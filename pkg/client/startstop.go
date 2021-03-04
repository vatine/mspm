package client

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Run a command (we expect this to be one of "start" or "stop"), if
// the path doesn't exist this isn't an error, it's "just" a library
// package.
func run(path string, out, errOut io.Writer) error {
	_, err := os.Stat(path)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path":  path,
		}).Info("run, non-existing path")
		return nil
	}

	cmd := exec.Command(path)
	cmd.Stdout = out
	cmd.Stderr = errOut

	return cmd.Run()
}

// Start the activated version of a package
func (c *Client) Start(pkgName string) error {
	log.WithFields(log.Fields{
		"package": pkgName,
	}).Debug("starting package")

	pkgSymlink := filepath.Join(c.mspmDir, pkgName)
	if !checkSymlink(pkgSymlink) {
		// This should be a symlink....
		err := fmt.Errorf("%s is supposed to be a symlink", pkgSymlink)
		log.WithFields(log.Fields{
			"package": pkgName,
			"pkgDir":  pkgSymlink,
			"action":  "directory-probe",
		}).Error("start")

		return err
	}

	startLink := filepath.Join(pkgSymlink, "start")
	return run(startLink, os.Stdout, os.Stderr)
}

// Stop the activated version of a package
func (c *Client) Stop(pkgName string) error {
	log.WithFields(log.Fields{
		"package": pkgName,
	}).Debug("stopping package")

	pkgSymlink := filepath.Join(c.mspmDir, pkgName)
	if !checkSymlink(pkgSymlink) {
		// This should be a symlink....
		err := fmt.Errorf("%s is supposed to be a symlink", pkgSymlink)
		log.WithFields(log.Fields{
			"package": pkgName,
			"pkgDir":  pkgSymlink,
			"action":  "directory-probe",
		}).Error("stop")

		return err
	}

	stopLink := filepath.Join(pkgSymlink, "stop")
	return run(stopLink, os.Stdout, os.Stderr)
}
