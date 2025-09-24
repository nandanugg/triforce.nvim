package riwayatpelatihanteknis

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatPelatihanTeknis struct {
	ID                     int64         `json:"id"`
	TipePelatihan          string        `json:"tipe_pelatihan"`
	JenisPelatihan         string        `json:"jenis_pelatihan"`
	NamaPelatihan          string        `json:"nama_pelatihan"`
	TanggalMulai           db.Date       `json:"tanggal_mulai"`
	TanggalSelesai         db.Date       `json:"tanggal_selesai"`
	Tahun                  *int          `json:"tahun"`
	Durasi                 pgtype.Float8 `json:"durasi"` // hour
	InstitusiPenyelenggara string        `json:"institusi_penyelenggara"`
	NomorSertifikat        string        `json:"nomor_sertifikat"`
}
