package typeutil

import "github.com/jackc/pgx/v5/pgtype"

// Cast safe convert interface into T
func Cast[T any](val any) T {
	v, _ := val.(T)
	return v
}

// ToPtr returns a pointer copy of value.
func ToPtr[T any](x T) *T {
	return &x
}

// PointerToPgtype convert pointer type to pgtype.
func PointerToPgtype[T any](ptr *T) any {
	if ptr == nil {
		switch any(*new(T)).(type) {
		case int16:
			return pgtype.Int2{Valid: false}
		case int32:
			return pgtype.Int4{Valid: false}
		case int64:
			return pgtype.Int8{Valid: false}
		case int:
			return pgtype.Int8{Valid: false}
		case string:
			return pgtype.Text{Valid: false}
		case bool:
			return pgtype.Bool{Valid: false}
		case float32:
			return pgtype.Float4{Valid: false}
		case float64:
			return pgtype.Float8{Valid: false}
		default:
			return nil
		}
	}

	switch v := any(*ptr).(type) {
	case int16:
		return pgtype.Int2{Int16: v, Valid: true}
	case int32:
		return pgtype.Int4{Int32: v, Valid: true}
	case int64:
		return pgtype.Int8{Int64: v, Valid: true}
	case int:
		return pgtype.Int8{Int64: int64(v), Valid: true}
	case string:
		return pgtype.Text{String: v, Valid: true}
	case bool:
		return pgtype.Bool{Bool: v, Valid: true}
	case float32:
		return pgtype.Float4{Float32: v, Valid: true}
	case float64:
		return pgtype.Float8{Float64: v, Valid: true}
	default:
		return nil
	}
}

// FromPtr returns the pointer value or empty.
func FromPtr[T any](x *T) (t T) {
	if x == nil {
		return t
	}
	return *x
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
