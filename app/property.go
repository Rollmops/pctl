package app

var PropertyMap = make(map[string]Property)

type Property interface {
	Name() string
	Value(*Process, bool) (string, error)
	FormattedSumValue(ProcessList) (string, error)
	GetComparator() PropertyComparator
}
