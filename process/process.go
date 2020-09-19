package process

import (
	"github.com/Rollmops/pctl/config"
	"os/exec"
)

type Process struct {

	config config.ProcessConfig
	cmd *exec.Cmd

}

func NewProcess(config config.ProcessConfig) Process {
	return Process{
		config: config,
		cmd: nil,
	}
}

func (p Process) Start() error {

	name := p.config.Cmd[0]

	var args []string
	if len(p.config.Cmd) > 1 {
		args = p.config.Cmd[1:]
	}

	p.cmd = exec.Command(name, args...)
	return p.cmd.Start()

}