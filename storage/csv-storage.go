package storage

import (
	"os"

	"github.com/gocarina/gocsv"
)

type CSVStorage[R any] struct {
	filePath string
}

func NewCSVStorage[R any](filePath string) *CSVStorage[R] {

	return &CSVStorage[R]{
		filePath: filePath,
	}
}

func (s *CSVStorage[R]) Save(rows []R) error {
	file, err := os.Create(s.filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	if err := gocsv.MarshalFile(&rows, file); err != nil {
		return err
	}

	return nil
}

func (s *CSVStorage[R]) Load() ([]R, error) {

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var out []R
	err = gocsv.UnmarshalFile(file, &out)
	if err != nil {
		return nil, err
	}

	return out, nil
}
