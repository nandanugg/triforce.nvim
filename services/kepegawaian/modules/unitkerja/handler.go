package unitkerja

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
	Data []anakUnitKerja    `json:"data"`
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

type adminGetRequest struct {
	ID string `param:"id"`
}

type adminGetResponse struct {
	Data *unitKerjaWithInduk `json:"data"`
}

func (h *handler) adminGet(c echo.Context) error {
	var req adminGetRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.get(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, adminGetResponse{Data: data})
}

type adminCreateRequest struct {
	DiatasanID    string  `json:"diatasan_id"`
	ID            string  `json:"id"`
	Nama          string  `json:"nama"`
	KodeInternal  string  `json:"kode_internal"`
	NamaJabatan   string  `json:"nama_jabatan"`
	PemimpinPNSID string  `json:"pemimpin_pns_id"`
	IsSatker      bool    `json:"is_satker"`
	UnorInduk     string  `json:"unor_induk"`
	ExpiredDate   db.Date `json:"expired_date"`
	Keterangan    string  `json:"keterangan"`
	Abbreviation  string  `json:"abbreviation"`
	Waktu         string  `json:"waktu"`
	JenisSatker   string  `json:"jenis_satker"`
	Peraturan     string  `json:"peraturan"`
}

type adminCreateResponse struct {
	Data *unitKerja `json:"data"`
}

func (h *handler) adminCreate(c echo.Context) error {
	var req adminCreateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.create(ctx, createParams{
		diatasanID:    req.DiatasanID,
		id:            req.ID,
		nama:          req.Nama,
		kodeInternal:  req.KodeInternal,
		namaJabatan:   req.NamaJabatan,
		pemimpinPNSID: req.PemimpinPNSID,
		isSatker:      req.IsSatker,
		unorInduk:     req.UnorInduk,
		expiredDate:   req.ExpiredDate,
		keterangan:    req.Keterangan,
		abbreviation:  req.Abbreviation,
		waktu:         req.Waktu,
		jenisSatker:   req.JenisSatker,
		peraturan:     req.Peraturan,
	})
	if err != nil {
		if db.IsPgErrorCode(err, db.PgErrUniqueViolation) {
			return echo.NewHTTPError(http.StatusConflict, "Data dengan ID ini sudah terdaftar")
		}
		slog.ErrorContext(ctx, "Error creating unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusCreated, adminCreateResponse{Data: data})
}

type adminUpdateRequest struct {
	ID            string  `param:"id"`
	DiatasanID    string  `json:"diatasan_id"`
	Nama          string  `json:"nama"`
	KodeInternal  string  `json:"kode_internal"`
	NamaJabatan   string  `json:"nama_jabatan"`
	PemimpinPNSID string  `json:"pemimpin_pns_id"`
	IsSatker      bool    `json:"is_satker"`
	UnorInduk     string  `json:"unor_induk"`
	ExpiredDate   db.Date `json:"expired_date"`
	Keterangan    string  `json:"keterangan"`
	Abbreviation  string  `json:"abbreviation"`
	Waktu         string  `json:"waktu"`
	JenisSatker   string  `json:"jenis_satker"`
	Peraturan     string  `json:"peraturan"`
}

type adminUpdateResponse struct {
	Data *unitKerja `json:"data"`
}

func (h *handler) adminUpdate(c echo.Context) error {
	var req adminUpdateRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.update(ctx, updateParams{
		diatasanID:    req.DiatasanID,
		id:            req.ID,
		nama:          req.Nama,
		kodeInternal:  req.KodeInternal,
		namaJabatan:   req.NamaJabatan,
		pemimpinPNSID: req.PemimpinPNSID,
		isSatker:      req.IsSatker,
		unorInduk:     req.UnorInduk,
		expiredDate:   req.ExpiredDate,
		keterangan:    req.Keterangan,
		abbreviation:  req.Abbreviation,
		waktu:         req.Waktu,
		jenisSatker:   req.JenisSatker,
		peraturan:     req.Peraturan,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error creating unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, adminUpdateResponse{Data: data})
}

type adminDeleteRequest struct {
	ID string `param:"id"`
}

func (h *handler) adminDelete(c echo.Context) error {
	var req adminDeleteRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.service.delete(ctx, req.ID)
	if err != nil {
		slog.ErrorContext(ctx, "Error deleting unit kerja.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
