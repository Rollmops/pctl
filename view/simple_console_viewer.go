package view

import (
	"fmt"
	"github.com/Rollmops/pctl/process"
	"io"
)

type SimpleConsoleViewer struct {
	writer io.Writer
}

func NewSimpleConsoleViewer(file io.Writer) SimpleConsoleViewer {
	return SimpleConsoleViewer{writer: file}
}

func (v *SimpleConsoleViewer) View(processes []process.Process) error {
	for _, p := range processes {
		line := fmt.Sprintf("%s: %s\n", p.Config.Name, p.Config.Cmd)
		if _, err := v.writer.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
