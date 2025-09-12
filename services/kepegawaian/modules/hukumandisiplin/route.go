package hukumandisiplin

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth echo.MiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-hukuman-disiplin", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/riwayat-hukuman-disiplin/:id/berkas", h.getBerkas, mwAuth)
}
