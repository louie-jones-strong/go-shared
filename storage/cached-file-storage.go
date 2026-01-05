package storage

import (
	"os"
	"time"

	"github.com/louie-jones-strong/go-shared/cache"
)

type CachedFileStorage[M any] struct {
	cache.BaseCacheScope[M]
	subStore Storage[M]

	filePath    string
	lastModTime *time.Time
}

func NewCachedFileStorage[M any](filePath string, subStore Storage[M]) *CachedFileStorage[M] {
	return &CachedFileStorage[M]{
		BaseCacheScope: cache.NewBaseCacheScope[M](),
		subStore:       subStore,
		filePath:       filePath,
		lastModTime:    nil,
	}
}

func (s *CachedFileStorage[M]) IsValid() bool {
	modTime, err := s.getModTime()
	if err != nil {
		modTime = nil
	}
	return s.isValid(modTime)
}

func (s *CachedFileStorage[M]) isValid(modTime *time.Time) bool {
	if !s.BaseCacheScope.IsValid() {
		return false
	}

	isValid := false
	if s.lastModTime != nil {
		isValid = modTime != nil && modTime.Equal(*s.lastModTime)
	}

	if !isValid {
		s.Clear()
	}
	return isValid
}

func (s *CachedFileStorage[M]) Clear() {
	s.lastModTime = nil
	s.BaseCacheScope.Clear()
}

func (s *CachedFileStorage[M]) ToString() string {
	return "CachedFileStorage: " + s.filePath
}

func (s *CachedFileStorage[M]) Save(obj M) error {
	s.Clear()
	return s.subStore.Save(obj)
}

func (s *CachedFileStorage[M]) Load() (M, error) {
	cs := cache.GetCacheService()
	cs.AddOpenScope(s)
	defer cs.CloseScope(s)

	var defaultOut M
	modTime, err := s.getModTime()
	if err != nil {
		return defaultOut, err
	}

	if s.isValid(modTime) {
		return s.GetValue(), err
	}

	output, err := s.subStore.Load()
	if err != nil {
		return defaultOut, err
	}

	s.lastModTime = modTime
	s.BaseCacheScope.SetValue(output)

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
