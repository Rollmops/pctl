package common_test

import (
	"github.com/Rollmops/pctl/common"
	"testing"
	"time"
)

var _testVariable = ""

func _delayedSetTestVariable(delay time.Duration, value string) {
	time.Sleep(delay)
	_testVariable = value
}

func TestWaitUntilTrue(t *testing.T) {
	go _delayedSetTestVariable(500*time.Millisecond, "Test1")

	if err := common.WaitUntilTrue(func() bool {
		return _testVariable == "Test1"
	}, 10*time.Millisecond, 30); err == nil {
		t.Fatal("did not expect variable to be set after 300ms")
	}

	if err := common.WaitUntilTrue(func() bool {
		return _testVariable == "Test1"
	}, 10*time.Millisecond, 30); err != nil {
		t.Fatal("did expect variable to be set after 600ms")
	}
}
