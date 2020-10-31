package app

import (
	"fmt"
)

type CpuPercentProperty struct{}

func init() {
	cpuPercentProperty := &CpuPercentProperty{}
	PropertyMap["cpu"] = cpuPercentProperty
	PropertyMap["cpu%"] = cpuPercentProperty
}

func (*CpuPercentProperty) Name() string {
	return "CPU %"
}

func (*CpuPercentProperty) Value(p *Process, _ bool) (string, error) {
	if p.RunningInfo != nil && p.IsRunning() {
		cpuPercent, err := p.RunningInfo.GopsutilProcess.CPUPercent()
		if err != nil {
			return "error", nil
		} else {
			return fmt.Sprintf("%.2f", cpuPercent), nil
		}
	}
	return "", nil
}

func (*CpuPercentProperty) FormattedSumValue(processList ProcessList) (string, error) {
	var cpuPercentSum float64
	for _, p := range processList {
		if p.RunningInfo != nil && p.IsRunning() {
			cpuPercent, err := p.RunningInfo.GopsutilProcess.CPUPercent()
			if err == nil {
				cpuPercentSum += cpuPercent
			}
		}
	}
	return fmt.Sprintf("Σ %.2f", cpuPercentSum), nil
}
func (*CpuPercentProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
