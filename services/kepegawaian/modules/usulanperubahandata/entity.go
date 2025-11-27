package usulanperubahandata

import "github.com/jackc/pgx/v5/pgtype"

type usulanPerubahanData struct {
	ID            int64             `json:"id"`
	JenisData     string            `json:"jenis_data"`
	Pegawai       *pegawai          `json:"pegawai,omitempty"`
	DataID        *pgtype.Text      `json:"data_id,omitempty"`
	PerubahanData map[string][2]any `json:"perubahan_data,omitempty"`
	Action        string            `json:"action,omitempty"`
	Status        string            `json:"status,omitempty"`
	Catatan       *string           `json:"catatan,omitempty"`
	CreatedAt     string            `json:"created_at,omitempty"`
}

type pegawai struct {
	NIP           string       `json:"nip"`
	Nama          string       `json:"nama"`
	GelarDepan    string       `json:"gelar_depan"`
	GelarBelakang string       `json:"gelar_belakang"`
	UnitKerja     string       `json:"unit_kerja"`
	Jabatan       *string      `json:"jabatan,omitempty"`
	Golongan      *string      `json:"golongan,omitempty"`
	Photo         *pgtype.Text `json:"photo,omitempty"`
	StatusPNS     *string      `json:"status_pns,omitempty"`
}

const (
	statusDiusulkan = "Diusulkan"
	statusDisetujui = "Disetujui"
	statusDitolak   = "Ditolak"
)

const (
	ActionCreate = "CREATE"
	ActionUpdate = "UPDATE"
	ActionDelete = "DELETE"
)
