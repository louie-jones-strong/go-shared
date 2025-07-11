package dataframe

import (
	"encoding/csv"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/louie-jones-strong/go-shared/storage"
)

type DataFrameStorage struct {
	filePath string
}

func NewDataFrameStorage(filePath string) *DataFrameStorage {
	return &DataFrameStorage{
		filePath: filePath,
	}
}

func (s *DataFrameStorage) Save(df *DataFrame) error {
	file, err := storage.OpenFileForWriting(s.filePath)
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

func (s *DataFrameStorage) Load() (*DataFrame, error) {

	file, err := storage.OpenFileForReading(s.filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.TrimLeadingSpace = true

	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	numRows := len(records)
	numCols := len(records[0])

	cols := make([]*series.Series, len(records[0]))
	for colIdx := 0; colIdx < numCols; colIdx++ {

		header := records[0][colIdx]

		strings := make([]string, numRows-1)
		for rowIdx := 1; rowIdx < numRows; rowIdx++ {
			strings[rowIdx-1] = records[rowIdx][colIdx]
		}

		colType, err := apptype.FindType(strings)
		if err != nil {
			return nil, err
		}

		var values any
		values = strings
		switch colType {
		case apptype.DateTime:
			values, err = apptype.ConvertArr(strings, apptype.StringToTime)
		case apptype.Int:
			values, err = apptype.ConvertArr(strings, apptype.StringToInt)
		case apptype.Float:
			values, err = apptype.ConvertArr(strings, apptype.StringToFloat)
		case apptype.Bool:
			values, err = apptype.ConvertArr(strings, apptype.StringToBool)
		}
		if err != nil {
			return nil, err
		}

		cols[colIdx] = series.BuildSeries(header, values)
	}

	df := New(
		cols,
	)

	return df, nil
}
