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
	if p.RunningInfo != nil && p.IsRunning() {
		createTime, err := p.RunningInfo.GopsutilProcess.CreateTime()
		if err != nil {
			uptime = err.Error()
		}
		uptimeInt := time.Now().Sub(time.Unix(createTime/1000, 0))
		uptime = DurationToString(uptimeInt)
	}
	return uptime, nil
}

func (*UptimeProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*UptimeProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
