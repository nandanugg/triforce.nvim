package riwayatpendidikan

import "errors"

var (
	errPegawaiNotFound = errors.New("pegawai not found")

	errTingkatPendidikanNotFound = errors.New("data tingkat pendidikan tidak ditemukan")
	errPendidikanNotFound        = errors.New("data pendidikan tidak ditemukan")
)
