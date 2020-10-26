package app

import gopsutil "github.com/shirou/gopsutil/process"

type Cache struct {
	runningInfoList []*RunningEnvironInfo
}

func (c *Cache) Refresh() error {
	c.runningInfoList = make([]*RunningEnvironInfo, 0)
	processIds, err := gopsutil.Pids()
	if err != nil {
		return err
	}
	for _, pid := range processIds {
		runningInfo, err := FindProcessRunningInfo(pid)
		if err != nil {
			return err
		}
		if runningInfo != nil {
			c.runningInfoList = append(c.runningInfoList, runningInfo)
		}
	}
	return nil
}

func (c *Cache) FindRunningInfoByGroupAndName(group string, name string) *RunningEnvironInfo {
	for _, r := range c.runningInfoList {
		if r.Config.Group == group && r.Config.Name == name {
			return r
		}
	}
	return nil
}
