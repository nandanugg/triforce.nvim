package jenissatker

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/jenis-satker", h.list, mwAuth(api.Kode_DataMaster_Public))

	e.Add(http.MethodGet, "/v1/admin/jenis-satker", h.list, mwAuth(api.Kode_DataMaster_Read))
	e.Add(http.MethodGet, "/v1/admin/jenis-satker/:id", h.adminGet, mwAuth(api.Kode_DataMaster_Read))

	e.Add(http.MethodPost, "/v1/admin/jenis-satker", h.adminCreate, mwAuth(api.Kode_DataMaster_Write))
	e.Add(http.MethodPut, "/v1/admin/jenis-satker/:id", h.adminUpdate, mwAuth(api.Kode_DataMaster_Write))
	e.Add(http.MethodDelete, "/v1/admin/jenis-satker/:id", h.adminDelete, mwAuth(api.Kode_DataMaster_Write))
}
