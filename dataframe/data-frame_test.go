package dataframe

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func TestUnit_print(t *testing.T) {

	tests := []struct {
		name        string
		df          *DataFrame
		showHeader  bool
		showTypes   bool
		showIndexes bool
		class       string
		expectedRes string
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
					series.New("name", series.String, []string{}),
				},
			),
			showHeader:  true,
			showTypes:   true,
			showIndexes: true,
			class:       "DataFrame",
			expectedRes: `[1x0] DataFrame
 name
 <string>`,
		},
		{
			name: "show all: all strings",
			df: New(
				[]*series.Series{
					series.New("name", series.String, []string{"Alice", "Bob"}),
					series.New("age", series.String, []string{"20", "22"}),
					series.New("grade", series.String, []string{"A", "B"}),
				},
			),
			showHeader:  true,
			showTypes:   true,
			showIndexes: true,
			class:       "DataFrame",
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
					series.New("name", series.String, []string{"Alice", "Bob"}),
					series.New("age", series.String, []string{"20", "22"}),
					series.New("grade", series.String, []string{"A", "B"}),
				},
			),
			showHeader:  true,
			showTypes:   false,
			showIndexes: false,
			expectedRes: `[3x2]
name  age grade
Alice 20  A
Bob   22  B`,
		},
		{
			name: "show only index: all strings",
			df: New(
				[]*series.Series{
					series.New("name", series.String, []string{"Alice", "Bob"}),
					series.New("age", series.String, []string{"20", "22"}),
					series.New("grade", series.String, []string{"A", "B"}),
				},
			),
			showHeader:  false,
			showTypes:   false,
			showIndexes: true,
			expectedRes: `[3x2]
0: Alice 20 A
1: Bob   22 B`,
		},
		{
			name: "show only types: all strings",
			df: New(
				[]*series.Series{
					series.New("name", series.String, []string{"Alice", "Bob"}),
					series.New("age", series.String, []string{"20", "22"}),
					series.New("grade", series.String, []string{"A", "B"}),
				},
			),
			showHeader:  false,
			showTypes:   true,
			showIndexes: false,
			expectedRes: `[3x2]
Alice    20       A
Bob      22       B
<string> <string> <string>`,
		},
		{
			name: "show all: 11 rows all strings",
			df: New(
				[]*series.Series{
					series.New("test", series.String, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}),
				},
			),
			showHeader:  true,
			showTypes:   true,
			showIndexes: true,
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
			name: "show all: 11 rows all ints",
			df: New(
				[]*series.Series{
					series.New("test", series.Int, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				},
			),
			showHeader:  true,
			showTypes:   true,
			showIndexes: true,
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
    <int>`,
		},
		{
			name: "show all: 11 rows of all types",
			df: New(
				[]*series.Series{
					series.New("A", series.String, []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}),
					series.New("B", series.Int, []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
					series.New("C", series.Float, []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
					series.New("D", series.Bool, []bool{true, false, true, false, true, false, true, false, true, false, true}),
				},
			),
			showHeader:  true,
			showTypes:   true,
			showIndexes: true,
			expectedRes: `[4x11]
    A        B     C       D
 0: 1        1     1       true
 1: 2        2     2       false
 2: 3        3     3       true
 3: 4        4     4       false
 4: 5        5     5       true
 5: 6        6     6       false
 6: 7        7     7       true
 7: 8        8     8       false
 8: 9        9     9       true
 9: 10       10    10      false
10: 11       11    11      true
    <string> <int> <float> <bool>`,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.df.Print(tc.showHeader, tc.showTypes, tc.showIndexes, tc.class)

			assert.Equal(t, tc.expectedRes, res)
		})
	}
}
