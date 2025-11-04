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
