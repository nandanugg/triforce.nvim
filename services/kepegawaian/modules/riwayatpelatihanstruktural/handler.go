package riwayatpelatihanstruktural

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

type listRequest struct {
	Limit  uint `query:"limit"`
	Offset uint `query:"offset"`
}

type listResponse struct {
	Data []riwayatPelatihanStruktural `json:"data"`
	Meta api.MetaPagination           `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.list(c.Request().Context(), api.CurrentUser(c).NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting list pelatihan struktural.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}
