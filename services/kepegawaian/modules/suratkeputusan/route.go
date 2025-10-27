package suratkeputusan

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, repo repository, db *pgxpool.Pool, mwAuth api.AuthMiddlewareFunc) {
	s := newService(repo, db)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/surat-keputusan", h.list, mwAuth(api.Kode_SuratKeputusan_Self))
	e.Add(http.MethodGet, "/v1/surat-keputusan/:id", h.get, mwAuth(api.Kode_SuratKeputusan_Self))
	e.Add(http.MethodGet, "/v1/surat-keputusan/:id/berkas", h.getBerkas, mwAuth(api.Kode_SuratKeputusan_Self))

	e.Add(http.MethodGet, "/v1/admin/surat-keputusan", h.listAdmin, mwAuth(api.Kode_SuratKeputusan_Read))
	e.Add(http.MethodGet, "/v1/admin/surat-keputusan/:id", h.getAdmin, mwAuth(api.Kode_SuratKeputusan_Read))
	e.Add(http.MethodGet, "/v1/admin/surat-keputusan/:id/berkas", h.getBerkasAdmin, mwAuth(api.Kode_SuratKeputusan_Read))

	e.Add(http.MethodGet, "/v1/koreksi-surat-keputusan", h.listKoreksi, mwAuth(api.Kode_SuratKeputusanApproval_Review))
	e.Add(http.MethodGet, "/v1/koreksi-surat-keputusan/:id", h.getDetailSuratKeputusan, mwAuth(api.Kode_SuratKeputusanApproval_Review))
	e.Add(http.MethodGet, "/v1/koreksi-surat-keputusan/:id/berkas", h.getBerkasAdmin, mwAuth(api.Kode_SuratKeputusanApproval_Review))
	e.Add(http.MethodGet, "/v1/koreksi-surat-keputusan/antrian", h.listKoreksiAntrian, mwAuth(api.Kode_SuratKeputusanApproval_Review))
	e.Add(http.MethodPost, "/v1/koreksi-surat-keputusan/:id", h.koreksiSuratKeputusan, mwAuth(api.Kode_SuratKeputusanApproval_Review))
}
