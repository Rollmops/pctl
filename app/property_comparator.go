package app

import (
	"fmt"
	"strings"
)

type PropertyComparator interface {
	Equal(interface{}, interface{}) (bool, error)
	Greater(interface{}, interface{}) (bool, error)
	GreaterEqual(interface{}, interface{}) (bool, error)
}

type StringPropertyComparator struct{}
type FloatPropertyComparator struct{}

func (c *StringPropertyComparator) Equal(value1 interface{}, value2 interface{}) (bool, error) {
	value1String := value1.(string)
	value2String := value2.(string)

	return strings.TrimSpace(strings.ToLower(value1String)) == strings.TrimSpace(strings.ToLower(value2String)), nil
}

func (c *StringPropertyComparator) Greater(_ interface{}, _ interface{}) (bool, error) {
	return false, fmt.Errorf("operator not supported")
}

func (c *StringPropertyComparator) GreaterEqual(_ interface{}, _ interface{}) (bool, error) {
	return false, fmt.Errorf("operator not supported")
}
