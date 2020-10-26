package app

import (
	"github.com/sirupsen/logrus"
	"time"
)

type AgentWatcher struct {
	processConfig *ProcessConfig
}

func (w *AgentWatcher) Start() {
	logrus.Debugf("Starting watcher for process '%s'", w.processConfig.Name)
	for {
		time.Sleep(1 * time.Second)
	}
}
