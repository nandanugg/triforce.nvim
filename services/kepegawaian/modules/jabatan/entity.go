package jabatan

type jabatan struct {
	Kode      string `json:"kode"`
	ID        int32  `json:"id"`
	Nama      string `json:"nama"`
	NamaFull  string `json:"nama_full"`
	Jenis     *int16 `json:"jenis"`
	Kelas     *int16 `json:"kelas"`
	Pensiun   *int16 `json:"pensiun"`
	KodeBkn   string `json:"kode_bkn"`
	NamaBkn   string `json:"nama_bkn"`
	Kategori  string `json:"kategori"`
	BknID     string `json:"bkn_id"`
	Tunjangan int64  `json:"tunjangan"`
}

type jabatanPublic struct {
	ID   string `json:"id"`   // kode_jabatan
	Nama string `json:"nama"` // nama_jabatan
}
