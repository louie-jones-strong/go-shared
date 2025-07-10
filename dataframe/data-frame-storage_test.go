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

func TestDataFrameStorage_LargeFile(t *testing.T) {
	storage := NewDataFrameStorage(path.Join("test-data", "electric_usage.csv"))

	// Test Loading
	res, err := storage.Load()
	assert.NoError(t, err)

	str := `[4x4560] DataFrame
      Consumption (kwh) Estimated Cost Inc. Tax (p) Start                     End
   0: 0.008             0.21                        2025-03-12T00:00:00Z      2025-03-12T00:30:00Z
   1: 0.008             0.21                        2025-03-12T00:30:00Z      2025-03-12T01:00:00Z
   2: 0.008             0.21                        2025-03-12T01:00:00Z      2025-03-12T01:30:00Z
   3: 0.008             0.21                        2025-03-12T01:30:00Z      2025-03-12T02:00:00Z
   4: 0.008             0.21                        2025-03-12T02:00:00Z      2025-03-12T02:30:00Z
      ...               ...                         ...                       ...
4557: 0.105             2.91                        2025-06-14T23:30:00+01:00 2025-06-15T00:00:00+01:00
4558: 0.091             2.52                        2025-06-15T00:00:00+01:00 2025-06-15T00:30:00+01:00
4559: 0.096             2.66                        2025-06-15T00:30:00+01:00 2025-06-15T01:00:00+01:00
      <float>           <float>                     <datetime>                <datetime>`
	assert.Equal(t, str, res.String())
}
