package maps

// ConvertMapToKVPList converts a map to a slice of Key-Value Pair (KVP) structs.
// Note: The order of values in the returned slice is not deterministic and may vary between runs.
func ConvertMapToKVPList[K comparable, V any](m map[K]V) []KVP[K, V] {
	kvpList := make([]KVP[K, V], 0, len(m))
	for k, v := range m {
		kvp := KVP[K, V]{key: k, value: v}
		kvpList = append(kvpList, kvp)
	}
	return kvpList
}

// GetMapKeys returns a slice of keys from the given map.
// Note: The order of values in the returned slice is not deterministic and may vary between runs.
func GetMapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// GetMapValues returns a slice of values from the given map.
// Note: The order of values in the returned slice is not deterministic and may vary between runs.
func GetMapValues[K comparable, V any](m map[K]V) []V {
	values := make([]V, 0, len(m))
	for _, v := range m {
		values = append(values, v)
	}
	return values
}
