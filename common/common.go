package common

import (
	"fmt"
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
