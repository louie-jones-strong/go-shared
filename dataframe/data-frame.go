package dataframe

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

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

// Print returns a easy to read tabular format of the data frame
func (df *DataFrame) Print(
	showHeaders bool,
	showTypes bool,
	showIndexes bool,
	class string,
) string {

	class = strings.TrimSpace(class)

	if df == nil {
		if class == "" {
			return "Nil"
		}
		return fmt.Sprintf("%v: Nil", class)
	}

	nRows := df.NumRows()
	nCols := df.NumColumns()
	if nCols == 0 {
		if class == "" {
			return "Empty"
		}
		return fmt.Sprintf("%v: Empty", class)
	}

	addRightPadding := func(s string, nChar int) string {
		if utf8.RuneCountInString(s) < nChar {
			return s + strings.Repeat(" ", nChar-utf8.RuneCountInString(s))
		}
		return s
	}

	addLeftPadding := func(s string, nChar int) string {
		if utf8.RuneCountInString(s) < nChar {
			return strings.Repeat(" ", nChar-utf8.RuneCountInString(s)) + s
		}
		return s
	}

	nMatRows := nRows
	nMatCols := nCols
	colOffset := 0
	rowOffset := 0
	if showHeaders {
		nMatRows++
		rowOffset = 1
	}
	if showTypes {
		nMatRows++
	}
	if showIndexes {
		nMatCols++
		colOffset = 1
	}

	matrix := make([][]string, nMatRows)
	maxWidths := make([]int, nMatCols)

	addToMatrix := func(r int, c int, s string) {
		length := len(s)

		if length > maxWidths[c] {
			maxWidths[c] = length
		}

		if len(matrix[r]) == 0 {
			matrix[r] = make([]string, nMatCols)
		}
		matrix[r][c] = s
	}

	// add headers to matrix
	if showHeaders {
		addToMatrix(0, 0, "")
		for c := 0; c < nCols; c++ {
			colName := df.columns[c].GetName()
			addToMatrix(0, c+colOffset, colName)
		}
	}

	// add data to matrix
	for r := 0; r < nRows; r++ {

		if showIndexes {
			addToMatrix(r+rowOffset, 0, strconv.Itoa(r)+":")
		}

		for c := 0; c < nCols; c++ {
			val, err := df.Get(c, r)
			if err != nil {
				panic(err)
			}
			addToMatrix(r+rowOffset, c+colOffset, val.ToString())
		}
	}

	// add types to matrix
	if showTypes {
		addToMatrix(nRows+rowOffset, 0, "")
		for c := 0; c < nCols; c++ {
			col := df.columns[c]
			typeText := fmt.Sprintf("<%v>", col.GetType())
			addToMatrix(nRows+rowOffset, c+colOffset, typeText)
		}
	}

	// build the text
	str := fmt.Sprintf("[%dx%d] %v", nCols, nRows, class)
	str = strings.TrimRight(str, " ")
	for r := 0; r < nMatRows; r++ {
		str += "\n"

		// add padding
		if showIndexes {
			matrix[r][0] = addLeftPadding(matrix[r][0], maxWidths[0])
		}
		for c := colOffset; c < nMatCols; c++ {
			matrix[r][c] = addRightPadding(matrix[r][c], maxWidths[c])
		}

		row := strings.Join(matrix[r], " ")
		str += strings.TrimRight(row, " ")
	}

	return str
}
