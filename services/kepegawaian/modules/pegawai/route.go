package pegawai

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, db *pgxpool.Pool, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo, db)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/pegawai/profil/:pns_id", h.getProfile)

	e.Add(http.MethodGet, "/v1/admin/pegawai", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodGet, "/v1/admin/pegawai/ppnpn", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodGet, "/v1/admin/pegawai/non_aktif", h.listAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodGet, "/v1/admin/pegawai/:nip", h.getAdmin, mwAuth(api.Kode_Pegawai_Read))
	e.Add(http.MethodPut, "/v1/admin/pegawai/:nip", h.putAdmin, mwAuth(api.Kode_Pegawai_Write))
}
