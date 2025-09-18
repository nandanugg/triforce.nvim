package golongan

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

type listRefGolonganResponse struct {
	Data []refGolongan      `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listRefGolongan(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	data, total, err := h.service.listRefGolongan(ctx, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list golongan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK,
		listRefGolonganResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}
