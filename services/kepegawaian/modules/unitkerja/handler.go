package unitkerja

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
	Data []unitKerja        `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listUnitKerjaRequest struct {
	Nama      string `query:"nama"`
	UnorInduk string `query:"unor_induk"`
	api.PaginationRequest
}

func (h *handler) listUnitKerja(c echo.Context) error {
	ctx := c.Request().Context()
	var req listUnitKerjaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listUnitKerja(ctx, listUnitKerjaParams{
		nama:      req.Nama,
		unorInduk: req.UnorInduk,
		limit:     req.Limit,
		offset:    req.Offset,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK,
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
