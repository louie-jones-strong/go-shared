package storage

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/dataframe"
	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func TestDataFrameStorage_CSV(t *testing.T) {
	storage := NewDataFrameStorage("df_test.csv")

	expectedData := dataframe.New(
		[]*series.Series{
			series.New("name", series.String, []string{"Alice", "Bob"}),
			series.New("age", series.String, []string{"20", "22"}),
			series.New("grade", series.String, []string{"A", "B"}),
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
