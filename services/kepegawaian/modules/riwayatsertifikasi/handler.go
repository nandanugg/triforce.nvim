package riwayatsertifikasi

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
	Data []riwayatSertifikasi `json:"data"`
	Meta api.MetaPagination   `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, api.CurrentUser(c).NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getBerkasRequest struct {
	ID int64 `param:"id"`
}

func (h *handler) getBerkas(c echo.Context) error {
	var req getBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, api.CurrentUser(c).NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat sertifikasi tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
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
		slog.ErrorContext(ctx, "Error getting list riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getBerkasAdminRequest struct {
	NIP string `param:"nip"`
	ID  int64  `param:"id"`
}

func (h *handler) getBerkasAdmin(c echo.Context) error {
	var req getBerkasAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, req.NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat sertifikasi tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type upsertParams struct {
	Tahun           int16  `json:"tahun"`
	NamaSertifikasi string `json:"nama_sertifikasi"`
	Deskripsi       string `json:"deskripsi"`
}

type adminCreateRequest struct {
	NIP string `param:"nip"`
	upsertParams
}

type adminCreateResponse struct {
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

func (h *handler) adminCreate(c echo.Context) error {
	var req adminCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, req)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "data pegawai tidak ditemukan")
		}
		slog.ErrorContext(ctx, "Error creating riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, adminCreateResponse{Data: struct {
		ID int64 `json:"id"`
	}{ID: data}})
}

type adminUpdateRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
	upsertParams
}

func (h *handler) adminUpdate(c echo.Context) error {
	var req adminUpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	success, err := h.service.update(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "Error updating riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminDeleteRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminDelete(c echo.Context) error {
	var req adminDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	success, err := h.service.delete(ctx, req.ID, req.NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminUploadBerkasRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminUploadBerkas(c echo.Context) error {
	var req adminUploadBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	fileBase64, _, err := api.GetFileBase64(c)
	if err != nil {
		return err
	}

	success, err := h.service.uploadBerkas(ctx, req.ID, req.NIP, fileBase64)
	if err != nil {
		slog.ErrorContext(ctx, "Error uploading berkas riwayat sertifikasi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
