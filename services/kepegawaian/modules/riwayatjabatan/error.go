package riwayatjabatan

import "errors"

var (
	errPegawaiNotFound = errors.New("pegawai not found")

	errJenisJabatanNotFound   = errors.New("data jenis jabatan tidak ditemukan")
	errJabatanNotFound        = errors.New("data jabatan tidak ditemukan")
	errSatuanKerjaNotFound    = errors.New("data satuan kerja tidak ditemukan")
	errUnitOrganisasiNotFound = errors.New("data unit organisasi tidak ditemukan")
)
