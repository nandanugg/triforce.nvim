package jenissatker

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
	Data []jenisSatker      `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listRequest struct {
	api.PaginationRequest
	Nama string `query:"nama"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, listParams{
		limit:  req.Limit,
		offset: req.Offset,
		nama:   req.Nama,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list jenis satker.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(
		http.StatusOK,
		listResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}

type adminGetRequest struct {
	ID int32 `param:"id"`
}

type adminGetResponse struct {
	Data *jenisSatker `json:"data"`
}

func (h *handler) adminGet(c echo.Context) error {
	var req adminGetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting jenis satker.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, adminGetResponse{Data: data})
}

type adminCreateRequest struct {
	Nama string `json:"nama"`
}

type adminCreateResponse struct {
	Data *jenisSatker `json:"data"`
}

func (h *handler) adminCreate(c echo.Context) error {
	var req adminCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, createParams{
		nama: req.Nama,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating jenis satker.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, adminCreateResponse{Data: data})
}

type adminUpdateRequest struct {
	ID   int32  `param:"id"`
	Nama string `json:"nama"`
}

type adminUpdateResponse struct {
	Data *jenisSatker `json:"data"`
}

func (h *handler) adminUpdate(c echo.Context) error {
	var req adminUpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.update(ctx, updateParams{
		id:   req.ID,
		nama: req.Nama,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating jenis satker.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, adminUpdateResponse{Data: data})
}

type adminDeleteRequest struct {
	ID int32 `param:"id"`
}

func (h *handler) adminDelete(c echo.Context) error {
	var req adminDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.service.delete(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting jenis satker.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
