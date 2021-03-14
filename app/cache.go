package app

import (
	"encoding/json"
	"fmt"
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"sync"
)

type Cache struct {
	runningInfoList []*RunningInfo
	runningInfoMap  map[int32]*RunningInfo
}

type RefreshResult struct {
	runningInfo *RunningInfo
	err         error
}

type RefreshChannelType chan *RefreshResult

func (c *Cache) refreshForPid(pid int32, channel RefreshChannelType, wg *sync.WaitGroup) {
	defer wg.Done()
	runningInfo, err := c.findTopProcessRunningInfo(pid)
	channel <- &RefreshResult{runningInfo, err}
}

func (c *Cache) collectRunningInfo() error {
	r := regexp.MustCompile(CurrentContext.GetProcessEnvironmentMarker() + `=(\{.+\})`)

	c.runningInfoMap = make(map[int32]*RunningInfo)
	out, err := exec.Command("ps", "e", "-ww").Output()
	if err != nil {
		return err
	}
	entries := strings.Split(string(out), "\n")
	for _, entry := range entries {
		entry = strings.Trim(entry, " ")

		if strings.Contains(entry, CurrentContext.GetProcessEnvironmentMarker()) {
			pid, err := strconv.Atoi(strings.Split(entry, " ")[0])
			if err != nil {
				return err
			}
			match := r.FindStringSubmatch(entry)
			if match == nil || len(match) != 2 {
				return fmt.Errorf("unable to find match in %s", entry)
			}
			var runningInfo RunningInfo
			err = json.Unmarshal([]byte(match[1]), &runningInfo)
			if err != nil {
				return err
			}
			processConfig := CurrentContext.Config.FindByGroupAndName(runningInfo.Config.Group, runningInfo.Config.Name)
			if processConfig == nil {
				return fmt.Errorf("unable to find running process %s with PID %d in config", runningInfo.Config.String(), pid)
			}
			runningInfo.Pid = int32(pid)
			err = runningInfo.SetDirty(processConfig)
			if err != nil {
				return err
			}
			c.runningInfoMap[int32(pid)] = &runningInfo
		}
	}
	return nil
}

func (c *Cache) Refresh() error {
	logrus.Debug("Start refreshing cache")
	err := c.collectRunningInfo()
	if err != nil {
		return err
	}
	c.runningInfoList = make([]*RunningInfo, 0)
	processIds, err := gopsutil.Pids()
	if err != nil {
		return err
	}
	logrus.Tracef("Found %d running PIDs", len(processIds))
	wg := sync.WaitGroup{}
	wg.Add(len(processIds))
	refreshResultChannel := make(chan *RefreshResult)
	go func() {
		wg.Wait()
		close(refreshResultChannel)
	}()
	for _, pid := range processIds {
		go c.refreshForPid(pid, refreshResultChannel, &wg)
	}
	for refreshChannelReturn := range refreshResultChannel {
		if refreshChannelReturn.err != nil {
			return refreshChannelReturn.err
		}
		if refreshChannelReturn.runningInfo != nil {
			c.runningInfoList = append(c.runningInfoList, refreshChannelReturn.runningInfo)
		}
	}
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

func (c *Cache) findTopProcessRunningInfo(pid int32) (*RunningInfo, error) {
	runningInfo := c.runningInfoMap[pid]
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
	runningInfoFromParent, err := c.findTopProcessRunningInfo(ppid)
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
