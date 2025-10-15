package template

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repo, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/template", h.list, mwAuth(api.Kode_DataMaster_Public))
	e.Add(http.MethodGet, "/v1/template/:id/berkas", h.getBerkas, mwAuth(api.Kode_DataMaster_Public))

	e.Add(http.MethodGet, "/v1/admin/template", h.list, mwAuth(api.Kode_DataMaster_Read))
	e.Add(http.MethodGet, "/v1/admin/template/:id", h.get, mwAuth(api.Kode_DataMaster_Read))
	e.Add(http.MethodGet, "/v1/admin/template/:id/berkas", h.getBerkas, mwAuth(api.Kode_DataMaster_Read))

	e.Add(http.MethodPost, "/v1/admin/template", h.create, mwAuth(api.Kode_DataMaster_Write))
	e.Add(http.MethodPut, "/v1/admin/template/:id", h.update, mwAuth(api.Kode_DataMaster_Write))
	e.Add(http.MethodDelete, "/v1/admin/template/:id", h.delete, mwAuth(api.Kode_DataMaster_Write))
}
