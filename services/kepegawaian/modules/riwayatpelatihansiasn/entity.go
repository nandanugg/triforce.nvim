package riwayatpelatihansiasn

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatPelatihanSIASN struct {
	ID                     int64       `json:"id"`
	JenisDiklat            string      `json:"jenis_diklat"`
	NamaDiklat             string      `json:"nama_diklat"`
	InstitusiPenyelenggara string      `json:"institusi_penyelenggara"`
	NomorSertifikat        string      `json:"nomor_sertifikat"`
	TanggalMulai           db.Date     `json:"tanggal_mulai"`
	TanggalSelesai         db.Date     `json:"tanggal_selesai"`
	Tahun                  pgtype.Int4 `json:"tahun"`
	Durasi                 pgtype.Int4 `json:"durasi"` // hour
}
