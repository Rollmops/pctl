package output

import (
	"fmt"
	"github.com/Rollmops/pctl/process"
	"io"
)

type SimpleConsoleOutput struct {
	writer io.Writer
}

func NewSimpleConsoleOutput(file io.Writer) SimpleConsoleOutput {
	return SimpleConsoleOutput{writer: file}
}

func (v *SimpleConsoleOutput) Write(processes []*process.Process) error {
	for _, p := range processes {
		isRunning := "no"
		if p.IsRunning() {
			isRunning = "yes"
		}
		line := fmt.Sprintf("%s: %s, running: %s\n", p.Config.Name, p.Config.Cmd, isRunning)
		if _, err := v.writer.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
