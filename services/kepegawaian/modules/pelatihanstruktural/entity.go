package pelatihanstruktural

import (
	"time"
)

type pelatihanStruktural struct {
	ID                    int64     `json:"id"`
	JenisDiklat           string    `json:"jenis_diklat"`
	NamaDiklat            string    `json:"nama_diklat"`
	IstitusiPenyelenggara string    `json:"istitusi_penyelenggara"`
	NomorSertifikat       string    `json:"nomor_sertifikat"`
	TanggalMulai          time.Time `json:"tanggal_mulai"`
	TanggalSelesai        time.Time `json:"tanggal_selesai"`
	Tahun                 int       `json:"tahun"`
	DurasiJam             int       `json:"durasi"` // hour
}
