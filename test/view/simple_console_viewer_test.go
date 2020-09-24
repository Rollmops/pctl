package view_test

import (
	"bytes"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	"github.com/Rollmops/pctl/view"
	"testing"
)

func TestSimpleConsoleViewer(t *testing.T) {
	processes := []*process.Process{
		process.NewProcess(
			&config.ProcessConfig{
				Name: "p1",
				Cmd:  []string{"sleep", "10"},
			}),
		process.NewProcess(
			&config.ProcessConfig{
				Name: "p2",
				Cmd:  []string{"ls", "-la"},
			}),
	}

	var w bytes.Buffer

	viewer := view.NewSimpleConsoleViewer(&w)
	//viewer := view.NewSimpleConsoleViewer(os.Stdout)

	err := viewer.View(processes)
	if err != nil {
		t.Fatal(err)
	}

	expectedOutput := `p1: [sleep 10], running: no
p2: [ls -la], running: no
`

	if w.String() != expectedOutput {
		t.Fatalf("%s != %s", expectedOutput, w.String())
	}
}
