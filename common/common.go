package common

import (
	"fmt"
	"github.com/davidscholberg/go-durationfmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func WaitUntilTrue(testFunction func() bool, interval time.Duration, attempts uint) error {
	var _attempt uint = 0
	for {
		if result := testFunction(); result == true {
			break
		}
		_attempt++
		if _attempt >= attempts {
			return fmt.Errorf("max attempts reached")
		}
		time.Sleep(interval)
	}
	return nil
}

func CompareStringSlices(slice1 []string, slice2 []string) bool {
	if len(slice1) != len(slice2) {
		return false
	}
	for index := range slice1 {
		if slice1[index] != slice2[index] {
			return false
		}
	}
	return true
}

func ExpandPath(path string) (string, error) {
	path = strings.Replace(path, "~", os.Getenv("HOME"), 1)
	path, err := filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		return "", err
	}

	return path, nil

}

func ByteCountIEC(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := uint64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %ciB",
		float64(b)/float64(div), "KMGTPE"[exp])
}

func DurationToString(d time.Duration) (string, error) {
	var format string
	if d > time.Hour*24 {
		format = "%dd %hh"
	} else if d > time.Hour {
		format = "%hh %mm"
	} else if d > time.Minute {
		format = "%mm %ss"
	} else {
		format = "%ss"
	}
	return durationfmt.Format(d, format)
}
