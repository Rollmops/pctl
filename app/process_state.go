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

func NewProcessStateMap(processes *[]*Process, dependencyGetter func(*Process) []string) *map[string]*ProcessState {
	processStateMap := make(map[string]*ProcessState)
	for _, p := range *processes {
		processStateMap = addToProcessStateMap(p, processes, processStateMap, dependencyGetter)
	}
	return &processStateMap
}

func addToProcessStateMap(p *Process, processes *[]*Process, processStateMap map[string]*ProcessState, dependencyGetter func(*Process) []string) map[string]*ProcessState {
	if processStateMap[p.Config.Name] == nil {
		processStateMap[p.Config.Name] = &ProcessState{
			Process: p,
			started: false,
		}
	}
	for _, d := range dependencyGetter(p) {
		processStateMap = addToProcessStateMap(findProcessByName(d, processes), processes, processStateMap, dependencyGetter)
		processStateMap[p.Config.Name].AddDependency(processStateMap[d])
	}
	return processStateMap
}

func findProcessByName(name string, processes *[]*Process) *Process {
	for _, p := range *processes {
		if p.Config.Name == name {
			return p
		}
	}
	return nil
}
