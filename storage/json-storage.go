package storage

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type JSONStorage[M any] struct {
	filePath string
}

func NewJSONStorage[M any](filePath string) *JSONStorage[M] {

	return &JSONStorage[M]{
		filePath: filePath,
	}
}

func (s *JSONStorage[M]) Save(obj M) error {
	start := time.Now()

	// Marshal the object to JSON
	data, err := json.MarshalIndent(obj, "", "  ")
	if err != nil {
		return err
	}

	// Open the file for writing (create if not exists, truncate if it does)
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the JSON data to the file
	if _, err := file.Write(data); err != nil {
		return err
	}

	log.Printf("Saving file: %v, took: %v", s.filePath, time.Since(start))
	return nil
}

func (s *JSONStorage[M]) Load() (M, error) {
	start := time.Now()

	var output M

	absPath, err := Check(s.filePath)
	if err != nil {
		return output, err
	}

	// Open the file for reading
	file, err := os.Open(absPath)
	if err != nil {
		return output, err
	}
	defer file.Close()

	// Decode the JSON data from the file
	err = json.NewDecoder(file).Decode(&output)
	if err != nil {
		return output, err
	}

	log.Printf("Loading file: %v, took: %v", s.filePath, time.Since(start))
	return output, nil
}
