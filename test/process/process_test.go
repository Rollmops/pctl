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
		Cmd: []string{ "rm", tmpFile.Name() },
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