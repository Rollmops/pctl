package app_test

import (
	"github.com/Rollmops/pctl/app"
	gopsutil "github.com/shirou/gopsutil/process"
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
	err = pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "stop", "Test1"})
	if err != nil {
		t.Fatal(err)
	}

	processes, err := gopsutil.Processes()
	if err != nil {
		t.Fatal(err)
	}
	for _, p := range processes {
		cmdline, _ := p.Cmdline()
		if cmdline == "sleep 10" {
			t.Fatal("Stopped process is still running")
		}
	}
}
