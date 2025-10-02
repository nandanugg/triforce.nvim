package riwayatpenghargaan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, r repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-penghargaan", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-penghargaan/:id/berkas", h.getBerkas, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-penghargaan", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-penghargaan/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))
}
