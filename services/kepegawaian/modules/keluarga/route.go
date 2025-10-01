package keluarga

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	dbRepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, sqlc dbRepo.Querier, mwAuth api.AuthMiddlewareFunc) {
	s := newService(sqlc)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/keluarga", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/keluarga", h.listAdmin, mwAuth(api.RoleAdmin))
}
