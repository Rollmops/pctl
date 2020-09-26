package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"os"
	"path"
	"testing"
)

func _loadCheckTestConfig(name string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	configPath := path.Join(cwd, "..", "fixtures", "check_test", name)
	return os.Setenv("PCTL_CONFIG_PATH", configPath)
}

func TestCheckPersistenceConfigDiscrepancyStillRunning(t *testing.T) {
	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should not be running")
	pctlApp := app.CreateCliApp(os.Stdout)

	assert.NoError(t, _loadCheckTestConfig("step1.yaml"))

	err := pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "Test1"})
	assert.NoError(t, err)
	assert.True(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be running")

	assert.NoError(t, _loadCheckTestConfig("step2.yaml"))

	out := test.CaptureLogOutput(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "info"}))
	})

	assert.Regexp(t, "level=error\\s+msg=\"Found tracked running process 'Test1' with pid \\d+ that could not be "+
		"found in config\"", out)
}
