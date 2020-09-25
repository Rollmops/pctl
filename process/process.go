package process

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/stop_strategy"
	"os/exec"
	"time"

	gopsutil "github.com/shirou/gopsutil/process"
)

type Process struct {
	Config *config.ProcessConfig
	info   *gopsutil.Process
	cmd    *exec.Cmd
}

func (p *Process) SynchronizeWithPid(pid int32) error {
	_p, err := gopsutil.NewProcess(pid)
	if err != nil {
		return nil
	}
	p.info = _p
	return nil
}

func (p *Process) IsRunning() bool {
	info, err := p.Info()
	if err != nil {
		return false
	}
	isRunning, err := info.IsRunning()
	return isRunning && err == nil
}

func (p *Process) Pid() (int32, error) {
	strategy := PidRetrieveStrategies[p.Config.PidRetrieveStrategyName]
	return strategy.Retrieve(p)
}

func (p *Process) Start() error {
	name := p.Config.Command[0]

	var args []string
	if len(p.Config.Command) > 1 {
		args = p.Config.Command[1:]
	}

	p.cmd = exec.Command(name, args...)
	return p.cmd.Start()

}

func (p *Process) Stop() error {
	_process, err := p.Info()
	if err != nil {
		return err
	}
	stopStrategy := stop_strategy.NewStopStrategyFromConfig(p.Config.StopStrategy)
	err = stopStrategy.Stop(p.Config.Name, _process)
	if err != nil {
		return err
	}
	err = common.WaitUntilTrue(func() bool {
		return !p.IsRunning()
	}, 100*time.Millisecond, 50)
	if err != nil {
		pid, _ := p.Pid()
		return fmt.Errorf("unable to stop process '%s' on PID %d", p.Config.Name, pid)
	}
	return nil
}

func (p *Process) Kill() error {
	return p.cmd.Process.Kill()
}

func (p *Process) Info() (*gopsutil.Process, error) {
	if p.info != nil {
		return p.info, nil
	}
	pid, err := p.Pid()
	if err != nil {
		return nil, err
	}
	info, err := gopsutil.NewProcess(pid)
	p.info = info
	return info, err
}
