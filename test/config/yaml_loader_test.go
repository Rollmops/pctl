package config_test

import (
	"os"
	"path"
	"path/filepath"
	"testing"

	"github.com/Rollmops/pctl/config"
)

var testDataDir string

func setUp() {
	cwd, _ := os.Getwd()
	testDataDir = path.Join(cwd, "..", "..", "test_data")
	_ = os.Setenv("TEST_DATA_DIR", testDataDir)
}

func TestLoadConfigOk(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "pctl.yml")
	yamlLoader, _ := config.NewYamlLoader(testConfigPath)

	_config, _ := yamlLoader.Load()

	if processCount := len(_config.Processes); processCount != 2 {
		t.Fatalf("Expected process count of 2, got %d", processCount)
	}
}

func TestCircleInclude(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "pctl_circle_include.yml")
	yamlLoader, _ := config.NewYamlLoader(testConfigPath)

	_, err := yamlLoader.Load()

	if err == nil {
		t.Fatalf("Expected error but got nil")
	}
}

func TestLoadConfigGlobIncludes(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "glob_test.yml")
	yamlLoader, _ := config.NewYamlLoader(testConfigPath)

	_config, _ := yamlLoader.Load()

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
