package app

import (
	"fmt"
)

type MemPercentageProperty struct{}

func init() {
	memPercentageProperty := &MemPercentageProperty{}
	PropertyMap["mem%"] = memPercentageProperty
	PropertyMap["mem-percentage"] = memPercentageProperty
}

func (*MemPercentageProperty) Name() string {
	return "Memory %"
}

func (*MemPercentageProperty) Value(p *Process, _ bool) (string, error) {
	if p.RunningInfo != nil && p.IsRunning() {
		memPercentage, err := p.RunningInfo.GopsutilProcess.MemoryPercent()
		if err != nil {
			return err.Error(), nil
		} else {
			return fmt.Sprintf("%.2f", memPercentage), nil
		}
	}
	return "", nil
}

func (*MemPercentageProperty) FormattedSumValue(processList ProcessList) (string, error) {
	var memPercentageSum float32
	for _, p := range processList {
		if p.RunningInfo != nil && p.IsRunning() {
			memPercentage, err := p.RunningInfo.GopsutilProcess.MemoryPercent()
			if err != nil {
				return "", err
			}
			memPercentageSum += memPercentage
		}
	}
	return fmt.Sprintf("Σ %.2f", memPercentageSum), nil
}
func (*MemPercentageProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
