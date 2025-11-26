package usulanperubahandata

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

type ServiceRoute struct {
	h *handler
}

type ServiceRouteInterface interface {
	Register(*echo.Echo, api.AuthMiddlewareFunc, ServiceInterface, string)
}

// RegisterRoutes registers the default routes for usulan perubahan data and returns a CustomRoute.
// Use CustomRoute.Register to add the "create" and "approve" routes for a specific jenis_data hook.
func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, repo repository, mwAuth api.AuthMiddlewareFunc) *ServiceRoute {
	s := newService(db, repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/usulan-perubahan-data/:jenis_data", h.myList, mwAuth(api.Kode_PegawaiPerubahanData_Request))
	e.Add(http.MethodPost, "/v1/usulan-perubahan-data/:jenis_data/:id/read", h.markAsRead, mwAuth(api.Kode_PegawaiPerubahanData_Request))
	e.Add(http.MethodDelete, "/v1/usulan-perubahan-data/:jenis_data/:id", h.delete, mwAuth(api.Kode_PegawaiPerubahanData_Request))

	e.Add(http.MethodGet, "/v1/admin/usulan-perubahan-data", h.adminList, mwAuth(api.Kode_PegawaiPerubahanData_Verify))
	e.Add(http.MethodGet, "/v1/admin/usulan-perubahan-data/:jenis_data/:id", h.adminDetail, mwAuth(api.Kode_PegawaiPerubahanData_Verify))
	e.Add(http.MethodPost, "/v1/admin/usulan-perubahan-data/:jenis_data/:id/reject", h.adminReject, mwAuth(api.Kode_PegawaiPerubahanData_Verify))

	return &ServiceRoute{h}
}

func (r *ServiceRoute) Register(e *echo.Echo, mwAuth api.AuthMiddlewareFunc, svc ServiceInterface, jenisData string) {
	e.Add(http.MethodPost, "/v1/usulan-perubahan-data/"+jenisData, r.h.create(svc, jenisData), mwAuth(api.Kode_PegawaiPerubahanData_Request))

	e.Add(http.MethodPost, "/v1/admin/usulan-perubahan-data/"+jenisData+"/:id/approve", r.h.adminApprove(svc, jenisData), mwAuth(api.Kode_PegawaiPerubahanData_Verify))
}
