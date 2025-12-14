package set

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUnit_Add(t *testing.T) {

	tests := []struct {
		name        string
		s           Set[int]
		itemToAdd   int
		expectedRes bool
	}{
		{
			name:        "Add new item",
			s:           New[int](),
			itemToAdd:   0,
			expectedRes: true,
		},
		{
			name:        "Add duplicate item",
			s:           NewWith(0),
			itemToAdd:   0,
			expectedRes: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			added := tc.s.Add(tc.itemToAdd)

			assert.Equal(t, tc.expectedRes, added)
		})
	}
}

func TestUnit_AddItems(t *testing.T) {
	tests := []struct {
		name        string
		s           Set[int]
		itemsToAdd  []int
		expectedRes int
	}{
		{
			name:        "Add two new items",
			s:           New[int](),
			itemsToAdd:  []int{2, 3},
			expectedRes: 2,
		},
		{
			name:        "Add items with duplicate",
			s:           NewWith(2),
			itemsToAdd:  []int{2, 3},
			expectedRes: 1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			added := tc.s.AddItems(tc.itemsToAdd...)
			assert.Equal(t, tc.expectedRes, added)
		})
	}
}

func TestUnit_Contains(t *testing.T) {
	tests := []struct {
		name        string
		s           Set[int]
		itemToCheck int
		expectedRes bool
	}{
		{
			name:        "Contains present item",
			s:           NewWith(5),
			itemToCheck: 5,
			expectedRes: true,
		},
		{
			name:        "Contains absent item",
			s:           NewWith(6),
			itemToCheck: 7,
			expectedRes: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			present := tc.s.Contains(tc.itemToCheck)
			assert.Equal(t, tc.expectedRes, present)
		})
	}
}

func TestUnit_Remove(t *testing.T) {
	tests := []struct {
		name         string
		s            Set[int]
		itemToRemove int
		expectedLen  int
	}{
		{
			name:         "Remove existing item",
			s:            NewWith(8, 9),
			itemToRemove: 9,
			expectedLen:  1,
		},
		{
			name:         "Remove non-existing item",
			s:            NewWith(8, 9),
			itemToRemove: 10,
			expectedLen:  2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			tc.s.Remove(tc.itemToRemove)
			assert.Equal(t, tc.expectedLen, tc.s.Len())
		})
	}
}

func TestUnit_Len(t *testing.T) {
	tests := []struct {
		name        string
		s           Set[int]
		expectedLen int
	}{
		{
			name:        "Length with two items",
			s:           NewWith(10, 11),
			expectedLen: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedLen, tc.s.Len())
		})
	}
}

func TestUnit_ToSlice(t *testing.T) {
	tests := []struct {
		name        string
		s           Set[int]
		expectedLen int
	}{
		{
			name:        "ToSlice length",
			s:           NewWith(4, 5),
			expectedLen: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			slice := tc.s.ToSlice()
			assert.Equal(t, tc.expectedLen, len(slice))
		})
	}
}

func TestUnit_Union(t *testing.T) {
	tests := []struct {
		name        string
		s1          Set[int]
		s2          Set[int]
		expectedOrd []int
	}{
		{
			name:        "Union disjoint",
			s1:          NewWith(1, 2),
			s2:          NewWith(3),
			expectedOrd: []int{1, 2, 3},
		},
		{
			name:        "Union overlapping",
			s1:          NewWith(1, 2),
			s2:          NewWith(2, 3),
			expectedOrd: []int{1, 2, 3},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.s1.Union(tc.s2)
			slice := SetToSortedSlice(res)
			assert.Equal(t, len(tc.expectedOrd), len(slice))
			for i := range slice {
				assert.Equal(t, tc.expectedOrd[i], slice[i])
			}
		})
	}
}

func TestUnit_Intersection(t *testing.T) {
	tests := []struct {
		name        string
		s1          Set[int]
		s2          Set[int]
		expectedOrd []int
	}{
		{
			name:        "Intersection non-empty",
			s1:          NewWith(1, 2),
			s2:          NewWith(2, 3),
			expectedOrd: []int{2},
		},
		{
			name:        "Intersection empty",
			s1:          NewWith(1, 4),
			s2:          NewWith(2, 3),
			expectedOrd: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.s1.Intersection(tc.s2)
			slice := SetToSortedSlice(res)
			assert.Equal(t, len(tc.expectedOrd), len(slice))
			for i := range slice {
				assert.Equal(t, tc.expectedOrd[i], slice[i])
			}
		})
	}
}

func TestUnit_Difference(t *testing.T) {
	tests := []struct {
		name        string
		s1          Set[int]
		s2          Set[int]
		expectedOrd []int
	}{
		{
			name:        "Difference removes common",
			s1:          NewWith(1, 2, 3),
			s2:          NewWith(2),
			expectedOrd: []int{1, 3},
		},
		{
			name:        "Difference when none common",
			s1:          NewWith(1, 2),
			s2:          NewWith(3, 4),
			expectedOrd: []int{1, 2},
		},
		{
			name:        "Difference all removed",
			s1:          NewWith(1, 2),
			s2:          NewWith(1, 2),
			expectedOrd: []int{},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			res := tc.s1.Difference(tc.s2)
			slice := SetToSortedSlice(res)
			assert.Equal(t, len(tc.expectedOrd), len(slice))
			for i := range slice {
				assert.Equal(t, tc.expectedOrd[i], slice[i])
			}
		})
	}
}
