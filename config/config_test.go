package config_test

import (
	"github.com/Rollmops/pctl/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

type TestLoader struct{}

func (t *TestLoader) Load(_ string) (*config.Config, error) {
	return nil, nil
}

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
	config.SuffixConfigLoaderMap["sfx"] = &TestLoader{}
	loader := config.GetLoaderFromPath("/path/to/config.sfx")
	assert.IsType(t, (*TestLoader)(nil), loader)
}

func TestValidateConfig(t *testing.T) {
	_config := config.Config{
		Processes: []*config.ProcessConfig{
			{
				Name:    "p1",
				Command: []string{"sleep 1"},
			},
			{
				Name:    "p1",
				Command: []string{"sleep 2"},
			},
		},
	}

	err := _config.Validate()

	if err == nil {
		t.Fatal("Expected failing _config validation")
	}
}

func TestInvalidCmdLength(t *testing.T) {
	_config := config.Config{
		Processes: []*config.ProcessConfig{
			{
				Name:    "p1",
				Command: []string{},
			},
		},
	}

	err := _config.Validate()
	if err == nil {
		t.Fatal("Expect error for process cmd length == 0")
	}

}
