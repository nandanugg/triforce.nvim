package unitkerja

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodPost, "/v1/admin/unit-kerja", h.adminCreate, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/unit-kerja/:id", h.adminGet, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPut, "/v1/admin/unit-kerja/:id", h.adminUpdate, mwAuth(api.RoleAdmin))
	e.Add(http.MethodDelete, "/v1/admin/unit-kerja/:id", h.adminDelete, mwAuth(api.RoleAdmin))

	e.Add(http.MethodGet, "/v1/unit-kerja", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/unit-kerja/akar", h.listAkar, mwAuth())
	e.Add(http.MethodGet, "/v1/unit-kerja/:id/anak", h.listAnak, mwAuth())
}
