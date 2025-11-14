package riwayathukumandisiplin

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type riwayatHukumanDisiplin struct {
	ID                  int64   `json:"id"`
	JenisHukuman        string  `json:"jenis_hukuman"`
	JenisHukumanID      int16   `json:"jenis_hukuman_id"`
	NamaGolongan        string  `json:"nama_golongan"`
	GolonganID          int16   `json:"golongan_id"`
	NamaPangkat         string  `json:"nama_pangkat"`
	NomorSK             string  `json:"nomor_sk"`
	TanggalSK           db.Date `json:"tanggal_sk"`
	TanggalMulai        db.Date `json:"tanggal_mulai"`
	TanggalAkhir        db.Date `json:"tanggal_akhir"`
	MasaTahun           int     `json:"masa_tahun"`
	MasaBulan           int     `json:"masa_bulan"`
	NomorPP             string  `json:"nomor_pp"`
	NomorSKPembatalan   string  `json:"nomor_sk_pembatalan"`
	TanggalSKPembatalan db.Date `json:"tanggal_sk_pembatalan"`
}
