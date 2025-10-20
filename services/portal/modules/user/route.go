package user

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/users", h.list, mwAuth(api.Kode_ManajemenAkses_Read))
	e.Add(http.MethodGet, "/v1/users/:nip", h.get, mwAuth(api.Kode_ManajemenAkses_Read))
}
