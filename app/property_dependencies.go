package app

import "strings"

type DependenciesProperty struct{}

func init() {
	dependenciesProperty := &DependenciesProperty{}
	PropertyMap["deps"] = dependenciesProperty
	PropertyMap["dependencies"] = dependenciesProperty
}

func (*DependenciesProperty) Name() string {
	return "Dependencies"
}

func (*DependenciesProperty) Value(p *Process, _ bool) (string, error) {
	return strings.Join(p.Config.DependsOn, ","), nil
}

func (*DependenciesProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*DependenciesProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
