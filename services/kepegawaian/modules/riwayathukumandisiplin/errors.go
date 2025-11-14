package riwayathukumandisiplin

import "errors"

var (
	errGolonganNotFound      = errors.New("data golongan tidak ditemukan")
	errJenisHukumanNotFound  = errors.New("data jenis hukuman tidak ditemukan")
	errPegawaiNotFound       = errors.New("data pegawai tidak ditemukan")
	errMasaHukumanTidakValid = errors.New("masa hukuman tidak valid")
)
