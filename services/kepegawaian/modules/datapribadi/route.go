package datapribadi

import (
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db *pgxpool.Pool, authMw echo.MiddlewareFunc) {
	r := newRepository(db)
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/data-pribadi", h.getDataPribadi, authMw)
	e.Add(http.MethodGet, "/v1/status-pernikahan", h.listStatusPernikahan, authMw)
}
