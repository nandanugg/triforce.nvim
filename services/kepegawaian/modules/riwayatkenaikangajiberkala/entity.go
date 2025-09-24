package riwayatkenaikangajiberkala

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatKenaikanGajiBerkala struct {
	ID                     int64       `json:"id"`
	IDGolongan             pgtype.Int4 `json:"id_golongan"`
	NamaGolongan           string      `json:"nama_golongan"`
	NamaGolonganPangkat    string      `json:"nama_golongan_pangkat"`
	TMTGolongan            db.Date     `json:"tmt_golongan"`
	MasaKerjaGolonganTahun pgtype.Int2 `json:"masa_kerja_golongan_tahun"`
	MasaKerjaGolonganBulan pgtype.Int2 `json:"masa_kerja_golongan_bulan"`
	NomorSK                string      `json:"nomor_sk"`
	TanggalSK              db.Date     `json:"tanggal_sk"`
	TMTKenaikanGajiBerkala db.Date     `json:"tmt_kenaikan_gaji_berkala"`
	TMTJabatan             db.Date     `json:"tmt_jabatan"`
	GajiPokok              pgtype.Int4 `json:"gaji_pokok"`
	Jabatan                string      `json:"jabatan"`
	Pendidikan             string      `json:"pendidikan"`
	TanggalLulus           db.Date     `json:"tanggal_lulus"`
	KantorPembayaran       string      `json:"kantor_pembayaran"`
	UnitKerjaIndukID       string      `json:"unit_kerja_induk_id"`
	UnitKerjaInduk         string      `json:"unit_kerja_induk"`
	Pejabat                string      `json:"pejabat"`
}
