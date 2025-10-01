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
	Data keluarga `json:"data"`
}

func (h *handler) list(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.list(ctx, api.CurrentUser(c).NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data keluarga.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{Data: data})
}

type listAdminRequest struct {
	NIP string `param:"nip"`
}

func (h *handler) listAdmin(c echo.Context) error {
	var req listAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	data, err := h.service.list(ctx, req.NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data keluarga.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{Data: data})
}
