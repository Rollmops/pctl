package app

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
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

type Info struct {
	Comment         string
	RunningCommand  []string
	DirtyCommand    bool
	DirtyMd5Hashes  []string
	GoPsutilProcess *gopsutil.Process
	Dirty           bool
}

type RunningEnvironInfo struct {
	Config    ProcessConfig
	Pid       int32
	Comment   string
	Md5Hashes map[string]string
}

type Process struct {
	Config *ProcessConfig
	Info   *Info
	Pid    int32
}

func (l *ProcessList) SyncRunningInfo() error {
	for _, p := range *l {
		err := p.SyncRunningInfo()
		if err != nil {
			return err
		}
	}
	return nil
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
	runningInfoStr, err := createRunningInfoJson(comment, p)
	if err != nil {
		return err
	}
	infoEnv := fmt.Sprintf("__PCTL_INFO__=%s", runningInfoStr)

	cmd := exec.Command("setsid", p.Config.Command...)
	cmd.Env = os.Environ()
	cmd.Env = append(cmd.Env, infoEnv)

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = CurrentContext.Cache.Refresh()
	if err != nil {
		return err
	}
	err = p.SyncRunningInfo()
	if err != nil {
		return err
	}
	if !p.IsRunning() {
		return fmt.Errorf("startup probe failed")
	}

	return nil
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
	return p.Info.GoPsutilProcess.SendSignal(syscall.SIGKILL)
}

func (p *Process) SyncRunningInfo() error {
	runningInfo := CurrentContext.Cache.FindRunningInfoByGroupAndName(p.Config.Group, p.Config.Name)
	return p.setRunningInfo(runningInfo)
}

func (p *Process) setRunningInfo(runningInfo *RunningEnvironInfo) error {
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
		p.Pid = runningInfo.Pid
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

func FindProcessRunningInfo(pid int32) (*RunningEnvironInfo, error) {
	runningInfo, err := findRunningEnvironInfoFromPid(pid)
	if err != nil {
		return nil, err
	}
	if runningInfo == nil {
		return nil, nil
	}

	proc, err := gopsutil.NewProcess(pid)
	if err != nil {
		return nil, err
	}
	ppid, err := proc.Ppid()
	if err != nil {
		return nil, err
	}
	runningInfoFromParent, err := FindProcessRunningInfo(ppid)
	if runningInfoFromParent != nil {
		return runningInfoFromParent, nil
	}
	return runningInfo, nil
}

func findRunningEnvironInfoFromPid(pid int32) (*RunningEnvironInfo, error) {
	envString := getProcessEnvironVariable(pid, "__PCTL_INFO__")
	if envString != "" {
		envJson := strings.SplitN(envString, "=", 2)[1]
		var runningInfo RunningEnvironInfo
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
		envContentString := string(envContent)
		if strings.Contains(envContentString, name) {
			envStrings := strings.Split(string(envContent), "\000")
			for _, envString := range envStrings {
				if strings.HasPrefix(envString, name) {
					return envString
				}
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
