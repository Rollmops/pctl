package config

import (
	"fmt"
)

func (c *Config) Validate() error {
	var processNames []string

	for _, p := range c.Processes {
		if _isInList(processNames, p.Name) {
			return fmt.Errorf("found duplicate process name %s", p.Name)
		}
		processNames = append(processNames, p.Name)
	}
	return nil
}

func _isInList(list []string, elem string) bool {
	for _, e := range list {
		if e == elem {
			return true
		}
	}
	return false
}
