package suratkeputusan

import (
	"time"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type logSuratKeputusan struct {
	Log       string    `json:"log"`
	Actor     string    `json:"actor"`
	Timestamp time.Time `json:"timestamp"`
}
type suratKeputusan struct {
	IDSK              string               `json:"id_sk"`
	KategoriSK        string               `json:"kategori_sk"`
	NoSK              string               `json:"no_sk"`
	TanggalSK         db.Date              `json:"tanggal_sk"`
	StatusSK          string               `json:"status_sk"`
	UnitKerja         string               `json:"unit_kerja"`
	NamaPemilik       string               `json:"nama_pemilik,omitempty"`
	NamaPenandaTangan string               `json:"nama_penandatangan,omitempty"`
	Logs              *[]logSuratKeputusan `json:"logs,omitempty"`
}

func statusSKText(statusSK int16) string {
	switch statusSK {
	case 0:
		return "Belum Dikoreksi"
	case 1:
		return "Sedang Dikoreksi"
	case 2:
		return "Sudah Dikoreksi & Dikembalikan"
	case 3:
		return "Sudah dikoreksi dan menunggu TTD"
	case 4:
		return "Sudah TTD"
	case 5:
		return "Sudah dikirim"
	default:
		return "Status Tidak Diketahui"
	}
}
