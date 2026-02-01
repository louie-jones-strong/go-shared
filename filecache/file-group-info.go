package filecache

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
	"github.com/louie-jones-strong/go-shared/storage"
)

const (
	DefaultFileKey string = "default"
)

type FileGroupInfo struct {
	Files map[string]string `json:"files"`

	// first time this file was created
	CreatedTimestamp int64 `json:"created-timestamp"`
	// last time this file was updated as a unix timestamp
	LastUpdatedTimestamp int64 `json:"last-updated"`

	cacheHolder iFileCache
}

func newFGI() *FileGroupInfo {

	fileInfo := &FileGroupInfo{
		Files:                map[string]string{},
		CreatedTimestamp:     time.Now().Unix(),
		LastUpdatedTimestamp: time.Now().Unix(),
	}
	return fileInfo
}

func (fi *FileGroupInfo) setFolder(cacheHolder iFileCache) {
	fi.cacheHolder = cacheHolder
}

func (fi *FileGroupInfo) GetCreatedTimestamp() time.Time {
	return time.Unix(fi.CreatedTimestamp, 0)
}

func (fi *FileGroupInfo) GetLastUpdated() time.Time {
	return time.Unix(fi.LastUpdatedTimestamp, 0)
}

func (fi *FileGroupInfo) IsValid(expireDuration time.Duration) bool {
	if fi == nil {
		return false
	}
	if len(fi.Files) == 0 {
		return false
	}

	if expireDuration < 0 {
		return true
	}

	expireThreshold := time.Now().UTC().Add(-expireDuration)

	lastUpdated := fi.GetLastUpdated()
	return lastUpdated.After(expireThreshold)
}

func (fi *FileGroupInfo) SaveFile(key string, ext string, content []byte) error {
	if fi == nil {
		return fmt.Errorf("SaveFile called with nil FileGroupInfo")
	}
	if fi.cacheHolder == nil {
		return fmt.Errorf("SaveFile called with nil cacheHolder")
	}

	unlock := fi.cacheHolder.wLock()
	defer unlock()

	fileName, exists := fi.Files[key]
	if !exists {
		fileName = uuid.New().String() + ext
		fi.Files[key] = fileName
	}
	filePath := filepath.Join(fi.cacheHolder.getFolderPath(), fileName)
	err := storage.WriteBytesToFile(filePath, content)
	if err != nil {
		return err
	}
	fi.LastUpdatedTimestamp = time.Now().Unix()

	err = fi.cacheHolder.saveManifest()
	if err != nil {
		return err
	}

	return nil
}

func (fi *FileGroupInfo) LoadFiles(keys ...string) (map[string][]byte, error) {
	if fi == nil {
		return nil, fmt.Errorf("LoadFiles called with nil FileGroupInfo")
	}
	if fi.cacheHolder == nil {
		return nil, fmt.Errorf("LoadFiles called with nil cacheHolder")
	}

	unlock := fi.cacheHolder.rLock()
	defer unlock()

	result := make(map[string][]byte)
	if len(keys) > 0 {
		for _, key := range keys {
			data, err := fi.LoadFile(key)
			if err != nil {
				return nil, err
			}
			result[key] = data
		}
	} else {
		for key := range fi.Files {
			data, err := fi.LoadFile(key)
			if err != nil {
				return nil, err
			}
			result[key] = data
		}
	}

	return result, nil
}

func (fi *FileGroupInfo) LoadFile(key string) ([]byte, error) {
	if fi == nil {
		return nil, fmt.Errorf("LoadFile called with nil FileGroupInfo")
	}
	if fi.cacheHolder == nil {
		return nil, fmt.Errorf("LoadFile called with nil cacheHolder")
	}

	unlock := fi.cacheHolder.rLock()
	defer unlock()

	fileName, exists := fi.Files[key]
	if !exists {
		return nil, fmt.Errorf("file with key %s does not exist", key)
	}

	filePath := filepath.Join(fi.cacheHolder.getFolderPath(), fileName)
	data, err := storage.ReadBytesFromFile(filePath)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (fi *FileGroupInfo) deleteFiles() error {
	if fi == nil {
		return fmt.Errorf("deleteFiles called with nil FileGroupInfo")
	}
	if fi.cacheHolder == nil {
		return fmt.Errorf("deleteFiles called with nil cacheHolder")
	}

	unlock := fi.cacheHolder.rLock()
	defer unlock()

	for _, fileName := range fi.Files {
		filePath := filepath.Join(fi.cacheHolder.getFolderPath(), fileName)
		err := os.Remove(filePath)
		if err != nil {
			return err
		}
	}
	return nil
}
