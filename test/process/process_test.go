package process_test

import (
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	"io/ioutil"
	"os"
	"testing"
	"time"
)

func TestProcessStart(t *testing.T) {

	tmpFile, err := ioutil.TempFile(os.TempDir(), "pctl_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}

	p := process.NewProcess(config.ProcessConfig{
		Name: "test",
		Cmd:  []string{"rm", tmpFile.Name()},
	})

	err = p.Start()
	if err != nil {
		t.Fatal(err)
	}

	attempts := 0
	for {
		attempts += 1

		_, err = os.Stat(tmpFile.Name())
		if err != nil {
			break
		}
		if attempts > 10 {
			t.Fatalf("Expect file %s to be removed after 1s.", tmpFile.Name())
		}

		time.Sleep(100 * time.Millisecond)
	}

}

func TestProcessInfo(t *testing.T) {

	p := process.NewProcess(config.ProcessConfig{
		Name: "test",
		Cmd:  []string{"sleep", "1"},
	})

	_, err := p.Info()
	if err == nil {
		t.Fatalf("Expect error on non started process")
	}

	err = p.Start()
	if err != nil {
		t.Fatal(err)
	}

	info, err := p.Info()
	if err != nil {
		t.Fatalf("Expect error on non started process")
	}

	cmdline, _ := info.Cmdline()
	if cmdline != "sleep 1" {
		t.Fatalf("%s != sleep 1", cmdline)
	}

}
