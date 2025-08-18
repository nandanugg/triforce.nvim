package penghargaan

import "gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"

type penghargaan struct {
	ID               int     `json:"id"`
	JenisPenghargaan string  `json:"jenis_penghargaan"`
	NamaPenghargaan  string  `json:"nama_penghargaan"`
	Deskripsi        string  `json:"deskripsi"`
	Tanggal          db.Date `json:"tanggal"`
}
