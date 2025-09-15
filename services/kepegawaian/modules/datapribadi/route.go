package datapribadi

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

func RegisterRoutes(e *echo.Echo, repo repository, authMw echo.MiddlewareFunc) {
	s := newService(repo)
	h := newHandler(s)

	e.Add(http.MethodGet, "/v1/data-pribadi", h.getDataPribadi, authMw)
}
