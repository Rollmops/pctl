package app

import "fmt"

type StatusProperty struct{}

func init() {
	statusProperty := &StatusProperty{}
	PropertyMap["status"] = statusProperty
	PropertyMap["state"] = statusProperty
}

func (*StatusProperty) Name() string {
	return "Status"
}

func (*StatusProperty) Value(p *Process, formatted bool) (string, error) {
	if p.IsRunning() {
		if formatted {
			return OkColor("Running"), nil
		} else {
			return "Running", nil
		}
	} else {
		if formatted {
			return FailedColor("Stopped"), nil
		} else {
			return "Stopped", nil
		}
	}
}

func (*StatusProperty) FormattedSumValue(processList ProcessList) (string, error) {
	runningCount := 0
	for _, p := range processList {
		if p.IsRunning() {
			runningCount++
		}
	}
	return fmt.Sprintf("Running: %d/%d", runningCount, len(processList)), nil
}

func (*StatusProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
