package storage

// Storage is an interface that defines the methods that a storage service must implement.
// This includes loading and saving of models.
type Storage[M any] interface {
	Load() (M, error)
	Save(M) error
}
