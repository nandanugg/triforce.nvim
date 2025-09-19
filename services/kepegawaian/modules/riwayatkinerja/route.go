package riwayatkinerja

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, mwAuth api.AuthMiddlewareFunc) {
	r := newRepository(db)
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/riwayat-kinerja", h.list, mwAuth())
}
