package riwayatkepangkatan

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, r repository, mwAuth api.AuthMiddlewareFunc) {
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-kepangkatan", h.list, mwAuth(api.Kode_Pegawai_Self))
	e.Add(http.MethodGet, "/v1/riwayat-kepangkatan/:id/berkas", h.getBerkas, mwAuth(api.Kode_Pegawai_Self))

	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-kepangkatan", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/riwayat-kepangkatan/:id/berkas", h.getBerkasAdmin, mwAuth(api.Kode_Pegawai_Read))

	e.Add(http.MethodPost, "/v1/admin/pegawai/:nip/riwayat-kepangkatan", h.adminCreate, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/riwayat-kepangkatan/:id", h.adminUpdate, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodDelete, "/v1/admin/pegawai/:nip/riwayat-kepangkatan/:id", h.adminDelete, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/riwayat-kepangkatan/:id/berkas", h.adminUploadBerkas, mwAuth(api.Kode_Pegawai_Write))
}
