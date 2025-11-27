package riwayatkepangkatan

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatKepangkatan struct {
	ID                        string      `json:"id"`
	IDJenisKP                 pgtype.Int4 `json:"id_jenis_kp"`
	NamaJenisKP               string      `json:"nama_jenis_kp"`
	IDGolongan                pgtype.Int4 `json:"id_golongan"`
	NamaGolongan              string      `json:"nama_golongan"`
	NamaGolonganPangkat       string      `json:"nama_golongan_pangkat"`
	TMTGolongan               db.Date     `json:"tmt_golongan"`
	SKNomor                   string      `json:"sk_nomor"`
	SKTanggal                 db.Date     `json:"sk_tanggal"`
	MKGolonganTahun           pgtype.Int2 `json:"mk_golongan_tahun"`
	MKGolonganBulan           pgtype.Int2 `json:"mk_golongan_bulan"`
	NoBKN                     string      `json:"no_bkn"`
	TanggalBKN                db.Date     `json:"tanggal_bkn"`
	JumlahAngkaKreditTambahan pgtype.Int4 `json:"jumlah_angka_kredit_tambahan"`
	JumlahAngkaKreditUtama    pgtype.Int4 `json:"jumlah_angka_kredit_utama"`
}

type usulanPerubahanData struct {
	JenisKPID                 [2]pgtype.Int4 `json:"jenis_kp_id"`
	NamaJenisKP               [2]pgtype.Text `json:"nama_jenis_kp"`
	GolonganID                [2]pgtype.Int2 `json:"golongan_id"`
	NamaGolongan              [2]pgtype.Text `json:"nama_golongan"`
	NamaGolonganPangkat       [2]pgtype.Text `json:"nama_golongan_pangkat"`
	TMTGolongan               [2]db.Date     `json:"tmt_golongan"`
	NomorSK                   [2]pgtype.Text `json:"nomor_sk"`
	TanggalSK                 [2]db.Date     `json:"tanggal_sk"`
	NomorBKN                  [2]pgtype.Text `json:"nomor_bkn"`
	TanggalBKN                [2]db.Date     `json:"tanggal_bkn"`
	MasaKerjaGolonganTahun    [2]pgtype.Int2 `json:"masa_kerja_golongan_tahun"`
	MasaKerjaGolonganBulan    [2]pgtype.Int2 `json:"masa_kerja_golongan_bulan"`
	JumlahAngkaKreditUtama    [2]pgtype.Int4 `json:"jumlah_angka_kredit_utama"`
	JumlahAngkaKreditTambahan [2]pgtype.Int4 `json:"jumlah_angka_kredit_tambahan"`
}
