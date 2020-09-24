package process

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	"os/exec"
	"syscall"
	"time"

	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
)

type Process struct {
	Config            *config.ProcessConfig
	info              *gopsutil.Process
	cmd               *exec.Cmd
	terminateStrategy TerminateStrategy
}

func NewProcess(config *config.ProcessConfig) *Process {
	return &Process{
		Config: config,
		cmd:    nil,
		terminateStrategy: &SignalTerminateStrategy{
			Signal: syscall.SIGTERM,
		},
	}
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
	name := p.Config.Cmd[0]

	var args []string
	if len(p.Config.Cmd) > 1 {
		args = p.Config.Cmd[1:]
	}

	p.cmd = exec.Command(name, args...)
	return p.cmd.Start()

}

func (p *Process) Terminate() error {
	err := p.terminateStrategy.Terminate(p)
	if err != nil {
		return err
	}
	err = common.WaitUntilTrue(func() bool {
		return p.IsRunning()
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
