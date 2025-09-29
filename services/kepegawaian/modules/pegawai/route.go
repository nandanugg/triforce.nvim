package pegawai

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/admin/pegawai", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip", h.getAdmin, mwAuth(api.RoleAdmin))

	e.Add(http.MethodGet, "/v1/pegawai/profil/:pns_id", h.getProfile)
}
