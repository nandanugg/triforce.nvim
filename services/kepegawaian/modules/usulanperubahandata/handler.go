package usulanperubahandata

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type handler struct {
	s *service
}

func newHandler(s *service) *handler {
	return &handler{s: s}
}

type myListRequest struct {
	JenisData string `param:"jenis_data"`
	api.PaginationRequest
}

type myListResponse struct {
	Data []usulanPerubahanData `json:"data"`
	Meta api.MetaPagination    `json:"meta"`
}

func (h *handler) myList(c echo.Context) error {
	var req myListRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.s.listUnread(ctx, api.CurrentUser(c).NIP, req.JenisData, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting my list usulan perubahan data.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, myListResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type adminListRequest struct {
	Nama        string `query:"nama"`
	NIP         string `query:"nip"`
	JenisData   string `query:"jenis_data"`
	KodeJabatan string `query:"kode_jabatan"`
	UnitKerjaID string `query:"unit_kerja_id"`
	GolonganID  *int16 `query:"golongan_id"`
	api.PaginationRequest
}

type adminListResponse struct {
	Data []usulanPerubahanData `json:"data"`
	Meta api.MetaPagination    `json:"meta"`
}

func (h *handler) adminList(c echo.Context) error {
	var req adminListRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.s.listPending(ctx, req.Limit, req.Offset, listPendingOptions{
		nama:        req.Nama,
		nip:         req.NIP,
		jenisData:   req.JenisData,
		unitKerjaID: req.UnitKerjaID,
		kodeJabatan: req.KodeJabatan,
		golonganID:  req.GolonganID,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting admin list usulan perubahan data.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, adminListResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type adminDetailRequest struct {
	JenisData string `param:"jenis_data"`
	ID        int64  `param:"id"`
}

type adminDetailResponse struct {
	Data *usulanPerubahanData `json:"data"`
}

func (h *handler) adminDetail(c echo.Context) error {
	var req adminDetailRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.s.detail(ctx, req.JenisData, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting admin detail usulan perubahan data.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}
	return c.JSON(http.StatusOK, adminDetailResponse{
		Data: data,
	})
}

type markAsReadRequest struct {
	ID        int64  `param:"id"`
	JenisData string `param:"jenis_data"`
}

func (h *handler) markAsRead(c echo.Context) error {
	var req markAsReadRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.s.markAsRead(ctx, api.CurrentUser(c).NIP, req.JenisData, req.ID)
	if err != nil {
		if errors.Is(err, errInvalidStateTransition) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "status transisi tidak valid, data tidak dapat ditandai sebagai telah dibaca")
		}

		slog.ErrorContext(ctx, "Error marking usulan perubahan data as read.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}
	return c.NoContent(http.StatusNoContent)
}

type deleteRequest struct {
	ID        int64  `param:"id"`
	JenisData string `param:"jenis_data"`
}

func (h *handler) delete(c echo.Context) error {
	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.s.delete(ctx, api.CurrentUser(c).NIP, req.JenisData, req.ID)
	if err != nil {
		if errors.Is(err, errInvalidStateTransition) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, "status transisi tidak valid, data tidak dapat dihapus")
		}

		slog.ErrorContext(ctx, "Error deleting usulan perubahan data.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}
	return c.NoContent(http.StatusNoContent)
}

type createRequest struct {
	Action string          `json:"action"`
	DataID string          `json:"data_id"`
	Data   json.RawMessage `json:"data"`
}

func (h *handler) create(svc ServiceInterface, jenisData string) func(echo.Context) error {
	return func(c echo.Context) error {
		var req createRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		if err := h.s.create(c.Request().Context(), svc, api.CurrentUser(c).NIP, jenisData, req); err != nil {
			var multiErr *api.MultiError
			if errors.As(err, &multiErr) {
				return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
			}
			if db.IsPgErrorCode(err, db.PgErrUniqueViolation) {
				return echo.NewHTTPError(http.StatusConflict, "data dengan id ini sudah diusulkan")
			}

			slog.ErrorContext(c.Request().Context(), "Error creating usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		return c.NoContent(http.StatusNoContent)
	}
}

type adminRejectRequest struct {
	JenisData string `param:"jenis_data"`
	ID        int64  `param:"id"`
	Catatan   string `json:"catatan"`
}

func (h *handler) adminReject(c echo.Context) error {
	var req adminRejectRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.s.reject(ctx, req.JenisData, req.ID, req.Catatan)
	if err != nil {
		if errors.Is(err, errInvalidStateTransition) {
			return echo.NewHTTPError(http.StatusUnprocessableEntity, `status transisi tidak valid, data tidak dapat ditolak`)
		}

		slog.ErrorContext(ctx, "Error rejecting usulan perubahan data.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}
	return c.NoContent(http.StatusNoContent)
}

type adminApproveRequest struct {
	ID int64 `param:"id"`
}

func (h *handler) adminApprove(svc ServiceInterface, jenisData string) func(echo.Context) error {
	return func(c echo.Context) error {
		var req adminApproveRequest
		if err := c.Bind(&req); err != nil {
			return err
		}

		ctx := c.Request().Context()
		found, err := h.s.approve(ctx, svc, jenisData, req.ID)
		if err != nil {
			if errors.Is(err, errInvalidStateTransition) {
				return echo.NewHTTPError(http.StatusUnprocessableEntity, `status transisi tidak valid, data tidak dapat disetujui`)
			}

			slog.ErrorContext(ctx, "Error approving usulan perubahan data.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError)
		}

		if !found {
			return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
		}
		return c.NoContent(http.StatusNoContent)
	}
}
