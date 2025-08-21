package jabatan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type jabatan struct {
	ID        string  `json:"id"`
	Jabatan   string  `json:"jabatan"`
	UnitKerja string  `json:"unit_kerja"`
	TMT       db.Date `json:"tmt"`
}

type jenisJabatan struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}
