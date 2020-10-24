package app

type Probe interface {
	Probe(*Process) (bool, error)
}

type NullProbe struct{}
type DefaultStartupProbe struct{}

var StartupProbes = make(map[string]Probe)
var ReadinessProbes = make(map[string]Probe)

func init() {
	ReadinessProbes[""] = &NullProbe{}
	StartupProbes[""] = &DefaultStartupProbe{}
}

func (p *NullProbe) Probe(_ *Process) (bool, error) {
	return true, nil
}

func (p *DefaultStartupProbe) Probe(process *Process) (bool, error) {
	err := CurrentContext.SyncRunningProcesses()
	if err != nil {
		return false, err
	}

	pByName, err := CurrentContext.GetProcessByName(process.Config.Name)
	if err != nil {
		return false, err
	}
	return pByName.IsRunning(), nil
}
