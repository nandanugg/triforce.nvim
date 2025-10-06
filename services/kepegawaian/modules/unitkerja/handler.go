package unitkerja

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

type listUnitKerjaResponse struct {
	Data []unitKerja        `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listUnitKerjaRequest struct {
	Nama      string `query:"nama"`
	UnorInduk string `query:"unor_induk"`
	api.PaginationRequest
}

func (h *handler) listUnitKerja(c echo.Context) error {
	ctx := c.Request().Context()
	var req listUnitKerjaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listUnitKerja(ctx, listUnitKerjaParams{
		nama:      req.Nama,
		unorInduk: req.UnorInduk,
		limit:     req.Limit,
		offset:    req.Offset,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK,
		listUnitKerjaResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}

type listAkarUnitKerjaResponse struct {
	Data []unitKerja        `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listAkarUnitKerjaRequest struct {
	api.PaginationRequest
}

func (h *handler) listAkarUnitKerja(c echo.Context) error {
	ctx := c.Request().Context()
	var req listAkarUnitKerjaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listAkarUnitKerja(ctx, listAkarUnitKerjaParams{
		limit:  req.Limit,
		offset: req.Offset,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list akar unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK,
		listAkarUnitKerjaResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}

type listAnakUnitKerjaResponse struct {
	Data []anakUnitKerja    `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listAnakUnitKerjaRequest struct {
	ID string `param:"id"`
	api.PaginationRequest
}

func (h *handler) listAnakUnitKerja(c echo.Context) error {
	ctx := c.Request().Context()
	var req listAnakUnitKerjaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listAnakUnitKerja(ctx, listAnakUnitKerjaParams{
		limit:  req.Limit,
		offset: req.Offset,
		id:     req.ID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list anak unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK,
		listAnakUnitKerjaResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}
