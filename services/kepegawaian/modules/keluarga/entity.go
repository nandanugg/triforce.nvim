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

var labelStatusSekolah = map[int16]statusSekolah{
	1: statusMasihSekolah,
	2: statusSudahBekerja,
}

func (s statusSekolah) toID() pgtype.Int2 {
	for status, label := range labelStatusSekolah {
		if s == label {
			return pgtype.Int2{Int16: status, Valid: true}
		}
	}
	return pgtype.Int2{}
}

type statusAnak string

const (
	statusAnakKandung statusAnak = "Kandung"
	statusAnakAngkat  statusAnak = "Angkat"
)

var labelStatusAnak = map[string]statusAnak{
	"1": statusAnakKandung,
	"2": statusAnakAngkat,
}

func (s statusAnak) toID() pgtype.Text {
	for status, label := range labelStatusAnak {
		if s == label {
			return pgtype.Text{String: status, Valid: true}
		}
	}
	return pgtype.Text{}
}

type hubunganOrangTua string

const (
	hubunganAyah hubunganOrangTua = "Ayah"
	hubunganIbu  hubunganOrangTua = "Ibu"
)

var labelHubunganOrangTua = map[int16]hubunganOrangTua{
	1: hubunganAyah,
	2: hubunganIbu,
}

func (s hubunganOrangTua) toID() pgtype.Int2 {
	for status, label := range labelHubunganOrangTua {
		if s == label {
			return pgtype.Int2{Int16: status, Valid: true}
		}
	}
	return pgtype.Int2{}
}

type hubunganPasangan string

const (
	hubunganIstri hubunganPasangan = "Istri"
	hubunganSuami hubunganPasangan = "Suami"
)

var labelHubunganPasangan = map[int16]hubunganPasangan{
	1: hubunganIstri,
	2: hubunganSuami,
}

func (s hubunganPasangan) toID() pgtype.Int2 {
	for status, label := range labelHubunganPasangan {
		if s == label {
			return pgtype.Int2{Int16: status, Valid: true}
		}
	}
	return pgtype.Int2{}
}
