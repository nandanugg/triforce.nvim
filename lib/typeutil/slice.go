package typeutil

// Map manipulates a slice and transforms it to a slice of another type.
func Map[T any, R any](collection []T, iteratee func(item T) R) []R {
	result := make([]R, 0, len(collection))
	for i := range collection {
		result = append(result, iteratee(collection[i]))
	}
	return result
}
