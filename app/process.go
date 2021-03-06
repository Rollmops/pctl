package app

import (
	"fmt"
	gopsutil "github.com/shirou/gopsutil/process"
	"os"
	"os/exec"
	"syscall"
	"time"
)

type ProcessList []*Process

type RunningInfo struct {
	Comment         string
	Config          ProcessConfig
	Pid             int32
	DirtyInfo       *DirtyInfo
	Md5Hashes       map[string]string
	GopsutilProcess *gopsutil.Process
}

type DirtyInfo struct {
	DirtyCommand   bool
	DirtyMd5Hashes []string
}

func (d *DirtyInfo) IsDirty() bool {
	return d.DirtyCommand || len(d.DirtyMd5Hashes) > 0
}

func (r *RunningInfo) SetDirty(processConfig *ProcessConfig) error {
	r.DirtyInfo = &DirtyInfo{
		DirtyCommand: !CompareStringSlices(processConfig.Command, r.Config.Command),
	}
	return nil
}

type Process struct {
	Config      *ProcessConfig
	RunningInfo *RunningInfo
}

func (l *ProcessList) SyncRunningInfo() error {
	for _, p := range *l {
		p.SyncRunningInfo()
	}
	return nil
}

func (p *Process) IsRunning() bool {
	if p.RunningInfo != nil {
		isRunning, err := p.RunningInfo.GopsutilProcess.IsRunning()
		if err != nil || !isRunning {
			return false
		} else {
			return true
		}
	}
	return false
}

func (p *Process) Start(comment string) (int32, error) {
	runningInfoStr, err := createRunningInfoJson(comment, p)
	if err != nil {
		return -1, err
	}
	infoEnv := fmt.Sprintf("__PCTL_INFO__=%s", runningInfoStr)

	cmd := exec.Command("setsid", p.Config.Command...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, infoEnv)

	for key, value := range p.Config.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	err = cmd.Start()
	return int32(cmd.Process.Pid), err
}

func (p *Process) WaitForStartup(pid int32) (bool, error) {
	startupProbe := p.Config.StartupProbe
	if startupProbe != nil {
		if p.RunningInfo == nil {
			p.RunningInfo = &RunningInfo{Pid: pid}
		} else {
			p.RunningInfo.Pid = pid
		}
		return startupProbe.Probe(p, 100*time.Millisecond)
	}
	return true, nil
}

func (p *Process) Stop() error {
	stopStrategy := NewStopStrategyFromConfig(p.Config.Stop)
	err := stopStrategy.Stop(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) WaitForStop() (bool, error) {
	timeout, err := p.Config.GetStopConfig().GetTimeout()
	if err != nil {
		return false, err
	}
	interval, err := p.Config.GetStopConfig().GetInterval()
	if err != nil {
		return false, err
	}
	stopped, err := WaitUntilTrue(func() (bool, error) {
		return !p.IsRunning(), nil
	}, timeout, interval)
	if err != nil {
		return false, err
	}
	return stopped, nil
}

func (p *Process) Kill() error {
	return p.RunningInfo.GopsutilProcess.SendSignal(syscall.SIGKILL)
}

func (p *Process) SyncRunningInfo() {
	p.RunningInfo = CurrentContext.Cache.FindRunningInfoByGroupAndName(p.Config.Group, p.Config.Name)
}
