package riwayatpenugasan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-penugasan", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-penugasan/:id/berkas", h.getBerkas, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-penugasan", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-penugasan/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))
}
