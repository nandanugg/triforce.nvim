package riwayatkenaikangajiberkala

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-kenaikan-gaji-berkala", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-kenaikan-gaji-berkala/:id/berkas", h.getBerkas, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-kenaikan-gaji-berkala", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-kenaikan-gaji-berkala/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))
}
