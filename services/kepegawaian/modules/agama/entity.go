package agama

import "github.com/jackc/pgx/v5/pgtype"

type agama struct {
	ID        int32              `json:"id"`
	Nama      string             `json:"nama"`
	CreatedAt pgtype.Timestamptz `json:"created_at"`
	UpdatedAt pgtype.Timestamptz `json:"updated_at"`
}
