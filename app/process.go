package app

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"syscall"
	"time"

	gopsutil "github.com/shirou/gopsutil/process"
)

type Info struct {
	Comment         string
	RunningCommand  []string
	DirtyCommand    bool
	DirtyMd5Hashes  []string
	GoPsutilProcess *gopsutil.Process
	Dirty           bool
}

type Process struct {
	Config *ProcessConfig
	Info   *Info
}

func (p *Process) IsRunning() bool {
	if p.Info != nil {
		isRunning, err := p.Info.GoPsutilProcess.IsRunning()
		if err != nil || !isRunning {
			return false
		} else {
			return true
		}
	}
	return false
}

func (p *Process) Start(comment string) error {
	runningInfoStr, err := _createRunningInfoJson(comment, p)
	if err != nil {
		return err
	}
	infoEnv := fmt.Sprintf("__PCTL_INFO__=%s", runningInfoStr)

	cmd := exec.Command("setsid", p.Config.Command...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, infoEnv)

	return cmd.Start()
}

func _createRunningInfoJson(comment string, p *Process) (string, error) {
	md5hashes, err := CreateFileHashesFromCommand(p.Config.Command)
	if err != nil {
		return "", err
	}
	runningEnvironInfo := RunningEnvironInfo{
		Config:    *p.Config,
		Comment:   comment,
		Md5Hashes: *md5hashes,
	}
	infoStr, err := json.Marshal(runningEnvironInfo)
	if err != nil {
		return "", err
	}
	return string(infoStr), nil
}

func (p *Process) WaitForStarted(maxWaitTime time.Duration, intervalDuration time.Duration) error {
	// TODO readiness probe implementation
	return nil
}

func (p *Process) Stop() error {
	stopStrategy := NewStopStrategyFromConfig(p.Config.StopStrategy)
	err := stopStrategy.Stop(p)
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) WaitForStop(maxWaitTime time.Duration, intervalDuration time.Duration) error {
	attempts := maxWaitTime / intervalDuration
	err := WaitUntilTrue(func() bool {
		return !p.IsRunning()
	}, intervalDuration, uint(attempts))
	if err != nil {
		pid := p.Info.GoPsutilProcess.Pid
		return fmt.Errorf("unable to stop process '%s' on PID %d", p.Config.Name, pid)
	}
	return nil
}

func (p *Process) Kill() error {
	return p.Info.GoPsutilProcess.SendSignal(syscall.SIGKILL)
}

func (p *Process) SetRunningInfo(runningInfo *RunningEnvironInfo) error {
	logrus.Tracef("Syncing process info for '%s'", p.Config.Name)
	if runningInfo != nil {
		dirtyHashes, err := collectDirtyHashes(&p.Config.Command, runningInfo)
		if err != nil {
			return err
		}
		dirtyCommand := !CompareStringSlices(runningInfo.Config.Command, p.Config.Command)
		gopsutilProcess, err := gopsutil.NewProcess(runningInfo.Pid)
		if err != nil {
			return err
		}
		p.Info = &Info{
			DirtyMd5Hashes:  *dirtyHashes,
			Comment:         runningInfo.Comment,
			RunningCommand:  runningInfo.Config.Command,
			DirtyCommand:    dirtyCommand,
			Dirty:           dirtyCommand || len(*dirtyHashes) > 0,
			GoPsutilProcess: gopsutilProcess,
		}
	}
	return nil
}

func collectDirtyHashes(command *[]string, runningInfo *RunningEnvironInfo) (*[]string, error) {
	logrus.Tracef("Collecting dirty file hashes from command '%v'", *command)
	var returnDirtyHashes []string
	md5hashes, err := CreateFileHashesFromCommand(*command)
	if err != nil {
		return nil, err
	}
	for arg, hash := range *md5hashes {
		runningHash := runningInfo.Md5Hashes[arg]
		if runningHash != hash {
			logrus.Tracef("Found dirty hash for arg '%s': %s != %s", arg, runningHash, hash)
			returnDirtyHashes = append(returnDirtyHashes, arg)
		}
	}
	return &returnDirtyHashes, nil
}
