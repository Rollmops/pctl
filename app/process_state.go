package app

type ProcessState struct {
	Process      *Process
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

func NewProcessStateMap(processes ProcessList, dependencyGetter func(*Process) []string) (*map[string]*ProcessState, error) {
	processStateMap := make(map[string]*ProcessState)
	var err error
	for _, p := range processes {
		processStateMap, err = addToProcessStateMap(p, processStateMap, dependencyGetter)
		if err != nil {
			return nil, err
		}
	}
	return &processStateMap, nil
}

func addToProcessStateMap(p *Process, processStateMap map[string]*ProcessState, dependencyGetter func(*Process) []string) (map[string]*ProcessState, error) {
	if processStateMap[p.Config.String()] == nil {
		processStateMap[p.Config.String()] = &ProcessState{
			Process: p,
			started: false,
		}
	}
	for _, d := range dependencyGetter(p) {
		dependencyConfigs, err := CurrentContext.Config.FindByGroupNameSpecifier(d)
		if err != nil {
			return nil, err
		}
		for _, dependencyConfig := range dependencyConfigs {
			dependencyProcess, err := CurrentContext.GetProcessByConfig(dependencyConfig)
			if err != nil {
				return nil, err
			}
			processStateMap, err = addToProcessStateMap(dependencyProcess, processStateMap, dependencyGetter)
			if err != nil {
				return nil, err
			}
			processStateMap[p.Config.String()].AddDependency(processStateMap[dependencyConfig.String()])
		}

	}
	return processStateMap, nil
}
