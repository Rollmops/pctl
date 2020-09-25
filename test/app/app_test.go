package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/test"
	"os"
	"path"
	"testing"
)

var testDataDir string

func init() {
	cwd, _ := os.Getwd()
	testDataDir = path.Join(cwd, "..", "fixtures", "integration.yaml")
	_ = os.Setenv("PCTL_CONFIG_PATH", testDataDir)
}

func TestStartStopCommand(t *testing.T) {
	pctlApp := app.CreateCliApp()

	err := pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "Test1"})
	if err != nil {
		t.Fatal(err)
	}

	if !test.IsCommandRunning("sleep 1234") {
		t.Fatal("'sleep 1234' should be running")
	}

	err = pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "stop", "Test1"})
	if err != nil {
		t.Fatal(err)
	}

	if test.IsCommandRunning("sleep 1234") {
		t.Fatal("'sleep 1234' should be stopped")
	}
}
