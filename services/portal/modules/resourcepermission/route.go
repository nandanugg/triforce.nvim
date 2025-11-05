package resourcepermission

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/resource-permissions/me", h.listMyResourcePermissions, mwAuth(api.Kode_Allow))

	e.Add(http.MethodGet, "/v1/resources", h.listResources, mwAuth(api.Kode_ManajemenAkses_Read))
}
