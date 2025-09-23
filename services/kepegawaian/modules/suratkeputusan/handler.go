package suratkeputusan

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
	Data []sk               `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listRequest struct {
	api.PaginationRequest
	StatusSK   *int32 `query:"status_sk"`
	KategoriSK string `query:"kategori_sk"`
	NoSK       string `query:"no_sk"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, listParams{
		Nip:        api.CurrentUser(c).NIP,
		Limit:      req.Limit,
		Offset:     req.Offset,
		StatusSK:   req.StatusSK,
		KategoriSK: req.KategoriSK,
		NoSK:       req.NoSK,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list sk pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getRequest struct {
	ID string `param:"id"`
}

type getResponse struct {
	Data sk `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, api.CurrentUser(c).NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting get sk pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{
		Data: *data,
	})
}

type getBerkasRequest struct {
	ID     string `param:"id"`
	Signed bool   `query:"signed"`
}

func (h *handler) getBerkas(c echo.Context) error {
	var req getBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, api.CurrentUser(c).NIP, req.ID, req.Signed)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas SK.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas SK tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}
