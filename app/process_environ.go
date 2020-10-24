package app

import (
	"encoding/json"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

type RunningEnvironInfo struct {
	Config    ProcessConfig
	Pid       int32
	Comment   string
	Md5Hashes map[string]string
}

func (c *Context) SyncRunningProcesses() error {
	runningProcesses, err := c.collectRunningProcesses()
	if err != nil {
		return err
	}
	c.RunningProcesses = runningProcesses
	return nil
}

func (c *Context) collectRunningProcesses() ([]*Process, error) {
	processIds, err := gopsutil.Pids()
	if err != nil {
		return nil, err
	}
	var runningProcesses []*Process
	for _, pid := range processIds {
		runningInfo, err := TraverseToTopParentWithRunningInfo(pid)
		if err != nil {
			return nil, err
		}
		if runningInfo != nil {
			config := c.Config.FindByName(runningInfo.Config.Name)
			if config == nil {
				logrus.Warningf("Unable to find config for running process '%s' with process id %d",
					runningInfo.Config.Name, pid)
				continue
			}
			p := Process{Config: config}
			err = p.SetRunningInfo(runningInfo)
			if err != nil {
				return nil, err
			}
			runningProcesses = append(runningProcesses, &p)
		}
	}
	return runningProcesses, nil
}

func TraverseToTopParentWithRunningInfo(pid int32) (*RunningEnvironInfo, error) {
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
	runningInfoFromParent, err := TraverseToTopParentWithRunningInfo(ppid)
	if runningInfoFromParent != nil {
		return runningInfoFromParent, nil
	}
	return runningInfo, nil
}

func findRunningEnvironInfoFromPid(pid int32) (*RunningEnvironInfo, error) {
	envFilePath := path.Join("/", "proc", strconv.Itoa(int(pid)), "environ")
	envContent, err := ioutil.ReadFile(envFilePath)
	if err == nil {
		envContentString := string(envContent)
		if strings.Contains(envContentString, "__PCTL_INFO__") {
			envStrings := strings.Split(string(envContent), "\000")
			for _, envString := range envStrings {
				if strings.Contains(envString, "__PCTL_INFO__") {
					envJson := strings.Join(strings.Split(envString, "=")[1:], "")
					var runningInfo RunningEnvironInfo
					err = json.Unmarshal([]byte(envJson), &runningInfo)
					if err != nil {
						return nil, err
					}
					runningInfo.Pid = pid
					return &runningInfo, nil
				}
			}
		}
	}
	return nil, nil
}
