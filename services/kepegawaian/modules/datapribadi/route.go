package datapribadi

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB, authMw echo.MiddlewareFunc) {
	r := newRepository(db)
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/data-pribadi", h.getDataPribadi, authMw)
}
