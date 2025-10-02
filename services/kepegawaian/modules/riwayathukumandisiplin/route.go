package riwayathukumandisiplin

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-hukuman-disiplin", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-hukuman-disiplin/:id/berkas", h.getBerkas, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-hukuman-disiplin", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-hukuman-disiplin/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))
}
