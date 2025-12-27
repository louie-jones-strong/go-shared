package storage

import (
	"encoding/json"
	"sync"
)

type JSONStorage[M any] struct {
	filePath string
	mu       sync.RWMutex
}

func NewJSONStorage[M any](filePath string) *JSONStorage[M] {

	return &JSONStorage[M]{
		filePath: filePath,
	}
}

func (s *JSONStorage[M]) Save(obj M) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Marshal the object to JSON
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	file, err := OpenFileForWriting(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	if _, err := file.Write(data); err != nil {
		return err
	}

	return nil
}

func (s *JSONStorage[M]) Load() (M, error) {
	var output M
	s.mu.RLock()
	defer s.mu.RUnlock()
	file, err := OpenFileForReading(s.filePath)
	if err != nil {
		return output, err
	}
	defer file.Close()

	// Decode the JSON data from the file
	err = json.NewDecoder(file).Decode(&output)
	if err != nil {
		return output, err
	}

	return output, nil
}
