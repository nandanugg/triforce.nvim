package pemberitahuan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type pemberitahuan struct {
	ID                 int64   `json:"id"`
	JudulBerita        string  `json:"judul_berita"`
	DeskripsiBerita    string  `json:"deskripsi_berita"`
	Status             string  `json:"status"`
	DiperbaruiOleh     string  `json:"diperbarui_oleh"`
	TerakhirDiperbarui db.Date `json:"terakhir_diperbarui"`
}
