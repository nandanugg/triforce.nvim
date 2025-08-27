package kinerja

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB, mwAuth echo.MiddlewareFunc) {
	r := newRepository(db)
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/kinerja", h.list, mwAuth)
}
