package data

import (
	"testing"
)

func TestModes(t *testing.T) {
	cases := []struct{m int32; e string}{
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
	cases := []struct{
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
	paths, _ := pathsUnderRoot("./testdata")

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
