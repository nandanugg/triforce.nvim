package template

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
	Data []template         `json:"data"`
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
		slog.ErrorContext(ctx, "Error getting list template.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

type getRequest struct {
	ID int32 `param:"id"`
}

type getResponse struct {
	Data *template `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting template.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{
		Data: data,
	})
}

type getBerkasRequest struct {
	ID int32 `param:"id"`
}

func (h *handler) getBerkas(c echo.Context) error {
	var req getBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas master template.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if len(blob) == 0 {
		return echo.NewHTTPError(http.StatusNotFound, "berkas master template tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type createRequest struct {
	Nama string `form:"nama"`
}

type createResponse struct {
	Data *template `json:"data"`
}

func (h *handler) create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	str, err := api.GetFileBase64(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, createParams{
		nama: req.Nama,
		file: str,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating template.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, createResponse{Data: data})
}

type updateRequest struct {
	Nama string `form:"nama"`
	ID   int32  `param:"id"`
}

type updateResponse struct {
	Data *template `json:"data"`
}

func (h *handler) update(c echo.Context) error {
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	str, err := api.GetFileBase64(c)
	if err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.update(ctx, req.ID, updateParams{
		nama: req.Nama,
		file: str,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error updating template.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, updateResponse{Data: data})
}

type deleteRequest struct {
	ID int32 `param:"id"`
}

func (h *handler) delete(c echo.Context) error {
	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	success, err := h.service.delete(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting template.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
