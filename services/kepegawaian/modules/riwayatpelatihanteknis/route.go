package riwayatpelatihanteknis

import (
	"net/http"

	"github.com/labstack/echo/v4"

	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, repo sqlc.Querier, mwAuth echo.MiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-teknis", h.list, mwAuth)
}
