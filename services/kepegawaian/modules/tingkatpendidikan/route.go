package tingkatpendidikan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/tingkat-pendidikan", h.listPublic, mwAuth())

	e.Add(http.MethodGet, "/v1/admin/tingkat-pendidikan", h.list, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/tingkat-pendidikan/:id", h.get, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPost, "/v1/admin/tingkat-pendidikan", h.create, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPut, "/v1/admin/tingkat-pendidikan/:id", h.update, mwAuth(api.RoleAdmin))
	e.Add(http.MethodDelete, "/v1/admin/tingkat-pendidikan/:id", h.delete, mwAuth(api.RoleAdmin))
}
