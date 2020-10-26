package app

import (
	"fmt"
	"strconv"
)

type RunningProperty struct{}

func init() {
	PropertyMap["running"] = &RunningProperty{}
}

func (*RunningProperty) Name() string {
	return "Running"
}

func (*RunningProperty) Value(p *Process, _ bool) (string, error) {
	return strconv.FormatBool(p.IsRunning()), nil
}

func (*RunningProperty) FormattedSumValue(processList ProcessList) (string, error) {
	runningCount := 0
	for _, p := range processList {
		if p.IsRunning() {
			runningCount++
		}
	}
	return fmt.Sprintf("%d/%d", runningCount, len(processList)), nil
}
func (*RunningProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
