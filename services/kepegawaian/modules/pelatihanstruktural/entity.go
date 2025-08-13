package pelatihanstruktural

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type pelatihanStruktural struct {
	NamaDiklat string  `json:"nama_diklat"`
	Nomor      *string `json:"nomor,omitempty"`
	Tanggal    db.Date `json:"tanggal"`
	Tahun      int     `json:"tahun"`
}
