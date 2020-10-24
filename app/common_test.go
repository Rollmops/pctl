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

	if err := app.WaitUntilTrue(func() bool {
		return _testVariable == "Test1"
	}, 10*time.Millisecond, 30); err == nil {
		t.Fatal("did not expect variable to be set after 300ms")
	}

	if err := app.WaitUntilTrue(func() bool {
		return _testVariable == "Test1"
	}, 10*time.Millisecond, 30); err != nil {
		t.Fatal("did expect variable to be set after 600ms")
	}
}

func TestDurationToString(t *testing.T) {
	s, err := app.DurationToString(10 * time.Millisecond)
	assert.NoError(t, err)
	assert.Equal(t, "0s", s)

	s, err = app.DurationToString(2 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "2s", s)

	s, err = app.DurationToString(1 * time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "60s", s)

	s, err = app.DurationToString(121 * time.Second)
	assert.NoError(t, err)
	assert.Equal(t, "2m 1s", s)

	s, err = app.DurationToString(121 * time.Minute)
	assert.NoError(t, err)
	assert.Equal(t, "2h 1m", s)

	s, err = app.DurationToString(25 * time.Hour)
	assert.NoError(t, err)
	assert.Equal(t, "1d 1h", s)
}

func TestByteCountIEC(t *testing.T) {
	assert.Equal(t, "10 B", app.ByteCountIEC(10))
	assert.Equal(t, "1.0 KiB", app.ByteCountIEC(1024))
	assert.Equal(t, "1.0 MiB", app.ByteCountIEC(1024*1024))
	assert.Equal(t, "1.0 GiB", app.ByteCountIEC(1024*1024*1024))
	assert.Equal(t, "1.0 TiB", app.ByteCountIEC(1024*1024*1024*1024))
}
