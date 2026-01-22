package dataframe

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/collections/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func TestUnit_print(t *testing.T) {

	tests := []struct {
		name          string
		df            *DataFrame
		showHeader    bool
		showTypes     bool
		showIndexes   bool
		numHeaderRows int
		numTailRows   int
		class         string
		expectedRes   string
	}{
		{
			name:        "nil",
			df:          nil,
			expectedRes: `Nil`,
		},
		{
			name:        "nil with class",
			df:          nil,
			class:       "DataFrame",
			expectedRes: `DataFrame: Nil`,
		},
		{
			name: "no columns without class",
			df: New(
				[]*series.Series{},
			),
			expectedRes: `Empty`,
		},
		{
			name: "no columns with class",
			df: New(
				[]*series.Series{},
			),
			class:       "DataFrame",
			expectedRes: `DataFrame: Empty`,
		},
		{
			name: "empty row",
			df: New(
				[]*series.Series{
					series.BuildSeries("name", []string{}),
				},
			),
			showHeader:    true,
			showTypes:     true,
			showIndexes:   true,
			numHeaderRows: 100,
			numTailRows:   100,
			class:         "DataFrame",
			expectedRes: `[1x0] DataFrame
 name
 <string>`,
		},
		{
			name: "show all: all strings",
			df: New(
				[]*series.Series{
					series.BuildSeries("name", []string{"Alice", "Bob"}),
					series.BuildSeries("age", []string{"20", "22"}),
					series.BuildSeries("grade", []string{"A", "B"}),
				},
			),
			showHeader:    true,
			showTypes:     true,
			showIndexes:   true,
			numHeaderRows: 100,
			numTailRows:   100,
			class:         "DataFrame",
			expectedRes: `[3x2] DataFrame
   name     age      grade
0: Alice    20       A
1: Bob      22       B
   <string> <string> <string>`,
		},
		{
			name: "show only header: all strings",
			df: New(
				[]*series.Series{
					series.BuildSeries("name", []string{"Alice", "Bob"}),
					series.BuildSeries("age", []string{"20", "22"}),
					series.BuildSeries("grade", []string{"A", "B"}),
				},
			),
			showHeader:    true,
			showTypes:     false,
			showIndexes:   false,
			numHeaderRows: 100,
			numTailRows:   100,
			expectedRes: `[3x2]
name  age grade
Alice 20  A
Bob   22  B`,
		},
		{
			name: "show only index: all strings",
			df: New(
				[]*series.Series{
					series.BuildSeries("name", []string{"Alice", "Bob"}),
					series.BuildSeries("age", []string{"20", "22"}),
					series.BuildSeries("grade", []string{"A", "B"}),
				},
			),
			showHeader:    false,
			showTypes:     false,
			showIndexes:   true,
			numHeaderRows: 100,
			numTailRows:   100,
			expectedRes: `[3x2]
0: Alice 20 A
1: Bob   22 B`,
		},
		{
			name: "show only types: all strings",
			df: New(
				[]*series.Series{
					series.BuildSeries("name", []string{"Alice", "Bob"}),
					series.BuildSeries("age", []string{"20", "22"}),
					series.BuildSeries("grade", []string{"A", "B"}),
				},
			),
			showHeader:    false,
			showTypes:     true,
			showIndexes:   false,
			numHeaderRows: 100,
			numTailRows:   100,
			expectedRes: `[3x2]
Alice    20       A
Bob      22       B
<string> <string> <string>`,
		},
		{
			name: "show all: 11 rows all strings",
			df: New(
				[]*series.Series{
					series.BuildSeries("test", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}),
				},
			),
			showHeader:    true,
			showTypes:     true,
			showIndexes:   true,
			numHeaderRows: 100,
			numTailRows:   100,
			expectedRes: `[1x11]
    test
 0: 1
 1: 2
 2: 3
 3: 4
 4: 5
 5: 6
 6: 7
 7: 8
 8: 9
 9: 10
10: 11
    <string>`,
		},
		{
			name: "show all: 11 rows of all types",
			df: New(
				[]*series.Series{
					series.BuildSeries("A", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}),
					series.BuildSeries("B", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
					series.BuildSeries("C", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
					series.BuildSeries("D", []bool{true, false, true, false, true, false, true, false, true, false, true}),
				},
			),
			showHeader:    true,
			showTypes:     true,
			showIndexes:   true,
			numHeaderRows: 100,
			numTailRows:   100,
			expectedRes: `[4x11]
    A        B       C       D
 0: 1        1       1       true
 1: 2        2       2       false
 2: 3        3       3       true
 3: 4        4       4       false
 4: 5        5       5       true
 5: 6        6       6       false
 6: 7        7       7       true
 7: 8        8       8       false
 8: 9        9       9       true
 9: 10       10      10      false
10: 11       11      11      true
    <string> <float> <float> <bool>`,
		},
		{
			name: "show all: 11 rows of all types only 5 header 3 tails",
			df: New(
				[]*series.Series{
					series.BuildSeries("A", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}),
					series.BuildSeries("B", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
					series.BuildSeries("C", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
					series.BuildSeries("D", []bool{true, false, true, false, true, false, true, false, true, false, true}),
				},
			),
			showHeader:    true,
			showTypes:     true,
			showIndexes:   true,
			numHeaderRows: 5,
			numTailRows:   3,
			expectedRes: `[4x11]
    A        B       C       D
 0: 1        1       1       true
 1: 2        2       2       false
 2: 3        3       3       true
 3: 4        4       4       false
 4: 5        5       5       true
    ...      ...     ...     ...
 8: 9        9       9       true
 9: 10       10      10      false
10: 11       11      11      true
    <string> <float> <float> <bool>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.df.Print(
				tc.showHeader,
				tc.showTypes,
				tc.showIndexes,
				tc.numHeaderRows,
				tc.numTailRows,
				tc.class,
			)

			assert.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestUnit_Describe(t *testing.T) {

	tests := []struct {
		name        string
		df          *DataFrame
		expectedRes *DataFrame
	}{
		{
			name:        "nil",
			df:          nil,
			expectedRes: nil,
		},
		{
			name: "nil columns",
			df: New(
				nil,
			),
			expectedRes: nil,
		},
		{
			name: "no columns",
			df: New(
				[]*series.Series{},
			),
			expectedRes: nil,
		},
		{
			name: "1 int column",
			df: New(
				[]*series.Series{
					series.BuildSeries("test", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				},
			),
			expectedRes: New(
				[]*series.Series{
					series.BuildSeries("index", []string{
						"count",
						"sum",
						"mean",
						"std",
						"min",
						// "25%",
						// "50%",
						// "75%",
						"max",
					}),
					series.BuildSeries("test", []float64{
						11,
						66,
						6,
						3.3166247903554,
						1,
						11,
					}),
				},
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.df.Describe()

			assert.Equal(t, tc.expectedRes, res)
		})
	}
}
