package stop_strategy

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/drone/envsubst"
	gopsutil "github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"strconv"
)

type ScriptStopStrategy struct {
	config.ScriptStopStrategyConfig
}

func (s *ScriptStopStrategy) Stop(c *config.ProcessConfig, p *gopsutil.Process) error {
	mapping := map[string]string{
		"pid":  strconv.Itoa(int(p.Pid)),
		"name": c.Name,
	}
	stopScriptPath, err := common.ExpandPath(s.Path)
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
