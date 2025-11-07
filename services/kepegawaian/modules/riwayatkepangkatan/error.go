package riwayatkepangkatan

import "errors"

var (
	errPegawaiNotFound = errors.New("pegawai not found")

	errJenisKenaikanPangkatNotFound = errors.New("data jenis kenaikan pangkat tidak ditemukan")
	errGolonganNotFound             = errors.New("data golongan tidak ditemukan")
)
