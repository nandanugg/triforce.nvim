package sertifikasi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth echo.MiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/sertifikasi", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/sertifikasi/:id/berkas", h.getBerkas, mwAuth)
}
