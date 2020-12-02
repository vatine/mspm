// Various internal data structures for MSPM
package data

import (
	"fmt"
	"sync"

	log "github.com/sirupsen/logrus"
)

// Represents a general MSPM package (that is, all versions and labels).
type Package struct {
	lock sync.Mutex
	name string
	versions map[string]*PackageVersion
	labels map[string]*PackageVersion
}

// Data for a specific version of a package.
type PackageVersion struct {
	Name string
	Version string
	Labels map[string]struct{}
	DataPath string
	fileMap map[string]fileInfo
}

type fileInfo struct {
	owner string
	mode  int32
}


type DataStore struct {
	lock sync.Mutex
	playground string
	store string
	packages map[string]*Package
}


// Set the label newLabel on the package-version designated by
// designator. This can either be the version number (a hash of the
// contents) or a label identifying the package you want to modify.
// If the label is alreday in use for a package-version, it will be
// removed from it.
func (p *Package) SetLabel(designator, newLabel string) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	return p.setLabel(designator, newLabel)
}

// Return the PackageVersion that corresponds to the requested
// designatoir, as well as a bool that is true if tereturn hre is a
// PackageVersion that corresponds to the requested designator.
//
// For this purpose, a "designator" is either the version string, or a label.
func (p *Package) GetVersion(designator string) (PackageVersion, bool) {
	p.lock.Lock()
	defer p.lock.Unlock()

	pv, ok := p.getVersion(designator)
	return *pv, ok
}

// Internal versio of GetVersion, this returns a *PackageVersion and
// expects to be called in an already-locked Package.
func (p *Package) getVersion(designator string) (*PackageVersion, bool) {
	pv, ok := p.versions[designator]
	if !ok {
		pv, ok = p.labels[designator]
	}

	return pv, ok
}

// Internal version of SetLabel, that doesn't perform any
// locking. This means it's safe to call from (and only from)
// functions that already hold the lock for a package.
func (p *Package) setLabel(designator, newLabel string) error {
	target, ok := p.getVersion(designator)
	if !ok {
		return fmt.Errorf("No package-version designated by %s", designator)
	}

	old, ok := p.labels[newLabel]
	if ok && old != target {
		// There is a package that has this label, let us
		// immediately get rid of it, so we do not have any
		// conflicts.
		delete(old.Labels, newLabel)
	}
	target.Labels[newLabel] = struct{}{}
	p.labels[newLabel] = target
	log.WithFields(log.Fields{
		"target": target,
		"old": old,
		"labels": target.Labels,
	}).Debug("setting labels")
	return nil
}

// Add a specific version of a Package. This will move the label
// "latest" to point at the new addition.
func (p *Package) AddVersion(pv PackageVersion) error {
	p.lock.Lock()
	defer p.lock.Unlock()

	version := pv.Version
	_, ok := p.versions[version]
	if ok {
		// Never modify a package...
		return fmt.Errorf("Package %s already has a version %s", pv.Name, version)
	}
	p.versions[version] = &pv
	return p.setLabel(version, "latest")
}

func newPackage(name string) *Package {
	p := new(Package)
	p.name = name
	p.versions = make(map[string]*PackageVersion)
	p.labels = make(map[string]*PackageVersion)

	return p
}

// Add a PackageVersion to the data store. As we already have the
// package name and version detail(s), we don't take them as extra
// parameters. If we happen to already have the specific version
// stored, we do nothing.
func (ds *DataStore) AddPackageVersion(pv PackageVersion) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	p, ok := ds.packages[pv.Name]

	if !ok {
		// No previous version of this package added, make it so
		p = newPackage(pv.Name)
		ds.packages[pv.Name] = p	}

	p.AddVersion(pv)
}

func (ds *DataStore) GetPackageVersion(pkg, designator string) (PackageVersion, error) {
	ds.lock.Lock()
	defer ds.lock.Unlock()

	p, ok := ds.packages[pkg]
	if !ok {
		// No package, we have an error!
		return PackageVersion{}, fmt.Errorf("No package named %s, designated %s", pkg, designator)
	}

	pv, ok := p.getVersion(designator)
	if !ok {
		return PackageVersion{}, fmt.Errorf("No package named %s, designated %s", pkg, designator)
	}

	return *pv, nil
}
