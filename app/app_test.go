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

func TestStartOrderCommand(t *testing.T) {
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
