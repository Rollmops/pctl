package config

import (
	"fmt"
)

func (c Config) validate() error {
	var processNames []string

	for _, p := range c.Processes {
		if isInList(processNames, p.Name) {
			return fmt.Errorf("Found duplicate process name %s", p.Name)
		}
		processNames = append(processNames, p.Name)
	}
	return nil
}

func isInList(list []string, elem string) bool {
	for _, e := range list {
		if e == elem {
			return true
		}
	}
	return false
}
