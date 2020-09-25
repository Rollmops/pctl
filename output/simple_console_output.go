package output

import (
	"fmt"
	"io"
	"os"
)

func init() {
	FormatMap["simple"] = NewSimpleConsoleOutput(os.Stdout)
}

type SimpleConsoleOutput struct {
	writer io.Writer
}

func NewSimpleConsoleOutput(file io.Writer) *SimpleConsoleOutput {
	return &SimpleConsoleOutput{writer: file}
}

func (v *SimpleConsoleOutput) Write(infoEntries []*InfoEntry) error {
	for _, e := range infoEntries {
		line := fmt.Sprintf("%s: [%s], running: %v, dirty: %v\n", e.Name, e.RunningCommand, e.IsRunning, e.ConfigCommandChanged)
		if _, err := v.writer.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
