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
	defer func() {
		assert.NoError(t, app.Run([]string{"pctl", "stop", "--nowait", "*"}))
	}()

	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should not be running")

	err := app.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "Test1"})
	assert.NoError(t, err)
	assert.True(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be running")

	err = app.Run([]string{"pctl", "--loglevel", "DEBUG", "stop", "--nowait", "Test1"})
	assert.NoError(t, err)
	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be stopped")
}

func TestStartCommand(t *testing.T) {
	defer func() {
		assert.NoError(t, app.Run([]string{"pctl", "stop", "--nowait", "*"}))
	}()

	assert.NoError(t, test.SetConfigEnvPath("start-order.yaml"))

	assert.False(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should not be running")

	err := app.Run([]string{"pctl", "--loglevel", "DEBUG", "start", "0"})
	assert.NoError(t, err)
	assert.True(t, test.IsCommandRunning("sleep 1234"), "'sleep 1234' should be running")
}

func TestStartWithDependencies(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("dependsOn.yaml"))

	defer func() {
		assert.NoError(t, app.Run([]string{"pctl", "stop", "--nowait", "*"}))
	}()

	out := test.CaptureStdout(func() {
		assert.NoError(t, app.Run([]string{"pctl", "--no-color", "start", "p1"}))
	})

	assert.True(t, test.IsCommandRunning("sleep 3456"), "'sleep 3456' should be running")
	assert.True(t, test.IsCommandRunning("sleep 4567"), "'sleep 4567' should be running")

	assert.Equal(t, `Started process 'p2'
Started process 'p1'
`, out)
}

func TestAppInfoCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))
	assert.NoError(t, app.Run([]string{"pctl", "info", "--format", "simple"}))
}

func TestAppInfoJsonCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))
	out := test.CaptureStdout(func() {
		assert.NoError(t, app.Run([]string{"pctl", "info", "--format", "json"}))
	})

	var obj interface{}
	assert.NoError(t, json.Unmarshal([]byte(out), &obj))
	assert.Equal(t, 2, len(obj.([]interface{})))
}

func TestAppInfoJsonFlatCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("integration.yaml"))
	out := test.CaptureStdout(func() {
		assert.NoError(t, app.Run([]string{"pctl", "info", "--format", "json-flat"}))
	})

	var obj interface{}
	assert.NoError(t, json.Unmarshal([]byte(out), &obj))
	assert.Equal(t, 2, len(obj.([]interface{})))
}

func TestValidatePersistenceConfigDiscrepancyStillRunning(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("check_test", "step1.yaml"))

	assert.False(t, test.IsCommandRunning("sleep 5555"), "'sleep 5555' should not be running")

	err := app.Run([]string{"pctl", "start", "Test1"})
	assert.NoError(t, err)
	assert.True(t, test.IsCommandRunning("sleep 5555"), "'sleep 5555' should be running")

	assert.NoError(t, test.SetConfigEnvPath("check_test", "step2.yaml"))

	err = app.Run([]string{"pctl", "info"})
	assert.Error(t, err)

	assert.Regexp(t, "found tracked running process 'Test1' with pid \\d+ that could not be found in config", err.Error())

}
