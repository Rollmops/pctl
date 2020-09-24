package process_test

import (
	"github.com/Rollmops/pctl/common"
	"github.com/Rollmops/pctl/config"
	"github.com/Rollmops/pctl/process"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"
)

var _testDataDir string

func init() {
	cwd, _ := os.Getwd()
	_testDataDir = path.Join(cwd, "..", "fixtures")
}

func TestDefaultPidRetrieveStrategy(t *testing.T) {
	tmpPidFilePath := "/tmp/pctl_test.pid"
	defer os.Remove(tmpPidFilePath)
	testScriptPath := filepath.Join(_testDataDir, "write_pid.sh")

	s := process.DefaultPidRetrieveStrategy{}
	c := &config.ProcessConfig{
		Name: "PidTest",
		Cmd:  []string{"bash", testScriptPath, tmpPidFilePath},
	}
	p := process.NewProcess(c)
	err := p.Start()
	if err != nil {
		t.Fatal(err)
	}

	if err = common.WaitUntilTrue(func() bool {
		if _, err := os.Stat(tmpPidFilePath); err == nil {
			return true
		}
		return false
	}, 100*time.Millisecond, 100); err != nil {
		t.Fatal(err)
	}
	pid, err := s.Retrieve(p)
	if err != nil {
		t.Fatal(err)
	}

	b, err := ioutil.ReadFile(tmpPidFilePath)
	if err != nil {
		t.Fatal(err)
	}

	writtenPid, err := strconv.Atoi(strings.TrimSpace(string(b)))
	if err != nil {
		t.Fatal(err)
	}

	if int32(writtenPid) != pid {
		t.Fatalf("%d != %d", writtenPid, pid)
	}

}
