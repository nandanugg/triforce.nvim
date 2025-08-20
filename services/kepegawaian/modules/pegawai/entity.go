package pegawai

type pegawai struct {
	ID            int    `json:"id"`
	NIP           string `json:"nip"`
	NamaPegawai   string `json:"nama_pegawai"`
	GelarDepan    string `json:"gelar_depan"`
	GelarBelakang string `json:"gelar_belakang"`
	Golongan      string `json:"golongan"`
	Jabatan       string `json:"jabatan"`
	UnitKerja     string `json:"unit_kerja"`
	StatusPegawai string `json:"status_pegawai"`
}
