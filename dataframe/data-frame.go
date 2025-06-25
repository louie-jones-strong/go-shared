package dataframe

import "github.com/louie-jones-strong/go-shared/dataframe/series"

type DataFrame struct {
	columns []*series.Series
}

func New(
	columns []*series.Series,
) *DataFrame {
	return &DataFrame{
		columns: columns,
	}
}

func (df *DataFrame) AddRows(rows [][]any) {
	for _, row := range rows {
		df.AddRow(row)
	}
}

func (df *DataFrame) AddRow(row []any) {
	if len(row) != len(df.columns) {
		return
	}

	for idx := range df.columns {
		df.columns[idx].Append(row[idx])
	}
}
