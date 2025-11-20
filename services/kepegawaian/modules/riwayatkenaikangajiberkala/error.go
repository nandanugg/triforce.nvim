package riwayatkenaikangajiberkala

import "errors"

var (
	errPegawaiNotFound   = errors.New("data pegawai tidak ditemukan")
	errGolonganNotFound  = errors.New("data golongan tidak ditemukan")
	errUnitKerjaNotFound = errors.New("data unit kerja tidak ditemukan")
)
