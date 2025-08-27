package kinerja

type kinerja struct {
	ID             int    `json:"id"`
	Tahun          int    `json:"tahun"`
	HasilKinerja   string `json:"hasil_kinerja"`
	PerilakuKerja  string `json:"perilaku_kerja"`
	KuadranKinerja string `json:"kuadran_kinerja"`
}
