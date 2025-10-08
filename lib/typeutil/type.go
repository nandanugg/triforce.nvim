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

// FromPtr returns the pointer value or empty.
func FromPtr[T any](x *T) (t T) {
	if x == nil {
		return t
	}
	return *x
}

// ValueOrNil returns a pointer to v if valid is true, otherwise nil.
func ValueOrNil[T any](v T, valid bool) *T {
	if valid {
		return &v
	}
	return nil
}

// Coalesce returns the first non-empty arguments. Arguments must be comparable.
func Coalesce[T comparable](values ...T) (t T) {
	for i := range values {
		if values[i] != t {
			return values[i]
		}
	}
	return t
}
