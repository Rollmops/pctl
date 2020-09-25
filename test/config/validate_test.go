package config_test

import "testing"
import "github.com/Rollmops/pctl/config"

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
