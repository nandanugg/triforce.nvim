package riwayatpekerjaan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type riwayatPekerjaan struct {
	ID               int     `json:"id"`
	PNSNIP           string  `json:"pns_nip"`
	JenisPerusahaan  string  `json:"jenis_perusahaan"`
	NamaPerusahaan   string  `json:"nama_perusahaan"`
	Sebagai          string  `json:"sebagai"`
	DariTanggal      db.Date `json:"dari_tanggal"`
	SampaiTanggal    db.Date `json:"sampai_tanggal"`
	PNSID            string  `json:"pns_id"`
	KeteranganBerkas string  `json:"keterangan_berkas"`
}
