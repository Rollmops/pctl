package config_test

import (
	"github.com/Rollmops/pctl/config"
	"os"
	"path"
	"path/filepath"
	"testing"
)

var testDataDir string

func init() {
	cwd, _ := os.Getwd()
	testDataDir = path.Join(cwd, "fixtures")
	_ = os.Setenv("TEST_DATA_DIR", testDataDir)
}

func TestLoadConfigOk(t *testing.T) {
	testConfigPath := path.Join(testDataDir, "pctl.yml")
	yamlLoader := config.NewYamlLoader()

	_config, _ := yamlLoader.Load(testConfigPath)

	if processCount := len(_config.Processes); processCount != 2 {
		t.Fatalf("Expected process count of 2, got %d", processCount)
	}
}

func TestCircleInclude(t *testing.T) {
	testConfigPath := path.Join(testDataDir, "pctl_circle_include.yml")
	yamlLoader := config.NewYamlLoader()

	_, err := yamlLoader.Load(testConfigPath)

	if err == nil {
		t.Fatalf("Expected error but got nil")
	}
}

func TestLoadConfigGlobIncludes(t *testing.T) {
	testConfigPath := path.Join(testDataDir, "glob_test.yml")
	yamlLoader := config.NewYamlLoader()

	_config, _ := yamlLoader.Load(testConfigPath)

	if processCount := len(_config.Processes); processCount != 3 {
		t.Fatalf("Expected process count of 3, got %d", processCount)
	}
}

func TestAbsPathLearning(t *testing.T) {
	absPath, _ := filepath.Abs("~")
	homePath := os.Getenv("HOME")

	if absPath == homePath {
		t.Fatalf("%s == %s", absPath, homePath)
	}
}
