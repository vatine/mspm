package client

import (
	"fmt"
	"os"
	"path"

	log "github.com/sirupsen/logrus"
)

// Activate the package with a label (or version).
// If the requested version is already active, leave things as they are.
func (c *Client) Activate(pkgName, label string) error {
	log.WithFields(log.Fields{
		"package name":  pkgName,
		"label/version": label,
	}).Debug("activate entered")

	version, err := c.matchLabelToVersion(pkgName, label)
	if err != nil {
		return err
	}

	fullName := fmt.Sprintf("%s-%s", pkgName, version)
	fullPath := path.Join(c.mspmDir, fullName)
	pkgStat, err := os.Lstat(fullPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"name":    pkgName,
			"label":   label,
			"version": version,
			"path":    fullPath,
		}).Error("Activate - failed to stat package path")
		return err
	}
	if !pkgStat.Mode().IsDir() {
		err := fmt.Errorf("expected %s to be a directory, mode is %s", fullPath, pkgStat.Mode())
		log.WithFields(log.Fields{
			"error":   err,
			"name":    pkgName,
			"version": version,
			"path":    fullPath,
		}).Error("Activate - package path is not a directory")
		return err
	}

	linkPath := path.Join(c.mspmDir, pkgName)
	_, err = os.Lstat(linkPath)
	if err == nil {
		err := os.Remove(linkPath)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"linkPath": linkPath,
				"label":    label,
			}).Error("Activate - failed removing symlink")
		}

	}

	linkPath2 := fmt.Sprintf("./%s", fullName)
	err = os.Symlink(linkPath2, linkPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":    err,
			"fullPath": fullPath,
			"linkPath": linkPath,
			"label":    label,
		}).Error("Activate - failed to create symlink")
	}
	return err
}

// Deactivate the package with a label (or version).
func (c *Client) Deactivate(pkgName, label string) error {
	log.WithFields(log.Fields{
		"package name":  pkgName,
		"label/version": label,
	}).Debug("deactivate entered")

	version, err := c.matchLabelToVersion(pkgName, label)
	if err != nil {
		return err
	}

	fullName := fmt.Sprintf("%s-%s", pkgName, version)
	fullPath := path.Join(c.mspmDir, fullName)
	pkgStat, err := os.Lstat(fullPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"name":    pkgName,
			"label":   label,
			"version": version,
			"path":    fullPath,
		}).Error("Deactivate - failed to stat package path")
		return err
	}
	if !pkgStat.Mode().IsDir() {
		err := fmt.Errorf("expected %s to be a directory, mode is %s", fullPath, pkgStat.Mode())
		log.WithFields(log.Fields{
			"error":   err,
			"name":    pkgName,
			"version": version,
			"path":    fullPath,
		}).Error("Deactivate - package path is not a directory")
		return err
	}

	linkPath := path.Join(c.mspmDir, pkgName)
	_, err = os.Lstat(linkPath)
	if err == nil {
		err := os.Remove(linkPath)
		if err != nil {
			log.WithFields(log.Fields{
				"error":    err,
				"linkPath": linkPath,
				"label":    label,
			}).Error("Deactivate - failed removing symlink")
		}
		return err
	}

	return nil
}
