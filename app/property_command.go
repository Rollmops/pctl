package app

import "strings"

type CommandProperty struct{}

func init() {
	commandProperty := &CommandProperty{}
	PropertyMap["cmd"] = commandProperty
	PropertyMap["command"] = commandProperty
}

func (*CommandProperty) Name() string {
	return "Command"
}

func (*CommandProperty) Value(p *Process, formatted bool) (string, error) {
	command := strings.Join(p.Config.Command, " ")
	if p.RunningInfo != nil {
		command = strings.Join(p.RunningInfo.Config.Command, " ")
		if p.RunningInfo.DirtyInfo.DirtyCommand && formatted {
			return FailedColor(command), nil
		}
	}
	return command, nil
}

func (*CommandProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*CommandProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
