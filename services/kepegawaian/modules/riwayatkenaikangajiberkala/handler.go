package riwayatkenaikangajiberkala

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
	Data []riwayatKenaikanGajiBerkala `json:"data"`
	Meta api.MetaPagination           `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, api.CurrentUser(c).NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat kenaikan gaji berkala.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

type getBerkasRequest struct {
	ID int64 `param:"id"`
}

func (h *handler) getBerkas(c echo.Context) error {
	var req getBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, api.CurrentUser(c).NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas riwayat kenaikan gaji berkala.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat kenaikan gaji berkala tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}
