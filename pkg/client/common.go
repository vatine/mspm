package client

import (
	"context"
	"fmt"
	"os"

	log "github.com/sirupsen/logrus"

	pb "github.com/vatine/mspm/pkg/protos"
)

// return true if the specified path points at a symlink.
func checkSymlink(dirPath string) bool {
	stat, err := os.Lstat(dirPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  dirPath,
		}).Error("checkSymlink")
		return false
	}
	return (stat.Mode() & os.ModeSymlink) == os.ModeSymlink
}

// Construct a full map of label -> version for a package
func (c *Client) fullLabelToVersion(pkgName string) (map[string]string, error) {
	rv := make(map[string]string)
	req := pb.PackageInformationRequest{PackageName: pkgName}
	resp, err := c.client.GetPackageInformation(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"pkgName": pkgName,
		}).Error("activate gRPC call")
		return rv, err
	}

	for _, pkgInfo := range resp.GetPackageData() {
		version := pkgInfo.GetVersion()
		for _, pkgLabel := range pkgInfo.GetLabel() {
			rv[pkgLabel] = version
		}
	}

	return rv, nil
}

// Return the version of a package that corresponds to a label/version.
func (c *Client) matchLabelToVersion(pkgName, label string) (string, error) {
	req := pb.PackageInformationRequest{PackageName: pkgName}
	resp, err := c.client.GetPackageInformation(context.Background(), &req)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"pkgName": pkgName,
			"label":   label,
		}).Error("activate gRPC call")
		return "", err
	}

	var version string

	for _, pkgInfo := range resp.GetPackageData() {
		if label == pkgInfo.GetVersion() {
			version = pkgInfo.GetVersion()
			break
		}
		for _, pkgLabel := range pkgInfo.GetLabel() {
			if label == pkgLabel {
				version = pkgInfo.GetVersion()
				break
			}
		}
	}

	if version == "" {
		log.WithFields(log.Fields{
			"pkgName": pkgName,
			"label":   label,
		}).Error("activate label/version not found")

		return "", fmt.Errorf("package %s version/label %s not found", pkgName, label)
	}

	return version, nil
}
