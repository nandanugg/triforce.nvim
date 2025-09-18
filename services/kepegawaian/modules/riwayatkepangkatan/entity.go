package riwayatkepangkatan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type riwayatKepangkatan struct {
	ID                        int32   `json:"id"`
	IDJenisKP                 int32   `json:"id_jenis_kp"`
	NamaJenisKP               string  `json:"nama_jenis_kp"`
	IDGolongan                int32   `json:"id_golongan"`
	NamaGolongan              string  `json:"nama_golongan"`
	NamaGolonganPangkat       string  `json:"nama_golongan_pangkat"`
	TMTGolongan               db.Date `json:"tmt_golongan"`
	SKNomor                   string  `json:"sk_nomor"`
	SKTanggal                 db.Date `json:"sk_tanggal"`
	MKGolonganTahun           int16   `json:"mk_golongan_tahun"`
	MKGolonganBulan           int16   `json:"mk_golongan_bulan"`
	NoBKN                     string  `json:"no_bkn"`
	TanggalBKN                db.Date `json:"tanggal_bkn"`
	JumlahAngkaKreditTambahan int16   `json:"jumlah_angka_kredit_tambahan"`
	JumlahAngkaKreditUtama    int16   `json:"jumlah_angka_kredit_utama"`
}
