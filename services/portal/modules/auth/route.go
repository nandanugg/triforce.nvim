package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
)

func RegisterRoutes(e *echo.Echo, keycloak config.Keycloak, client *http.Client) {
	s := newService(keycloak, client)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/auth/login", h.login)
	e.Add(http.MethodGet, "/v1/auth/logout", h.logout)
	e.Add(http.MethodPost, "/v1/auth/exchange-token", h.exchangeToken)
	e.Add(http.MethodPost, "/v1/auth/refresh-token", h.refreshToken)
}
