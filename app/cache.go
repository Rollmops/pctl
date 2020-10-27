package app

import (
	gopsutil "github.com/shirou/gopsutil/process"
	"github.com/sirupsen/logrus"
	"sync"
)

type Cache struct {
	runningInfoList []*RunningInfo
}

func refreshForPid(pid int32, c *chan *RunningInfo, wg *sync.WaitGroup) {
	runningInfo, _ := FindProcessRunningInfo(pid)
	if runningInfo != nil {
		*c <- runningInfo
	}
	wg.Done()
}

func (c *Cache) Refresh() error {
	c.runningInfoList = make([]*RunningInfo, 0)
	processIds, err := gopsutil.Pids()
	if err != nil {
		return err
	}
	wg := sync.WaitGroup{}
	wg.Add(len(processIds))
	channel := make(chan *RunningInfo)
	go func() {
		wg.Wait()
		close(channel)
	}()
	for _, pid := range processIds {
		go refreshForPid(pid, &channel, &wg)
	}
	for runningInfo := range channel {
		c.runningInfoList = append(c.runningInfoList, runningInfo)
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
			logrus.Tracef("Found dirty hash for arg '%s': %s != %s", arg, runningHash, hash)
			returnDirtyHashes = append(returnDirtyHashes, arg)
		}
	}
	return returnDirtyHashes, nil
}
