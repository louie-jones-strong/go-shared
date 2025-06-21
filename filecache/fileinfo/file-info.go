package fileinfo

import (
	"path/filepath"
	"time"
)

type FileInfo struct {
	FileName string `json:"file-name"`
	// last time this file was updated as a unix timestamp
	LastUpdatedTimestamp int64 `json:"last-updated"`
}

func New(
	fileName string,
) *FileInfo {

	fileInfo := &FileInfo{
		FileName:             fileName,
		LastUpdatedTimestamp: 0,
	}

	fileInfo.UpdateLastUpdated()
	return fileInfo
}

func (fi *FileInfo) UpdateLastUpdated() {
	fi.LastUpdatedTimestamp = time.Now().Unix()
}

func (fi *FileInfo) GetLastUpdated() time.Time {
	return time.Unix(fi.LastUpdatedTimestamp, 0)
}

func (fi *FileInfo) IsValid(expireDuration time.Duration) bool {

	if expireDuration < 0 {
		return true
	}

	expireThreshold := time.Now().UTC().Add(-expireDuration)

	lastUpdated := fi.GetLastUpdated()
	return lastUpdated.After(expireThreshold)
}

func (fi *FileInfo) GetFilePath(folderPath string) string {
	return filepath.Join(folderPath, fi.FileName)
}
