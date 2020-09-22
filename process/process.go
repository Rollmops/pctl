package process

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"os/exec"
)

type Process struct {
	Config config.ProcessConfig
	info   *gopsutil.Process
	cmd    *exec.Cmd
}

func NewProcess(config config.ProcessConfig) Process {
	return Process{
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
	info, err := gopsutil.NewProcess(int32(p.cmd.Process.Pid))
	p.info = info
	return info, err
}
