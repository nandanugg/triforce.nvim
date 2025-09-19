package riwayatkepangkatan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, r repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-kepangkatan", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-kepangkatan/:id/berkas", h.getBerkas, mwAuth())
}
