package suratkeputusan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/admin/surat-keputusan", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/surat-keputusan/:id", h.getAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/surat-keputusan/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))

	e.Add(http.MethodGet, "/v1/surat-keputusan", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/surat-keputusan/:id", h.get, mwAuth())
	e.Add(http.MethodGet, "/v1/surat-keputusan/:id/berkas", h.getBerkas, mwAuth())
}
