package pelatihanfungsional

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type pelatihanFungsional struct {
	ID                     string   `json:"id"`
	JenisDiklat            string   `json:"jenis_diklat"`
	NamaDiklat             string   `json:"nama_diklat"`
	Tanggal                *db.Date `json:"tanggal,omitempty"`
	Tahun                  *string  `json:"tahun,omitempty"`
	InstitusiPenyelenggara string   `json:"institusi_penyelenggara"`
	NomorSertifikat        string   `json:"nomor_sertifikat"`
}
