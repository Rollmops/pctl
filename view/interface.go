package view

import "github.com/Rollmops/pctl/process"

type Viewer interface {
	View([]process.Process) error
}
