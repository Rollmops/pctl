package app

import "fmt"

type VmsProperty struct{}

func init() {
	PropertyMap["vms"] = &VmsProperty{}
}

func (*VmsProperty) Name() string {
	return "Vms"
}

func (*VmsProperty) Value(p *Process, _ bool) (string, error) {
	var vms string
	if p.RunningInfo != nil && p.IsRunning() {
		memoryInfo, err := p.RunningInfo.GopsutilProcess.MemoryInfo()
		if err != nil {
			vms = "error"
		} else {
			vms = ByteCountIEC(memoryInfo.VMS)
		}
	}
	return vms, nil
}

func (*VmsProperty) FormattedSumValue(processList ProcessList) (string, error) {
	var vmsSum uint64
	for _, p := range processList {
		if p.RunningInfo != nil && p.IsRunning() {
			memoryInfo, err := p.RunningInfo.GopsutilProcess.MemoryInfo()
			if err == nil {
				vmsSum += memoryInfo.VMS
			}
		}
	}
	return fmt.Sprintf("Σ %s", ByteCountIEC(vmsSum)), nil
}
func (*VmsProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
