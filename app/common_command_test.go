package app

import (
	"github.com/Rollmops/pctl/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelpCommand(t *testing.T) {

	helpText1 := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "--help"}))
	})

	helpText2 := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "-h"}))
	})

	helpTextStart := `NAME:
   pctl - process control

USAGE:
   `
	helpTextEnd := `[global options] command [command options] [arguments...]

COMMANDS:
   version  show the version information
   start    start process(es)
   restart  restart process(es)
   stop     stop process(es)
   kill     kill process(es)
   info     show info for all or specified processes
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --context value, -c value   allows you to run independent pctl contexts in parallel [$PCTL_CONTEXT]
   --loglevel value, -L value  level: trace,debug,info,warn,warning,error,fatal,panic (default: "warning") [$PCTL_LOG_LEVEL]
   --no-color                  do not use colors (default: false) [$PCTL_NO_COLOR]
   --help, -h                  show help (default: false)
`

	assert.Contains(t, helpText1, helpTextStart)
	assert.Contains(t, helpText1, helpTextEnd)
	assert.Contains(t, helpText2, helpTextStart)
	assert.Contains(t, helpText2, helpTextEnd)

}

func TestVersionCommand(t *testing.T) {
	assert.NoError(t, test.SetConfigEnvPath("pctl.yml"))
	versionText := test.CaptureStdout(func() {
		assert.NoError(t, Run([]string{"pctl", "version"}))
	})

	assert.Equal(t, "pctl __VERSION_PLACEHOLDER__\n", versionText)
}
