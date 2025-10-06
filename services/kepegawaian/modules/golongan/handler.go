package golongan

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

type listPublicResponse struct {
	Data []golonganPublic   `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listPublic(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	data, total, err := h.service.listPublic(ctx, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list golongan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listPublicResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

type listResponse struct {
	Data []golongan         `json:"data"`
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
		slog.ErrorContext(ctx, "Error getting list golongan.", "error", err)
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
	Data *golongan `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting golongan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{Data: data})
}

type createRequest struct {
	Nama        string `json:"nama"`
	NamaPangkat string `json:"nama_pangkat"`
	Nama2       string `json:"nama_2"`
	Gol         int16  `json:"gol"`
	GolPppk     string `json:"gol_pppk"`
}

type createResponse struct {
	Data *golongan `json:"data"`
}

func (h *handler) create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, createParams{
		nama:        req.Nama,
		namaPangkat: req.NamaPangkat,
		nama2:       req.Nama2,
		gol:         req.Gol,
		golPppk:     req.GolPppk,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating golongan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, createResponse{Data: data})
}

type updateRequest struct {
	ID          int32  `param:"id"`
	Nama        string `json:"nama"`
	NamaPangkat string `json:"nama_pangkat"`
	Nama2       string `json:"nama_2"`
	Gol         int16  `json:"gol"`
	GolPppk     string `json:"gol_pppk"`
}

type updateResponse struct {
	Data *golongan `json:"data"`
}

func (h *handler) update(c echo.Context) error {
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.update(ctx, req.ID, updateParams{
		nama:        req.Nama,
		namaPangkat: req.NamaPangkat,
		nama2:       req.Nama2,
		gol:         req.Gol,
		golPppk:     req.GolPppk,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error updating golongan.", "error", err)
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
		slog.ErrorContext(ctx, "Error deleting golongan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
