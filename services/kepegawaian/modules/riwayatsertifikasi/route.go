package riwayatsertifikasi

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-sertifikasi", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-sertifikasi/:id/berkas", h.getBerkas, mwAuth())
}
