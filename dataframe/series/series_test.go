package series

import (
	"testing"

	"github.com/louie-jones-strong/go-shared/dataframe/apptype"
	"github.com/stretchr/testify/assert"
)

func TestUnit_Order(t *testing.T) {

	tests := []struct {
		name        string
		s           *Series
		reverse     bool
		expectedRes []int
	}{
		{
			name:        "empty",
			s:           New("test", apptype.Int, []int{}),
			expectedRes: nil,
		},
		{
			name:        "ordered: acceding",
			s:           New("test", apptype.Int, []int{1, 2, 3}),
			expectedRes: []int{0, 1, 2},
		},
		{
			name:        "reverse ordered: acceding",
			s:           New("test", apptype.Int, []int{3, 2, 1}),
			expectedRes: []int{2, 1, 0},
		},
		{
			name:        "ordered: descending",
			s:           New("test", apptype.Int, []int{1, 2, 3}),
			reverse:     true,
			expectedRes: []int{2, 1, 0},
		},
		{
			name:        "reverse ordered: descending",
			s:           New("test", apptype.Int, []int{3, 2, 1}),
			reverse:     true,
			expectedRes: []int{0, 1, 2},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := tc.s.Order(tc.reverse)

			assert.Equal(t, tc.expectedRes, res)
		})
	}
}
