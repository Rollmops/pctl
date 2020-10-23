package config

import (
	"fmt"
	"regexp"
	"strings"
)

func defaultCompare(value1 string, value2 string) bool {
	return value1 == value2
}

var filterPatternRegex *regexp.Regexp

func init() {
	filterPatternRegex = regexp.MustCompile(`([.a-zA-Z0-9_-]+)([!=<>]+)([.a-zA-Z0-9_-]+)`)
}

func (p *ProcessConfig) IsRelevantForFilter(filterPattern string) (bool, error) {
	match := filterPatternRegex.FindStringSubmatch(filterPattern)

	if len(match) != 4 {
		return false, fmt.Errorf("invalid filter pattern: '%s'", filterPattern)
	}

	fieldValue := ""
	field := match[1]
	value := match[3]

	if strings.HasPrefix(field, "metadata.") {
		fieldValue = p.Metadata[strings.Split(field, ".")[1]]
	}

	switch operator := match[2]; operator {
	case "=":
		return defaultCompare(fieldValue, value), nil
	case "==":
		return defaultCompare(fieldValue, value), nil
	case "!=":
		return !defaultCompare(fieldValue, value), nil
	}

	return false, nil

}
