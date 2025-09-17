package tingkatpendidikan

import (
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	service *service
}

func newHandler(s *service) *handler {
	return &handler{service: s}
}

type listResponse struct {
	Data []tingkatPendidikan `json:"data"`
}

func (h *handler) list(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.list(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list tingkat pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
	})
}
