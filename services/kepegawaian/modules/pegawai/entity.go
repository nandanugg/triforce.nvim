package pegawai

type profile struct {
	NIPLama        string   `json:"nip_lama"`
	NIPBaru        string   `json:"nip_baru"`
	GelarDepan     string   `json:"gelar_depan"`
	GelarBelakang  string   `json:"gelar_belakang"`
	Nama           string   `json:"nama"`
	Pangkat        string   `json:"pangkat"`
	Golongan       string   `json:"golongan"`
	Jabatan        string   `json:"jabatan"`
	UnitOrganisasi []string `json:"unit_organisasi"`
}
