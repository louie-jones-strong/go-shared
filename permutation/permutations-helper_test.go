package permutation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_CalcNumPermutations(t *testing.T) {

	tests := []struct {
		name        string
		axisOptions [][]int
		expectedRes int
	}{
		{
			name:        "nil axis options",
			axisOptions: nil,
			expectedRes: 0,
		},
		{
			name:        "empty axis options",
			axisOptions: [][]int{},
			expectedRes: 0,
		},
		{
			name: "nil axis",
			axisOptions: [][]int{
				nil,
			},
			expectedRes: 0,
		},
		{
			name: "empty axis",
			axisOptions: [][]int{
				{},
			},
			expectedRes: 0,
		},
		{
			name: "single axis with 3 options",
			axisOptions: [][]int{
				{1, 2, 3},
			},
			expectedRes: 3,
		},
		{
			name: "axis with 3 options and axis with 0 options",
			axisOptions: [][]int{
				{1, 2, 3},
				{},
			},
			expectedRes: 0,
		},
		{
			name: "axis with 3 options and axis with 2 options",
			axisOptions: [][]int{
				{1, 2, 3},
				{4, 5},
			},
			expectedRes: 6,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res := CalcNumPermutations(tc.axisOptions)

			assert.Equal(t, tc.expectedRes, res)
		})
	}
}

func TestUnit_SelectPermutation(t *testing.T) {

	tests := []struct {
		name        string
		axisOptions [][]int
		idx         int
		expectedRes []int
		expectedErr error
	}{
		{
			name:        "nil axis options",
			axisOptions: nil,
			idx:         0,
			expectedRes: nil,
			expectedErr: fmt.Errorf("idx 0 is out of bounds for permutation space of size 0"),
		},
		{
			name:        "empty axis options",
			axisOptions: [][]int{},
			idx:         0,
			expectedRes: nil,
			expectedErr: fmt.Errorf("idx 0 is out of bounds for permutation space of size 0"),
		},
		{
			name: "nil axis",
			axisOptions: [][]int{
				nil,
			},
			idx:         0,
			expectedRes: nil,
			expectedErr: fmt.Errorf("idx 0 is out of bounds for permutation space of size 0"),
		},
		{
			name: "empty axis",
			axisOptions: [][]int{
				{},
			},
			idx:         0,
			expectedRes: nil,
			expectedErr: fmt.Errorf("idx 0 is out of bounds for permutation space of size 0"),
		},
		{
			name: "single axis with 3 options and idx 0",
			axisOptions: [][]int{
				{1, 2, 3},
			},
			idx:         0,
			expectedRes: []int{1},
		},
		{
			name: "single axis with 3 options and idx 2",
			axisOptions: [][]int{
				{1, 2, 3},
			},
			idx:         2,
			expectedRes: []int{3},
		},
		{
			name: "single axis with 3 options and idx 3 (out of bounds)",
			axisOptions: [][]int{
				{1, 2, 3},
			},
			idx:         3,
			expectedRes: nil,
			expectedErr: fmt.Errorf("idx 3 is out of bounds for permutation space of size 3"),
		},
		{
			name: "axis with 3 options and axis with 0 options",
			axisOptions: [][]int{
				{1, 2, 3},
				{},
			},
			idx:         0,
			expectedRes: nil,
			expectedErr: fmt.Errorf("idx 0 is out of bounds for permutation space of size 0"),
		},
		{
			name: "axis with 3 options and axis with 2 options",
			axisOptions: [][]int{
				{1, 2, 3},
				{4, 5},
			},
			idx:         0,
			expectedRes: []int{1, 4},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			res, err := SelectPermutation(tc.axisOptions, tc.idx)

			if tc.expectedErr != nil {
				assert.EqualError(t, err, tc.expectedErr.Error())
				assert.Zero(t, res)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedRes, res)
			}
		})
	}
}
