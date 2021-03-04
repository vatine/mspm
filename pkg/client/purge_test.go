package client

import (
	"testing"
)

func TestSplit(t *testing.T) {
	testcases := []struct {
		path    string
		pkg     string
		version string
		err     bool
	}{
		{"blah/blah-beef", "blah", "beef", false},
		{"blah/blah-beef", "blub", "", true},
	}

	for ix, tc := range testcases {
		_, seenVersion, err := splitPackageVersion(tc.path, tc.pkg)
		if err != nil {
			if !tc.err {
				t.Errorf("Case #%d, saw error %v, none expcted.", ix, err)
			}
		} else {
			if tc.err {
				t.Errorf("Case #%d, error expected, none seen.", ix)
			}
			if tc.version != seenVersion {
				t.Errorf("Case #%d, saw version %s, want %s", ix, seenVersion, tc.version)
			}
		}
	}
}

func TestListAll(t *testing.T) {
	t1, e1 := listAllPossibleVersions("./testdata", "fake1")
	if e1 != nil {
		t.Errorf("Unexepected error %v", e1)
		return
	}

	expected := []string{"deadbeef", "feedf00d"}
	for ix, seen := range t1 {
		want := expected[ix]
		if seen != want {
			t.Errorf("Postition #%d, saw %s, want %s", ix, seen, want)
		}
	}

	t2, e2 := listAllPossibleVersions("./testdata", "fake2")
	if e2 != nil {
		// This is expected
		return
	}

	if len(t2) > 0 {
		t.Errorf("Elements returned, %v. None were expected.", t2)
	}
}
