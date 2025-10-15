package template

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repo, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/template", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/template/:id/berkas", h.getBerkas, mwAuth())

	e.Add(http.MethodGet, "/v1/admin/template", h.list, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/template/:id", h.get, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/template/:id/berkas", h.getBerkas, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPost, "/v1/admin/template", h.create, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPut, "/v1/admin/template/:id", h.update, mwAuth(api.RoleAdmin))
	e.Add(http.MethodDelete, "/v1/admin/template/:id", h.delete, mwAuth(api.RoleAdmin))
}
