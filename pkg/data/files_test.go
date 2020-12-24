package data

import (
	"fmt"
	"testing"
)

func TestModes(t *testing.T) {
	cases := []struct {
		m int32
		e string
	}{
		{0752, "rwxr-x-w-"}, {1, "--------x"},
	}

	for ix, c := range cases {
		got := modes(c.m)
		want := c.e
		if got != want {
			t.Errorf("Case #%d, got %s, want %s", ix, got, want)
		}
	}

}

func TestFileInfo(t *testing.T) {
	cases := []struct {
		f fileInfo
		e string
	}{
		{fileInfo{"owner", 0777}, "«owner»«rwxrwxrwx»"},
		{fileInfo{"owner", 0525}, "«owner»«r-x-w-r-x»"},
	}

	for ix, c := range cases {
		got := c.f.forHash()
		want := c.e
		if got != want {
			t.Errorf("Case #%d, got %s, want %s", ix, got, want)
		}
	}
}

func TestListFiles(t *testing.T) {
	paths, _ := pathsUnderRoot("./testdata/hash")

	expected := []string{"dir1/", "dir1/f1", "dir1/f2", "dir2/", "dir2/dir21/", "dir2/f1", "dir3/"}

	if len(expected) != len(paths) {
		t.Errorf("Mismatched path lengths, saw %v, want %v", paths, expected)
		return
	}
	for ix, want := range expected {
		got := paths[ix]
		if got != want {
			t.Errorf("Path element %d, saw %s want %s", ix, got, want)
		}
	}
}

func TestHash(t *testing.T) {
	pv := &PackageVersion{
		Name:     "testpackage",
		DataPath: "./testdata/hash",
		Labels:   make(map[string]struct{}),
		fileMap:  make(map[string]fileInfo),
	}
	pv.fileMap["dir1/"] = fileInfo{"root", 0755}
	pv.fileMap["dir2/"] = fileInfo{"root", 0755}
	pv.fileMap["dir2/dir21/"] = fileInfo{"root", 0755}
	pv.fileMap["dir3/"] = fileInfo{"root", 0755}
	pv.fileMap["dir1/f1"] = fileInfo{"root", 0644}
	pv.fileMap["dir1/f2"] = fileInfo{"root", 0644}
	pv.fileMap["dir2/f1"] = fileInfo{"root", 0644}

	h, err := pv.hash()
	if err != nil {
		t.Errorf("Failed to hash, saw error %s", err)
	}
	got := fmt.Sprintf("%x", h)
	want := "3c855f31642549cd0edc9af183c56a61a16da04df477f3c8697f51d82a37ac11d814bc37de9e34a0031ce1450f41ba039053fa52f9a252fc164ffb0c3da286ce"

	if got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}
