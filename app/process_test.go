package app_test

import (
	"github.com/Rollmops/pctl/app"
	"github.com/stretchr/testify/assert"
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

	p := &app.Process{Config: &app.ProcessConfig{
		Name:    "test",
		Command: []string{"rm", tmpFile.Name()},
	}}

	err = p.Start("")
	if err != nil {
		t.Fatal(err)
	}

	if err = app.WaitUntilTrue(func() bool {
		if _, err = os.Stat(tmpFile.Name()); err != nil {
			return true
		}
		return false
	}, 100*time.Millisecond, 10); err != nil {
		t.Fatalf("Expect file %s to be removed after 1s.", tmpFile.Name())
	}
}

func TestProcessIsRunning(t *testing.T) {
	p := &app.Process{Config: &app.ProcessConfig{
		Name:    "test",
		Command: []string{"sleep", "1"},
	}}

	assert.False(t, p.IsRunning())

	err := p.Start("")
	assert.NoError(t, err)

	if err = app.WaitUntilTrue(func() bool {
		return p.IsRunning()
	}, 100*time.Millisecond, 10); err != nil {
		t.Fatal("Expect process to be running")
	}

	// unfortunately the process hangs in a defunct state after sleep 1 exited (also with releasing it)
}
