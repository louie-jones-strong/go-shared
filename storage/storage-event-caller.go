package storage

type EventCallableStorageItem interface {
	OnLoad() error
	OnSave() error
}

type StorageEventCaller[M EventCallableStorageItem] struct {
	subStore Storage[M]
}

func NewStorageEventCaller[M EventCallableStorageItem](subStore Storage[M]) *StorageEventCaller[M] {
	return &StorageEventCaller[M]{
		subStore: subStore,
	}
}

func (s *StorageEventCaller[M]) Save(obj M) error {
	err := obj.OnSave()
	if err != nil {
		return err
	}
	return s.subStore.Save(obj)
}

func (s *StorageEventCaller[M]) Load() (M, error) {
	var output M

	obj, err := s.subStore.Load()
	if err != nil {
		return output, err
	}

	err = obj.OnLoad()
	if err != nil {
		return output, err
	}

	return obj, nil
}
