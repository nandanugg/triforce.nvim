package auth

import (
	"errors"
	"fmt"
)

type httpStatusError struct {
	code    int
	message []byte
}

func (e *httpStatusError) Error() string {
	return fmt.Sprintf("status=%d: %s", e.code, e.message)
}

var errUserNotFound = errors.New("user tidak ditemukan")
