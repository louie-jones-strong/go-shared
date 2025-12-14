package set

// Set is a generic set data structure for comparable types.
type Set[T comparable] map[T]struct{}

// New creates a new empty set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// NewWith creates a new set initialized with the provided items.
func NewWith[T comparable](items ...T) Set[T] {
	s := New[T]()
	s.AddItems(items...)
	return s
}

// Contains checks if the set contains the specified item.
func (s Set[T]) Contains(item T) bool {
	_, exists := s[item]
	return exists
}

// AddItems adds multiple items to the set and returns the count of newly added items.
func (s Set[T]) AddItems(items ...T) int {
	count := 0
	for _, item := range items {
		if s.Add(item) {
			count++
		}
	}
	return count
}

// Add adds an item to the set and returns true if the item was not already present.
func (s Set[T]) Add(item T) bool {
	if s.Contains(item) {
		return false
	}
	s[item] = struct{}{}
	return true
}

// Remove removes an item from the set.
func (s Set[T]) Remove(item T) {
	delete(s, item)
}

// Len returns the number of items in the set.
func (s Set[T]) Len() int {
	return len(s)
}

// ToSlice converts the set to a slice of items.
// Note the order is not deterministic, if you it to be deterministic use SetToSortedSlice.
func (s Set[T]) ToSlice() []T {
	slice := make([]T, 0, len(s))
	for item := range s {
		slice = append(slice, item)
	}
	return slice
}

// Union returns a new set that is the union of the current set and another set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s {
		result.Add(item)
	}
	for item := range other {
		result.Add(item)
	}
	return result
}

// Intersection returns a new set that is the intersection of the current set and another set.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s {
		if other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}

// Difference returns a new set that is the difference of the current set and another set.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	result := New[T]()
	for item := range s {
		if !other.Contains(item) {
			result.Add(item)
		}
	}
	return result
}
