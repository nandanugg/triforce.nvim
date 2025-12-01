package riwayatpenghargaan

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatPenghargaan struct {
	ID               int     `json:"id"`
	JenisPenghargaan string  `json:"jenis_penghargaan"`
	NamaPenghargaan  string  `json:"nama_penghargaan"`
	Deskripsi        string  `json:"deskripsi"`
	Tanggal          db.Date `json:"tanggal"`
}

type usulanPerubahanData struct {
	JenisPenghargaan [2]pgtype.Text `json:"jenis_penghargaan"`
	NamaPenghargaan  [2]pgtype.Text `json:"nama_penghargaan"`
	Deskripsi        [2]pgtype.Text `json:"deskripsi"`
	Tanggal          [2]db.Date     `json:"tanggal"`
}

type JenisPenghargaan string

const (
	JenisPenghargaanInternational             JenisPenghargaan = "Internasional"
	JenisPenghargaanUnitKerjaEselonDuaKebawah JenisPenghargaan = "Unit Kerja (eselon 2 ke bawah)"
	JenisPenghargaanUnitUtama                 JenisPenghargaan = "Unit Utama"
	JenisPenghargaanNasional                  JenisPenghargaan = "Nasional"
	JenisPenghargaanInstansional              JenisPenghargaan = "Instansional (Kementerian/Lembaga)"
)

func (s *service) validateJenisPenghargaan(jenis string) (JenisPenghargaan, bool) {
	switch jenis {
	case
		string(JenisPenghargaanInternational),
		string(JenisPenghargaanUnitKerjaEselonDuaKebawah),
		string(JenisPenghargaanUnitUtama),
		string(JenisPenghargaanNasional),
		string(JenisPenghargaanInstansional):
		return JenisPenghargaan(jenis), true
	}

	return "", false
}
