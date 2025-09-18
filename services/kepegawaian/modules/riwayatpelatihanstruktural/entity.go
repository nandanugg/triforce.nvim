package riwayatpelatihanstruktural

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type riwayatPelatihanStruktural struct {
	ID                     string  `json:"id"`
	JenisDiklat            string  `json:"jenis_diklat"`
	NamaDiklat             string  `json:"nama_diklat"`
	InstitusiPenyelenggara string  `json:"institusi_penyelenggara"`
	NomorSertifikat        string  `json:"nomor_sertifikat"`
	TanggalMulai           db.Date `json:"tanggal_mulai"`
	TanggalSelesai         db.Date `json:"tanggal_selesai"`
	Tahun                  *int16  `json:"tahun"`
	Durasi                 int     `json:"durasi"` // hour
}
