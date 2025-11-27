package pegawai

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"unicode/utf8"

	"github.com/labstack/echo/v4"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
)

type handler struct {
	service *service
}

func newHandler(s *service) *handler {
	return &handler{service: s}
}

type profileRequest struct {
	PNSID string `param:"pns_id"`
}

type profileResponse struct {
	Data *profile `json:"data"`
}

func (h *handler) getProfile(c echo.Context) error {
	var req profileRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	pnsID, err := base64.RawURLEncoding.DecodeString(req.PNSID)
	if err != nil || !utf8.Valid(pnsID) {
		// treat as not found
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	ctx := c.Request().Context()
	data, err := h.service.getProfileByPNSID(ctx, string(pnsID))
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data profil pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, profileResponse{
		Data: data,
	})
}

type listAdminRequest struct {
	api.PaginationRequest
	Keyword    string `query:"keyword"`
	UnitID     string `query:"unit_id"`
	GolonganID int32  `query:"golongan_id"`
	JabatanID  string `query:"jabatan_id"`
	Status     string `query:"status"`
	Tipe       string `query:"tipe"`
}

type listAdminResponse struct {
	Data []pegawai          `json:"data"`
	Meta api.MetaPagination `json:"meta"`
}

func (h *handler) listAdmin(c echo.Context) error {
	var req listAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	var data []pegawai
	var count uint
	var err error
	params := adminListPegawaiParams{
		limit:      req.Limit,
		offset:     req.Offset,
		keyword:    req.Keyword,
		unitID:     req.UnitID,
		golonganID: req.GolonganID,
		jabatanID:  req.JabatanID,
		status:     req.Status,
	}

	switch req.Tipe {
	case "", tipePegawaiAktif:
		data, count, err = h.service.adminListPegawai(ctx, params)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting data list pegawai aktif.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	case tipePegawaiPPNPN:
		data, count, err = h.service.adminListPegawaiPPNPN(ctx, params)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting data list pegawai ppnpn.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	case tipePegawaiNonAktif:
		data, count, err = h.service.adminListPegawaiNonAktif(ctx, params)
		if err != nil {
			slog.ErrorContext(ctx, "Error getting data list pegawai nonaktif.", "error", err)
			return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
		}
	}

	return c.JSON(http.StatusOK, listAdminResponse{
		Data: data,
		Meta: api.MetaPagination{
			Total:  count,
			Limit:  req.Limit,
			Offset: req.Offset,
		},
	})
}

type getAdminRequest struct {
	NIP string `param:"nip"`
}

type getAdminResponse struct {
	Data pegawaiDetail `json:"data"`
}

func (h *handler) getAdmin(c echo.Context) error {
	var req getAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	data, err := h.service.adminGetPegawai(ctx, req.NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data detil pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if data == nil {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.JSON(http.StatusOK, getAdminResponse{Data: *data})
}
