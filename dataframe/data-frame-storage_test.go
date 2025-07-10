package dataframe

import (
	"path"
	"testing"
	"time"

	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func TestDataFrameStorage_CSV(t *testing.T) {
	storage := NewDataFrameStorage(path.Join("test-data", "df_test.csv"))

	dateTime := time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC)
	expectedData := New(
		[]*series.Series{
			series.New("name", []string{"Alice", "Bob"}),
			series.New("age", []int{20, 22}),
			series.New("marks", []float64{100.1, 60}),
			series.New("grade", []string{"A", "B"}),
			series.New("date", []time.Time{dateTime, dateTime}),
		},
	)

	// Test Saving
	err := storage.Save(expectedData)
	assert.NoError(t, err)

	// Test Loading
	res, err := storage.Load()
	assert.NoError(t, err)
	assert.Equal(t, expectedData, res)
}
