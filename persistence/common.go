package persistence

import (
	"fmt"
	"os"
	"path/filepath"
)

func GetStateFilePath() (string, error) {
	home := os.Getenv("HOME")
	if home == "" {
		return "", fmt.Errorf("unable to retrieve HOME env var")
	}
	stateFilePath := filepath.Join(home, "var", "pctl")

	if _, err := os.Stat(stateFilePath); err != nil {
		err = os.MkdirAll(stateFilePath, os.ModePerm)
		if err != nil {
			return "", err
		}
	}
	return stateFilePath, nil
}
