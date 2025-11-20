package riwayatkenaikangajiberkala

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
	Data []riwayatKenaikanGajiBerkala `json:"data"`
	Meta api.MetaPagination           `json:"meta"`
}

func (h *handler) list(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, api.CurrentUser(c).NIP, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list riwayat kenaikan gaji berkala.", "error", err)
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
		slog.ErrorContext(ctx, "Error getting berkas riwayat kenaikan gaji berkala.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat kenaikan gaji berkala tidak ditemukan")
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
		slog.ErrorContext(ctx, "Error getting list riwayat kenaikan gaji berkala.", "error", err)
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
		slog.ErrorContext(ctx, "Error getting berkas riwayat kenaikan gaji berkala.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas riwayat kenaikan gaji berkala tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type upsertParams struct {
	GolonganID             int32   `json:"golongan_id"`
	TMTGolongan            db.Date `json:"tmt_golongan"`
	MasaKerjaGolonganTahun int16   `json:"masa_kerja_golongan_tahun"`
	MasaKerjaGolonganBulan int16   `json:"masa_kerja_golongan_bulan"`
	NomorSK                string  `json:"nomor_sk"`
	TanggalSK              db.Date `json:"tanggal_sk"`
	GajiPokok              int32   `json:"gaji_pokok"`
	Jabatan                string  `json:"jabatan"`
	TMTJabatan             db.Date `json:"tmt_jabatan"`
	TMTKenaikanGajiBerkala db.Date `json:"tmt_kenaikan_gaji_berkala"`
	Pendidikan             string  `json:"pendidikan"`
	TanggalLulus           db.Date `json:"tanggal_lulus"`
	KantorPembayaran       string  `json:"kantor_pembayaran"`
	Pejabat                string  `json:"pejabat"`
	UnitKerjaIndukID       string  `json:"unit_kerja_induk_id"`
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

func (h *handler) createAdmin(c echo.Context) error {
	var req adminCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	id, err := h.service.create(ctx, req)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, errPegawaiNotFound.Error())
		}

		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(c.Request().Context(), "Error admin creating riwayat kenaikan gaji berkala.", "error", err)
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

func (h *handler) updateAdmin(c echo.Context) error {
	var req adminUpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.service.update(ctx, req)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, errPegawaiNotFound.Error())
		}

		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(c.Request().Context(), "Error admin updating riwayat kenaikan gaji berkala.", "error", err)
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

func (h *handler) deleteAdmin(c echo.Context) error {
	var req adminDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.service.delete(ctx, req.NIP, req.ID)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, errPegawaiNotFound.Error())
		}
		slog.ErrorContext(c.Request().Context(), "Error admin deleting riwayat kenaikan gaji berkala.", "error", err)
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
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, errPegawaiNotFound.Error())
		}
		slog.ErrorContext(c.Request().Context(), "Error admin uploading berkas riwayat kenaikan gaji berkala.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
