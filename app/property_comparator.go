package app

import (
	"fmt"
	"github.com/minio/minio/pkg/wildcard"
	"strings"
)

type PropertyComparator interface {
	Equal(string, string) (bool, error)
	Greater(string, string) (bool, error)
}

type StringPropertyComparator struct{}
type FloatPropertyComparator struct{}

func (c *StringPropertyComparator) Equal(value1 string, value2 string) (bool, error) {
	preparedValue1 := strings.TrimSpace(strings.ToLower(value1))
	preparedValue2 := strings.TrimSpace(strings.ToLower(value2))

	return wildcard.Match(preparedValue2, preparedValue1), nil
}

func (c *StringPropertyComparator) Greater(_ string, _ string) (bool, error) {
	return false, fmt.Errorf("operator not supported")
}
