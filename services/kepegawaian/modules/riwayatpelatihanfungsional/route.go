package riwayatpelatihanfungsional

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db repository, mwAuth echo.MiddlewareFunc) {
	s := newService(db)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-fungsional", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-fungsional/:id/berkas", h.getBerkas, mwAuth)
}
