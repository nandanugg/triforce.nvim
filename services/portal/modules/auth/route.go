package auth

import (
	"crypto/rsa"
	"database/sql"
	"net/http"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB, keycloak config.Keycloak, client *http.Client, privateKey *rsa.PrivateKey, keyfunc jwt.Keyfunc) {
	r := newRepository(db)
	s := newService(r, keycloak, client, privateKey, keyfunc)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/auth/login", h.login)
	e.Add(http.MethodGet, "/v1/auth/logout", h.logout)
	e.Add(http.MethodPost, "/v1/auth/exchange-token", h.exchangeToken)
	e.Add(http.MethodPost, "/v1/auth/refresh-token", h.refreshToken)
}
