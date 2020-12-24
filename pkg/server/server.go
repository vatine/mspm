package server

import (
	"context"
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	
	"github.com/vatine/mspm/pkg/data"
	pb "github.com/vatine/mspm/pkg/protos"
)

type Server struct {
	pb.UnimplementedMspmServer
	dataStore *data.DataStore
}

func packageInformationFromPackageVersion(pv data.PackageVersion) *pb.PackageInformation {
	rv := new(pb.PackageInformation)

	rv.PackageName = pv.Name
	rv.Version = pv.Version
	for _, label := range pv.GetAllLabels() {
		rv.Label = append(rv.Label, label)
	}

	return rv
}

func NewServer(playground, store string) *Server {
	ds := data.NewDataStore(playground, store)

	return &Server{dataStore: ds}
}

func (s *Server) SetLabels(ctx context.Context, in *pb.SetLabelRequest, opts ...grpc.CallOption) (*pb.PackageInformation, error) {
	pkgName := in.GetPackageName()
	version := in.GetVersion()
	log.WithFields(log.Fields{
		"name": pkgName,
		"version": version,
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
			"ix": ix,
			"name": pkgName,
			"version": version,
			"label": label,
		}).Debug("SetLabels - setting label")
		err := s.dataStore.SetLabel(pkgName, version, label)
		if err != nil {
			rErr = err
			log.WithFields(log.Fields{
				"err": err,
				"name": pkgName,
				"version": version,
				"label": label,
			}).Error("SetLabels - setting label")
		}
	}

	pv, err := s.dataStore.GetPackageVersion(pkgName, version)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"package": pkgName,
			"version": version,
		}).Error("SetLabels - unexpected missing")
		return nil, err
	}
	return packageInformationFromPackageVersion(pv), rErr
}

func (s *Server) GetPackageInformation(in *pb.PackageInformationRequest, opts ...grpc.CallOption) (*pb.PackageInformationResponse, error) {
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

func (s *Server) UploadPackage(in *pb.NewPackage, opts ...grpc.CallOption) (*pb.PackageInformation, error) {
	name := in.GetPackageName()
	if name == "" {
		return nil, fmt.Errorf("No package name specified.")
	}

	pv, err := s.dataStore.NewPackageVersion(name)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name": name,
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
