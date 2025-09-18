package typeutil

// Cast safe convert interface into T
func Cast[T any](val any) T {
	v, _ := val.(T)
	return v
}

// ToPtr returns a pointer copy of value.
func ToPtr[T any](x T) *T {
	return &x
}
