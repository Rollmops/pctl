package config_test

import (
	"github.com/Rollmops/pctl/config"
	"testing"
)

func TestFindByName(t *testing.T) {
	_config := config.Config{
		Processes: []config.ProcessConfig{
			{
				Name: "p1",
				Cmd:  []string{"sleep", "10"},
			},
			{
				Name: "p2",
				Cmd:  []string{"ls", "-la"},
			},
		},
	}

	if p := _config.FindByName("p1"); p.Cmd[0] != "sleep" {
		t.Fatalf("sleep != %s", p.Cmd[0])
	}
	if p := _config.FindByName("p2"); p.Cmd[1] != "-la" {
		t.Fatalf("-la != %s", p.Cmd[1])
	}
	if p := _config.FindByName("NOT_THERE"); p != nil {
		t.Fatal("Expected process config NOT_THERE to be nil")
	}

}
