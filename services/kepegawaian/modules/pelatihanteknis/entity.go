package pelatihanteknis

import (
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type pelatihanTeknis struct {
	ID                     int64   `json:"id"`
	TipePelatihan          string  `json:"tipe_pelatihan"`          // Tipe Pelatihan (Sertifikat/Non-sertifikat)
	JenisPelatihan         string  `json:"jenis_pelatihan"`         // Jenis Pelatihan (Workshop/Seminar/Kursus/Lainnya)
	NamaPelatihan          string  `json:"nama_pelatihan"`          // Nama Pelatihan
	TanggalMulai           db.Date `json:"tanggal_mulai"`           // Tanggal Mulai
	TanggalSelesai         db.Date `json:"tanggal_selesai"`         // Tanggal Selesai
	Tahun                  *int    `json:"tahun"`                   // Tahun
	Durasi                 int     `json:"durasi"`                  // Durasi (Jam)
	InstitusiPenyelenggara string  `json:"institusi_penyelenggara"` // Institusi Penyelenggara
	NomorSertifikat        string  `json:"nomor_sertifikat"`        // Nomor Sertifikat
}
