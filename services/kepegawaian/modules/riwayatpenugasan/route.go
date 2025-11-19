package riwayatpenugasan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-penugasan", h.list, mwAuth(api.Kode_Pegawai_Self))
	e.Add(http.MethodGet, "/v1/riwayat-penugasan/:id/berkas", h.getBerkas, mwAuth(api.Kode_Pegawai_Self))

	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-penugasan", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-penugasan/:id/berkas", h.getBerkasAdmin, mwAuth(api.Kode_Pegawai_Read))

	e.Add(http.MethodPost, "/v1/admin/pegawai/:nip/riwayat-penugasan", h.adminCreate, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/riwayat-penugasan/:id", h.adminUpdate, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodDelete, "/v1/admin/pegawai/:nip/riwayat-penugasan/:id", h.adminDelete, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/riwayat-penugasan/:id/berkas", h.adminUploadBerkas, mwAuth(api.Kode_Pegawai_Write))
}
