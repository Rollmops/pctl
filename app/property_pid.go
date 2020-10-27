package app

import (
	"strconv"
)

type PidProperty struct{}

func init() {
	PropertyMap["pid"] = &PidProperty{}
}

func (*PidProperty) Name() string {
	return "Pid"
}

func (*PidProperty) Value(p *Process, _ bool) (string, error) {
	if p.IsRunning() {
		return strconv.Itoa(int(p.RunningInfo.Pid)), nil
	}
	return "", nil
}

func (*PidProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*PidProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
