package app

import (
	"github.com/sirupsen/logrus"
	"time"
)

type AgentWatcher struct {
	processConfig *ProcessConfig
	runningInfo   *RunningInfo
}

func (w *AgentWatcher) Start() {
	logrus.Infof("Starting watcher for process %s", w.processConfig)
	for {
		if w.runningInfo == nil {
			w.runningInfo = CurrentContext.Cache.FindRunningInfoByGroupAndName(w.processConfig.Group, w.processConfig.Name)
		}

		if w.runningInfo != nil {
			isRunning, err := w.runningInfo.GopsutilProcess.IsRunning()
			if err != nil {
				logrus.Warningf(err.Error())
			}
			if !isRunning {
				logrus.Warningf("Process %s stopped")
			}
		}
		time.Sleep(1 * time.Second)
	}
}
