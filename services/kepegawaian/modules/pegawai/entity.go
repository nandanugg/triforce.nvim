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

type pegawai struct {
	NIP           string `json:"nip"`
	GelarDepan    string `json:"gelar_depan"`
	GelarBelakang string `json:"gelar_belakang"`
	Nama          string `json:"nama"`
	Golongan      string `json:"golongan"`
	Jabatan       string `json:"jabatan"`
	UnitKerja     string `json:"unit_kerja"`
	Status        string `json:"status"`
}

const (
	statusPNSMPP   = "Masa Persiapan Pensiun"
	statusPNSAktif = "Aktif"
)

func getStatusHukum(params string) []string {
	switch params {
	case "PNS", "CPNS":
		return []string{statusPNSAktif}
	case "MPP":
		return []string{statusPNSMPP}
	default:
		return []string{statusPNSAktif, statusPNSMPP}
	}
}
