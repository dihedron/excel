package collections

import (
	"cmp"
	"iter"
	"slices"
)

// MapKeys extracts the keys from a map.
func MapKeys[K comparable, V any](m map[K]V) []K {
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	return keys
}

// SortedMapIter returns an iterator yielding map entries in key-sorted order.
// The key type 'K' must be ordered (e.g., strings, integers) so it can be sorted.
func SortedMap[K cmp.Ordered, V any](m map[K]V) iter.Seq2[K, V] {
	// 1. Extract the keys from the map
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// 2. Sort the keys
	slices.Sort(keys)

	// 3. Return the iterator function
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			// Yield passes the key and value to the range loop block.
			// If yield returns false (e.g., the loop breaks), we stop iterating.
			if !yield(k, m[k]) {
				return
			}
		}
	}
}

// SortedMapIterFunc returns an iterator yielding map entries ordered by a custom comparison function.
// The key type 'K' only needs to be comparable (the standard map key constraint).
func SortedMapFunc[K comparable, V any](m map[K]V, cmpFunc func(a, b K) int) iter.Seq2[K, V] {
	// 1. Extract the keys from the map
	keys := make([]K, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}

	// 2. Sort the keys using the provided custom comparison function
	slices.SortFunc(keys, cmpFunc)

	// 3. Return the iterator function
	return func(yield func(K, V) bool) {
		for _, k := range keys {
			if !yield(k, m[k]) {
				return
			}
		}
	}
}
