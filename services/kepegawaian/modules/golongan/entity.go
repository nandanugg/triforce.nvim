package golongan

type golonganPublic struct {
	ID          int32  `json:"id"`
	Nama        string `json:"nama"`
	NamaPangkat string `json:"nama_pangkat"`
}

type golongan struct {
	ID          int32  `json:"id"`
	Nama        string `json:"nama"`
	NamaPangkat string `json:"nama_pangkat"`
	Nama2       string `json:"nama_2"`
	Gol         int16  `json:"gol"`
	GolPppk     string `json:"gol_pppk"`
}
