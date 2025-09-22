package riwayatasesmenninebox

import "github.com/jackc/pgx/v5/pgtype"

type riwayatAsesmenNineBox struct {
	ID         int         `json:"id"`
	Tahun      pgtype.Int2 `json:"tahun"`
	Kesimpulan string      `json:"kesimpulan"`
}
