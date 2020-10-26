package app

import "strconv"

type NiceProperty struct{}

func init() {
	PropertyMap["nice"] = &NiceProperty{}
}

func (*NiceProperty) Name() string {
	return "Nice"
}

func (*NiceProperty) Value(p *Process, _ bool) (string, error) {
	if p.Info != nil && p.IsRunning() {
		nice, err := p.Info.GoPsutilProcess.Nice()
		if err != nil {
			return err.Error(), nil
		} else {
			return strconv.Itoa(int(nice)), nil
		}
	}
	return "", nil
}

func (*NiceProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*NiceProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
