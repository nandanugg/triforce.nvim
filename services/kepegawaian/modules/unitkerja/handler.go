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

type listResponse struct {
	Data []unitKerjaPublic  `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listRequest struct {
	Nama      string `query:"nama"`
	UnorInduk string `query:"unor_induk"`
	api.PaginationRequest
}

func (h *handler) list(c echo.Context) error {
	ctx := c.Request().Context()
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.list(ctx, listParams{
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

type listAkarResponse struct {
	Data []unitKerjaPublic  `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listAkarRequest struct {
	api.PaginationRequest
}

func (h *handler) listAkar(c echo.Context) error {
	ctx := c.Request().Context()
	var req listAkarRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listAkar(ctx, listAkarParams{
		limit:  req.Limit,
		offset: req.Offset,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list akar unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK,
		listAkarResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}

type listAnakResponse struct {
	Data []anakUnitKerja    `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listAnakRequest struct {
	ID string `param:"id"`
	api.PaginationRequest
}

func (h *handler) listAnak(c echo.Context) error {
	ctx := c.Request().Context()
	var req listAnakRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.listAnak(ctx, listAnakParams{
		limit:  req.Limit,
		offset: req.Offset,
		id:     req.ID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list anak unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	return c.JSON(http.StatusOK,
		listAnakResponse{
			Data: data,
			Meta: api.MetaPagination{
				Limit:  req.Limit,
				Offset: req.Offset,
				Total:  uint(total),
			},
		},
	)
}

// type getRequest struct {
// 	ID int32 `param:"id"`
// }

// type getResponse struct {
// 	Data *unitKerja `json:"data"`
// }

// func (h *handler) get(c echo.Context) error {
// 	var req getRequest
// 	if err := c.Bind(&req); err != nil {
// 		return err
// 	}

// 	ctx := c.Request().Context()
// 	data, err := h.service.get(ctx, req.ID)
// 	if err != nil {
// 		slog.ErrorContext(ctx, "Error getting golongan.", "error", err)
// 		return echo.NewHTTPError(http.StatusInternalServerError)
// 	}

// 	if data == nil {
// 		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
// 	}

// 	return c.JSON(http.StatusOK, getResponse{Data: data})
// }
