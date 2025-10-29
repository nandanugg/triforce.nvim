package pemberitahuan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/admin/pemberitahuan", h.list, mwAuth(api.Kode_Pemberitahuan_Read))
	e.Add(http.MethodPost, "/v1/admin/pemberitahuan", h.create, mwAuth(api.Kode_Pemberitahuan_Write))
	e.Add(http.MethodPut, "/v1/admin/pemberitahuan/:id", h.update, mwAuth(api.Kode_Pemberitahuan_Write))
	e.Add(http.MethodDelete, "/v1/admin/pemberitahuan/:id", h.delete, mwAuth(api.Kode_Pemberitahuan_Write))

	e.Add(http.MethodGet, "/v1/pemberitahuan", h.listPublic, mwAuth(api.Kode_Pemberitahuan_Public))
}
