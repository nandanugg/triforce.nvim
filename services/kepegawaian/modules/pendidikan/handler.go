package pendidikan

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
	Nama string `query:"nama"`
	api.PaginationRequest
}

type listResponse struct {
	Data []pendidikan       `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, req.Limit, req.Offset, req.Nama)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list pendidikan.", "error", err)
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
	Data *pendidikan `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{Data: data})
}

type createRequest struct {
	Nama                string `json:"nama"`
	TingkatPendidikanID int32  `json:"tingkat_pendidikan_id"`
}

type createResponse struct {
	Data *pendidikan `json:"data"`
}

func (h *handler) create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, createParams{
		nama:                req.Nama,
		tingkatPendidikanID: req.TingkatPendidikanID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, createResponse{Data: data})
}

type updateRequest struct {
	ID                  string `param:"id"`
	Nama                string `json:"nama"`
	TingkatPendidikanID int32  `json:"tingkat_pendidikan_id"`
}

type updateResponse struct {
	Data *pendidikan `json:"data"`
}

func (h *handler) update(c echo.Context) error {
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.update(ctx, req.ID, updateParams{
		nama:                req.Nama,
		tingkatPendidikanID: req.TingkatPendidikanID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error updating pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, updateResponse{Data: data})
}

type deleteRequest struct {
	ID string `param:"id"`
}

func (h *handler) delete(c echo.Context) error {
	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	success, err := h.service.delete(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
