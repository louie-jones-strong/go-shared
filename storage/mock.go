package storage

type MockStorage[M any] struct {
	Data      M
	SaveError error
	LoadError error
}

// Save saves the data to the storage
func (s *MockStorage[M]) Save(_ M) error {
	return s.SaveError
}

// Load loads the data from the storage
func (s *MockStorage[M]) Load() (M, error) {

	return s.Data, s.LoadError
}
