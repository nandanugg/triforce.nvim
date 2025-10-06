package unitkerja

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/unit-kerja", h.listUnitKerja, mwAuth())
	e.Add(http.MethodGet, "/v1/unit-kerja/akar", h.listAkarUnitKerja, mwAuth())
	e.Add(http.MethodGet, "/v1/unit-kerja/:id/anak", h.listAnakUnitKerja, mwAuth())
}
