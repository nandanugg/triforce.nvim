package riwayatpelatihansiasn

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	dbRepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, db dbRepo.Querier, mwAuth api.AuthMiddlewareFunc) {
	s := newService(db)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-siasn", h.list, mwAuth())
	e.Add(http.MethodGet, "/v1/riwayat-pelatihan-siasn/:id/berkas", h.getBerkas, mwAuth())
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-pelatihan-siasn", h.listAdmin, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-pelatihan-siasn/:id/berkas", h.getBerkasAdmin, mwAuth(api.RoleAdmin))
}
