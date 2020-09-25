package stop_strategy

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"strconv"
)

type ScriptStopStrategy struct {
	config.ScriptStopStrategyConfig
}

func (s *ScriptStopStrategy) Stop(name string, p *gopsutil.Process) error {
	stopScriptPath, err := common.ExpandPath(s.Path)
	if err != nil {
		return err
	}
	args := []string{name, strconv.Itoa(int(p.Pid))}
	args = append(args, s.Args...)
	cmd := exec.Command(stopScriptPath, args...)
	if s.ForwardStdout {
		cmd.Stdout = os.Stdout
	}
	if s.ForwardStderr {
		cmd.Stderr = os.Stderr
	}
	return cmd.Run()
}
