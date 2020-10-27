package app

type UsernameProperty struct{}

func init() {
	usernameProperty := &UsernameProperty{}
	PropertyMap["username"] = usernameProperty
	PropertyMap["user"] = usernameProperty
}

func (*UsernameProperty) Name() string {
	return "Username"
}

func (*UsernameProperty) Value(p *Process, _ bool) (string, error) {
	if p.RunningInfo != nil && p.IsRunning() {
		username, err := p.RunningInfo.GopsutilProcess.Username()
		if err != nil {
			return err.Error(), nil
		} else {
			return username, nil
		}
	}
	return "", nil
}

func (*UsernameProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*UsernameProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
