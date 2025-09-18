package riwayatpelatihanstruktural

import (
	"net/http"

	"github.com/labstack/echo/v4"

	dbRepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, db dbRepo.Querier, mwAuth echo.MiddlewareFunc) {
	s := newService(db)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-struktural", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-struktural/:id/berkas", h.getBerkas, mwAuth)
}
