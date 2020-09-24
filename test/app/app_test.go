package app_test

import (
	"github.com/Rollmops/pctl/app"
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

func TestStartCommand(t *testing.T) {
	pctlApp := app.CreateCliApp()

	err := pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "Test1"})
	if err != nil {
		t.Fatal(err)
	}
}
