package keluarga

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type keluarga struct {
	OrangTua []orangTua `json:"orang_tua"`
	Pasangan []pasangan `json:"pasangan"`
	Anak     []anak     `json:"anak"`
}

type orangTua struct {
	ID          int32  `json:"id"`
	Nama        string `json:"nama"`
	NIK         string `json:"nik"`
	Agama       string `json:"agama"`
	Hubungan    string `json:"hubungan"`
	StatusHidup string `json:"status_hidup"`
}

type pasangan struct {
	ID               int64   `json:"id"`
	Nama             string  `json:"nama"`
	NIK              string  `json:"nik"`
	StatusPNS        string  `json:"status_pns"`
	Agama            string  `json:"agama"`
	StatusNikah      string  `json:"status_nikah"`
	TanggalMenikah   db.Date `json:"tanggal_menikah"`
	TanggalMeninggal db.Date `json:"tanggal_meninggal"`
	TanggalCerai     db.Date `json:"tanggal_cerai"`
	TanggalLahir     db.Date `json:"tanggal_lahir"`
	AkteNikah        string  `json:"akte_nikah"`
	AkteMeninggal    string  `json:"akte_meninggal"`
	AkteCerai        string  `json:"akte_cerai"`
	Karsus           string  `json:"karsus"`
}

type anak struct {
	ID            int64   `json:"id"`
	Nama          string  `json:"nama"`
	NIK           string  `json:"nik"`
	JenisKelamin  string  `json:"jenis_kelamin"`
	StatusAnak    string  `json:"status_anak"`
	StatusSekolah string  `json:"status_sekolah"`
	TanggalLahir  db.Date `json:"tanggal_lahir"`
	NamaOrangTua  string  `json:"nama_orang_tua"`
	AnakKe        int64   `json:"anak_ke"`
}
