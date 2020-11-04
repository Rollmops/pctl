package app_test

import (
	"github.com/Rollmops/pctl/app"
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

	if finished, _ := app.WaitUntilTrue(func() (bool, error) {
		return _testVariable == "Test1", nil
	}, 300*time.Millisecond, 10*time.Millisecond); finished == true {
		t.Fatal("did not expect variable to be set after 300ms")
	}

	if finished, _ := app.WaitUntilTrue(func() (bool, error) {
		return _testVariable == "Test1", nil
	}, 300*time.Millisecond, 10*time.Millisecond); finished == false {
		t.Fatal("did expect variable to be set after 600ms")
	}
}

func TestDurationToString(t *testing.T) {
	s := app.DurationToString(10 * time.Millisecond)
	assert.Equal(t, "0s", s)

	s = app.DurationToString(2 * time.Second)
	assert.Equal(t, "2s", s)

	s = app.DurationToString(1 * time.Minute)
	assert.Equal(t, "60s", s)

	s = app.DurationToString(121 * time.Second)
	assert.Equal(t, "2m 1s", s)

	s = app.DurationToString(121 * time.Minute)
	assert.Equal(t, "2h 1m", s)

	s = app.DurationToString(25 * time.Hour)
	assert.Equal(t, "1d 1h", s)
}

func TestByteCountIEC(t *testing.T) {
	assert.Equal(t, "10 B", app.ByteCountIEC(10))
	assert.Equal(t, "1.0 KiB", app.ByteCountIEC(1024))
	assert.Equal(t, "1.0 MiB", app.ByteCountIEC(1024*1024))
	assert.Equal(t, "1.0 GiB", app.ByteCountIEC(1024*1024*1024))
	assert.Equal(t, "1.0 TiB", app.ByteCountIEC(1024*1024*1024*1024))
}
