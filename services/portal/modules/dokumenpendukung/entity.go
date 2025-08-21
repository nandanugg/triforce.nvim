package dokumenpendukung

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type dokumenPendukung struct {
	ID                 int64    `json:"id"`
	NamaTombol         string   `json:"nama_tombol"`
	NamaHalaman        string   `json:"nama_halaman"`
	DiperbaruiOleh     *string  `json:"diperbarui_oleh,omitempty"`
	TerakhirDiperbarui *db.Date `json:"terakhir_diperbarui,omitempty"`
	Status             string   `json:"status"`
}
