package riwayatjabatan

type riwayatJabatan struct {
	ID                      int64  `json:"id"`
	JenisJabatan            string `json:"jenis_jabatan"`
	NamaJabatan             string `json:"nama_jabatan"`
	TmtJabatan              string `json:"tmt_jabatan"`
	NoSk                    string `json:"no_sk"`
	TanggalSk               string `json:"tanggal_sk"`
	SatuanKerja             string `json:"satuan_kerja"`
	UnitOrganisasi          string `json:"unit_organisasi"`
	StatusPlt               bool   `json:"status_plt"`
	KelasJabatan            string `json:"kelas_jabatan"`
	PeriodeJabatanStartDate string `json:"periode_jabatan_start_date"`
	PeriodeJabatanEndDate   string `json:"periode_jabatan_end_date"`
}
