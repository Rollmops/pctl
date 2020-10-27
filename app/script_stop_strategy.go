package app

import (
	"github.com/drone/envsubst"
	"os"
	"os/exec"
	"strconv"
)

type ScriptStopStrategyConfig struct {
	Path          string
	Args          []string
	ForwardStdout bool `yaml:"forwardStdout"`
	ForwardStderr bool `yaml:"forwardStderr"`
}

type ScriptStopStrategy struct {
	ScriptStopStrategyConfig
}

func (s *ScriptStopStrategy) Stop(process *Process) error {
	mapping := map[string]string{
		"pid":  strconv.Itoa(int(process.RunningInfo.GopsutilProcess.Pid)),
		"name": process.Config.Name,
	}
	stopScriptPath := ExpandPath(s.Path)

	var substArgs []string
	for _, arg := range s.Args {
		substArg, err := envsubst.Eval(arg, func(s string) string {
			return mapping[s]
		})
		if err != nil {
			return err
		}
		substArgs = append(substArgs, substArg)
	}
	cmd := exec.Command(stopScriptPath, substArgs...)
	if s.ForwardStdout {
		cmd.Stdout = os.Stdout
	}
	if s.ForwardStderr {
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
