package agama

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/agama", h.list, mwAuth())

	e.Add(http.MethodGet, "/v1/admin/agama", h.list, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPost, "/v1/admin/agama", h.create, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/agama/:id", h.get, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPut, "/v1/admin/agama/:id", h.update, mwAuth(api.RoleAdmin))
	e.Add(http.MethodDelete, "/v1/admin/agama/:id", h.delete, mwAuth(api.RoleAdmin))
}
