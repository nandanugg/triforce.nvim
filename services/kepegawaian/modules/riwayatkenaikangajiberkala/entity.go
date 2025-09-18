package riwayatkenaikangajiberkala

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type riwayatKenaikanGajiBerkala struct {
	ID        int64   `json:"id"`
	TMTKGB    db.Date `json:"tmt_kgb"`
	NoSK      string  `json:"no_sk"`
	TglSK     db.Date `json:"tgl_sk"`
	GolRuang  string  `json:"gol_ruang"`
	GajiPokok string  `json:"gaji_pokok"`
}
