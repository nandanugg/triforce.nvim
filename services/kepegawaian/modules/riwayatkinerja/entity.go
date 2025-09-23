package riwayatkinerja

import "github.com/jackc/pgx/v5/pgtype"

type riwayatKinerja struct {
	ID             int32       `json:"id"`
	Tahun          pgtype.Int4 `json:"tahun"`
	HasilKinerja   string      `json:"hasil_kinerja"`
	PerilakuKerja  string      `json:"perilaku_kerja"`
	KuadranKinerja string      `json:"kuadran_kinerja"`
}
