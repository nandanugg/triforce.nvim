package keluarga

import (
	"errors"
)

var (
	errPegawaiNotFound = errors.New("pegawai not found")

	errAgamaNotFound            = errors.New("data agama tidak ditemukan")
	errStatusPernikahanNotFound = errors.New("data status pernikahan tidak ditemukan")
	errPasanganOrangTuaNotFound = errors.New("data pasangan orang tua tidak ditemukan")
)
