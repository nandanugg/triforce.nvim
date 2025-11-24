package jenisjabatan

import (
	"errors"
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
	Data []jenisJabatan     `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list jenis jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

type adminGetRequest struct {
	ID int32 `param:"id"`
}

type adminGetResponse struct {
	Data *jenisJabatan `json:"data"`
}

func (h *handler) adminGet(c echo.Context) error {
	var req adminGetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting jenis jabatan.", "error", err)
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
	Data *jenisJabatan `json:"data"`
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
		slog.ErrorContext(ctx, "Error creating jenis jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, adminCreateResponse{Data: data})
}

type adminUpdateRequest struct {
	ID   int32  `param:"id"`
	Nama string `json:"nama"`
}

type adminUpdateResponse struct {
	Data *jenisJabatan `json:"data"`
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
		slog.ErrorContext(ctx, "Error creating jenis jabatan.", "error", err)
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
		if errors.Is(err, errJenisJabatanReferenced) {
			return echo.NewHTTPError(http.StatusBadRequest, errJenisJabatanReferenced.Error())
		}
		slog.ErrorContext(ctx, "Error deleting jenis jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
