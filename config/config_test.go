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
	testDataDir = path.Join(cwd, "..", "test_data")
	os.Setenv("TEST_DATA_DIR", testDataDir)

}

func TestLoadConfigOk(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "pctl.yml")

	config, _ := config.LoadConfig(testConfigPath)

	if processCount := len(config.Processes); processCount != 2 {
		t.Fatalf("Expected process count of 2, got %d", processCount)
	}
}

func TestCircleInclude(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "pctl_circle_include.yml")

	_, err := config.LoadConfig(testConfigPath)

	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

}

func TestLoadConfigGlobIncludes(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "glob_test.yml")

	config, _ := config.LoadConfig(testConfigPath)

	if processCount := len(config.Processes); processCount != 3 {
		t.Fatalf("Expected process count of 3, got %d", processCount)
	}
}

func TestDuplicateProcessNames(t *testing.T) {
	setUp()
	testConfigPath := path.Join(testDataDir, "duplicate_process_name.yml")

	_, err := config.LoadConfig(testConfigPath)

	if err == nil {
		t.Fatalf("Expected duplicate process name error")
	}
}

func TestAbsPathLearning(t *testing.T) {
	absPath, _ := filepath.Abs("~")
	homePath := os.Getenv("HOME")

	if absPath == homePath {
		t.Fatalf("%s == %s", absPath, homePath)
	}

}
