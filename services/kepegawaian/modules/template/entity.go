package template

import "time"

type template struct {
	ID        int32     `json:"id"`
	Nama      string    `json:"nama"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
