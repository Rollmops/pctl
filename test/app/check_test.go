package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCheckPersistenceConfigDiscrepancyStillRunning(t *testing.T) {
	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should not be running")
	pctlApp := app.CreateCliApp()

	assert.NoError(t, test.SetConfigEnvPath("check_test", "step1.yaml"))

	err := pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "Test1"})
	assert.NoError(t, err)
	assert.True(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be running")

	assert.NoError(t, test.SetConfigEnvPath("check_test", "step2.yaml"))

	out := test.CaptureLogOutput(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "info"}))
	})

	assert.Regexp(t, "level=error\\s+msg=\"Found tracked running process 'Test1' with pid \\d+ that could not be "+
		"found in config\"", out)
}
