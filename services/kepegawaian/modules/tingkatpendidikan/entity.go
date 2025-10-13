package tingkatpendidikan

import "github.com/jackc/pgx/v5/pgtype"

type tingkatPendidikanPublic struct {
	ID   int32  `json:"id"`
	Nama string `json:"nama"`
}

type tingkatPendidikan struct {
	Nama             string       `json:"nama"`
	Abbreviation     pgtype.Text  `json:"abbreviation"`
	ID               int32        `json:"id"`
	GolonganID       pgtype.Int4  `json:"golongan_id"`
	GolonganAwalID   pgtype.Int4  `json:"golongan_awal_id"`
	Tingkat          pgtype.Int2  `json:"tingkat"`
	NamaGolongan     *pgtype.Text `json:"nama_golongan,omitempty"`
	NamaGolonganAwal *pgtype.Text `json:"nama_golongan_awal,omitempty"`
}
