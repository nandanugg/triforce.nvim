package keluarga

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type keluarga struct {
	OrangTua []orangTua `json:"orang_tua"`
	Pasangan []pasangan `json:"pasangan"`
	Anak     []anak     `json:"anak"`
}

type orangTua struct {
	ID               int32            `json:"id"`
	Nama             string           `json:"nama"`
	NIK              string           `json:"nik"`
	AgamaID          pgtype.Int2      `json:"agama_id"`
	Agama            string           `json:"agama"`
	Hubungan         hubunganOrangTua `json:"hubungan"`
	StatusHidup      string           `json:"status_hidup"`
	TanggalMeninggal db.Date          `json:"tanggal_meninggal"`
	AkteMeninggal    string           `json:"akte_meninggal"`
}

type pasangan struct {
	ID                 int64       `json:"id"`
	Nama               string      `json:"nama"`
	NIK                string      `json:"nik"`
	StatusPNS          string      `json:"status_pns"`
	AgamaID            pgtype.Int2 `json:"agama_id"`
	Agama              string      `json:"agama"`
	StatusPernikahanID pgtype.Int2 `json:"status_pernikahan_id"`
	StatusNikah        string      `json:"status_nikah"`
	TanggalMenikah     db.Date     `json:"tanggal_menikah"`
	TanggalMeninggal   db.Date     `json:"tanggal_meninggal"`
	TanggalCerai       db.Date     `json:"tanggal_cerai"`
	TanggalLahir       db.Date     `json:"tanggal_lahir"`
	AkteNikah          string      `json:"akte_nikah"`
	AkteMeninggal      string      `json:"akte_meninggal"`
	AkteCerai          string      `json:"akte_cerai"`
	Karsus             string      `json:"karsus"`
}

type anak struct {
	ID                 int64         `json:"id"`
	Nama               string        `json:"nama"`
	NIK                string        `json:"nik"`
	JenisKelamin       string        `json:"jenis_kelamin"`
	StatusAnak         statusAnak    `json:"status_anak"`
	StatusSekolah      statusSekolah `json:"status_sekolah"`
	StatusPernikahanID pgtype.Int2   `json:"status_pernikahan_id"`
	StatusPernikahan   string        `json:"status_pernikahan"`
	TanggalLahir       db.Date       `json:"tanggal_lahir"`
	PasanganOrangTuaID pgtype.Int8   `json:"pasangan_orang_tua_id"`
	NamaOrangTua       string        `json:"nama_orang_tua"`
	AgamaID            pgtype.Int2   `json:"agama_id"`
	Agama              string        `json:"agama"`
	AnakKe             pgtype.Int2   `json:"anak_ke"`
}

type statusSekolah string

const (
	statusMasihSekolah statusSekolah = "Masih Sekolah"
	statusSudahBekerja statusSekolah = "Sudah Bekerja"
)

func (s statusSekolah) toID() pgtype.Int2 {
	switch s {
	case statusMasihSekolah:
		return pgtype.Int2{Int16: 1, Valid: true}
	case statusSudahBekerja:
		return pgtype.Int2{Int16: 2, Valid: true}
	default:
		return pgtype.Int2{}
	}
}

func labelStatusSekolah(statusID pgtype.Int2) statusSekolah {
	switch statusID.Int16 {
	case 1:
		return statusMasihSekolah
	case 2:
		return statusSudahBekerja
	default:
		return ""
	}
}

type statusAnak string

const (
	statusAnakKandung statusAnak = "Kandung"
	statusAnakAngkat  statusAnak = "Angkat"
)

func (s statusAnak) toID() pgtype.Text {
	switch s {
	case statusAnakKandung:
		return pgtype.Text{String: "1", Valid: true}
	case statusAnakAngkat:
		return pgtype.Text{String: "2", Valid: true}
	default:
		return pgtype.Text{}
	}
}

func labelStatusAnak(statusID pgtype.Text) statusAnak {
	switch statusID.String {
	case "1":
		return statusAnakKandung
	case "2":
		return statusAnakAngkat
	default:
		return ""
	}
}

type hubunganOrangTua string

const (
	hubunganAyah hubunganOrangTua = "Ayah"
	hubunganIbu  hubunganOrangTua = "Ibu"
)

func (s hubunganOrangTua) toID() pgtype.Int2 {
	switch s {
	case hubunganAyah:
		return pgtype.Int2{Int16: 1, Valid: true}
	case hubunganIbu:
		return pgtype.Int2{Int16: 2, Valid: true}
	default:
		return pgtype.Int2{}
	}
}

func labelHubunganOrangTua(hubunganID pgtype.Int2) hubunganOrangTua {
	switch hubunganID.Int16 {
	case 1:
		return hubunganAyah
	case 2:
		return hubunganIbu
	default:
		return ""
	}
}

type hubunganPasangan string

const (
	hubunganIstri hubunganPasangan = "Istri"
	hubunganSuami hubunganPasangan = "Suami"
)

func (s hubunganPasangan) toID() pgtype.Int2 {
	switch s {
	case hubunganIstri:
		return pgtype.Int2{Int16: 1, Valid: true}
	case hubunganSuami:
		return pgtype.Int2{Int16: 2, Valid: true}
	default:
		return pgtype.Int2{}
	}
}
