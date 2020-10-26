package app

type GroupProperty struct{}

func init() {
	PropertyMap["group"] = &GroupProperty{}
}

func (*GroupProperty) Name() string {
	return "Group"
}

func (*GroupProperty) Value(p *Process, _ bool) (string, error) {
	return p.Config.Group, nil
}

func (*GroupProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*GroupProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
