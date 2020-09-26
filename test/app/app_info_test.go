package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppInfoCommand(t *testing.T) {
	pctlApp := app.CreateCliApp()
	out := test.CaptureStdout(func() {
		assert.NoError(t, pctlApp.Run([]string{"pctl", "info"}))
	})
	expectedInfoOut := `Test1: ["sleep","1234"], running: false, dirty: false
Test2: ["sleep","2345"], running: false, dirty: false
`
	assert.Equal(t, expectedInfoOut, out)
}
