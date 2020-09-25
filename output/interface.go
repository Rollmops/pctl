package output

import "github.com/Rollmops/pctl/process"

type Output interface {
	Write([]process.Process) error
}
