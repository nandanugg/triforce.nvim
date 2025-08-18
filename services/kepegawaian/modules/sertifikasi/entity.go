package sertifikasi

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type sertifikasi struct {
	ID              int64   `json:"id"`
	NamaSertifikasi string  `json:"nama_sertifikasi"`
	Tahun           int64   `json:"tahun"`
	TanggalUpload   db.Date `json:"tanggal_upload"`
}
