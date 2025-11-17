package riwayatjabatan

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type handler struct {
	service *service
}

func newHandler(s *service) *handler {
	return &handler{service: s}
}

type listResponse struct {
	Data []riwayatJabatan   `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	ctx := c.Request().Context()
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	data, total, err := h.service.list(ctx, api.CurrentUser(c).NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
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
		slog.ErrorContext(ctx, "Error getting berkas riwayat jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat jabatan tidak ditemukan")
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
		slog.ErrorContext(ctx, "Error getting list riwayat jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
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
		slog.ErrorContext(ctx, "Error getting berkas riwayat jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat jabatan tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type upsertParams struct {
	JenisJabatanID          *int32  `json:"jenis_jabatan_id"`
	JabatanID               string  `json:"id_jabatan"`
	SatuanKerjaID           string  `json:"satuan_kerja_id"`
	UnitOrganisasiID        string  `json:"unit_organisasi_id"`
	TMTJabatan              db.Date `json:"tmt_jabatan"`
	NoSK                    string  `json:"no_sk"`
	TanggalSK               db.Date `json:"tanggal_sk"`
	StatusPlt               *bool   `json:"status_plt"`
	PeriodeJabatanStartDate db.Date `json:"periode_jabatan_start_date"`
	PeriodeJabatanEndDate   db.Date `json:"periode_jabatan_end_date"`
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

	id, err := h.service.create(c.Request().Context(), req.NIP, req.upsertParams)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "data pegawai tidak ditemukan")
		}

		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(c.Request().Context(), "Error admin creating riwayat jabatan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var resp adminCreateResponse
	resp.Data.ID = id
	return c.JSON(http.StatusCreated, resp)
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

	found, err := h.service.update(c.Request().Context(), req.ID, req.NIP, req.upsertParams)
	if err != nil {
		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(c.Request().Context(), "Error admin updating riwayat jabatan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
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

	found, err := h.service.delete(c.Request().Context(), req.ID, req.NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin deleting riwayat jabatan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
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

	fileBase64, _, err := api.GetFileBase64(c)
	if err != nil {
		return err
	}

	found, err := h.service.uploadBerkas(c.Request().Context(), req.ID, req.NIP, fileBase64)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin uploading berkas riwayat jabatan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}
	return c.NoContent(http.StatusNoContent)
}
