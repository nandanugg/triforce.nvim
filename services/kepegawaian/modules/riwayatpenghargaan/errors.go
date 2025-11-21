package riwayatpenghargaan

import (
	"errors"
	"fmt"
)

var ErrJenisPenghargaanInvalid = errors.New("jenis penghargaan tidak valid")

func NewError(base error, message string, args ...any) error {
	if len(args) > 0 {
		message = fmt.Sprintf(message, args...)
	}
	return fmt.Errorf("%w: %s", base, message)
}
