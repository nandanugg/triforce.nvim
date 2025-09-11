package pelatihanstruktural

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

type handler struct {
	service *service
}

func newHandler(s *service) *handler {
	return &handler{service: s}
}

type listResponse struct {
	Data []pelatihanStruktural `json:"data"`
}

func (h *handler) list(c echo.Context) error {
	data, err := h.service.list(c.Request().Context(), api.CurrentUser(c).NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting list pelatihan struktural.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
	})
}
