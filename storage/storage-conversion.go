package storage

type StorageConversion[M any, D any] struct {
	subStore         Storage[D]
	serializerFunc   func(M) (D, error)
	deserializerFunc func(D) (M, error)
}

func NewStorageConversion[M any, D any](
	subStore Storage[D],
	serializerFunc func(M) (D, error),
	deserializerFunc func(D) (M, error),
) *StorageConversion[M, D] {
	return &StorageConversion[M, D]{
		subStore:         subStore,
		serializerFunc:   serializerFunc,
		deserializerFunc: deserializerFunc,
	}
}

func (s *StorageConversion[M, D]) Save(model M) error {

	data, err := s.serializerFunc(model)
	if err != nil {
		return err
	}

	return s.subStore.Save(data)
}

func (s *StorageConversion[M, D]) Load() (M, error) {
	data, err := s.subStore.Load()
	if err != nil {
		var zero M
		return zero, err
	}
	model, err := s.deserializerFunc(data)
	if err != nil {
		var zero M
		return zero, err
	}

	return model, nil
}
