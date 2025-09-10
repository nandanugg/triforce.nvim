package jeniskenaikanpangkat

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

type listJenisKPResponse struct {
	Data []jenisKp          `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listJenisKPRequest struct {
	Limit  uint `query:"limit"`
	Offset uint `query:"offset"`
}

func (h *handler) listJenisKP(c echo.Context) error {
	var req listJenisKPRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	data, total, err := h.service.listJenisKP(ctx, listJenisKPParams(req))
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list jenis kenaikan pangkat.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(
		http.StatusOK,
		listJenisKPResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}
