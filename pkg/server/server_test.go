package server

import (
	"testing"

	pb "github.com/vatine/mspm/pkg/protos"
)

func TestAssignment(t *testing.T) {
	s := Server{}
	var ms pb.MspmServer
	ms = &s
	if _, ok := ms.(*Server); !ok {
		t.Errorf("Failed to convert")
	}
}
