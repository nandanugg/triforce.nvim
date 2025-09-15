package riwayatpendidikan

type riwayatPendidikan struct {
	ID                   int32  `json:"id"`
	JenjangPendidikan    string `json:"jenjang_pendidikan"`
	Pendidikan           string `json:"jurusan"`
	NamaSekolah          string `json:"nama_sekolah"`
	TahunLulus           int16  `json:"tahun_lulus"`
	NomorIjazah          string `json:"nomor_ijazah"`
	GelarDepan           string `json:"gelar_depan"`
	GelarBelakang        string `json:"gelar_belakang"`
	TugasBelajar         string `json:"tugas_belajar"`
	KeteranganPendidikan string `json:"keterangan_pendidikan"`
}

var tugasBelajar = map[int16]string{
	1: "Tugas Belajar",
	2: "Izin Belajar",
}
