package client

import (
	"strings"
	"testing"
)

func TestCheckSymlink(t *testing.T) {
	testcases := []struct {
		path string
		ok   bool
	}{
		{"./testdata/fake1", true},
		{"./snurgle", false},
		{"./testdata/fake1-deadbeef", false},
	}

	for ix, tc := range testcases {
		saw := checkSymlink(tc.path)
		want := tc.ok

		if saw != want {
			t.Errorf("Case #%d, saw %v, want %v", ix, saw, want)
		}
	}
}

func TestRun(t *testing.T) {
	testcases := []struct {
		path    string
		wantOut string
		wantErr string
		errSeen bool
	}{
		{"./testdata/fake1/start", "start fake1-deadbeef\n", "", false},
		{"./testdata/fake1/stop", "stop fake1-deadbeef\n", "", false},
		{"./testdata/fake1-deadbeef/stop", "stop fake1-deadbeef\n", "", false},
		{"./testdata/fake1-deadbeef/doesnotexist", "", "", false},
		{"./testdata/fake1-feedf00d/start", "", "", true},
		{"./testdata/fake1-feedf00d/stop", "", "", true},
	}

	for ix, tc := range testcases {
		o := strings.Builder{}
		e := strings.Builder{}
		err := run(tc.path, &o, &e)
		errSeen := (err != nil)

		switch {
		case errSeen && !tc.errSeen:
			t.Errorf("Case #%d, saw error %s, expected none", ix, err)
		case !errSeen && tc.errSeen:
			t.Errorf("Case #%d, saw no error, expected one.", ix)
		}

		if out := o.String(); out != tc.wantOut {
			t.Errorf("Case #%d, output seen «%s», want «%s»", ix, out, tc.wantOut)
		}

		if err := e.String(); err != tc.wantErr {
			t.Errorf("Case #%d, output seen «%s», want «%s»", ix, err, tc.wantErr)
		}
	}
}
