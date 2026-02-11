package maps

// KVP represents a Key-Value Pair.
type KVP[K comparable, V any] struct {
	key   K
	value V
}

// Key returns the key of the Key-Value Pair.
func (kvp KVP[K, V]) Key() K {
	return kvp.key
}

// Value returns the value of the Key-Value Pair.
func (kvp KVP[K, V]) Value() V {
	return kvp.value
}
