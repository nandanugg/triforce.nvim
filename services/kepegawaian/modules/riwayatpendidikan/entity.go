package riwayatpendidikan

import "github.com/jackc/pgx/v5/pgtype"

type riwayatPendidikan struct {
	ID                   int32         `json:"id"`
	TingkatPendidikanID  pgtype.Int2   `json:"tingkat_pendidikan_id"`
	JenjangPendidikan    string        `json:"jenjang_pendidikan"`
	PendidikanID         pgtype.Text   `json:"pendidikan_id"`
	Pendidikan           string        `json:"jurusan"`
	NamaSekolah          string        `json:"nama_sekolah"`
	TahunLulus           pgtype.Int2   `json:"tahun_lulus"`
	NomorIjazah          string        `json:"nomor_ijazah"`
	GelarDepan           string        `json:"gelar_depan"`
	GelarBelakang        string        `json:"gelar_belakang"`
	TugasBelajar         statusBelajar `json:"tugas_belajar"`
	KeteranganPendidikan string        `json:"keterangan_pendidikan"`
}

type usulanPerubahanData struct {
	TingkatPendidikanID [2]pgtype.Int2 `json:"tingkat_pendidikan_id"`
	TingkatPendidikan   [2]pgtype.Text `json:"jenjang_pendidikan"`
	PendidikanID        [2]pgtype.Text `json:"pendidikan_id"`
	Pendidikan          [2]pgtype.Text `json:"jurusan"`
	NamaSekolah         [2]pgtype.Text `json:"nama_sekolah"`
	TahunLulus          [2]pgtype.Int2 `json:"tahun_lulus"`
	NomorIjazah         [2]pgtype.Text `json:"nomor_ijazah"`
	GelarDepan          [2]pgtype.Text `json:"gelar_depan"`
	GelarBelakang       [2]pgtype.Text `json:"gelar_belakang"`
	TugasBelajar        [2]pgtype.Text `json:"tugas_belajar"`
	NegaraSekolah       [2]pgtype.Text `json:"negara_sekolah"`
}

type statusBelajar string

const (
	statusTugasBelajar statusBelajar = "Tugas Belajar"
	statusIzinBelajar  statusBelajar = "Izin Belajar"
)

var labelStatusBelajar = map[int16]statusBelajar{
	1: statusTugasBelajar,
	2: statusIzinBelajar,
}

func (s statusBelajar) toID() pgtype.Int2 {
	for status, label := range labelStatusBelajar {
		if s == label {
			return pgtype.Int2{Int16: status, Valid: true}
		}
	}
	return pgtype.Int2{}
}
