package datapribadi

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type dataPribadi struct {
	ID                       int     `json:"id"`
	NIP                      string  `json:"nip"`
	NIPBaru                  string  `json:"nip_baru"`
	NIK                      string  `json:"nik"`
	TempatLahir              string  `json:"tempat_lahir"`
	TanggalLahir             db.Date `json:"tanggal_lahir"`
	EmailDikbud              string  `json:"email_dikbud"`
	EmailLain                string  `json:"email_lain"`
	Alamat                   string  `json:"alamat"`
	Nama                     string  `json:"nama"`
	GelarDepan               string  `json:"gelar_depan"`
	GelarBelakang            string  `json:"gelar_belakang"`
	NomorHP                  string  `json:"nomor_hp"`
	Agama                    *string `json:"agama,omitempty"`
	Photo                    string  `json:"photo"`
	JenisKelamin             string  `json:"jenis_kelamin"`
	TingkatPendidikan        string  `json:"tingkat_pendidikan"`
	Pendidikan               string  `json:"pendidikan"`
	MasaKerja                string  `json:"masa_kerja"`
	LokasiKerja              string  `json:"lokasi_kerja"`
	JenisPegawai             *string `json:"jenis_pegawai,omitempty"`
	PangkatGolonganAktif     string  `json:"pangkat_golongan_aktif"`
	GolonganRuangAwal        string  `json:"golongan_ruang_awal"`
	GolonganRuangTerakhir    string  `json:"golongan_ruang_terakhir"`
	TMTGolongan              db.Date `json:"tmt_golongan"`
	Jabatan                  string  `json:"jabatan"`
	UnitKerja                *string `json:"unit_kerja,omitempty"`
	GajiPokok                string  `json:"gaji_pokok"`
	TMTASN                   db.Date `json:"tmt_asn"`
	SKASN                    string  `json:"sk_asn"`
	StatusPNS                string  `json:"status_pns"`
	KartuPegawai             *string `json:"kartu_pegawai"`
	StatusPerkawinan         *string `json:"status_perkawinan,omitempty"`
	NomorSuratDokter         string  `json:"nomor_surat_dokter"`
	TanggalSuratDokter       db.Date `json:"tanggal_surat_dokter"`
	NomorSuratBebasNarkoba   string  `json:"nomor_surat_bebas_narkoba"`
	TanggalSuratBebasNarkoba db.Date `json:"tanggal_surat_bebas_narkoba"`
	NomorCatatanPolisi       string  `json:"nomor_catatan_polisi"`
	TanggalCatatanPolisi     db.Date `json:"tanggal_catatan_polisi"`
	AkteKelahiran            string  `json:"akte_kelahiran"`
	NomorBPJS                string  `json:"nomor_bpjs"`
	NPWP                     string  `json:"npwp"`
	TanggalNPWP              db.Date `json:"tanggal_npwp"`
	NomorDarurat             string  `json:"nomor_darurat"`
}

type statusPernikahan struct {
	ID   string `json:"id"`
	Nama string `json:"nama"`
}
