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
	IDSK                 string               `json:"id_sk"`
	KategoriSK           string               `json:"kategori_sk"`
	NoSK                 string               `json:"no_sk"`
	TanggalSK            db.Date              `json:"tanggal_sk"`
	StatusSK             string               `json:"status_sk"`
	UnitKerja            string               `json:"unit_kerja"`
	NamaPemilik          string               `json:"nama_pemilik,omitempty"`
	NIPPemilik           string               `json:"nip_pemilik,omitempty"`
	NamaPenandaTangan    string               `json:"nama_penandatangan,omitempty"`
	JabatanPenandaTangan string               `json:"jabatan_penandatangan,omitempty"`
	Logs                 *[]logSuratKeputusan `json:"logs,omitempty"`
}

type statusSK int16

const (
	statusSKBelumDikoreksi  statusSK = 0
	statusSKSedangDikoreksi statusSK = 1
	statusSKDikembalikan    statusSK = 2
	statusSKSudahDikoreksi  statusSK = 3
	statusSKSudahTtd        statusSK = 4
	statusSKSudahKirim      statusSK = 5
)

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

type koreksiSuratKeputusan struct {
	IDSK         string                   `json:"id_sk"`
	KategoriSK   string                   `json:"kategori_sk"`
	NoSK         string                   `json:"no_sk"`
	TanggalSK    db.Date                  `json:"tanggal_sk"`
	UnitKerja    string                   `json:"unit_kerja"`
	NamaPemilik  string                   `json:"nama_pemilik"`
	NIPPemilik   string                   `json:"nip_pemilik"`
	ListKorektor []korektorSuratKeputusan `json:"list_korektor,omitempty"`
	Aksi         *string                  `json:"aksi,omitempty"`
}

type korektorSuratKeputusan struct {
	KorektorKe     int16  `json:"korektor_ke"`
	Nama           string `json:"nama"`
	NIP            string `json:"nip"`
	GelarDepan     string `json:"gelar_depan"`
	GelarBelakang  string `json:"gelar_belakang"`
	StatusKoreksi  string `json:"status_koreksi"`
	CatatanKoreksi string `json:"catatan_koreksi"`
}

var statusKoreksiMap = map[string]int32{
	"Sudah Dikoreksi": 1,
	"Belum Dikoreksi": 2,
}

func statusKoreksiValue(status string) *int32 {
	if val, ok := statusKoreksiMap[status]; ok {
		return &val
	}
	return nil
}

type antrianKoreksiSuratKeputusan struct {
	KategoriSK string `json:"kategori_sk"`
	Jumlah     int64  `json:"jumlah"`
}

type statusKoreksiSK int16

const (
	statusKoreksiBelumDikoreksi statusKoreksiSK = 0
	statusKoreksiSudahDikoreksi statusKoreksiSK = 1
	statusKoreksiDikembalikan   statusKoreksiSK = 3
)

func (s statusKoreksiSK) sudahDikoreksi() bool {
	return s == statusKoreksiSudahDikoreksi
}

type statusTtd int16

const (
	statusTtdBelumTtd     statusTtd = 0
	statusTtdSudahTtd     statusTtd = 1
	statusTtdDikembalikan statusTtd = 3
)

func (s statusTtd) belumTtd() bool {
	return s == statusTtdBelumTtd
}

type statusKorektorSK int16

const (
	statusKorektorSKBelumDikoreksi statusKorektorSK = 2
	statusKorektorSKSudahDikoreksi statusKorektorSK = 1
	statusKorektorSKDikembalikan   statusKorektorSK = 3
)

func (s statusKorektorSK) sudahDikoreksi() bool {
	return s == statusKorektorSKSudahDikoreksi
}

func (s statusKorektorSK) belumDikoreksi() bool {
	return s == statusKorektorSKBelumDikoreksi
}

func (s statusKorektorSK) dikembalikan() bool {
	return s == statusKorektorSKDikembalikan
}

func (s statusKorektorSK) String() string {
	switch s {
	case statusKorektorSKBelumDikoreksi:
		return "Belum Dikoreksi"
	case statusKorektorSKSudahDikoreksi:
		return "Sudah Dikoreksi"
	case statusKorektorSKDikembalikan:
		return "Dikembalikan"
	}
	return ""
}

type suratKeputusanRiwayatMessage string

const (
	diteruskanKePenandatangan       suratKeputusanRiwayatMessage = "SK di teruskan ke penandatangan"
	diteruskanKeKorektorSelanjutnya suratKeputusanRiwayatMessage = "SK di teruskan ke korektor " // korektor n
	dikembalikan                    suratKeputusanRiwayatMessage = "Koreksi SK (SK dikembalikan)"
)
