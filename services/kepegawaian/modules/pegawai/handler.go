package pegawai

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
	Cari       string `query:"cari"`
	UnitID     string `query:"unit_id"`
	GolonganID int64  `query:"golongan_id"`
	JabatanID  string `query:"jabatan_id"`
	Status     string `query:"status"`

	api.PaginationRequest
}

type listResponse struct {
	Data []pegawai          `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	if api.CurrentUser(c).Role != api.RoleAdmin {
		return echo.NewHTTPError(http.StatusForbidden)
	}

	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, uint64(req.Limit), uint64(req.Offset), listOptions{
		cari:       req.Cari,
		unitID:     req.UnitID,
		golonganID: req.GolonganID,
		jabatanID:  req.JabatanID,
		status:     req.Status,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

func (h *handler) listStatusPegawai(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.listStatusPegawai(ctx)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list status pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, map[string]any{"data": data})
}
