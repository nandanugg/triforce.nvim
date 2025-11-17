package auth

import "github.com/jackc/pgx/v5/pgtype"

const (
	sourceZimbra   = "zimbra"
	sourceKeycloak = "keycloak"
)

type token struct {
	AccessToken      string `json:"access_token"`
	ExpiresIn        int    `json:"expires_in"`
	RefreshToken     string `json:"refresh_token"`
	RefreshExpiresIn int    `json:"refresh_expires_in"`
	IDToken          string `json:"id_token"`
}

type user struct {
	id     pgtype.UUID
	source string
	nip    string
}
