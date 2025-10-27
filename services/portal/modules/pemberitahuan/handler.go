package pemberitahuan

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5/pgtype"
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
	Data []pemberitahuan    `json:"data"`
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
		slog.ErrorContext(ctx, "Error getting list pemberitahuan", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type createRequest struct {
	JudulBerita     string    `json:"judul_berita"`
	DeskripsiBerita string    `json:"deskripsi_berita"`
	Pinned          bool      `json:"pinned"`
	DiterbitkanPada time.Time `json:"diterbitkan_pada"`
	DitarikPada     time.Time `json:"ditarik_pada"`
}
type createUpdateResponse struct {
	Data *pemberitahuan `json:"data"`
}

func (h *handler) create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	usr := api.CurrentUser(c)
	ctx := c.Request().Context()

	data, err := h.service.create(ctx, createPemberitahuanParams{
		JudulBerita:      req.JudulBerita,
		DeskripsiBerita:  req.DeskripsiBerita,
		Pinned:           req.Pinned,
		DiterbitkanPada:  pgtype.Timestamptz{Time: req.DiterbitkanPada, Valid: true},
		DitarikPada:      pgtype.Timestamptz{Time: req.DitarikPada, Valid: true},
		DiperbaharuiOleh: usr.NIP,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating pemberitahuan", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusCreated, createUpdateResponse{Data: data})
}

type updateRequest struct {
	ID int64 `param:"id"`
	createRequest
}

func (h *handler) update(c echo.Context) error {
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	usr := api.CurrentUser(c)
	ctx := c.Request().Context()
	data, err := h.service.update(ctx, req.ID, updatePemberitahuanParams{
		JudulBerita:      req.JudulBerita,
		DeskripsiBerita:  req.DeskripsiBerita,
		Pinned:           req.Pinned,
		DiterbitkanPada:  pgtype.Timestamptz{Time: req.DiterbitkanPada, Valid: true},
		DitarikPada:      pgtype.Timestamptz{Time: req.DitarikPada, Valid: true},
		DiperbaharuiOleh: usr.NIP,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error updating pemberitahuan", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}
	return c.JSON(http.StatusOK, createUpdateResponse{Data: data})
}

type deleteRequest struct {
	ID int64 `param:"id"`
}

func (h *handler) delete(c echo.Context) error {
	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	success, err := h.service.delete(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting pemberitahuan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
