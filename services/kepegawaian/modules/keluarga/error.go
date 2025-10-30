package keluarga

import "errors"

var (
	errPegawaiNotFound          = errors.New("pegawai not found")
	errAgamaNotFound            = errors.New("agama not found")
	errStatusPernikahanNotFound = errors.New("status pernikahan not found")
	errPasanganOrangTuaNotFound = errors.New("pasangan orang tua not found")
)
