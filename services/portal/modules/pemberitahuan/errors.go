package pemberitahuan

import (
	"errors"
	"fmt"
)

// ErrConflict is the sentinel error used for type comparison.
var ErrConflict = errors.New("pemberitahuan conflict")

func NewError(base error, message string, args ...interface{}) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return fmt.Errorf("%w: %s", base, message)
}
