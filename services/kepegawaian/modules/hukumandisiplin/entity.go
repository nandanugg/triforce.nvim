package hukumandisiplin

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type hukumanDisiplin struct {
	ID           int64   `json:"id"`
	JenisHukuman string  `json:"jenis_hukuman"`
	NomorSK      string  `json:"nomor_sk"`
	TanggalSK    db.Date `json:"tanggal_sk"`
	TanggalMulai db.Date `json:"tanggal_mulai"`
	MasaTahun    int     `json:"masa_tahun"`
	MasaBulan    int     `json:"masa_bulan"`
}
