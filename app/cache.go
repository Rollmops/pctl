package app

import (
	"encoding/json"
	"fmt"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"
)

type Cache struct {
	runningInfoList []*RunningInfo
}

type RefreshResult struct {
	runningInfo *RunningInfo
	err         error
}

type RefreshChannelType chan *RefreshResult

func (c *Cache) refreshForPid(pid int32, channel RefreshChannelType, wg *sync.WaitGroup) {
	defer wg.Done()
	runningInfo, err := FindTopProcessRunningInfo(pid)
	channel <- &RefreshResult{runningInfo, err}
}

func (c *Cache) Refresh() error {
	logrus.Debug("Start refreshing cache")
	start := time.Now()
	c.runningInfoList = make([]*RunningInfo, 0)
	processIds, err := gopsutil.Pids()
	if err != nil {
		return err
	}
	logrus.Tracef("Found %d running PIDs", len(processIds))
	wg := sync.WaitGroup{}
	wg.Add(len(processIds))
	refreshResultChannel := make(chan *RefreshResult, len(CurrentContext.Config.ProcessConfigs))
	go func() {
		wg.Wait()
		close(refreshResultChannel)
	}()
	for _, pid := range processIds {
		go c.refreshForPid(pid, refreshResultChannel, &wg)
	}
	for refreshChannelReturn := range refreshResultChannel {
		if refreshChannelReturn.err != nil {
			err = refreshChannelReturn.err
		}
		if refreshChannelReturn.runningInfo != nil {
			c.runningInfoList = append(c.runningInfoList, refreshChannelReturn.runningInfo)
		}
	}
	elapsed := time.Since(start)
	logrus.Debugf("Refreshing cache took %s", elapsed)
	return err
}

func (c *Cache) FindRunningInfoByGroupAndName(group string, name string) *RunningInfo {
	for _, r := range c.runningInfoList {
		if r.Config.Group == group && r.Config.Name == name {
			return r
		}
	}
	return nil
}

func collectDirtyHashes(command *[]string, runningInfo *RunningInfo) ([]string, error) {
	logrus.Tracef("Collecting dirty file hashes from command '%v'", *command)
	var returnDirtyHashes []string
	md5hashes, err := CreateFileHashesFromCommand(*command)
	if err != nil {
		return nil, err
	}
	for arg, hash := range *md5hashes {
		runningHash := runningInfo.Md5Hashes[arg]
		if runningHash != hash {
			logrus.Tracef("Found dirty hash for arg %s: %s != %s", arg, runningHash, hash)
			returnDirtyHashes = append(returnDirtyHashes, arg)
		}
	}
	return returnDirtyHashes, nil
}

func FindTopProcessRunningInfo(pid int32) (*RunningInfo, error) {
	runningInfo, err := FindRunningEnvironInfoFromPid(pid)
	if err != nil {
		return nil, err
	}
	if runningInfo == nil {
		return nil, nil
	}

	gopsutilProcess, err := gopsutil.NewProcess(pid)
	if err != nil {
		// do not return the error here since the pid might have been gone already
		return nil, nil
	}
	ppid, _ := gopsutilProcess.Ppid()
	if ppid == -1 {
		return runningInfo, nil
	}
	runningInfoFromParent, err := FindTopProcessRunningInfo(ppid)
	if err != nil {
		return nil, err
	}
	if runningInfoFromParent != nil {
		return runningInfoFromParent, nil
	}

	runningInfo.GopsutilProcess = gopsutilProcess
	runningInfo.DirtyInfo.DirtyMd5Hashes, _ = collectDirtyHashes(&runningInfo.Config.Command, runningInfo)
	return runningInfo, nil
}

func FindRunningEnvironInfoFromPid(pid int32) (*RunningInfo, error) {
	envString := getProcessEnvironVariable(pid, "__PCTL_INFO__")
	if envString != "" {
		envJson := strings.SplitN(envString, "=", 2)[1]
		var runningInfo RunningInfo
		err := json.Unmarshal([]byte(envJson), &runningInfo)
		if err != nil {
			return nil, err
		}
		processConfig := CurrentContext.Config.FindByGroupAndName(runningInfo.Config.Group, runningInfo.Config.Name)
		if processConfig == nil {
			return nil, fmt.Errorf("unable to find running process %s with PID %d in config", runningInfo.Config.String(), pid)
		}
		runningInfo.Pid = pid
		return &runningInfo, runningInfo.SetDirty(processConfig)
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
