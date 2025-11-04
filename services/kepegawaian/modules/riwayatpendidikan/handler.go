package riwayatpendidikan

import (
	"errors"
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
	Data []riwayatPendidikan `json:"data"`
	Meta api.MetaPagination  `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, api.CurrentUser(c).NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type listAdminRequest struct {
	NIP string `param:"nip"`
	api.PaginationRequest
}

func (h *handler) listAdmin(c echo.Context) error {
	var req listAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, req.NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
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
	mimeType, blob, err := h.service.getBerkas(ctx, api.CurrentUser(c).NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas riwayat pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat pendidikan tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type getBerkasAdminRequest struct {
	NIP string `param:"nip"`
	ID  int32  `param:"id"`
}

func (h *handler) getBerkasAdmin(c echo.Context) error {
	var req getBerkasAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, req.NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas riwayat pendidikan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat pendidikan tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type upsertParams struct {
	TingkatPendidikanID int16         `json:"tingkat_pendidikan_id"`
	PendidikanID        *string       `json:"pendidikan_id"`
	NamaSekolah         string        `json:"nama_sekolah"`
	TahunLulus          int16         `json:"tahun_lulus"`
	NomorIjazah         string        `json:"nomor_ijazah"`
	GelarDepan          string        `json:"gelar_depan"`
	GelarBelakang       string        `json:"gelar_belakang"`
	TugasBelajar        statusBelajar `json:"tugas_belajar"`
	NegaraSekolah       string        `json:"negara_sekolah"`
}

type adminCreateRequest struct {
	NIP string `param:"nip"`
	upsertParams
}

type adminCreateResponse struct {
	Data struct {
		ID int32 `json:"id"`
	} `json:"data"`
}

func (h *handler) adminCreate(c echo.Context) error {
	var req adminCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	id, err := h.service.create(c.Request().Context(), req.NIP, req.upsertParams)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "data pegawai tidak ditemukan")
		}

		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(c.Request().Context(), "Error admin creating riwayat pendidikan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var resp adminCreateResponse
	resp.Data.ID = id
	return c.JSON(http.StatusCreated, resp)
}

type adminUpdateRequest struct {
	ID  int32  `param:"id"`
	NIP string `param:"nip"`
	upsertParams
}

func (h *handler) adminUpdate(c echo.Context) error {
	var req adminUpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.update(c.Request().Context(), req.ID, req.NIP, req.upsertParams)
	if err != nil {
		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(c.Request().Context(), "Error admin updating riwayat pendidikan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminDeleteRequest struct {
	ID  int32  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminDelete(c echo.Context) error {
	var req adminDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.delete(c.Request().Context(), req.ID, req.NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin deleting riwayat pendidikan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminUploadBerkasRequest struct {
	ID  int32  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminUploadBerkas(c echo.Context) error {
	var req adminUploadBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	fileBase64, _, err := api.GetFileBase64(c)
	if err != nil {
		return err
	}

	found, err := h.service.uploadBerkas(c.Request().Context(), req.ID, req.NIP, fileBase64)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin uploading berkas riwayat pendidikan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
