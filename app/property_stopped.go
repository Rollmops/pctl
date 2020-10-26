package app

import (
	"fmt"
	"strconv"
)

type StoppedProperty struct{}

func init() {
	PropertyMap["stopped"] = &StoppedProperty{}
}

func (*StoppedProperty) Name() string {
	return "Running"
}

func (*StoppedProperty) Value(p *Process, _ bool) (string, error) {
	return strconv.FormatBool(!p.IsRunning()), nil
}

func (*StoppedProperty) FormattedSumValue(processList ProcessList) (string, error) {
	stoppedCount := 0
	for _, p := range processList {
		if !p.IsRunning() {
			stoppedCount++
		}
	}
	return fmt.Sprintf("%d/%d", stoppedCount, len(processList)), nil
}
func (*StoppedProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
