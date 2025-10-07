package unitkerja

type unitKerjaPublic struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

type anakUnitKerja struct {
	ID      string `json:"id"`
	Nama    string `json:"nama"`
	HasAnak bool   `json:"has_anak"`
}
