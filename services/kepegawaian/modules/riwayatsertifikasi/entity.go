package riwayatsertifikasi

type riwayatSertifikasi struct {
	ID              int64  `json:"id"`
	NamaSertifikasi string `json:"nama_sertifikasi"`
	Tahun           int64  `json:"tahun"`
	Deskripsi       string `json:"deskripsi"`
}
