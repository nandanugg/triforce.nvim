package jenisjabatan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/jenis-jabatan", h.list, mwAuth())

	e.Add(http.MethodGet, "/v1/admin/jenis-jabatan", h.list, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPost, "/v1/admin/jenis-jabatan", h.adminCreate, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/jenis-jabatan/:id", h.adminGet, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPut, "/v1/admin/jenis-jabatan/:id", h.adminUpdate, mwAuth(api.RoleAdmin))
	e.Add(http.MethodDelete, "/v1/admin/jenis-jabatan/:id", h.adminDelete, mwAuth(api.RoleAdmin))
}
