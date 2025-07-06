package dataframe

import (
	"path"
	"testing"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func TestDataFrameStorage_CSV(t *testing.T) {
	storage := NewDataFrameStorage(path.Join("test-data", "df_test.csv"))

	expectedData := New(
		[]*series.Series{
			series.New("name", apptype.String, []string{"Alice", "Bob"}),
			series.New("age", apptype.String, []string{"20", "22"}),
			series.New("grade", apptype.String, []string{"A", "B"}),
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
