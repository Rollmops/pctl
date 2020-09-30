package app_test

import (
	"encoding/json"
	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStartStopCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))

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
		assert.NoError(t, pctlApp.Run([]string{"pctl", "--no-color-Output", "start", "p1"}))
	})

	assert.True(t, test.IsCommandRunning("sleep 3456"), "'sleep 3456' should be running")
	assert.True(t, test.IsCommandRunning("sleep 4567"), "'sleep 4567' should be running")

	assert.Equal(t, out, `Starting dependency 'p2' ... Ok
Starting process 'p1' ... Ok
`)
}

func TestAppInfoCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))
	pctlApp := app.CreateCliApp()
	out := test.CaptureStdout(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "info"}))
	})
	expectedInfoOut := `Test1: ["sleep","1234"], running: false, dirty: false
Test2: ["sleep","2345"], running: false, dirty: false
`
	assert.Equal(t, expectedInfoOut, out)
}

func TestAppInfoJsonCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))
	pctlApp := app.CreateCliApp()
	out := test.CaptureStdout(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "info", "--format", "json"}))
	})

	var obj interface{}
	assert.NoError(t, json.Unmarshal([]byte(out), &obj))
	assert.Equal(t, 2, len(obj.([]interface{})))
}

func TestAppInfoJsonFlatCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))
	pctlApp := app.CreateCliApp()
	out := test.CaptureStdout(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "info", "--format", "json-flat"}))
	})

	var obj interface{}
	assert.NoError(t, json.Unmarshal([]byte(out), &obj))
	assert.Equal(t, 2, len(obj.([]interface{})))
}

func TestValidatePersistenceConfigDiscrepancyStillRunning(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("check_test", "step1.yaml"))

	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should not be running")
	pctlApp := app.CreateCliApp()

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
