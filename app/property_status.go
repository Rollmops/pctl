package app

import (
	"fmt"
	"strings"
)

type StatusProperty struct{}

func init() {
	statusProperty := &StatusProperty{}
	PropertyMap["status"] = statusProperty
	PropertyMap["state"] = statusProperty
}

func (*StatusProperty) Name() string {
	return "State"
}

func (*StatusProperty) Value(p *Process, formatted bool) (string, error) {
	var status []string
	if p.IsRunning() {
		if formatted {
			status = append(status, OkColor("Running"))
		} else {
			status = append(status, "Running")
		}
	} else {
		if formatted {
			status = append(status, FailedColor("Stopped"))
		} else {
			status = append(status, "Stopped")
		}
	}
	if p.RunningInfo != nil {
		if p.RunningInfo.DirtyInfo.IsDirty() {
			if formatted {
				status = append(status, WarningColor("Dirty"))
			} else {
				status = append(status, "Dirty")
			}

		}
	}
	return strings.Join(status, " | "), nil
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
