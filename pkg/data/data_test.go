package data

import (
	"fmt"
	"sync"
	"testing"
)

func newPackageVersion(name, version string) PackageVersion {
	return PackageVersion{
		Name:     name,
		Labels:   make(map[string]struct{}),
		DataPath: "dataPath",
		Version:  version,
	}
}

func TestSetLabel(t *testing.T) {
	// This actually tries to provoke data races

	var wg sync.WaitGroup
	var pvs []*PackageVersion
	p := newPackage("foo")
	c := make(chan struct{})

	tests := 30
	wg.Add(tests)

	// Set everything up for data racing
	for n := 0; n < tests; n++ {
		go func(n int) {
			<-c
			v := fmt.Sprintf("%d", n)
			d1 := fmt.Sprintf("v%d", n)
			d2 := fmt.Sprintf("v%d.0", n)
			err := p.SetLabel(v, d1)
			if err != nil {
				t.Errorf("n=%d, v=%s, d1=%s, failed to set label, err: %s", n, v, d1, err)
			}
			err = p.SetLabel(d1, d2)
			if err != nil {
				t.Errorf("n=%d, d1=%s, d2=%s, failed to set label, err: %s", n, d1, d2, err)
			}
			wg.Done()
		}(n)

		v := fmt.Sprintf("%d", n)
		pv := newPackageVersion("foo", v)
		pvs = append(pvs, &pv)
		p.versions[v] = &pv
	}

	// Unleash the goroutines

	close(c)
	wg.Wait()

	for n, pv := range pvs {
		wantV := fmt.Sprintf("%d", n)
		wantD1 := fmt.Sprintf("v%d", n)
		wantD2 := fmt.Sprintf("v%d.0", n)

		if pv.Version != wantV {
			t.Errorf("PV %d, version is %s, want %s", n, pv.Version, wantV)
		}
		if _, ok := pv.Labels[wantD1]; !ok {
			t.Errorf("PV %d, %v missing label %s", n, pv.Labels, wantD1)
		}
		if _, ok := pv.Labels[wantD2]; !ok {
			t.Errorf("PV %d, %v missing label %s", n, pv.Labels, wantD2)
		}
	}
}
