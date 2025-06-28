package storage

import (
	"encoding/csv"
	"os"

	"github.com/louie-jones-strong/go-shared/dataframe"
	"github.com/louie-jones-strong/go-shared/dataframe/series"
)

type DataFrameStorage struct {
	filePath string
}

func NewDataFrameStorage(filePath string) *DataFrameStorage {
	return &DataFrameStorage{
		filePath: filePath,
	}
}

func (s *DataFrameStorage) Save(df *dataframe.DataFrame) error {
	file, err := os.OpenFile(s.filePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)

	records := make([][]string, df.NumRows()+1)
	records[0] = make([]string, df.NumColumns())

	for c := 0; c < df.NumColumns(); c++ {
		col, err := df.GetColumn(c)
		if err != nil {
			return err
		}
		records[0][c] = col.GetName()
	}

	for r := 0; r < df.NumRows(); r++ {
		records[r+1] = make([]string, df.NumColumns())
		for c := 0; c < df.NumColumns(); c++ {
			item, err := df.Get(c, r)
			if err != nil {
				return err
			}

			records[r+1][c] = item.ToString()
		}
	}

	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}

func (s *DataFrameStorage) Load() (*dataframe.DataFrame, error) {

	file, err := os.Open(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	numRows := len(records)
	numCols := len(records[0])

	cols := make([]*series.Series, len(records[0]))
	for colIdx := 0; colIdx < numCols; colIdx++ {

		header := records[0][colIdx]

		values := make([]string, numRows-1)
		for rowIdx := 1; rowIdx < numRows; rowIdx++ {
			values[rowIdx-1] = records[rowIdx][colIdx]
		}

		cols[colIdx] = series.New(header, series.String, values)
	}

	df := dataframe.New(
		cols,
	)

	return df, nil
}
