package riwayatsertifikasi

import "github.com/jackc/pgx/v5/pgtype"

type riwayatSertifikasi struct {
	ID              int64       `json:"id"`
	NamaSertifikasi string      `json:"nama_sertifikasi"`
	Tahun           pgtype.Int8 `json:"tahun"`
	Deskripsi       string      `json:"deskripsi"`
}
