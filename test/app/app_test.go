package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartStopCommand(t *testing.T) {
	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should not be running")
	pctlApp := app.CreateCliApp()

	err := pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "Test1"})
	assert.NoError(t, err)
	assert.True(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be running")

	err = pctlApp.Run([]string{"pctl", "--loglevel", "DEBUG", "stop", "--nowait", "Test1"})
	assert.NoError(t, err)
	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be stopped")
}

func TestStartWithDependencies(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("dependsOn.yaml"))

	pctlApp := app.CreateCliApp()

	out := test.CaptureStdout(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "--no-color-output", "start", "p1"}))
	})

	assert.True(t, test.IsCommandRunning("sleep 3456"), "'sleep 3456' should be running")
	assert.True(t, test.IsCommandRunning("sleep 4567"), "'sleep 4567' should be running")

	assert.Equal(t, out, `Starting dependency 'p2' ... Ok
Starting process 'p1' ... Ok
`)
}
