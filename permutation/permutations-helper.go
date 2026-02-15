package permutation

import "fmt"

// CalcNumPermutations calculates the total number of permutations given the options for each axis.
func CalcNumPermutations[T any](
	axisOptions [][]T,
) int {
	if len(axisOptions) == 0 {
		return 0
	}

	numPermutations := 1
	for _, options := range axisOptions {
		numPermutations *= len(options)
	}
	return numPermutations
}

// SelectPermutation selects one permutation out of the permutation space defined by axisOptions.
// e.g. if axisOptions is [[1, 2, 3], [4, 5]], then the permutation space is:
// - [1, 4]
// - [1, 5]
// - [2, 4]
// - [2, 5]
// - [3, 4]
// - [3, 5]
// and if idx is 2, then the selected permutation is [2, 4].
func SelectPermutation[T any](
	axisOptions [][]T,
	idx int,
) ([]T, error) {

	maxPermutations := CalcNumPermutations(axisOptions)

	if idx < 0 || idx >= maxPermutations {
		return nil, fmt.Errorf("idx %d is out of bounds for permutation space of size %d", idx, maxPermutations)
	}

	items := make([]T, len(axisOptions))
	for i, options := range axisOptions {

		items[i] = options[idx%len(options)]
		idx /= len(options)
	}
	return items, nil
}
