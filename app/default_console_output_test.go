package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/jedib0t/go-pretty/table"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestDefaultConsoleOutput(t *testing.T) {
	r, w, _ := os.Pipe()
	o := app.DefaultConsoleOutput{
		Writer: w,
		Style:  table.StyleDefault,
	}

	processes := []*app.Process{
		{
			Config: &app.ProcessConfig{
				Name:    "Process1",
				Command: []string{"sleep", "1"},
			},
		},
	}
	assert.NoError(t, o.Write(processes))
	assert.NoError(t, w.Close())
	out, err := ioutil.ReadAll(r)
	assert.NoError(t, err)

	expected := `+-------------------------------------------------------------+
|    Name      Status        Pid  Uptime  Rss    Vms  Command |
+-------------------------------------------------------------+
| 1  Process1  Stopped                                        |
+-------------------------------------------------------------+
|              Running: 0/1               Σ 0 B               |
+-------------------------------------------------------------+
`

	assert.Equal(t, expected, string(out))
}
