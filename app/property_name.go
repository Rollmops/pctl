package app

type NameProperty struct{}

func init() {
	PropertyMap["name"] = &NameProperty{}
}

func (*NameProperty) Name() string {
	return "Name"
}

func (*NameProperty) Value(p *Process, _ bool) (string, error) {
	return p.Config.Name, nil
}

func (*NameProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*NameProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
