package penghargaan

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, r repository, mwAuth echo.MiddlewareFunc) {
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-penghargaan", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/riwayat-penghargaan/:id/berkas", h.getBerkas, mwAuth)
}
