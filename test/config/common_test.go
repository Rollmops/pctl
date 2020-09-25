package config_test

import (
	"github.com/Rollmops/pctl/config"
	"testing"
)

func TestFindByName(t *testing.T) {
	_config := config.Config{
		Processes: []*config.ProcessConfig{
			{
				Name:    "p1",
				Command: []string{"sleep", "10"},
			},
			{
				Name:    "p2",
				Command: []string{"ls", "-la"},
			},
		},
	}

	if p := _config.FindByName("p1"); p.Command[0] != "sleep" {
		t.Fatalf("sleep != %s", p.Command[0])
	}
	if p := _config.FindByName("p2"); p.Command[1] != "-la" {
		t.Fatalf("-la != %s", p.Command[1])
	}
	if p := _config.FindByName("NOT_THERE"); p != nil {
		t.Fatal("Expected process config NOT_THERE to be nil")
	}

}

func TestGetLoaderFromPathYaml(t *testing.T) {
	loader := config.GetLoaderFromPath("/path/to/config.yaml")
	if _, ok := loader.(*config.YamlLoader); !ok {
		t.Fatal("expected YamlLoader")
	}
}

func TestGetLoaderFromPathYml(t *testing.T) {
	loader := config.GetLoaderFromPath("/path/to/config.yml")
	if _, ok := loader.(*config.YamlLoader); !ok {
		t.Fatal("expected YamlLoader")
	}
}
