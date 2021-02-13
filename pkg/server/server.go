package server

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"

	"github.com/vatine/mspm/pkg/data"
	pb "github.com/vatine/mspm/pkg/protos"
)

type Server struct {
	pb.UnimplementedMspmServer
	dataStore *data.DataStore
}

// Convert a data.PackageVersion to a pb.PackageInformation, as this
// will be useful in a few cases.
func packageInformationFromPackageVersion(pv data.PackageVersion) *pb.PackageInformation {
	rv := new(pb.PackageInformation)

	rv.PackageName = pv.Name
	rv.Version = pv.Version
	for _, label := range pv.GetAllLabels() {
		rv.Label = append(rv.Label, label)
	}

	return rv
}

// Create a new Server data structure, populate it with paths to the
// playground (temp storage) and primary storage directories.
func NewServer(playground, store string) *Server {
	ds := data.NewDataStore(playground, store)

	return &Server{dataStore: ds}
}

// Set labels on a specific version of a package.
func (s *Server) SetLabels(ctx context.Context, in *pb.SetLabelRequest) (*pb.PackageInformation, error) {
	pkgName := in.GetPackageName()
	version := in.GetVersion()
	log.WithFields(log.Fields{
		"name":      pkgName,
		"version":   version,
		"newLabels": in.Label,
	}).Debug("SetLabels")

	if pkgName == "" {
		log.Error("SetLabels - missing package name")
		return nil, fmt.Errorf("No package name specified")
	}

	if version == "" {
		log.WithFields(log.Fields{
			"name": pkgName,
		}).Error("SetLabels - missing version designator")
		return nil, fmt.Errorf("No version designator specified")
	}

	var rErr error

	for ix, label := range in.GetLabel() {
		log.WithFields(log.Fields{
			"ix":      ix,
			"name":    pkgName,
			"version": version,
			"label":   label,
		}).Debug("SetLabels - setting label")
		err := s.dataStore.SetLabel(pkgName, version, label)
		if err != nil {
			rErr = err
			log.WithFields(log.Fields{
				"err":     err,
				"name":    pkgName,
				"version": version,
				"label":   label,
			}).Error("SetLabels - setting label")
		}
	}

	pv, err := s.dataStore.GetPackageVersion(pkgName, version)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"package": pkgName,
			"version": version,
		}).Error("SetLabels - unexpected missing")
		return nil, err
	}
	return packageInformationFromPackageVersion(pv), rErr
}

// Get information on a specific package.
func (s *Server) GetPackageInformation(ctx context.Context, in *pb.PackageInformationRequest) (*pb.PackageInformationResponse, error) {
	rv := new(pb.PackageInformationResponse)

	name := in.GetPackageName()
	if name == "" {
		return rv, fmt.Errorf("No package named %s", name)
	}
	pvs, ok := s.dataStore.GetPackageVersions(name)
	if !ok {
		log.WithFields(log.Fields{
			"name": name,
		}).Warning("GetPackageInformation - package not found")
	}

	for _, pv := range pvs {
		rv.PackageData = append(rv.PackageData, packageInformationFromPackageVersion(pv))
	}

	return rv, nil
}

// Receive a new version of a package.
func (s *Server) UploadPackage(ctx context.Context, in *pb.NewPackage) (*pb.PackageInformation, error) {
	name := in.GetPackageName()
	if name == "" {
		return nil, fmt.Errorf("No package name specified.")
	}

	pv, err := s.dataStore.NewPackageVersion(name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  name,
		}).Error("UploadPackage")
		return nil, err
	}

	for _, file := range in.GetFiles() {
		if strings.HasSuffix(file.GetName(), "/") {
			pv.AddDir(file)
		}
	}

	return packageInformationFromPackageVersion(pv), nil
}

// Send a specific version of a package to a client.
func (s *Server) GetPackage(ctx context.Context, in *pb.GetPackageRequest) (*pb.GetPackageResponse, error) {
	name := in.GetPackageName()
	labelDes := in.GetDesignator()

	if name == "" {
		log.Error("GetPackage called with empty name")
		return nil, fmt.Errorf("No name specified")
	}
	if labelDes == "" {
		log.WithFields(log.Fields{
			"name": name,
		}).Error("GetPackage, blank version designator")
		return nil, fmt.Errorf("No designator specified")
	}

	pv, err := s.dataStore.GetPackageVersion(name, labelDes)
	if err != nil {
		log.WithFields(log.Fields{
			"error":      err,
			"name":       name,
			"designator": labelDes,
		}).Error("GetPackage fetching packageversion")
		return nil, err
	}

	resp := pb.GetPackageResponse{
		PackageData: packageInformationFromPackageVersion(pv),
	}

	return &resp, nil
}
