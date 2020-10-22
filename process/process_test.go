package process_test

import (
	"fmt"
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"testing"
	"time"
)

var _testDataDir string

func init() {
	cwd, _ := os.Getwd()
	_testDataDir = path.Join(cwd, "fixtures")
}

func TestProcessStart(t *testing.T) {
	tmpFile, err := ioutil.TempFile(os.TempDir(), "pctl_test.*.txt")
	if err != nil {
		t.Fatal(err)
	}

	p := &process.Process{Config: &config.ProcessConfig{
		Name:    "test",
		Command: []string{"rm", tmpFile.Name()},
	}}

	err = p.Start("")
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

func TestProcessIsRunning(t *testing.T) {
	p := &process.Process{Config: &config.ProcessConfig{
		Name:    "test",
		Command: []string{"sleep", "1"},
	}}

	assert.False(t, p.IsRunning())

	err := p.Start("")
	assert.NoError(t, err)

	if err = common.WaitUntilTrue(func() bool {
		return p.IsRunning()
	}, 100*time.Millisecond, 10); err != nil {
		t.Fatal("Expect process to be running")
	}

	// unfortunately the process hangs in a defunct state after sleep 1 exited (also with releasing it)
}

func TestProcessInfo(t *testing.T) {

	var cmdline string
	var seconds = 0

	for {
		seconds += 1
		p := &process.Process{Config: &config.ProcessConfig{
			Name:    "test",
			Command: []string{"sleep", strconv.Itoa(seconds)},
		}}

		_, err := p.GetPsutilProcess()
		if err == nil {
			t.Fatalf("Expect error on non started process")
		}

		err = p.Start("")
		if err != nil {
			t.Fatal(err)
		}

		info, err := p.GetPsutilProcess()
		if err != nil {
			t.Fatalf("Expect psutilProcess of started process")
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
