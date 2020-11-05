package app

import (
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInfoCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name      Pid  State        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1  km     spserver       Stopped      ┃
┃ 2         p1             Stopped      ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/2 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandStopped(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "--stopped"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name      Pid  State        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1  km     spserver       Stopped      ┃
┃ 2         p1             Stopped      ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/2 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandRunning(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))
	err := Run([]string{"pctl", "--no-color", "info", "--running"})
	assert.EqualError(t, err, "no matching processes found")
}

func TestInfoCommandGroupSpecifier(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "km:"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name      Pid  State        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1  km     spserver       Stopped      ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/1 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandAllSpecifier(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", ":"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name      Pid  State        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1  km     spserver       Stopped      ┃
┃ 2         p1             Stopped      ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/2 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandWildcardGroupSpecifier(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "k*:"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name      Pid  State        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1  km     spserver       Stopped      ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/1 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandWildcardNameSpecifier(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "p*"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name  Pid  State        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1         p1         Stopped      ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                      Running: 0/1 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}
func TestInfoCommandOnlyColumnCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "-c", "command"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━┓
┃    Command        ┃
┣━━━━━━━━━━━━━━━━━━━┫
┃ 1  sleep infinity ┃
┃ 2  sleep infinity ┃
┣━━━━━━━━━━━━━━━━━━━┫
┃                   ┃
┗━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandPlusColumnCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "-c", "+,command"}))
	})

	expectedOutput := `┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃    Group  Name      Pid  State         Command        ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃ 1  km     spserver       Stopped       sleep infinity ┃
┃ 2         p1             Stopped       sleep infinity ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/2                 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}
