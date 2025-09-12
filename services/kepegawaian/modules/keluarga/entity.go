package keluarga

import "time"

type keluarga struct {
	OrangTua []orangTua `json:"orang_tua"`
	Pasangan []pasangan `json:"pasangan"`
	Anak     []anak     `json:"anak"`
}

type orangTua struct {
	Nama            *string `json:"nama"`
	Nik             *string `json:"nik"`
	Agama           *string `json:"agama"`
	Hubungan        string  `json:"hubungan"`
	StatusMeninggal string  `json:"tanggal_meninggal"`
}

type pasangan struct {
	Nama           *string    `json:"nama"`
	Nik            *string    `json:"nik"`
	StatusPNS      string     `json:"status_pns"`
	Agama          *string    `json:"agama"`
	StatusNikah    string     `json:"status_nikah"`
	TanggalMenikah *time.Time `json:"tanggal_menikah"`
}

type anak struct {
	Nama          *string    `json:"nama"`
	Nik           *string    `json:"nik"`
	JenisKelamin  string     `json:"jenis_kelamin"`
	StatusAnak    string     `json:"status_anak"`
	StatusSekolah string     `json:"status_sekolah"`
	TanggalLahir  *time.Time `json:"tanggal_lahir"`
	AnakKe        *int64     `json:"anak_ke"`
}
