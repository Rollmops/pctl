package app

import (
	"fmt"
	"strings"
)

type EnvProperty struct{}

func init() {
	PropertyMap["env"] = &EnvProperty{}
}

func (*EnvProperty) Name() string {
	return "Environment"
}

func (*EnvProperty) Value(p *Process, _ bool) (string, error) {
	var envStrings []string
	for key, value := range p.Config.Env {
		envStrings = append(envStrings, fmt.Sprintf("%s=%s", key, value))
	}
	return strings.Join(envStrings, ","), nil
}

func (*EnvProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*EnvProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
