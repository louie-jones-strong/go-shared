package filecache

import (
	"errors"
	"maps"
	"os"
	"time"

	"github.com/louie-jones-strong/go-shared/filecache/fileinfo"
	"github.com/louie-jones-strong/go-shared/logger"
	"github.com/louie-jones-strong/go-shared/storage"

	"github.com/google/uuid"
)

type FileCache[K comparable] struct {
	manifestStore  storage.Storage[map[K]*fileinfo.FileInfo]
	itemFolderPath string

	manifest map[K]*fileinfo.FileInfo
}

func New[K comparable](
	manifestPath string,
	itemFolderPath string,
) *FileCache[K] {
	return &FileCache[K]{
		manifestStore: storage.NewJSONStorage[map[K]*fileinfo.FileInfo](
			manifestPath,
		),
		itemFolderPath: itemFolderPath,
		manifest:       nil,
	}
}

func (fc *FileCache[K]) CleanupExpiredItems(expireDuration time.Duration) (int, error) {

	manifest := fc.getManifest()
	if manifest == nil {
		return 0, nil
	}

	defer fc.saveManifest()

	keys := maps.Keys(manifest)

	numRemoved := 0
	for key := range keys {
		fi := manifest[key]

		if fi.IsValid(expireDuration) {
			continue
		}

		// delete the file info from the manifest
		delete(manifest, key)

		// remove the file
		filePath := fi.GetFilePath(fc.itemFolderPath)
		err := os.Remove(filePath)
		if err != nil {
			return numRemoved, err
		}

		numRemoved++
	}

	return numRemoved, nil
}

func (fc *FileCache[K]) TryLoadFileWithExpire(key K, expireDuration time.Duration) ([]byte, error) {
	fi := fc.tryGetFileInfoWithExpire(key, expireDuration)
	if fi == nil {
		return nil, nil
	}

	filePath := fi.GetFilePath(fc.itemFolderPath)

	data, err := storage.ReadBytesFromFile(filePath)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func (fc *FileCache[K]) TryLoadFile(key K) ([]byte, error) {
	return fc.TryLoadFileWithExpire(key, -1)
}

func (fc *FileCache[K]) SaveFileWithExt(key K, data []byte, ext string) error {

	fi := fc.getOrCreateFileInfo(key, ext)

	filePath := fi.GetFilePath(fc.itemFolderPath)

	err := storage.WriteBytesToFile(filePath, data)
	if err != nil {
		return err
	}

	fi.UpdateLastUpdated()

	err = fc.saveManifest()

	return err
}

func (fc *FileCache[K]) SaveFile(key K, data []byte) error {
	return fc.SaveFileWithExt(key, data, "")
}

func (fc *FileCache[K]) getManifest() map[K]*fileinfo.FileInfo {
	if fc.manifest == nil {
		manifest, err := fc.manifestStore.Load()
		if err != nil {
			logger.Debug("error loading manifest store")
		}

		if manifest == nil {
			manifest = map[K]*fileinfo.FileInfo{}
		}

		fc.manifest = manifest
	}

	return fc.manifest
}

func (fc *FileCache[K]) getOrCreateFileInfo(key K, ext string) *fileinfo.FileInfo {

	manifest := fc.getManifest()

	fi, found := manifest[key]
	if found {
		return fi
	}
	fileName := uuid.New().String()
	fileName += ext

	fi = fileinfo.New(fileName)

	manifest[key] = fi

	return fi
}

func (fc *FileCache[K]) tryGetFileInfo(key K) *fileinfo.FileInfo {

	manifest := fc.getManifest()

	fi, found := manifest[key]
	if !found {
		return nil
	}

	return fi
}

func (fc *FileCache[K]) tryGetFileInfoWithExpire(key K, expireDuration time.Duration) *fileinfo.FileInfo {

	fi := fc.tryGetFileInfo(key)
	if fi == nil {
		return nil
	}

	if !fi.IsValid(expireDuration) {
		return nil
	}

	return fi
}

func (fc *FileCache[K]) saveManifest() error {
	if fc.manifest == nil {
		return errors.New("Cannot save file cache with nil manifest")
	}

	err := fc.manifestStore.Save(fc.manifest)
	if err != nil {
		return err
	}

	return nil
}
