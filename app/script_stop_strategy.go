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
		"pid":  strconv.Itoa(int(process.Info.GoPsutilProcess.Pid)),
		"name": process.Config.Name,
	}
	stopScriptPath, err := ExpandPath(s.Path)
	if err != nil {
		return err
	}

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
