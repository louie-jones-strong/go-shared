package storage

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Student struct {
	Name  string `csv:"name"`
	Age   int    `csv:"age"`
	Grade string `csv:"grade"`
}

func TestCSVStorage(t *testing.T) {
	storage := NewCSVStorage[*Student]("test.csv")

	// Test Save
	students := []*Student{
		{Name: "Alice", Age: 20, Grade: "A"},
		{Name: "Bob", Age: 22, Grade: "B"},
	}
	err := storage.Save(students)
	assert.NoError(t, err)

	// Test Load
	loadedStudents, err := storage.Load()
	assert.NoError(t, err)
	assert.Equal(t, students, loadedStudents)

	// Clean up
	os.Remove("test_students.csv")
}
