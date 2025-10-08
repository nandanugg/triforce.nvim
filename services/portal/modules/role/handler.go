package role

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

type handler struct {
	svc *service
}

func newHandler(s *service) *handler {
	return &handler{svc: s}
}

type listResponse struct {
	Data []role             `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.svc.list(c.Request().Context(), req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting list roles.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getRequest struct {
	ID int16 `param:"id"`
}

type getResponse struct {
	Data *role `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, err := h.svc.get(c.Request().Context(), req.ID)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting detail role.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{
		Data: data,
	})
}

type createRequest struct {
	Nama                  string  `json:"nama"`
	Deskripsi             string  `json:"deskripsi"`
	IsDefault             bool    `json:"is_default"`
	ResourcePermissionIDs []int32 `json:"resource_permission_ids"`
}

type createResponse struct {
	Data struct {
		ID int16 `json:"id"`
	} `json:"data"`
}

func (h *handler) create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	id, err := h.svc.create(c.Request().Context(), createParams{
		nama:                  req.Nama,
		deskripsi:             req.Deskripsi,
		isDefault:             req.IsDefault,
		resourcePermissionIDs: req.ResourcePermissionIDs,
	})
	if err != nil {
		if errors.Is(err, errResourcePermissionNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data resource permission tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error creating role.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var resp createResponse
	resp.Data.ID = id
	return c.JSON(http.StatusCreated, resp)
}

type updateRequest struct {
	ID                    int16    `param:"id"`
	Nama                  *string  `json:"nama"`
	Deskripsi             *string  `json:"deskripsi"`
	IsDefault             *bool    `json:"is_default"`
	ResourcePermissionIDs *[]int32 `json:"resource_permission_ids"`
}

func (h *handler) update(c echo.Context) error {
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.svc.update(c.Request().Context(), req.ID, updateOptions{
		nama:                  req.Nama,
		deskripsi:             req.Deskripsi,
		isDefault:             req.IsDefault,
		resourcePermissionIDs: req.ResourcePermissionIDs,
	})
	if err != nil {
		if errors.Is(err, errResourcePermissionNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data resource permission tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error updating role.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
