package app

type CwdProperty struct{}

func init() {
	PropertyMap["cwd"] = &CwdProperty{}
}

func (*CwdProperty) Name() string {
	return "Cwd"
}

func (*CwdProperty) Value(p *Process, _ bool) (string, error) {
	if p.Info != nil && p.IsRunning() {
		cwd, err := p.Info.GoPsutilProcess.Cwd()
		if err != nil {
			return err.Error(), nil
		} else {
			return cwd, nil
		}
	}
	return "", nil
}

func (*CwdProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*CwdProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
