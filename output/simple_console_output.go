package output

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func init() {
	FormatMap["simple"] = &SimpleConsoleOutput{}
}

type SimpleConsoleOutput struct {
	writer io.Writer
}

func (o *SimpleConsoleOutput) SetWriter(writer *os.File) {
	o.writer = writer
}

func (o *SimpleConsoleOutput) Write(infoEntries []*InfoEntry) error {
	for _, e := range infoEntries {

		b, err := json.Marshal(e.RunningCommand)
		if err != nil {
			return err
		}

		line := fmt.Sprintf("%s: %s, running: %o, dirty: %o\n", e.Name, string(b), e.IsRunning, e.ConfigCommandChanged)
		if _, err := o.writer.Write([]byte(line)); err != nil {
			return err
		}
	}
	return nil
}
