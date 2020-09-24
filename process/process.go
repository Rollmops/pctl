package process

import (
	"fmt"
	"os/exec"

	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
)

type Process struct {
	Config config.ProcessConfig
	info   *gopsutil.Process
	cmd    *exec.Cmd
}

func NewProcess(config config.ProcessConfig) *Process {
	return &Process{
		Config: config,
		cmd:    nil,
	}
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

func (p *Process) Info() (*gopsutil.Process, error) {
	if p.cmd == nil {
		return nil, fmt.Errorf("command to yet started")
	}
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
