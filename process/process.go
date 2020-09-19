package process

import (
	"fmt"
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"os/exec"
)

type Process struct {
	config config.ProcessConfig
	cmd    *exec.Cmd
}

func NewProcess(config config.ProcessConfig) Process {
	return Process{
		config: config,
		cmd:    nil,
	}
}

func (p *Process) Start() error {

	name := p.config.Cmd[0]

	var args []string
	if len(p.config.Cmd) > 1 {
		args = p.config.Cmd[1:]
	}

	p.cmd = exec.Command(name, args...)
	return p.cmd.Start()

}

func (p *Process) Info() (*gopsutil.Process, error) {

	if p.cmd == nil {
		return nil, fmt.Errorf("command to yet started")
	}

	return gopsutil.NewProcess(int32(p.cmd.Process.Pid))

}
