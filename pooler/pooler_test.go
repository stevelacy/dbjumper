package pooler

import (
	"errors"
	"github.com/stevelacy/dbjumper/pkg"
	"testing"
)

func TestGetRemote(t *testing.T) {
	cfg := dbjumper.Config{
		ListenAddress: "127.0.0.1:5432",
		Instances:     map[string]dbjumper.Instance{},
	}
	inst1 := dbjumper.Instance{
		ConnCount: 40,
		Address:   "127.0.0.1:1234",
	}
	inst2 := dbjumper.Instance{
		ConnCount: 5,
		Address:   "127.0.0.1:1235",
	}
	inst3 := dbjumper.Instance{
		ConnCount: 8,
		Address:   "127.0.0.1:1236",
	}

	cfg.Instances["inst1"] = inst1
	cfg.Instances["inst2"] = inst2
	cfg.Instances["inst3"] = inst3
	found, err := getAvailableInstance(&cfg)
	if err != nil {
		t.Error(err)
	}
	if found.Address != inst2.Address {
		t.Error(errors.New("incorrect instance returned"))
	}
}

func TestGetRemoteNone(t *testing.T) {
	cfg := dbjumper.Config{
		ListenAddress: "127.0.0.1:5432",
		Instances:     map[string]dbjumper.Instance{},
	}
	_, err := getAvailableInstance(&cfg)
	if err.Error() != "0 instances configured" {
		t.Error("Incorrect number of instances returned")
	}
}
