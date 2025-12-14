package set

import (
	"sort"

	"golang.org/x/exp/constraints"
)

// ToSliceOrdered converts a set of ordered values to a sorted slice.
// It requires that the element type be ordered (supports < operator).
func SetToSortedSlice[T constraints.Ordered](s Set[T]) []T {
	slice := s.ToSlice()
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	return slice
}
