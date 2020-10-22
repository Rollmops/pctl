package process

import (
	"encoding/json"
	"github.com/Rollmops/pctl/config"
	gopsutil "github.com/shirou/gopsutil/process"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
)

type RunningEnvironInfo struct {
	Config  config.ProcessConfig
	Pid     int32
	Comment string
}

var _pctlInfoMap = make(map[string]*RunningEnvironInfo)

func FindRunningInfo(name string) (*RunningEnvironInfo, error) {
	runningEnvironInfo := _pctlInfoMap[name]
	if runningEnvironInfo != nil {
		return runningEnvironInfo, nil
	}

	processIds, _ := gopsutil.Pids()
	for _, pid := range processIds {
		runningInfo, err := traverseToTopParentWithRunningInfo(pid, name)

		if err != nil {
			return nil, err
		}
		if runningInfo != nil {
			_pctlInfoMap[runningInfo.Config.Name] = runningInfo
			return runningInfo, nil
		}
	}
	return nil, nil
}

func traverseToTopParentWithRunningInfo(pid int32, name string) (*RunningEnvironInfo, error) {
	runningInfo, err := findRunningEnvironInfoFromPid(pid)
	if err != nil {
		return nil, err
	}
	if runningInfo == nil || runningInfo.Config.Name != name {
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
	runningInfoFromParent, err := traverseToTopParentWithRunningInfo(ppid, name)
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
