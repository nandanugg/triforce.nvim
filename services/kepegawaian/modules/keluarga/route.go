package keluarga

import (
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	dbRepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

func RegisterRoutes(e *echo.Echo, sqlc dbRepo.Querier, mwAuth api.AuthMiddlewareFunc) {
	s := newService(sqlc)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/keluarga", h.list, mwAuth(api.Kode_Pegawai_Self))

	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip/keluarga", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))

	e.Add(http.MethodPost, "/v1/admin/pegawai/:nip/orang-tua", h.adminCreateOrangTua, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/orang-tua/:id", h.adminUpdateOrangTua, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodDelete, "/v1/admin/pegawai/:nip/orang-tua/:id", h.adminDeleteOrangTua, mwAuth(api.Kode_Pegawai_Write))

	e.Add(http.MethodPost, "/v1/admin/pegawai/:nip/pasangan", h.adminCreatePasangan, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/pasangan/:id", h.adminUpdatePasangan, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodDelete, "/v1/admin/pegawai/:nip/pasangan/:id", h.adminDeletePasangan, mwAuth(api.Kode_Pegawai_Write))

	e.Add(http.MethodPost, "/v1/admin/pegawai/:nip/anak", h.adminCreateAnak, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip/anak/:id", h.adminUpdateAnak, mwAuth(api.Kode_Pegawai_Write))
	e.Add(http.MethodDelete, "/v1/admin/pegawai/:nip/anak/:id", h.adminDeleteAnak, mwAuth(api.Kode_Pegawai_Write))
}
