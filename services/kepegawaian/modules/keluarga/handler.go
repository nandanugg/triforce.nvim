package keluarga

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
	Data keluarga `json:"data"`
}

func (h *handler) list(c echo.Context) error {
	ctx := c.Request().Context()
	data, err := h.service.list(ctx, api.CurrentUser(c).NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data keluarga.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{Data: data})
}

type listAdminRequest struct {
	NIP string `param:"nip"`
}

func (h *handler) listAdmin(c echo.Context) error {
	var req listAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}
	ctx := c.Request().Context()
	data, err := h.service.list(ctx, req.NIP)
	if err != nil {
		slog.ErrorContext(ctx, "Error admin getting data keluarga.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, listResponse{Data: data})
}

type upsertOrangTuaParams struct {
	Nama             string           `json:"nama"`
	NIK              string           `json:"nik"`
	AgamaID          *int16           `json:"agama_id"`
	Hubungan         hubunganOrangTua `json:"hubungan"`
	TanggalMeninggal db.Date          `json:"tanggal_meninggal"`
	AkteMeninggal    string           `json:"akte_meninggal"`
}

type adminCreateOrangTuaRequest struct {
	NIP string `param:"nip"`
	upsertOrangTuaParams
}

type adminCreateOrangTuaResponse struct {
	Data struct {
		ID int32 `json:"id"`
	} `json:"data"`
}

func (h *handler) adminCreateOrangTua(c echo.Context) error {
	var req adminCreateOrangTuaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	id, err := h.service.createOrangTua(c.Request().Context(), req.NIP, req.upsertOrangTuaParams)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "data pegawai tidak ditemukan")
		}
		if errors.Is(err, errAgamaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data agama tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error admin creating orang tua pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var resp adminCreateOrangTuaResponse
	resp.Data.ID = id
	return c.JSON(http.StatusCreated, resp)
}

type adminUpdateOrangTuaRequest struct {
	ID  int32  `param:"id"`
	NIP string `param:"nip"`
	upsertOrangTuaParams
}

func (h *handler) adminUpdateOrangTua(c echo.Context) error {
	var req adminUpdateOrangTuaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.updateOrangTua(c.Request().Context(), req.ID, req.NIP, req.upsertOrangTuaParams)
	if err != nil {
		if errors.Is(err, errAgamaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data agama tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error admin updating orang tua pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminDeleteOrangTuaRequest struct {
	ID  int32  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminDeleteOrangTua(c echo.Context) error {
	var req adminDeleteOrangTuaRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.deleteOrangTua(c.Request().Context(), req.ID, req.NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin deleting orang tua pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type upsertPasanganParams struct {
	Nama               string           `json:"nama"`
	NIK                string           `json:"nik"`
	IsPNS              bool             `json:"is_pns"`
	TanggalLahir       db.Date          `json:"tanggal_lahir"`
	NoKarsus           string           `json:"no_karsus"`
	AgamaID            *int16           `json:"agama_id"`
	StatusPernikahanID int16            `json:"status_pernikahan_id"`
	Hubungan           hubunganPasangan `json:"hubungan"`
	TanggalMenikah     db.Date          `json:"tanggal_menikah"`
	AkteNikah          string           `json:"akte_nikah"`
	TanggalMeninggal   db.Date          `json:"tanggal_meninggal"`
	AkteMeninggal      string           `json:"akte_meninggal"`
	TanggalCerai       db.Date          `json:"tanggal_cerai"`
	AkteCerai          string           `json:"akte_cerai"`
}

type adminCreatePasanganRequest struct {
	NIP string `param:"nip"`
	upsertPasanganParams
}

type adminCreatePasanganResponse struct {
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

func (h *handler) adminCreatePasangan(c echo.Context) error {
	var req adminCreatePasanganRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	id, err := h.service.createPasangan(c.Request().Context(), req.NIP, req.upsertPasanganParams)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "data pegawai tidak ditemukan")
		}
		if errors.Is(err, errAgamaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data agama tidak ditemukan")
		}
		if errors.Is(err, errStatusPernikahanNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data status pernikahan tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error admin creating pasangan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var resp adminCreatePasanganResponse
	resp.Data.ID = id
	return c.JSON(http.StatusCreated, resp)
}

type adminUpdatePasanganRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
	upsertPasanganParams
}

func (h *handler) adminUpdatePasangan(c echo.Context) error {
	var req adminUpdatePasanganRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.updatePasangan(c.Request().Context(), req.ID, req.NIP, req.upsertPasanganParams)
	if err != nil {
		if errors.Is(err, errAgamaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data agama tidak ditemukan")
		}
		if errors.Is(err, errStatusPernikahanNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data status pernikahan tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error admin updating pasangan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminDeletePasanganRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminDeletePasangan(c echo.Context) error {
	var req adminDeletePasanganRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.deletePasangan(c.Request().Context(), req.ID, req.NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin deleting pasangan pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type upsertAnakParams struct {
	Nama               string        `json:"nama"`
	NIK                string        `json:"nik"`
	JenisKelamin       string        `json:"jenis_kelamin"`
	TanggalLahir       db.Date       `json:"tanggal_lahir"`
	PasanganOrangTuaID int64         `json:"pasangan_orang_tua_id"`
	AgamaID            *int16        `json:"agama_id"`
	StatusPernikahanID int16         `json:"status_pernikahan_id"`
	StatusAnak         statusAnak    `json:"status_anak"`
	StatusSekolah      statusSekolah `json:"status_sekolah"`
	AnakKe             *int16        `json:"anak_ke"`
}

type adminCreateAnakRequest struct {
	NIP string `param:"nip"`
	upsertAnakParams
}

type adminCreateAnakResponse struct {
	Data struct {
		ID int64 `json:"id"`
	} `json:"data"`
}

func (h *handler) adminCreateAnak(c echo.Context) error {
	var req adminCreateAnakRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	id, err := h.service.createAnak(c.Request().Context(), req.NIP, req.upsertAnakParams)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return echo.NewHTTPError(http.StatusNotFound, "data pegawai tidak ditemukan")
		}
		if errors.Is(err, errAgamaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data agama tidak ditemukan")
		}
		if errors.Is(err, errStatusPernikahanNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data status pernikahan tidak ditemukan")
		}
		if errors.Is(err, errPasanganOrangTuaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data pasangan orang tua tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error admin creating anak pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	var resp adminCreateAnakResponse
	resp.Data.ID = id
	return c.JSON(http.StatusCreated, resp)
}

type adminUpdateAnakRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
	upsertAnakParams
}

func (h *handler) adminUpdateAnak(c echo.Context) error {
	var req adminUpdateAnakRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.updateAnak(c.Request().Context(), req.ID, req.NIP, req.upsertAnakParams)
	if err != nil {
		if errors.Is(err, errAgamaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data agama tidak ditemukan")
		}
		if errors.Is(err, errStatusPernikahanNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data status pernikahan tidak ditemukan")
		}
		if errors.Is(err, errPasanganOrangTuaNotFound) {
			return echo.NewHTTPError(http.StatusBadRequest, "data pasangan orang tua tidak ditemukan")
		}

		slog.ErrorContext(c.Request().Context(), "Error admin updating anak pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}

type adminDeleteAnakRequest struct {
	ID  int64  `param:"id"`
	NIP string `param:"nip"`
}

func (h *handler) adminDeleteAnak(c echo.Context) error {
	var req adminDeleteAnakRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	found, err := h.service.deleteAnak(c.Request().Context(), req.ID, req.NIP)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error admin deleting anak pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
