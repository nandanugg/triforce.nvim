package user

import (
	"time"

	"github.com/google/uuid"
)

type user struct {
	NIP      string    `json:"nip"`
	Profiles []profile `json:"profiles"`
	Roles    []role    `json:"roles"`
}

type role struct {
	ID        int16  `json:"id"`
	Nama      string `json:"nama"`
	IsDefault bool   `json:"is_default"`
	IsAktif   bool   `json:"is_aktif"`
}

type profile struct {
	ID          uuid.UUID  `json:"id"`
	Source      string     `json:"source"`
	Nama        *string    `json:"nama"`
	Email       *string    `json:"email"`
	LastLoginAt *time.Time `json:"last_login_at"`
}
