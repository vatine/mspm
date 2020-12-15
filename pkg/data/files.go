// This provides the file interface of our backend storage. What we
// expect here is atht we call NewPackageVersion when we have a new
// package, then we fill it with files and eventually we call
// pv.Done(), at which point we compute the aggregate hash sum, stamp
// it with a version and add it as the "latest" version to the main
// package.
package data

import (
	"archive/tar"
	"crypto/sha512"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"

	log "github.com/sirupsen/logrus"
	
	pb "github.com/vatine/mspm/pkg/protos"
)

func (ds DataStore) NewPackageVersion(name string) (PackageVersion, error) {
	tdPath := filepath.Join(ds.playground, "tmp", name)
	dataPath, err := ioutil.TempDir(tdPath, "tmp-")
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"tdPath": tdPath,
			"name": name,
		}).Error("creating tempdir")
	}
	return PackageVersion{
		Name: name,
		Labels: make(map[string]struct{}),
		DataPath: dataPath,
		fileMap: make(map[string]fileInfo),
	}, err
}

func (pv *PackageVersion) AddDir(pvFile pb.File) error {
	targetPath := filepath.Join(pv.DataPath, pvFile.Name)

	pv.fileMap[pvFile.Name] = fileInfo{pvFile.Owner, pvFile.Mode}
	return os.Mkdir(targetPath, os.ModeDir | 0777)
}

// Add a file to the on-disk temporary storage of a file.
func (pv PackageVersion) AddFile(pvFile pb.File) error {
	targetPath := filepath.Join(pv.DataPath, pvFile.Name)

	out, err := os.Create(targetPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name": pv.Name,
			"targetPath": targetPath,
		}).Error("opening PackageVersion file")
		return err
	}

	written := 0
	for written < len(pvFile.Contents) {
		n, err := out.Write(pvFile.Contents[written:])
		if err != nil {
			if n > 0 {
				log.WithFields(log.Fields{
					"error": err,
					"name": pv.Name,
					"n": n,
					"written": written,
				}).Warning("non-full write")
				written += n
				continue
			}
			log.WithFields(log.Fields{
				"error": err,
				"name": pv.Name,
				"written": written,
			}).Error("failed writing data")
			return err
		}
		written += n
	}

	pv.fileMap[pvFile.Name] = fileInfo{pvFile.Owner, pvFile.Mode}
	return nil
}

// Returns a slice of os.FileInfo, sorted asciibetically after name
func pathsUnderRoot(root string) ([]string, error) {
	rv, err := pathsUnderRootInternal(root, "")

	return rv, err
}

func pathsUnderRootInternal(root, offset string) ([]string, error) {
	target := filepath.Join(root, offset)

	fi, err := os.Stat(target)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"root": root,
			"offset": offset,
			"target": target,
		}).Error("stating initial directory")
		return []string{}, err
	}

	if !fi.IsDir() {
		log.WithFields(log.Fields{
			"target": target,
		}).Warning("unexpected non-directory")
		return []string{offset}, nil
	}

	var rv []string
	dir, err := os.Open(target)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"root": root,
			"offset": offset,
			"target": target,
		}).Error("opening initial directory")
		return rv, err
	}
	defer dir.Close()

	names, err := dir.Readdirnames(-1)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"target": target,
		}).Error("listing files")
	}

	sort.Strings(names)
	for _, name := range names {
		fi, err := os.Stat(filepath.Join(target, name))
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"target": target,
				"name": name,
			}).Error("stating internal file")
			return rv, err
		}
		if fi.IsDir() {
			thisName := fmt.Sprintf("%s/", filepath.Join(offset, name))
			if name != "." && name != ".." {
				rv = append(rv, thisName)
				internals, err := pathsUnderRootInternal(root, thisName)
				rv = append(rv, internals...)
				if err != nil {
					return rv, nil
				}
			}
		} else {
			rv = append(rv, filepath.Join(offset, name))
		}
	}

	return rv, nil
}

// Return a "mode-bit" representation for 3 bits
func mode(in int32) string {
	rv := []byte("---")
	data := []struct{m int32; s byte}{{4, 'r'},{2, 'w'},{1, 'x'}}

	for i, d := range data {
		if (in & d.m) == d.m {
			rv[i] = d.s
		}
	}

	return string(rv)
}

// Return a string that is the textual representation of a file-mode
func modes(i int32) string {
	return fmt.Sprintf("%s%s%s", mode(i >> 6), mode(i >> 3), mode(i))
}

func (f fileInfo) forHash() string {
	return fmt.Sprintf("«%s»«%s»", f.owner, modes(f.mode))
	
}

func (pv *PackageVersion) hash() ([]byte, error) {
	hash := sha512.New()

	paths, err := pathsUnderRoot(pv.DataPath)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"path": pv.DataPath,
		}).Error("listing files")
		return []byte{}, err
	}

	for ix, name := range paths {
		fi, ok := pv.fileMap[name]
		if !ok {
			log.WithFields(log.Fields{
				"pv.Name": pv.Name,
				"name": name,
			}).Error("missing fileIfo for name")
			return hash.Sum(nil), fmt.Errorf("File %s is unknown", name)
		}
		fmt.Fprintf(hash, "«%d»«%s»%s", ix, name, fi.forHash())
		if !strings.HasSuffix(name, "/") {
			func() {
				f, err := os.Open(filepath.Join(pv.DataPath, name))
				if err != nil {
					return
				}
				defer f.Close()
				io.Copy(hash, f)
			}()
		}
	}

	return hash.Sum(nil), nil
}

func (pv *PackageVersion) Finish() error {
	var saveErr error
	hash, err := pv.hash()

	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name": pv.Name,
		}).Error("hashing packageVersion")
		return err
	}

	pv.Version = fmt.Sprintf("%x", hash)

	files, _ := pathsUnderRoot(pv.DataPath)
	outName := filepath.Join(pv.DataPath, "../..", fmt.Sprintf("%s-%s.tgz", pv.Name, pv.Version))

	out, err := os.Create(outName)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"pv": pv.Name,
			"version": pv.Version,
			"filename": outName,
		}).Error("opening archive file")
		return err
	}
	defer out.Close()

	zipper := gzip.NewWriter(out)

	tarball := tar.NewWriter(zipper)
	defer tarball.Close()

	for _, fname := range files {
		fsname := filepath.Join(pv.DataPath, fname)
		tarname := fmt.Sprintf("%s-%s/%s", pv.Name, pv.Version, fname)
		fi, err := os.Stat(fsname)
		if err != nil {
			log.WithFields(log.Fields{
				"error": err,
				"fsname": fsname,
			}).Error("archiving stat failed")
			return err
		}
		func () {
			pvMode := pv.fileMap[fname].mode & 0777
			tarHdr, _ := tar.FileInfoHeader(fi, "")
			tarHdr.Name = tarname
			tarHdr.Mode = (tarHdr.Mode & 0xFFFE00) | int64(pvMode)
			tarHdr.Uid = 0
			tarHdr.Gid = 0
			
			in, err := os.Open(fsname)
			if err != nil {
				log.WithFields(log.Fields{
					"error": err,
					"fsname": fsname,
				}).Error("failed to open input file")
			}
			defer in.Close()
			tarball.WriteHeader(tarHdr)
			io.Copy(tarball, in)
		}()
	}
	
	return saveErr
}
