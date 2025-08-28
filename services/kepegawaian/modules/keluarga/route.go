package keluarga

import (
	"database/sql"
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, db *sql.DB, mwAuth echo.MiddlewareFunc) {
	r := newRepository(db)
	s := newService(r)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/keluarga", h.list, mwAuth)
	e.Add(http.MethodGet, "/v1/keluarga/anak", h.listAnak, mwAuth)
	e.Add(http.MethodGet, "/v1/keluarga/orang-tua", h.listOrangTua, mwAuth)
	e.Add(http.MethodGet, "/v1/keluarga/pasangan", h.listPasangan, mwAuth)
}
