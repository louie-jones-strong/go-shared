package storage

type EventCallableStorageItem interface {
	OnLoad() error
	OnSave() error
}

type StorageEventCaller[M EventCallableStorageItem] struct {
	subStore Storage[[]M]
}

func NewStorageEventCaller[M EventCallableStorageItem](subStore Storage[[]M]) *StorageEventCaller[M] {
	return &StorageEventCaller[M]{
		subStore: subStore,
	}
}

func (s *StorageEventCaller[M]) Save(items []M) error {
	for _, item := range items {
		err := item.OnSave()
		if err != nil {
			return err
		}
	}
	return s.subStore.Save(items)
}

func (s *StorageEventCaller[M]) Load() ([]M, error) {

	items, err := s.subStore.Load()
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		err := item.OnLoad()
		if err != nil {
			return nil, err
		}
	}

	return items, nil
}
