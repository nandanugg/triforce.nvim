package typeutil

// Cast safe convert interface into T
func Cast[T any](val any) T {
	v, _ := val.(T)
	return v
}
