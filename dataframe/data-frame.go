package dataframe

import (
	"fmt"

	"github.com/louie-jones-strong/go-shared/dataframe/series"
)

type DataFrame struct {
	nameMap map[string]int
	columns []*series.Series
}

func New(
	columns []*series.Series,
) *DataFrame {
	return &DataFrame{
		columns: columns,
	}
}

func (df *DataFrame) NumColumns() int {
	return len(df.columns)
}

func (df *DataFrame) NumRows() int {
	if df.NumColumns() == 0 {
		return 0
	}

	return df.columns[0].Len()
}

func (df *DataFrame) AddColumns(columns []*series.Series) error {
	for _, column := range columns {
		err := df.AddColumn(column)
		if err != nil {
			return err
		}
	}
	return nil
}

func (df *DataFrame) AddColumn(column *series.Series) error {
	if column.Len() != df.NumRows() {
		return fmt.Errorf("Add Column called with column length not matching the dataframe's length")
	}

	df.columns = append(df.columns, column)
	return nil
}

func (df *DataFrame) AddRows(rows [][]any) error {
	for _, row := range rows {
		err := df.AddRow(row)
		if err != nil {
			return err
		}
	}
	return nil
}

func (df *DataFrame) AddRow(row []any) error {
	if len(row) != len(df.columns) {
		return fmt.Errorf("Add row called with row length not matching the dataframe's length")
	}

	for idx := range df.columns {
		df.columns[idx].Append(row[idx])
	}
	return nil
}

func (df *DataFrame) GetRow(rowIdx int) ([]series.Element, error) {
	if rowIdx < 0 || rowIdx >= df.NumRows() {
		return nil, fmt.Errorf(
			"GetRow called with rowIdx: %v out or range %v - %v",
			rowIdx,
			0,
			df.NumRows(),
		)
	}

	row := make([]series.Element, df.NumColumns())
	for c := 0; c < df.NumColumns(); c++ {
		row[c] = df.columns[c].Elem(rowIdx)
	}

	return row, nil
}

func (df *DataFrame) GetColumnByName(columnName string) (*series.Series, error) {
	columnIdx, found := df.nameMap[columnName]
	if !found {
		return nil, fmt.Errorf("column name not found in dataframe: %v", columnName)
	}

	return df.GetColumn(columnIdx)
}

func (df *DataFrame) GetColumn(columnIdx int) (*series.Series, error) {
	if columnIdx < 0 || columnIdx >= df.NumColumns() {
		return nil, fmt.Errorf(
			"GetColumnByIdx called with columnIdx: %v out or range %v - %v",
			columnIdx,
			0,
			df.NumColumns(),
		)
	}

	return df.columns[columnIdx], nil
}

func (df *DataFrame) GetByName(columnName string, rowIdx int) (series.Element, error) {
	col, err := df.GetColumnByName(columnName)
	if err != nil {
		return nil, err
	}

	if rowIdx < 0 || rowIdx >= df.NumRows() {
		return nil, fmt.Errorf(
			"Get called with rowIdx: %v out or range %v - %v",
			rowIdx,
			0,
			df.NumRows(),
		)
	}

	item := col.Elem(rowIdx)
	return item, nil
}

func (df *DataFrame) Get(columnIdx int, rowIdx int) (series.Element, error) {
	col, err := df.GetColumn(columnIdx)
	if err != nil {
		return nil, err
	}

	if rowIdx < 0 || rowIdx >= df.NumRows() {
		return nil, fmt.Errorf(
			"Get called with rowIdx: %v out or range %v - %v",
			rowIdx,
			0,
			df.NumRows(),
		)
	}

	item := col.Elem(rowIdx)
	return item, nil
}
