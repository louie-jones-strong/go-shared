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

// ConvertMapToKVPList converts a map to a slice of Key-Value Pair (KVP) structs.
func ConvertMapToKVPList[K comparable, V any](m map[K]V) []KVP[K, V] {
	kvpList := make([]KVP[K, V], 0, len(m))
	for k, v := range m {
		kvp := KVP[K, V]{key: k, value: v}
		kvpList = append(kvpList, kvp)
	}
	return kvpList
}

// GetMapKeys returns a slice of keys from the given map.
func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetMapValues returns a slice of values from the given map.
func GetMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
