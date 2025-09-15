package pendidikanformal

type pendidikanFormal struct {
	ID                   int    `json:"id"`
	JenjangPendidikan    string `json:"jenjang_pendidikan"`
	Pendidikan           string `json:"jurusan"`
	NamaSekolah          string `json:"nama_sekolah"`
	TahunLulus           string `json:"tahun_lulus"`
	NomorIjazah          string `json:"nomor_ijazah"`
	GelarDepan           string `json:"gelar_depan"`
	GelarBelakang        string `json:"gelar_belakang"`
	TugasBelajar         string `json:"tugas_belajar"`
	KeteranganPendidikan string `json:"keterangan_pendidikan"`
}
