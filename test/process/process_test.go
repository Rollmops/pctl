package process_test

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	"io/ioutil"
	"os"
	"strconv"
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

	if err = common.WaitUntilTrue(func() bool {
		if _, err = os.Stat(tmpFile.Name()); err != nil {
			return true
		}
		return false
	}, 100*time.Millisecond, 10); err != nil {
		t.Fatalf("Expect file %s to be removed after 1s.", tmpFile.Name())
	}
}

func TestProcessInfo(t *testing.T) {

	var cmdline string
	var seconds = 0

	for {
		seconds += 1
		p := process.NewProcess(config.ProcessConfig{
			Name: "test",
			Cmd:  []string{"sleep", strconv.Itoa(seconds)},
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
			t.Fatalf("Expect info of started process")
		}

		cmdline, err = info.Cmdline()
		if err == nil && cmdline != "" {
			break
		} else {
			if seconds > 5 {
				t.Fatalf("Expect running process")
			}
		}
	}

	expected := fmt.Sprintf("sleep %d", seconds)
	if cmdline != expected {
		t.Fatalf("%s != %s", cmdline, expected)
	}

}
