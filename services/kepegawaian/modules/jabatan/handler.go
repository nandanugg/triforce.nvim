package jabatan

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

type listRequest struct {
	Nama string `query:"nama"`
	api.PaginationRequest
}

type listResponse struct {
	Data []jabatanPublic    `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listJabatan(c echo.Context) error {
	var req listRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, total, err := h.service.listJabatan(ctx, req.Nama, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

type listAdminRequest struct {
	Keyword string `query:"keyword"`
	api.PaginationRequest
}

type listAdminResponse struct {
	Data []jabatan          `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listAdmin(c echo.Context) error {
	var req listAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()

	data, total, err := h.service.listAdmin(ctx, req.Keyword, req.Limit, req.Offset)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting list jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listAdminResponse{
		Data: data,
		Meta: api.MetaPagination{Limit: req.Limit, Offset: req.Offset, Total: uint(total)},
	})
}

type getRequest struct {
	ID int32 `param:"id"`
}

type getResponse struct {
	Data *jabatan `json:"data"`
}

func (h *handler) get(c echo.Context) error {
	var req getRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getResponse{Data: data})
}

type createRequest struct {
	KodeJabatan      string `json:"kode_jabatan"`
	NamaJabatan      string `json:"nama_jabatan"`
	NamaJabatanFull  string `json:"nama_jabatan_full"`
	JenisJabatan     *int16 `json:"jenis_jabatan"`
	Kelas            *int16 `json:"kelas"`
	Pensiun          *int16 `json:"pensiun"`
	KodeBkn          string `json:"kode_bkn"`
	NamaJabatanBkn   string `json:"nama_jabatan_bkn"`
	KategoriJabatan  string `json:"kategori_jabatan"`
	BknID            string `json:"bkn_id"`
	TunjanganJabatan int64  `json:"tunjangan_jabatan"`
}

type createResponse struct {
	Data *jabatan `json:"data"`
}

func (h *handler) create(c echo.Context) error {
	var req createRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, createParams{
		kodeJabatan:      req.KodeJabatan,
		namaJabatan:      req.NamaJabatan,
		namaJabatanFull:  req.NamaJabatanFull,
		jenisJabatan:     req.JenisJabatan,
		kelas:            req.Kelas,
		pensiun:          req.Pensiun,
		kodeBkn:          req.KodeBkn,
		namaJabatanBkn:   req.NamaJabatanBkn,
		kategoriJabatan:  req.KategoriJabatan,
		bknID:            req.BknID,
		tunjanganJabatan: req.TunjanganJabatan,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, createResponse{Data: data})
}

type updateRequest struct {
	ID               int32  `param:"id"`
	KodeJabatan      string `json:"kode_jabatan"`
	NamaJabatan      string `json:"nama_jabatan"`
	NamaJabatanFull  string `json:"nama_jabatan_full"`
	JenisJabatan     *int16 `json:"jenis_jabatan"`
	Kelas            *int16 `json:"kelas"`
	Pensiun          *int16 `json:"pensiun"`
	KodeBkn          string `json:"kode_bkn"`
	NamaJabatanBkn   string `json:"nama_jabatan_bkn"`
	KategoriJabatan  string `json:"kategori_jabatan"`
	BknID            string `json:"bkn_id"`
	TunjanganJabatan int64  `json:"tunjangan_jabatan"`
}

type updateResponse struct {
	Data *jabatan `json:"data"`
}

func (h *handler) update(c echo.Context) error {
	var req updateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.update(ctx, req.ID, updateParams{
		kodeJabatan:      req.KodeJabatan,
		namaJabatan:      req.NamaJabatan,
		namaJabatanFull:  req.NamaJabatanFull,
		jenisJabatan:     req.JenisJabatan,
		kelas:            req.Kelas,
		pensiun:          req.Pensiun,
		kodeBkn:          req.KodeBkn,
		namaJabatanBkn:   req.NamaJabatanBkn,
		kategoriJabatan:  req.KategoriJabatan,
		bknID:            req.BknID,
		tunjanganJabatan: req.TunjanganJabatan,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error updating jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, updateResponse{Data: data})
}

type deleteRequest struct {
	ID int32 `param:"id"`
}

func (h *handler) delete(c echo.Context) error {
	var req deleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	success, err := h.service.delete(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting jabatan.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}
	if !success {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
