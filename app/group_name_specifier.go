package app

import (
	"fmt"
	"github.com/minio/minio/pkg/wildcard"
	"regexp"
)

var groupNameSpecifierPatternRegex *regexp.Regexp

func init() {
	groupNameSpecifierPatternRegex = regexp.MustCompile(`([.a-zA-Z0-9_-]*)(:?)([.a-zA-Z0-9_-]*)`)
}

type GroupNameSpecifier struct {
	name  string
	group string
}

func NewGroupNameSpecifier(specifier string) (*GroupNameSpecifier, error) {
	match := groupNameSpecifierPatternRegex.FindStringSubmatch(specifier)
	if len(match) != 4 {
		return nil, fmt.Errorf("invalid name specifier pattern: '%s'", specifier)
	}
	group := match[1]
	name := match[3]
	if match[2] == "" {
		name = group
		group = "*"
	}
	if name == "" {
		name = "*"
	}
	if group == "" {
		group = "*"
	}

	return &GroupNameSpecifier{
		group: group,
		name:  name,
	}, nil
}

func (s *GroupNameSpecifier) IsMatchingGroupAndName(group string, name string) bool {
	return wildcard.Match(s.group, group) && wildcard.Match(s.name, name)
}
