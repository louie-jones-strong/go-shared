package condition

type MockCondition[C any] struct {
	value bool
	err   error
}

// NewMockCondition creates a new MockCondition with the specified value and error.
func NewMockCondition[C any](value bool, err error) *MockCondition[C] {
	return &MockCondition[C]{
		value: value,
		err:   err,
	}
}

// Evaluate returns the predefined value and error of the MockCondition.
func (c *MockCondition[C]) Evaluate(_ C) (bool, error) {
	return c.value, c.err
}
