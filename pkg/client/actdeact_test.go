package client

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"google.golang.org/grpc"

	pb "github.com/vatine/mspm/pkg/protos"
)

// The fakeActDeactServer is designed to do two things:
//  1. Help us juggle a temp directory for the testing
//  2. Give us desired responses to package information responses
type fakeActDeactServer struct {
	pb.UnimplementedMspmServer
	tmpDir string
	pvMap  map[string]map[string][]string
}

func newFakeActDeact(t *testing.T) *fakeActDeactServer {
	var rv fakeActDeactServer

	dirname, err := ioutil.TempDir("./tempdir", "actdeact")
	if err != nil {
		t.Fatalf("Failed to create test directory")
	}

	rv.tmpDir = dirname
	rv.pvMap = make(map[string]map[string][]string)

	return &rv
}

func (f *fakeActDeactServer) SetLabels(ctx context.Context, in *pb.SetLabelRequest, opts ...grpc.CallOption) (*pb.PackageInformation, error) {
	return nil, nil
}

func (f *fakeActDeactServer) GetPackageInformation(ctx context.Context, in *pb.PackageInformationRequest, opts ...grpc.CallOption) (*pb.PackageInformationResponse, error) {
	name := in.GetPackageName()
	rv := new(pb.PackageInformationResponse)

	submap, ok := f.pvMap[name]
	if ok {
		for version, labels := range submap {
			tmp := new(pb.PackageInformation)
			tmp.PackageName = name
			tmp.Version = version
			tmp.Label = labels
			rv.PackageData = append(rv.PackageData, tmp)
		}
	}

	return rv, nil
}

func (f *fakeActDeactServer) UploadPackage(ctx context.Context, in *pb.NewPackage, opts ...grpc.CallOption) (*pb.PackageInformation, error) {
	return nil, nil
}

func (f *fakeActDeactServer) GetPackage(ctx context.Context, in *pb.GetPackageRequest, opts ...grpc.CallOption) (*pb.GetPackageResponse, error) {
	return nil, nil
}

func (f *fakeActDeactServer) addPackage(name, version string, labels ...string) {
	fullName := fmt.Sprintf("%s-%s", name, version)
	os.Mkdir(path.Join(f.tmpDir, fullName), 0777)

	subMap, ok := f.pvMap[name]
	if !ok {
		subMap = make(map[string][]string)
		f.pvMap[name] = subMap
	}
	f.pvMap[name][version] = labels
}

func (f *fakeActDeactServer) tearDown() {
	os.RemoveAll(f.tmpDir)
}

func TestActivateWrongPackage(t *testing.T) {
	fs := newFakeActDeact(t)
	c := new(Client)
	c.client = fs
	c.mspmDir = fs.tmpDir
	defer fs.tearDown()

	fs.addPackage("bob", "deadbeef", "latest")

	err := c.Activate("alice", "latest")

	if err != nil {
	} else {
		t.Errorf("Expected error, saw none")
	}
}

func TestActivateRightPackage(t *testing.T) {
	fs := newFakeActDeact(t)
	c := new(Client)
	c.client = fs
	c.mspmDir = fs.tmpDir
	defer fs.tearDown()

	fs.addPackage("bob", "deadbeef", "latest")
	fs.addPackage("bob", "f00dbeef", "banjo", "kazooie")

	testcases := []struct {
		label string
		err   bool
	}{
		{"latest", false}, {"deadbeef", false}, {"f00dbeef", false},
		{"banjo", false}, {"dexter", true}, {"kazooie", false},
	}

	for ix, tc := range testcases {
		err := c.Activate("bob", tc.label)

		switch {
		case err != nil && !tc.err:
			t.Errorf("Case #%d, saw error %s, expected none", ix, err)
		case err != nil && tc.err:
			// All good
		case tc.err:
			t.Errorf("Case #%d, saw no error, expected one.", ix)
		}
	}
}
