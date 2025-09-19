package riwayatpenugasan

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth echo.MiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-penugasan", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/riwayat-penugasan/:id/berkas", h.getBerkas, mwAuth)
}
