package filecache

import (
	"errors"
	"sync"
	"time"

	"github.com/louie-jones-strong/go-shared/collections/maps"
	"github.com/louie-jones-strong/go-shared/logger"
	"github.com/louie-jones-strong/go-shared/storage"
)

type iFileCache interface {
	saveManifest() error
	getFolderPath() string
}

type FileCache[K comparable] struct {
	manifestStore  storage.Storage[map[K]*FileGroupInfo]
	itemFolderPath string

	manifest map[K]*FileGroupInfo
	mu       sync.RWMutex
}

func New[K comparable](
	manifestPath string,
	itemFolderPath string,
) *FileCache[K] {
	return &FileCache[K]{
		manifestStore: storage.NewJSONStorage[map[K]*FileGroupInfo](
			manifestPath,
		),
		itemFolderPath: itemFolderPath,
		manifest:       nil,
	}
}

func (fc *FileCache[K]) getFolderPath() string {
	return fc.itemFolderPath
}

func (fc *FileCache[K]) getManifest() map[K]*FileGroupInfo {
	if fc.manifest == nil {
		manifest, err := fc.manifestStore.Load()
		if err != nil {
			logger.Debug("error loading manifest store")
		}

		if manifest == nil {
			manifest = map[K]*FileGroupInfo{}
		}

		for _, v := range manifest {
			v.setFolder(fc)
		}

		fc.manifest = manifest
	}

	return fc.manifest
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

func (fc *FileCache[K]) CleanupExpiredItems(expireDuration time.Duration) (int, error) {

	manifest := fc.getManifest()
	if manifest == nil {
		return 0, nil
	}

	keysToRemove := make([]K, 0)
	for key := range manifest {
		fi := manifest[key]
		if !fi.IsValid(expireDuration) {
			keysToRemove = append(keysToRemove, key)
		}
	}
	return fc.RemoveFiles(keysToRemove...)
}

func (fc *FileCache[K]) RemoveFiles(keysToRemove ...K) (int, error) {
	manifest := fc.getManifest()
	if manifest == nil {
		logger.Debug("RemoveFiles called but manifest is nil")
		return 0, nil
	}

	fc.mu.Lock()
	defer fc.saveManifest()
	defer fc.mu.Unlock()

	numRemoved := 0
	for _, key := range keysToRemove {
		fi, found := manifest[key]
		if !found {
			logger.Debug("FileCache: No file info found for key during removal: %v", key)
			continue
		}

		// delete the file info from the manifest
		delete(manifest, key)
		err := fi.DeleteFiles()
		if err != nil {
			return numRemoved, err
		}
		numRemoved++
	}

	return numRemoved, nil
}

func (fc *FileCache[K]) GetItems() []maps.KVP[K, *FileGroupInfo] {
	manifest := fc.getManifest()
	return maps.ConvertMapToKVPList(manifest)
}

func (fc *FileCache[K]) TryGetFileInfo(key K) *FileGroupInfo {

	manifest := fc.getManifest()

	fi, found := manifest[key]
	if found {
		return fi
	}
	return nil
}

func (fc *FileCache[K]) GetOrCreateFileInfo(key K) *FileGroupInfo {

	manifest := fc.getManifest()

	fi, found := manifest[key]
	if found {
		return fi
	}

	fi = newFGI()
	fi.setFolder(fc)
	manifest[key] = fi

	return fi
}
