package dataframe

import (
	"testing"
	"time"

	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/louie-jones-strong/go-shared/dataframe/series/elements/element"
	"github.com/stretchr/testify/assert"
)

func Benchmark(b *testing.B) {

	for i := 0; i < b.N; i++ {

		colA := series.BuildSeries("A", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"})
		var dTime time.Time
		df := New(
			[]*series.Series{
				colA,
				series.BuildSeries("B", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				series.BuildSeries("C", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				series.BuildSeries("D", []bool{true, false, true, false, true, false, true, false, true, false, true}),
				series.BuildSeries("D_2", []bool{true, false, true, false, true, false, true, false, true, false, true}),
				series.BuildSeries("E", []time.Time{dTime, dTime, dTime, dTime, dTime, dTime, dTime, dTime, dTime, dTime, dTime}),
			},
		)
		assert.NotNil(b, df)

		col1, err := df.GetColumn(0)
		assert.NoError(b, err)
		assert.Equal(b, colA, col1)

		col1Clone := col1.Clone()
		assert.Equal(b, colA, col1Clone)

		col1Clone.Rename("A_Clone")

		err = df.AddColumn(col1Clone)
		assert.NoError(b, err)

		err = df.DropColumn("D_2")
		assert.NoError(b, err)

		describeDf := df.Describe()
		assert.NotNil(b, describeDf)

		str := describeDf.String()

		expectedStr := `[7x6] DataFrame
   index    A               B               C               D                  E                A_Clone
0: count    11              11              11              11                 11               11
1: sum      NaN             66              66              NaN                -6.834915648e+11 NaN
2: mean     6               6               6               0.5454545454545454 -6.21355968e+10  6
3: std      3.3166247903554 3.3166247903554 3.3166247903554 0.5222329678670935 0                3.3166247903554
4: min      NaN             1               1               NaN                -6.21355968e+10  NaN
5: max      NaN             11              11              NaN                -6.21355968e+10  NaN
   <string> <float>         <float>         <float>         <float>            <float>          <float>`
		assert.Equal(b, expectedStr, str)

		delegate := func(item element.IElement) (element.IElement, error) {
			newElem := element.NewFloatElement(item.ToFloat() * 10)
			return newElem, nil
		}
		col1Clone.ApplyInPlace(delegate)

		describeDf = df.Describe()
		assert.NotNil(b, describeDf)
		str = describeDf.String()

		expectedStr = `[7x6] DataFrame
   index    A               B               C               D                  E                A_Clone
0: count    11              11              11              11                 11               11
1: sum      NaN             66              66              NaN                -6.834915648e+11 NaN
2: mean     6               6               6               0.5454545454545454 -6.21355968e+10  60
3: std      3.3166247903554 3.3166247903554 3.3166247903554 0.5222329678670935 0                33.166247903554
4: min      NaN             1               1               NaN                -6.21355968e+10  NaN
5: max      NaN             11              11              NaN                -6.21355968e+10  NaN
   <string> <float>         <float>         <float>         <float>            <float>          <float>`
		assert.Equal(b, expectedStr, str)

	}
}
