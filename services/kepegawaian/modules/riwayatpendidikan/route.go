package riwayatpendidikan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, sqlc sqlc.Querier, mwAuth api.AuthMiddlewareFunc) {
	s := newService(sqlc)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-pendidikan", h.list, mwAuth(api.Kode_Pegawai_Self))
	e.Add(http.MethodGet, "/v1/riwayat-pendidikan/:id/berkas", h.getBerkas, mwAuth(api.Kode_Pegawai_Self))

	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-pendidikan", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-pendidikan/:id/berkas", h.getBerkasAdmin, mwAuth(api.Kode_Pegawai_Read))
}
