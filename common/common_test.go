package common_test

import (
	"github.com/Rollmops/pctl/common"
	"github.com/stretchr/testify/assert"
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

func TestDurationToString(t *testing.T) {
	s, err := common.DurationToString(10 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, "0s", s)

	s, err = common.DurationToString(2 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "2s", s)

	s, err = common.DurationToString(1 * time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "60s", s)

	s, err = common.DurationToString(121 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "2m 1s", s)

	s, err = common.DurationToString(121 * time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "2h 1m", s)

	s, err = common.DurationToString(25 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, "1d 1h", s)
}

func TestByteCountIEC(t *testing.T) {
	assert.Equal(t, "10 B", common.ByteCountIEC(10))
	assert.Equal(t, "1.0 KiB", common.ByteCountIEC(1024))
	assert.Equal(t, "1.0 MiB", common.ByteCountIEC(1024*1024))
	assert.Equal(t, "1.0 GiB", common.ByteCountIEC(1024*1024*1024))
	assert.Equal(t, "1.0 TiB", common.ByteCountIEC(1024*1024*1024*1024))
}
