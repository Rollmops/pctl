package app

type Probe interface {
	Probe(*Process) (bool, error)
}

type NullProbe struct{}
type DefaultStartupProbe struct{}

var ReadinessProbes = make(map[string]Probe)

func init() {
	ReadinessProbes[""] = &NullProbe{}
}

func (p *NullProbe) Probe(_ *Process) (bool, error) {
	return true, nil
}
