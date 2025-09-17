package pelatihanfungsional

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type pelatihanFungsional struct {
	ID                     string  `json:"id"`
	JenisDiklat            string  `json:"jenis_diklat"`
	NamaDiklat             string  `json:"nama_diklat"`
	TanggalSelesai         db.Date `json:"tanggal_selesai"`
	TanggalMulai           db.Date `json:"tanggal_mulai"`
	Durasi                 int32   `json:"durasi"`
	InstitusiPenyelenggara string  `json:"institusi_penyelenggara"`
	NomorSertifikat        string  `json:"nomor_sertifikat"`
	Tahun                  int16   `json:"tahun"`
}
