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
┃ 2  g1     p1             Stopped      ┃
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
┃ 2  g1     p1             Stopped      ┃
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
┃ 2  g1     p1             Stopped      ┃
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
┃ 1  g1     p1         Stopped      ┃
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
┃ 2  g1     p1             Stopped       sleep infinity ┃
┣━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┫
┃                          Running: 0/2                 ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
`
	assert.Equal(t, expectedOutput, infoOutput)
}

func TestInfoCommandAllColumnsRunning(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))
	defer func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "stop", "--nowait", ":"}))
	}()

	assert.NoError(t, Run([]string{"pctl", "--no-color", "start", "--comment", "test comment", "km:spserver"}))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "-c", "status,stopped,running,vms,cpu%,mem%,gids,uptime,nice,comment,dirty,rss"}))
	})
	assert.Regexp(t, "State\\s+Stopped\\s+Running\\s+Vms\\s+CPU %\\s+Memory %\\s+Gids\\s+Uptime\\s+Nice\\s+Comment\\s+Dirty\\s+Rss", infoOutput)
	assert.Regexp(t, "1\\s+Running\\s+false\\s+true\\s+\\S+\\s+\\w+\\s+\\d+\\.\\d+\\s+\\S+\\s+\\d\\S+\\s+\\d+.+test comment", infoOutput)
	assert.Regexp(t, "2\\s+Stopped\\s+true\\s+false", infoOutput)
}

func TestInfoCommandAllColumns(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))
	defer func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "stop", "--nowait", ":"}))
	}()

	assert.NoError(t, Run([]string{"pctl", "--no-color", "start", "km:spserver"}))

	infoOutput := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--no-color", "info", "-c", "name,group,cwd,deps,user,metadata,command,env"}))
	})

	assert.Regexp(t, "Name\\s+Group\\s+Cwd\\s+Dependencies\\s+Username\\s+Metadata\\s+Command\\s+Environment", infoOutput)
	assert.Regexp(t, "\\s+2\\s+p1\\s+g1\\s+spserver\\s+key=value\\s+sleep infinity", infoOutput)
}
