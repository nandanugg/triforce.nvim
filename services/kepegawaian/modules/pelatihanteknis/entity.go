package pelatihanteknis

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type pelatihanTeknis struct {
	ID              int64   `json:"id"`
	TipeDiklat      string  `json:"tipe_diklat"`
	JenisDiklat     string  `json:"jenis_diklat"`
	NamaDiklat      string  `json:"nama_diklat"`
	JumlahJam       int     `json:"jumlah_jam"`
	TanggalDiklat   db.Date `json:"tanggal_diklat"`
	NomorSertifikat string  `json:"nomor_sertifikat"`
}
