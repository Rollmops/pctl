package app

import "time"

type UptimeProperty struct{}

func init() {
	PropertyMap["uptime"] = &UptimeProperty{}
}

func (*UptimeProperty) Name() string {
	return "Uptime"
}

func (*UptimeProperty) Value(p *Process, _ bool) (string, error) {
	var uptime string
	if p.Info != nil && p.IsRunning() {
		createTime, err := p.Info.GoPsutilProcess.CreateTime()
		if err != nil {
			uptime = "error"
		}
		uptimeInt := time.Now().Sub(time.Unix(createTime/1000, 0))
		uptime, err = DurationToString(uptimeInt)
		if err != nil {
			uptime = "error"
		}
	}
	return uptime, nil
}

func (*UptimeProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*UptimeProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
