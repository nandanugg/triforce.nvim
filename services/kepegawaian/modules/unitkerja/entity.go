package unitkerja

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type unitKerjaPublic struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}

type anakUnitKerja struct {
	ID      string `json:"id"`
	Nama    string `json:"nama"`
	HasAnak bool   `json:"has_anak"`
}

type unitKerja struct {
	ID                   string  `json:"id"`
	No                   int32   `json:"no"`
	KodeInternal         string  `json:"kode_internal"`
	Nama                 string  `json:"nama"`
	EselonID             string  `json:"eselon_id"`
	CepatKode            string  `json:"cepat_kode"`
	NamaJabatan          string  `json:"nama_jabatan"`
	NamaPejabat          string  `json:"nama_pejabat"`
	DiatasanID           string  `json:"diatasan_id"`
	InstansiID           string  `json:"instansi_id"`
	PemimpinPnsID        string  `json:"pemimpin_pns_id"`
	JenisUnorID          string  `json:"jenis_unor_id"`
	UnorInduk            string  `json:"unor_induk"`
	JumlahIdealStaff     int16   `json:"jumlah_ideal_staff"`
	Order                int32   `json:"order"`
	IsSatker             bool    `json:"is_satker"`
	Eselon1              string  `json:"eselon_1"`
	Eselon2              string  `json:"eselon_2"`
	Eselon3              string  `json:"eselon_3"`
	Eselon4              string  `json:"eselon_4"`
	ExpiredDate          db.Date `json:"expired_date"`
	Keterangan           string  `json:"keterangan"`
	JenisSatker          string  `json:"jenis_satker"`
	Abbreviation         string  `json:"abbreviation"`
	UnorIndukPenyetaraan string  `json:"unor_induk_penyetaraan"`
	JabatanID            string  `json:"jabatan_id"`
	Waktu                string  `json:"waktu"`
	Peraturan            string  `json:"peraturan"`
	Remark               string  `json:"remark"`
	Aktif                bool    `json:"aktif"`
	EselonNama           string  `json:"eselon_nama"`
}
