package keluarga

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
	Data []keluarga `json:"data"`
}

func (h *handler) list(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.list(ctx, api.CurrentUser(c).ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list keluarga.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{Data: data})
}

func (h *handler) listOrangTua(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.listOrangTua(ctx, api.CurrentUser(c).ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list orang tua.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{Data: data})
}
