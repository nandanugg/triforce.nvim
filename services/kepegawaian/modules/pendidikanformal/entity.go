package pendidikanformal

type pendidikanFormal struct {
	ID                   int    `json:"id"`
	JenjangPendidikan    string `json:"jenjang_pendidikan"`
	NamaSekolah          string `json:"nama_sekolah"`
	Jurusan              string `json:"jurusan"`
	KeteranganPendidikan string `json:"keterangan_pendidikan"`
	TahunLulus           string `json:"tahun_lulus"`
	NomorIjazah          string `json:"nomor_ijazah"`
}
