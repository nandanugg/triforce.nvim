package resourcepermission

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

type listMyResourcePermissionsResponse struct {
	Data []string `json:"data"`
}

func (h *handler) listMyResourcePermissions(c echo.Context) error {
	data, err := h.svc.listResourcePermissionsByNip(c.Request().Context(), api.CurrentUser(c).NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting list my resource permissions.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listMyResourcePermissionsResponse{
		Data: data,
	})
}

type listResourcesResponse struct {
	Data []resource         `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listResources(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.svc.listResources(c.Request().Context(), req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error getting list resources.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResourcesResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}
