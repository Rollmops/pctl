package app

type CommentProperty struct{}

func init() {
	PropertyMap["comment"] = &CommentProperty{}
}

func (*CommentProperty) Name() string {
	return "Comment"
}

func (*CommentProperty) Value(p *Process, _ bool) (string, error) {
	if p.RunningInfo != nil {
		return p.RunningInfo.Comment, nil
	}
	return "", nil
}

func (*CommentProperty) FormattedSumValue(_ ProcessList) (string, error) {
	return "", nil
}
func (*CommentProperty) GetComparator() PropertyComparator {
	return &StringPropertyComparator{}
}
