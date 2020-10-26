package app

import (
	"strconv"
	"strings"
)

type GidsProperty struct{}

func init() {
	PropertyMap["gids"] = &GidsProperty{}
}

func (*GidsProperty) Name() string {
	return "Gids"
}

func (*GidsProperty) Value(p *Process, _ bool) (string, error) {
	if p.Info != nil && p.IsRunning() {
		gids, err := p.Info.GoPsutilProcess.Gids()
		if err != nil {
			return err.Error(), nil
		} else {
			var gidsString []string
			for _, gid := range gids {
				gidsString = append(gidsString, strconv.Itoa(int(gid)))
			}
			return strings.Join(gidsString, ","), nil
		}
	}
	return "", nil
}

func (*GidsProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*GidsProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
