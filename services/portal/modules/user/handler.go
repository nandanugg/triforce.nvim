package user

import (
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

type listRequest struct {
	NIP    string `query:"nip"`
	RoleID int16  `query:"role_id"`
	api.PaginationRequest
}

type listResponse struct {
	Data []user             `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	opts := listOptions{
		nip:    req.NIP,
		roleID: req.RoleID,
	}
	data, total, err := h.svc.list(c.Request().Context(), opts, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting list users.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getRequest struct {
	NIP string `param:"nip"`
}

type getResponse struct {
	Data *user `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, err := h.svc.get(c.Request().Context(), req.NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting detail user.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{
		Data: data,
	})
}
