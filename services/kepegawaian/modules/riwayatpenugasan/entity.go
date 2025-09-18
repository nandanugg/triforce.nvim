package riwayatpenugasan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type riwayatPenugasan struct {
	ID               int64   `json:"id"`
	TipeJabatan      string  `json:"tipe_jabatan"`
	DeskripsiJabatan string  `json:"deskripsi_jabatan"`
	TanggalMulai     db.Date `json:"tanggal_mulai"`
	TanggalSelesai   db.Date `json:"tanggal_selesai"`
}
