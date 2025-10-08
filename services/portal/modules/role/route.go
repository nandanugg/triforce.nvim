package role

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

	e.Add(http.MethodGet, "/v1/roles", h.list, mwAuth(api.RoleAdmin))
	e.Add(http.MethodGet, "/v1/roles/:id", h.get, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPost, "/v1/roles", h.create, mwAuth(api.RoleAdmin))
	e.Add(http.MethodPatch, "/v1/roles/:id", h.update, mwAuth(api.RoleAdmin))
}
