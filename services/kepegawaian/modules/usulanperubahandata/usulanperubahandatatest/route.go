package usulanperubahandatatest

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
)

type ServiceRoute struct {
	h *handler
}

func NewServiceRoute(db *pgxpool.Pool) *ServiceRoute {
	return &ServiceRoute{
		h: &handler{db, sqlc.New(db)},
	}
}

func (r *ServiceRoute) Register(e *echo.Echo, mwAuth api.AuthMiddlewareFunc, svc usulanperubahandata.ServiceInterface, jenisData string) {
	usulanperubahandata.RegisterRoutes(e, r.h.db, r.h.repo, mwAuth)

	e.Add(http.MethodPost, "/v1/usulan-perubahan-data/"+jenisData, r.h.createThenApprove(svc, jenisData), mwAuth(api.Kode_PegawaiPerubahanData_Request))
}
