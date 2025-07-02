package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

func createDir(filePath string) error {
	dir := filepath.Dir(filePath)

	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return fmt.Errorf("Error creating directory: %v", err)
	}
	return nil
}

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

func OpenFileForReading(filePath string) (*os.File, error) {
	filePath, err := Check(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}

	return file, nil
}

func OpenFileForWriting(filePath string) (*os.File, error) {

	err := createDir(filePath)
	if err != nil {
		return nil, err
	}

	file, err := os.OpenFile(filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	return file, nil
}
