package dataframe

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/dataframe/series"
	"github.com/stretchr/testify/assert"
)

func Benchmark(b *testing.B) {

	for i := 0; i < b.N; i++ {
		df := New(
			[]*series.Series{
				series.New("A", []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "10", "11"}),
				series.New("B", []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				series.New("C", []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}),
				series.New("D", []bool{true, false, true, false, true, false, true, false, true, false, true}),
				series.New("D_2", []bool{true, false, true, false, true, false, true, false, true, false, true}),
			},
		)
		assert.NotNil(b, df)

		err := df.DropColumn("D_2")
		assert.NoError(b, err)

		describeDf := df.Describe()
		assert.NotNil(b, describeDf)

		describeDf.Print(true, true, true, "Describe")
	}
}
