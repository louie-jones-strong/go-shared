package dataframe

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func Benchmark(b *testing.B) {

	for i := 0; i < b.N; i++ {

		colA := series.BuildSeries("A", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"})
		df := New(
			[]*series.Series{
				colA,
				series.BuildSeries("B", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				series.BuildSeries("C", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				series.BuildSeries("D", []bool{true, false, true, false, true, false, true, false, true, false, true}),
				series.BuildSeries("D_2", []bool{true, false, true, false, true, false, true, false, true, false, true}),
			},
		)
		assert.NotNil(b, df)

		col1, err := df.GetColumn(0)
		assert.NoError(b, err)
		assert.Equal(b, colA, col1)

		col1Clone := col1.Clone()
		assert.Equal(b, colA, col1Clone)

		err = df.DropColumn("D_2")
		assert.NoError(b, err)

		describeDf := df.Describe()
		assert.NotNil(b, describeDf)

		str := describeDf.String()

		expectedStr := `[5x6] DataFrame
   index    A               B               C               D
0: count    11              11              11              11
1: sum      NaN             66              66              NaN
2: mean     6               6               6               0.5454545454545454
3: std      3.3166247903554 3.3166247903554 3.3166247903554 0.5222329678670935
4: min      NaN             1               1               NaN
5: max      NaN             11              11              NaN
   <string> <float>         <float>         <float>         <float>`
		assert.Equal(b, expectedStr, str)
	}
}
