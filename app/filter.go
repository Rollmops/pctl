package app

import (
	"fmt"
	"regexp"
	"strings"
)

var filterPatternRegex *regexp.Regexp

func init() {
	filterPatternRegex = regexp.MustCompile(`([.a-zA-Z0-9_-]+)([!=<>]+)([.a-zA-Z0-9_-]+)`)
}

type Filter struct {
	field    string
	operator string
	value    string
}

func NewFilter(pattern string) (*Filter, error) {
	match := filterPatternRegex.FindStringSubmatch(pattern)
	if len(match) != 4 {
		return nil, fmt.Errorf("invalid filter pattern: '%s'", pattern)
	}
	return &Filter{
		field:    match[1],
		operator: match[2],
		value:    match[3],
	}, nil
}

func (f *Filter) IsMatchingProcess(process *Process) (bool, error) {

	fieldValue := ""

	var comparator PropertyComparator
	if strings.HasPrefix(f.field, "metadata.") {
		fieldValue = process.Config.Metadata[strings.Split(f.field, ".")[1]]
		comparator = &StringPropertyComparator{}
	} else if f.field == "group" {
		fieldValue = process.Config.Group
		comparator = &StringPropertyComparator{}
	} else if strings.HasPrefix(f.field, "state.") || strings.HasPrefix(f.field, "property.") {
		propertyId := strings.Split(f.field, ".")[1]
		property := PropertyMap[propertyId]
		if property == nil {
			return false, fmt.Errorf("property/state '%s' not available", propertyId)
		}
		var err error
		fieldValue, err = property.Value(process, false)
		if err != nil {
			return false, err
		}
		comparator = property.GetComparator()
	} else {
		return false, fmt.Errorf("invalid filter field: '%s'", f.field)
	}

	switch f.operator {
	case "=":
		return comparator.Equal(fieldValue, f.value)
	case "==":
		return comparator.Equal(fieldValue, f.value)
	case "!=":
		isEqual, err := comparator.Equal(fieldValue, f.value)
		if err != nil {
			return false, err
		}
		return !isEqual, nil
	case ">":
		return comparator.Greater(fieldValue, f.value)
	}

	return false, nil
}
