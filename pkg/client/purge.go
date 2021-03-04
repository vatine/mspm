package client

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	log "github.com/sirupsen/logrus"
)

func splitPackageVersion(filePath, pkgName string) (string, string, error) {
	targetName := path.Base(filePath)
	if !strings.HasPrefix(targetName, pkgName) {
		err := fmt.Errorf("malformed symlink, expected target %s to start with %s", targetName, pkgName)
		log.WithFields(log.Fields{
			"error":      err,
			"pkgName":    pkgName,
			"targetName": targetName,
		}).Error("splitPackageVersion")
		return "", "", err
	}

	version := targetName[1+len(pkgName) : len(targetName)]

	return pkgName, version, nil
}

// Return the package and version of what a symlink is pointing at.
func resolveSymlink(link string) (string, string, error) {
	log.WithFields(log.Fields{
		"link": link,
	}).Debug("resolveSymlink called")
	if !checkSymlink(link) {
		err := fmt.Errorf("link %s is not a symlink", link)
		return "", "", err
	}

	resolved, err := os.Readlink(link)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"link":  link,
		}).Error("resolveSymlínk - os.ReadLink")
		return "", "", err
	}

	pkgName := path.Base(link)

	return splitPackageVersion(resolved, pkgName)
}

// Return a list of version IDs for a specific package.
func listAllPossibleVersions(base, pkgName string) ([]string, error) {
	log.WithFields(log.Fields{
		"base":    base,
		"pkgName": pkgName,
	}).Debug("listAllPossibleVersions")
	wildcard := fmt.Sprintf("%s/%s-*", base, pkgName)

	files, err := filepath.Glob(wildcard)

	if err != nil {
		log.WithFields(log.Fields{
			"error":   err,
			"base":    base,
			"pkgName": pkgName,
		}).Error("listAllPossibleVersions - glob")
		return []string{}, err
	}

	var rv []string
	for _, name := range files {
		_, version, err := splitPackageVersion(name, pkgName)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"name":  name,
			}).Error("listAllPossibleVersions ")
			continue
		}
		rv = append(rv, version)
	}

	return rv, nil
}

// Delete downloaded package/version, either "all but active", or
// "specified labels, not active".
func (c *Client) Purge(pkgName string, labels ...string) error {
	avoid := make(map[string]bool)
	zap := make(map[string]bool)
	pkgSymlink := path.Join(c.mspmDir, pkgName)

	if checkSymlink(pkgSymlink) {
		_, version, err := resolveSymlink(pkgSymlink)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"pkgName": pkgName,
				"symlink": pkgSymlink,
			}).Error("Purge, failed to resolve symlink")
			return err
		}

		avoid[version] = true
	}

	if len(labels) == 0 {
		allVersions, err := listAllPossibleVersions(c.mspmDir, pkgName)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"pkgName": pkgName,
			}).Error("Purge - failed to get existing versions")
			return err
		}
		for _, version := range allVersions {
			zap[version] = true
		}
	} else {
		versionMap, err := c.fullLabelToVersion(pkgName)
		if err != nil {
			log.WithFields(log.Fields{
				"error":   err,
				"pkgName": pkgName,
			}).Error("Purge - failed to get full label->version map")
			return err
		}

		for _, label := range labels {
			version, ok := versionMap[label]
			if ok {
				zap[version] = true
			} else {
				log.WithFields(log.Fields{
					"pkgName": pkgName,
					"label":   label,
				}).Warn("Purge - label is not set")
			}
		}
	}
	return c.purgeInner(pkgName, avoid, zap)
}

// Loop through the versions in zap, if there is a package with that
// version, delete it (unless it is in the avoid map).
func (c *Client) purgeInner(pkgName string, avoid, zap map[string]bool) error {
	log.WithFields(log.Fields{
		"pkgName": pkgName,
		"avoid":   avoid,
		"zap":     zap,
	}).Debug("purgeInner entered")
	for version, _ := range zap {
		dirName := fmt.Sprintf("%s/%s-%s", c.mspmDir, pkgName, version)
		if !avoid[version] {
			_, err := os.Stat(dirName)
			if err != nil {
				log.WithFields(log.Fields{
					"pkgName": pkgName,
					"version": version,
					"path":    dirName,
				}).Warning("purgeInner - unable to purge version")
				continue
			}

			err = os.RemoveAll(dirName)
			if err != nil {
				log.WithFields(log.Fields{
					"error":   err,
					"pkgName": pkgName,
					"version": version,
					"path":    dirName,
				}).Error("purgeInner - errored out deleting directory")
				return err
			}
		} else {
			log.WithFields(log.Fields{
				"pkgName": pkgName,
				"version": version,
				"path":    dirName,
			}).Info("purgeInner - this is the active version")
		}
	}

	return nil
}
