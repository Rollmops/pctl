package app

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

func defaultCompare(value1 string, value2 string, caseInsensitive bool) bool {
	if caseInsensitive {
		return strings.ToLower(value1) == strings.ToLower(value2)
	}
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
	caseInsensitive := false

	if strings.HasPrefix(f.field, "metadata.") {
		fieldValue = process.Config.Metadata[strings.Split(f.field, ".")[1]]
	} else if strings.HasPrefix(f.field, "state.") {
		caseInsensitive = true
		filterField := strings.Split(f.field, ".")[1]
		switch filterField {
		case "pid":
			if process.Info != nil {
				fieldValue = strconv.Itoa(int(process.Info.GoPsutilProcess.Pid))
			} else {
				fieldValue = "-1"
			}
		case "running":
			fieldValue = strconv.FormatBool(process.IsRunning())
		case "stopped":
			fieldValue = strconv.FormatBool(!process.IsRunning())
		case "dirty":
			if process.Info != nil {
				fieldValue = strconv.FormatBool(process.Info.Dirty)
			} else {
				fieldValue = "false"
			}
		case "dirtyConfig":
			if process.Info != nil {
				fieldValue = strconv.FormatBool(process.Info.DirtyCommand)
			} else {
				fieldValue = "false"
			}
		case "dirtyMd5":
			if process.Info != nil {
				fieldValue = strconv.FormatBool(len(process.Info.DirtyMd5Hashes) > 0)
			} else {
				fieldValue = "false"
			}
		}
	}

	switch f.operator {
	case "=":
		return defaultCompare(fieldValue, f.value, caseInsensitive), nil
	case "==":
		return defaultCompare(fieldValue, f.value, caseInsensitive), nil
	case "!=":
		return !defaultCompare(fieldValue, f.value, caseInsensitive), nil
	}

	return false, nil
}
