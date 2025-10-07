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

// ValueOrNil returns a pointer to v if valid is true, otherwise nil.
func ValueOrNil[T any](v T, valid bool) *T {
	if valid {
		return &v
	}
	return nil
}
