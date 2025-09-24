package pegawai

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo repository) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/pegawai/profil/:pns_id", h.getProfile)
}
