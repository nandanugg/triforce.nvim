package riwayatjabatan

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatJabatan struct {
	ID                      int64       `json:"id"`
	JenisJabatan            string      `json:"jenis_jabatan"`
	IDJabatan               pgtype.Text `json:"id_jabatan"`
	NamaJabatan             string      `json:"nama_jabatan"`
	TmtJabatan              db.Date     `json:"tmt_jabatan"`
	NoSk                    string      `json:"no_sk"`
	TanggalSk               db.Date     `json:"tanggal_sk"`
	SatuanKerja             string      `json:"satuan_kerja"`
	UnitOrganisasi          string      `json:"unit_organisasi"`
	StatusPlt               bool        `json:"status_plt"`
	KelasJabatan            string      `json:"kelas_jabatan"`
	PeriodeJabatanStartDate db.Date     `json:"periode_jabatan_start_date"`
	PeriodeJabatanEndDate   db.Date     `json:"periode_jabatan_end_date"`
}
