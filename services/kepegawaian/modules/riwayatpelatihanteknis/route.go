package riwayatpelatihanteknis

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, repo sqlc.Querier, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-teknis", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-teknis/:id/berkas", h.getBerkas, mwAuth())
}
