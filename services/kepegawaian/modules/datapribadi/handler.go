package datapribadi

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

type getDataPribadiResponse struct {
	Data dataPribadi `json:"data"`
}

func (h *handler) getDataPribadi(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.get(ctx, api.CurrentUser(c).NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data pribadi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getDataPribadiResponse{Data: *data})
}
