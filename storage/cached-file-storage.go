package storage

import (
	"os"
	"time"

	"github.com/louie-jones-strong/go-shared/logger"
)

type CachedFileStorage[M any] struct {
	subStore Storage[M]

	filePath      string
	lastModTime   *time.Time
	cachedContent M
}

func NewCachedFileStorage[M any](filePath string, subStore Storage[M]) *CachedFileStorage[M] {
	var defaultOut M
	return &CachedFileStorage[M]{
		subStore:      subStore,
		filePath:      filePath,
		lastModTime:   nil,
		cachedContent: defaultOut,
	}
}

func (s *CachedFileStorage[M]) Save(obj M) error {
	var output M
	s.cachedContent = output
	s.lastModTime = nil
	return s.subStore.Save(obj)
}

func (s *CachedFileStorage[M]) Load() (M, error) {

	var defaultOut M
	modTime, err := s.getModTime()
	if err != nil {
		return defaultOut, err
	}

	if s.lastModTime != nil && modTime.Equal(*s.lastModTime) {
		logger.Debug("Cache HIT for: %v", s.filePath)
		return s.cachedContent, nil
	}
	logger.Debug("Cache MISS for: %v", s.filePath)

	output, err := s.subStore.Load()
	if err != nil {
		return defaultOut, err
	}

	s.lastModTime = modTime
	s.cachedContent = output

	return output, nil
}

func (s *CachedFileStorage[M]) getModTime() (*time.Time, error) {
	fileInfo, err := os.Stat(s.filePath)
	if err != nil {
		return nil, err
	}

	modificationTime := fileInfo.ModTime()
	return &modificationTime, nil
}
