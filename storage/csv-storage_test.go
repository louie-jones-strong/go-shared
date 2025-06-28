package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	Name  string `csv:"name"`
	Age   int    `csv:"age"`
	Grade string `csv:"grade"`
}

func TestCSVStorage(t *testing.T) {
	storage := NewCSVStorage[*Student]("csv_test.csv")

	expectedData := []*Student{
		{Name: "Alice", Age: 20, Grade: "A"},
		{Name: "Bob", Age: 22, Grade: "B"},
	}

	// Test Saving
	err := storage.Save(expectedData)
	assert.NoError(t, err)

	// Test Loading
	loadedStudents, err := storage.Load()
	assert.NoError(t, err)
	assert.Equal(t, expectedData, loadedStudents)

}
