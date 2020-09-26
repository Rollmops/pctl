package app_test

import (
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAppInfoCommand(t *testing.T) {
	infoOut, err := test.StartAppAndGetStdout([]string{"pctl", "-L", "debug", "info"})
	assert.NoError(t, err)
	expectedInfoOut := `Test1: ["sleep","1234"], running: %!o(bool=false), dirty: %!o(bool=false)
Test2: ["sleep","2345"], running: %!o(bool=false), dirty: %!o(bool=false)
`
	assert.Equal(t, expectedInfoOut, infoOut)
}
