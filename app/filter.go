package app

import (
	"fmt"
	"regexp"
	"strings"
)

func defaultCompare(value1 string, value2 string) bool {
	return value1 == value2
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

var filterPatternRegex *regexp.Regexp

func init() {
	filterPatternRegex = regexp.MustCompile(`([.a-zA-Z0-9_-]+)([!=<>]+)([.a-zA-Z0-9_-]+)`)
}

func (f *Filter) IsMatchingProcess(process *Process) (bool, error) {

	fieldValue := ""

	if strings.HasPrefix(f.field, "metadata.") {
		fieldValue = process.Config.Metadata[strings.Split(f.field, ".")[1]]
	} else if strings.HasPrefix(f.field, "state.") {
	}

	switch f.operator {
	case "=":
		return defaultCompare(fieldValue, f.value), nil
	case "==":
		return defaultCompare(fieldValue, f.value), nil
	case "!=":
		return !defaultCompare(fieldValue, f.value), nil
	}

	return false, nil

}
