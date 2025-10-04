package riwayatpelatihanstruktural

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatPelatihanStruktural struct {
	ID         string        `json:"id"`
	NamaDiklat string        `json:"nama_diklat"`
	Tanggal    db.Date       `json:"tanggal"`
	Tahun      pgtype.Int2   `json:"tahun"`
	Nomor      string        `json:"nomor"`
	Lama       pgtype.Float4 `json:"lama"` // hour
}
