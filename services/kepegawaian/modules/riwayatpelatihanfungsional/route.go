package riwayatpelatihanfungsional

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, db repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(db)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-fungsional", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-fungsional/:id/berkas", h.getBerkas, mwAuth())
}
