package suratkeputusan

import (
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
	Data []suratKeputusan   `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listRequest struct {
	api.PaginationRequest
	StatusSK   []int32 `query:"status_sk"`
	KategoriSK string  `query:"kategori_sk"`
	NoSK       string  `query:"no_sk"`
}

func (h *handler) list(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.list(ctx, listParams{
		nip:          api.CurrentUser(c).NIP,
		limit:        req.Limit,
		offset:       req.Offset,
		listStatusSK: req.StatusSK,
		kategoriSK:   req.KategoriSK,
		noSK:         req.NoSK,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list sk pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getRequest struct {
	ID string `param:"id"`
}

type getResponse struct {
	Data *suratKeputusan `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, api.CurrentUser(c).NIP, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting get sk pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{
		Data: data,
	})
}

type getBerkasRequest struct {
	ID     string `param:"id"`
	Signed bool   `query:"signed"`
}

func (h *handler) getBerkas(c echo.Context) error {
	var req getBerkasRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkas(ctx, api.CurrentUser(c).NIP, req.ID, req.Signed)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas SK.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas SK tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type listAdminResponse struct {
	Data []suratKeputusan   `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

type listAdminRequest struct {
	api.PaginationRequest
	UnitKerjaID    string  `query:"unit_kerja_id"`
	NamaPemilik    string  `query:"nama_pemilik"`
	NIP            string  `query:"nip"`
	GolonganID     int32   `query:"golongan_id"`
	JabatanID      string  `query:"jabatan_id"`
	KategoriSK     string  `query:"kategori_sk"`
	TanggalSKMulai db.Date `query:"tanggal_sk_mulai"`
	TanggalSKAkhir db.Date `query:"tanggal_sk_akhir"`
	StatusSK       []int32 `query:"status_sk"`
}

func (h *handler) listAdmin(c echo.Context) error {
	var req listAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.listAdmin(ctx, listAdminParams{
		Limit:          req.Limit,
		Offset:         req.Offset,
		UnitKerjaID:    req.UnitKerjaID,
		NamaPemilik:    req.NamaPemilik,
		NIP:            req.NIP,
		GolonganID:     req.GolonganID,
		JabatanID:      req.JabatanID,
		KategoriSK:     req.KategoriSK,
		TanggalSKMulai: req.TanggalSKMulai,
		TanggalSKAkhir: req.TanggalSKAkhir,
		ListStatusSK:   req.StatusSK,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list sk pegawai by admin.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listAdminResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getAdminRequest struct {
	ID string `param:"id"`
}

type getAdminResponse struct {
	Data *suratKeputusan `json:"data"`
}

func (h *handler) getAdmin(c echo.Context) error {
	var req getAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.getAdmin(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data sk pegawai by admin.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getAdminResponse{
		Data: data,
	})
}

type getBerkasAdminRequest struct {
	ID     string `param:"id"`
	Signed bool   `query:"signed"`
}

func (h *handler) getBerkasAdmin(c echo.Context) error {
	var req getBerkasAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	mimeType, blob, err := h.service.getBerkasAdmin(ctx, req.ID, req.Signed)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting berkas SK by admin.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if blob == nil {
		return echo.NewHTTPError(http.StatusNotFound, "berkas SK tidak ditemukan")
	}

	c.Response().Header().Set("Content-Disposition", "inline")
	return c.Blob(http.StatusOK, mimeType, blob)
}

type listKoreksiRequest struct {
	api.PaginationRequest
	UnitKerjaID   string `query:"unit_kerja_id"`
	NamaPemilik   string `query:"nama_pemilik"`
	NipPemilik    string `query:"nip_pemilik"`
	GolonganID    int32  `query:"golongan_id"`
	JabatanID     string `query:"jabatan_id"`
	KategoriSK    string `query:"kategori_sk"`
	NoSK          string `query:"no_sk"`
	StatusKoreksi string `query:"status"`
}

type listKoreksiResponse struct {
	Data []koreksiSuratKeputusan `json:"data"`
	Meta api.MetaPagination      `json:"meta"`
}

func (h *handler) listKoreksi(c echo.Context) error {
	var req listKoreksiRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.listKoreksi(ctx, listKoreksiParams{
		limit:       req.Limit,
		offset:      req.Offset,
		unitKerjaID: req.UnitKerjaID,
		namaPemilik: req.NamaPemilik,
		nip:         api.CurrentUser(c).NIP,
		nipPemilik:  req.NipPemilik,
		golonganID:  req.GolonganID,
		jabatanID:   req.JabatanID,
		kategoriSK:  req.KategoriSK,
		noSK:        req.NoSK,
		status:      req.StatusKoreksi,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list koreksi surat keputusan belum dikoreksi.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listKoreksiResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type listKoreksiAntrianResponse struct {
	Data []antrianKoreksiSuratKeputusan `json:"data"`
	Meta api.MetaPagination             `json:"meta"`
}

func (h *handler) listKoreksiAntrian(c echo.Context) error {
	var req api.PaginationRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.listKoreksiAntrian(ctx, listKoreksiAntrianParams{
		limit:  req.Limit,
		offset: req.Offset,
		nip:    api.CurrentUser(c).NIP,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list koreksi surat keputusan antrian.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listKoreksiAntrianResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: total},
	})
}

type getDetailSuratKeputusanRequest struct {
	ID string `param:"id"`
}

type getDetailSuratKeputusanResponse struct {
	Data *koreksiSuratKeputusan `json:"data"`
}

func (h *handler) getDetailSuratKeputusan(c echo.Context) error {
	var req getDetailSuratKeputusanRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.getDetailSuratKeputusan(ctx, req.ID, api.CurrentUser(c).NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting detail surat keputusan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getDetailSuratKeputusanResponse{
		Data: data,
	})
}

type koreksiSuratKeputusanRequest struct {
	ID             string `param:"id"`
	StatusKoreksi  string `json:"status_koreksi"`
	CatatanKoreksi string `json:"catatan_koreksi"`
}

func (h *handler) koreksiSuratKeputusan(c echo.Context) error {
	var req koreksiSuratKeputusanRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	success, err := h.service.koreksiSuratKeputusan(ctx, koreksiSuratKeputusanParams{
		id:             req.ID,
		statusKoreksi:  req.StatusKoreksi,
		catatanKoreksi: req.CatatanKoreksi,
		nip:            api.CurrentUser(c).NIP,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error koreksi surat keputusan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !success {
		return echo.NewHTTPError(http.StatusBadRequest, "Surat keputusan tidak dapat dikoreksi")
	}

	return c.NoContent(http.StatusNoContent)
}
