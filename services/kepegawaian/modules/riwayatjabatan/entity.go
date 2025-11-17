package riwayatjabatan

import (
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type riwayatJabatan struct {
	ID                      int64       `json:"id"`
	JenisJabatanID          pgtype.Int4 `json:"jenis_jabatan_id"`
	JenisJabatan            string      `json:"jenis_jabatan"`
	IDJabatan               pgtype.Text `json:"id_jabatan"`
	NamaJabatan             string      `json:"nama_jabatan"`
	TmtJabatan              db.Date     `json:"tmt_jabatan"`
	NoSk                    string      `json:"no_sk"`
	TanggalSk               db.Date     `json:"tanggal_sk"`
	SatuanKerjaID           pgtype.Text `json:"satuan_kerja_id"`
	SatuanKerja             string      `json:"satuan_kerja"`
	UnitOrganisasiID        pgtype.Text `json:"unit_organisasi_id"`
	UnitOrganisasi          string      `json:"unit_organisasi"`
	StatusPlt               bool        `json:"status_plt"`
	KelasJabatanID          pgtype.Int4 `json:"kelas_jabatan_id"`
	KelasJabatan            string      `json:"kelas_jabatan"`
	PeriodeJabatanStartDate db.Date     `json:"periode_jabatan_start_date"`
	PeriodeJabatanEndDate   db.Date     `json:"periode_jabatan_end_date"`
}
