package app

import "fmt"

type ProcessGraphNode struct {
	Process      *Process
	dependencies []*ProcessGraphNode
	started      bool
	stopped      bool
	stopErr      *error
	startErr     *error
}

// Is ready to start, when all dependencies are started
func (c *ProcessGraphNode) IsReadyToStart() (bool, error) {
	for _, d := range c.dependencies {
		if d.startErr != nil {
			return false, *d.startErr
		}
		if !d.started {
			return false, nil
		}
	}
	return true, nil
}

func (c *ProcessGraphNode) IsReadyToStop() (bool, error) {
	for _, d := range c.dependencies {
		if d.stopErr != nil {
			return false, *d.stopErr
		}
		if !d.stopped {
			return false, nil
		}
	}
	return true, nil
}

func (c *ProcessGraphNode) AddDependency(d *ProcessGraphNode) {
	for _, startDep := range c.dependencies {
		if d == startDep {
			return
		}
	}
	c.dependencies = append(c.dependencies, d)
}

func NewProcessGraphMap(processes ProcessList, dependencyGetter func(*Process) []string) (*map[string]*ProcessGraphNode, error) {
	processStateMap := make(map[string]*ProcessGraphNode)
	var err error
	for _, p := range processes {
		processStateMap, err = addToProcessGraphMap(p, processStateMap, dependencyGetter)
		if err != nil {
			return nil, err
		}
	}
	return &processStateMap, nil
}

func addToProcessGraphMap(p *Process, processStateMap map[string]*ProcessGraphNode, dependencyGetter func(*Process) []string) (map[string]*ProcessGraphNode, error) {
	if processStateMap[p.Config.String()] == nil {
		processStateMap[p.Config.String()] = &ProcessGraphNode{
			Process: p,
			started: false,
		}
	}
	for _, d := range dependencyGetter(p) {
		dependencyConfigs, err := CurrentContext.Config.FindByGroupNameSpecifier(d)
		if err != nil {
			return nil, err
		}
		if len(dependencyConfigs) == 0 {
			return nil, fmt.Errorf("unable to find processes for specifier '%s'", d)
		}
		for _, dependencyConfig := range dependencyConfigs {
			dependencyProcess, err := CurrentContext.GetProcessByConfig(dependencyConfig)
			if err != nil {
				return nil, err
			}
			processStateMap, err = addToProcessGraphMap(dependencyProcess, processStateMap, dependencyGetter)
			if err != nil {
				return nil, err
			}
			processStateMap[p.Config.String()].AddDependency(processStateMap[dependencyConfig.String()])
		}

	}
	return processStateMap, nil
}
