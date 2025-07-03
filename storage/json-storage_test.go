package storage

import (
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Save(t *testing.T) {
	storage := NewJSONStorage[*Student](path.Join("test-data", "json_test.json"))

	expectedData := &Student{
		Name:  "Alice",
		Age:   20,
		Grade: "A",
	}

	// Test Saving
	err := storage.Save(expectedData)
	assert.NoError(t, err)

	// Test Loading
	loadedStudents, err := storage.Load()
	assert.NoError(t, err)
	assert.Equal(t, expectedData, loadedStudents)
}
