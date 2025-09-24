package typeutil

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(item T) R) []R {
	result := make([]R, 0, len(collection))
	for i := range collection {
		result = append(result, iteratee(collection[i]))
	}
	return result
}

// UniqMap manipulates a slice and transforms it to a slice of another type with unique values.
func UniqMap[T any, R comparable](collection []T, iteratee func(item T, index int) R) []R {
	result := make([]R, 0, len(collection))
	seen := make(map[R]struct{}, len(collection))

	for i := range collection {
		r := iteratee(collection[i], i)
		if _, ok := seen[r]; !ok {
			result = append(result, r)
			seen[r] = struct{}{}
		}
	}
	return result
}

// SliceToMap returns a map containing key-value pairs provided by transform function applied to elements of the given slice.
// If any of two pairs would have the same key the last one gets added to the map.
// The order of keys in returned map is not specified and is not guaranteed to be the same from the original array.
func SliceToMap[T any, K comparable, V any](collection []T, transform func(item T) (K, V)) map[K]V {
	result := make(map[K]V, len(collection))

	for i := range collection {
		k, v := transform(collection[i])
		result[k] = v
	}

	return result
}
