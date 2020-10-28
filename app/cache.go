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
)

type Cache struct {
	runningInfoList []*RunningInfo
}

type refreshChannelType struct {
	runningInfo *RunningInfo
	err         error
}

func (c *Cache) refreshForPid(pid int32, channel chan *refreshChannelType, wg *sync.WaitGroup) {
	defer wg.Done()
	runningInfo, _ := c.FindProcessRunningInfo(pid)
	if runningInfo != nil {
		channel <- &refreshChannelType{runningInfo, nil}
	}
}

func (c *Cache) Refresh() error {
	c.runningInfoList = make([]*RunningInfo, 0)
	processIds, err := gopsutil.Pids()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(processIds))
	channel := make(chan *refreshChannelType)
	go func() {
		wg.Wait()
		close(channel)
	}()
	for _, pid := range processIds {
		go c.refreshForPid(pid, channel, &wg)
	}
	for refreshChannelReturn := range channel {
		if refreshChannelReturn.err != nil {
			return err
		}
		c.runningInfoList = append(c.runningInfoList, refreshChannelReturn.runningInfo)
	}

	return nil
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

func (c *Cache) FindProcessRunningInfo(pid int32) (*RunningInfo, error) {
	runningInfo, err := c.findRunningEnvironInfoFromPid(pid)
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
	runningInfoFromParent, _ := c.FindProcessRunningInfo(ppid)
	if runningInfoFromParent != nil {
		return runningInfoFromParent, nil
	}

	runningInfo.GopsutilProcess = gopsutilProcess
	// TODO do this in agent
	//runningInfo.DirtyMd5Hashes, _ = collectDirtyHashes(&runningInfo.Config.Command, runningInfo)
	return runningInfo, nil
}

func (c *Cache) findRunningEnvironInfoFromPid(pid int32) (*RunningInfo, error) {
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
			return nil, fmt.Errorf("unable to find running process %s in config", runningInfo.Config)
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
