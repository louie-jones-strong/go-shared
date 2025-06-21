package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

func Check(path string) (string, error) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return "", err
	}

	_, statErr := os.Stat(absPath)
	if os.IsNotExist(statErr) {
		return absPath, fmt.Errorf("Could not find file with path: %v", absPath)
	}

	return absPath, nil
}
