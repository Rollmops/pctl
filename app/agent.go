package app

import (
	"github.com/facebookgo/pidfile"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
)

func FindAgentProcess() (*gopsutil.Process, error) {
	pid, err := pidfile.Read()
	if err != nil {
		return nil, nil
	}
	p, _ := gopsutil.NewProcess(int32(pid))
	isRunning, _ := p.IsRunning()
	if !isRunning {
		return nil, nil
	}
	return p, nil
}

func CheckRunningAgentProcess() error {
	agentProcess, err := FindAgentProcess()
	if err != nil {
		return err
	}
	if agentProcess != nil {
		logrus.Fatalf("Agent process already running on PID %d", agentProcess.Pid)
	}
	return nil
}
