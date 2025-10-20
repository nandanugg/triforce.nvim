package user

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, repo sqlcRepository, mwAuth api.AuthMiddlewareFunc) {
	r := newRepository(db, repo)
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/users", h.list, mwAuth(api.Kode_ManajemenAkses_Read))
	e.Add(http.MethodGet, "/v1/users/:nip", h.get, mwAuth(api.Kode_ManajemenAkses_Read))

	e.Add(http.MethodPatch, "/v1/users/:nip", h.update, mwAuth(api.Kode_ManajemenAkses_Write))
}
