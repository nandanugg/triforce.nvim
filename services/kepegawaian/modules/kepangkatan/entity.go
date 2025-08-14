package kepangkatan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type kepangkatan struct {
	ID              string  `json:"id"`
	Pangkat         string  `json:"pangkat"`
	Golongan        string  `json:"golongan"`
	TMT             db.Date `json:"tmt"`
	MKGolonganTahun string  `json:"mk_golongan_tahun"`
	MKGolonganBulan string  `json:"mk_golongan_bulan"`
}
