package kepangkatan

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
	Data []kepangkatan      `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, listRiwayatParams{
		Limit:  req.Limit,
		Offset: req.Offset,
		PnsNip: api.CurrentUser(c).NIP,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat kepangkatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}
