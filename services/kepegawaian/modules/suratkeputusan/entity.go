package suratkeputusan

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"

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
	StatusSK     string                   `json:"status_sk,omitempty"`
	ListKorektor []korektorSuratKeputusan `json:"list_korektor,omitempty"`
	Aksi         *string                  `json:"aksi,omitempty"`
	KorektorKe   int                      `json:"korektor_ke,omitempty"`
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
	dikembalikanOlehPenandatangan   suratKeputusanRiwayatMessage = "SK dikembalikan oleh penandatangan"
	ditandatangani                  suratKeputusanRiwayatMessage = "Berhasil ditandatangan"
)

type statusTandaTangan string

const (
	statusTandaTanganBelumDitandatangan statusTandaTangan = "Belum Ditandatangan"
	statusTandaTanganSudahDitandatangan statusTandaTangan = "Sudah Ditandatangan"
)

func (s statusTandaTangan) value() int32 {
	switch s {
	case statusTandaTanganBelumDitandatangan:
		return 0
	case statusTandaTanganSudahDitandatangan:
		return 1
	}
	return -1
}

func (s statusTandaTangan) valid() bool {
	return s == statusTandaTanganBelumDitandatangan || s == statusTandaTanganSudahDitandatangan
}

type antreanSK struct {
	KategoriSK string `json:"kategori_sk"`
	Jumlah     int64  `json:"jumlah"`
}

type statusTandatanganRequest string

const (
	statusTandatanganRequestTandatangan  statusTandatanganRequest = "Tandatangan"
	statusTandatanganRequestDikembalikan statusTandatanganRequest = "Dikembalikan"
)

func (s statusTandatanganRequest) tandaTangan() bool {
	return s == statusTandatanganRequestTandatangan
}

func (s statusTandatanganRequest) dikembalikan() bool {
	return s == statusTandatanganRequestDikembalikan
}

type statusKoreksiRequest string

const (
	statusKoreksiRequestSudahDikoreksi statusKoreksiRequest = "Sudah Dikoreksi"
	statusKoreksiRequestBelumDikoreksi statusKoreksiRequest = "Belum Dikoreksi"
)

func (s statusKoreksiRequest) sudahDikoreksi() bool {
	return s == statusKoreksiRequestSudahDikoreksi
}

func (s statusKoreksiRequest) belumDikoreksi() bool {
	return s == statusKoreksiRequestBelumDikoreksi
}

func (s statusKoreksiRequest) value() []int32 {
	if s.sudahDikoreksi() {
		return []int32{
			int32(statusKorektorSK(statusKorektorSKSudahDikoreksi)),
			int32(statusKorektorSK(statusKorektorSKDikembalikan)),
		}
	}

	if s.belumDikoreksi() {
		return []int32{
			int32(statusKorektorSK(statusKorektorSKBelumDikoreksi)),
		}
	}

	return []int32{}
}

type letakTTDsk int16

const (
	letakTTDskKiriBawah   letakTTDsk = 0
	letakTTDskTengahBawah letakTTDsk = 1
	letakTTDskKananBawah  letakTTDsk = 2
)

func (s letakTTDsk) koordinat() (xAxis float64, yAxis float64) {
	switch s {
	case letakTTDskKiriBawah:
		return 5, 10
	case letakTTDskTengahBawah:
		return 150, 10
	case letakTTDskKananBawah:
		return 300, 10
	}
	return 600, 70
}

type errorMessage string

const (
	errorMessageBelumAdaTTD      errorMessage = "Data tanda tangan digital anda belum ada"
	errorMessageNIKNotFound      errorMessage = "NIK tidak ditemukan"
	errorMessageBukanPegawaiTTD  errorMessage = "Anda tidak memiliki izin untuk menandatangani surat keputusan ini."
	errorMessageStatusTtdInvalid errorMessage = "Surat keputusan ini belum siap untuk ditandatangani."
)

func (e errorMessage) Error() string {
	return string(e)
}

type statusLogSuratKeputusan string

const (
	statusLogSuratKeputusanGagal    statusLogSuratKeputusan = "Gagal"
	statusLogSuratKeputusanBerhasil statusLogSuratKeputusan = "Berhasil"
)

var labelStatusLogSuratKeputusan = map[int16]statusLogSuratKeputusan{
	1: statusLogSuratKeputusanGagal,
	2: statusLogSuratKeputusanBerhasil,
}

func (s statusLogSuratKeputusan) toID() pgtype.Int2 {
	for status, label := range labelStatusLogSuratKeputusan {
		if s == label {
			return pgtype.Int2{Int16: status, Valid: true}
		}
	}
	return pgtype.Int2{}
}

type logBSRESuratKeputusan struct {
	FileID     string                  `json:"file_id"`
	NIK        string                  `json:"nik"`
	Nama       string                  `json:"nama"`
	Keterangan string                  `json:"keterangan"`
	Status     statusLogSuratKeputusan `json:"status"`
	Tanggal    db.Date                 `json:"tanggal"`
}
