package jabatan

import "github.com/jackc/pgx/v5/pgtype"

type jabatan struct {
	Kode      string      `json:"kode"`
	ID        int32       `json:"id"`
	Nama      string      `json:"nama"`
	NamaFull  string      `json:"nama_full"`
	NamaJenis *string     `json:"nama_jenis,omitempty"`
	Jenis     pgtype.Int2 `json:"jenis"`
	Kelas     pgtype.Int2 `json:"kelas"`
	Pensiun   pgtype.Int2 `json:"pensiun"`
	KodeBkn   string      `json:"kode_bkn"`
	NamaBkn   string      `json:"nama_bkn"`
	Kategori  string      `json:"kategori"`
	BknID     string      `json:"bkn_id"`
	Tunjangan int64       `json:"tunjangan"`
}

type jabatanPublic struct {
	ID   string `json:"id"`   // kode_jabatan
	Nama string `json:"nama"` // nama_jabatan
}
