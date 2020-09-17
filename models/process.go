package models

type Process struct {
	Name string
}

func (p Process) ToString() string {
	return "Process(" + p.Name + ")"
}
