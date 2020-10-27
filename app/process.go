package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
	"syscall"
	"time"

	gopsutil "github.com/shirou/gopsutil/process"
)

type ProcessList []*Process

type RunningInfo struct {
	Comment         string
	Config          ProcessConfig
	DirtyMd5Hashes  []string
	DirtyCommand    bool
	Dirty           bool
	Pid             int32
	Md5Hashes       map[string]string
	GopsutilProcess *gopsutil.Process
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

func (p *Process) Start(comment string) error {
	runningInfoStr, err := createRunningInfoJson(comment, p)
	if err != nil {
		return err
	}
	infoEnv := fmt.Sprintf("__PCTL_INFO__=%s", runningInfoStr)

	cmd := exec.Command("setsid", p.Config.Command...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, infoEnv)

	for key, value := range p.Config.Env {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", key, value))
	}
	return cmd.Start()
}

func (p *Process) WaitForReady() error {
	readinessProbe := ReadinessProbes[p.Config.ReadinessProbe]
	// TODO maxWaitTime and IntervalDuration from config
	maxWaitTime := 5 * time.Second
	intervalDuration := 100 * time.Millisecond
	attempts := maxWaitTime / intervalDuration
	err := WaitUntilTrue(func() (bool, error) {
		running, err := readinessProbe.Probe(p)
		if err != nil {
			return false, err
		}
		return running, nil
	}, intervalDuration, uint(attempts))
	if err != nil {
		return err
	}
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

func (p *Process) WaitForStop(timeout time.Duration, intervalDuration time.Duration) error {
	attempts := timeout / intervalDuration
	err := WaitUntilTrue(func() (bool, error) {
		return p.IsRunning(), nil
	}, intervalDuration, uint(attempts))
	if err != nil {
		return err
	}
	return nil
}

func (p *Process) Kill() error {
	return p.RunningInfo.GopsutilProcess.SendSignal(syscall.SIGKILL)
}

func (p *Process) SyncRunningInfo() {
	p.RunningInfo = CurrentContext.Cache.FindRunningInfoByGroupAndName(p.Config.Group, p.Config.Name)
}

func FindProcessRunningInfo(pid int32) (*RunningInfo, error) {
	runningInfo, err := findRunningEnvironInfoFromPid(pid)
	if err != nil {
		return nil, err
	}
	if runningInfo == nil {
		return nil, nil
	}

	gopsutilProcess, err := gopsutil.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	ppid, err := gopsutilProcess.Ppid()
	if err != nil {
		return nil, err
	}
	runningInfoFromParent, err := FindProcessRunningInfo(ppid)
	if runningInfoFromParent != nil {
		return runningInfoFromParent, nil
	}
	runningInfo.GopsutilProcess = gopsutilProcess
	//runningInfo.DirtyMd5Hashes, _ = collectDirtyHashes(&runningInfo.Config.Command, runningInfo)
	return runningInfo, nil
}

func findRunningEnvironInfoFromPid(pid int32) (*RunningInfo, error) {
	envString := getProcessEnvironVariable(pid, "__PCTL_INFO__")
	if envString != "" {
		envJson := strings.SplitN(envString, "=", 2)[1]
		var runningInfo RunningInfo
		err := json.Unmarshal([]byte(envJson), &runningInfo)
		if err != nil {
			return nil, err
		}
		runningInfo.Pid = pid
		return &runningInfo, nil
	}
	return nil, nil
}

func getProcessEnvironVariable(pid int32, name string) string {
	envFilePath := path.Join("/", "proc", strconv.Itoa(int(pid)), "environ")
	envContent, err := ioutil.ReadFile(envFilePath)
	if err == nil {
		envStrings := strings.Split(string(envContent), "\000")
		for _, envString := range envStrings {
			if strings.HasPrefix(envString, name) {
				return envString
			}
		}
	}
	return ""
}

func createRunningInfoJson(comment string, p *Process) (string, error) {
	md5hashes, err := CreateFileHashesFromCommand(p.Config.Command)
	if err != nil {
		return "", err
	}
	runningEnvironInfo := RunningInfo{
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
