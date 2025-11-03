package api

import "strings"

type MultiError struct {
	errs []error
}

func NewMultiError(errs []error) *MultiError {
	return &MultiError{errs}
}

func (e *MultiError) Error() string {
	var msg strings.Builder
	for _, e := range e.errs {
		if msg.Len() != 0 {
			msg.WriteString(" | ")
		}
		msg.WriteString(e.Error())
	}
	return msg.String()
}
