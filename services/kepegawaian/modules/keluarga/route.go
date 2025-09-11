package keluarga

import (
	"net/http"

	"github.com/labstack/echo/v4"

	dbRepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, sqlc dbRepo.Querier, mwAuth echo.MiddlewareFunc) {
	s := newService(sqlc)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/keluarga", h.list, mwAuth)
}
