package riwayatpelatihanstruktural

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, repo sqlc.Querier, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-struktural", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-struktural/:id/berkas", h.getBerkas, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-pelatihan-struktural", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-pelatihan-struktural/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))
}
