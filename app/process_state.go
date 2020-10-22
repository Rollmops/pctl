package app

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
)

type ProcessState struct {
	Process      *process.Process
	dependencies []*ProcessState
	started      bool
	stopped      bool
}

// Is ready to start, when all dependencies are started
func (c *ProcessState) IsReadyToStart() bool {
	for _, d := range c.dependencies {
		if !d.started {
			return false
		}
	}
	return true
}

func (c *ProcessState) IsReadyToStop() bool {
	for _, d := range c.dependencies {
		if !d.stopped {
			return false
		}
	}
	return true
}

func (c *ProcessState) AddDependency(d *ProcessState) {
	for _, startDep := range c.dependencies {
		if d == startDep {
			return
		}
	}
	c.dependencies = append(c.dependencies, d)
}

func NewFromProcessConfigs(processConfigs *[]*config.ProcessConfig, dependencyGetter func(*config.ProcessConfig) []string) *map[string]*ProcessState {
	processStateMap := make(map[string]*ProcessState)
	for _, p := range *processConfigs {
		processStateMap = addToProcessStateMap(p, processStateMap, dependencyGetter)
	}
	return &processStateMap
}

func addToProcessStateMap(c *config.ProcessConfig, processStateMap map[string]*ProcessState, dependencyGetter func(*config.ProcessConfig) []string) map[string]*ProcessState {
	if processStateMap[c.Name] == nil {
		processStateMap[c.Name] = &ProcessState{
			Process: &process.Process{Config: c},
			started: false,
		}
	}
	for _, d := range dependencyGetter(c) {
		processStateMap = addToProcessStateMap(CurrentContext.Config.FindByName(d), processStateMap, dependencyGetter)
		processStateMap[c.Name].AddDependency(processStateMap[d])
	}
	return processStateMap
}
