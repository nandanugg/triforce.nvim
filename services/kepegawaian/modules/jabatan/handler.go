package jabatan

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
	Data []jabatan          `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listRiwayatJabatanResponse struct {
	Data []riwayatJabatan   `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listJabatan(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.listJabatan(ctx, listParams(req))
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

func (h *handler) listRiwayatJabatan(c echo.Context) error {
	ctx := c.Request().Context()
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listRiwayatJabatan(ctx, listRiwayatJabatanParams{
		Limit:  req.Limit,
		Offset: req.Offset,
		NIP:    "41",
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listRiwayatJabatanResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}
