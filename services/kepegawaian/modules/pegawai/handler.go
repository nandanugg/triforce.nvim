package pegawai

import (
	"encoding/base64"
	"errors"
	"log/slog"
	"net/http"
	"unicode/utf8"

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
	data, count, err := h.service.adminListPegawai(ctx, adminListPegawaiParams{
		limit:      req.Limit,
		offset:     req.Offset,
		keyword:    req.Keyword,
		unitID:     req.UnitID,
		golonganID: req.GolonganID,
		jabatanID:  req.JabatanID,
		status:     req.Status,
	})
	if err != nil {
		slog.ErrorContext(ctx, "Error getting data list pegawai aktif.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
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

type updatePegawaiParams struct {
	Nama                     string  `json:"nama"`
	GelarDepan               string  `json:"gelar_depan"`
	GelarBelakang            string  `json:"gelar_belakang"`
	NIP                      string  `json:"nip"`
	JenisKelamin             string  `json:"jenis_kelamin"`
	TanggalLahir             db.Date `json:"tanggal_lahir"`
	NIK                      string  `json:"nik"`
	NomorKK                  string  `json:"nomor_kk"`
	TempatLahirID            string  `json:"tempat_lahir_id"`
	TingkatPendidikanID      *int16  `json:"tingkat_pendidikan_id"`
	PendidikanID             string  `json:"pendidikan_id"`
	StatusPernikahanID       *int16  `json:"status_pernikahan_id"`
	AgamaID                  *int16  `json:"agama_id"`
	JenisPegawaiID           *int16  `json:"jenis_pegawai_id"`
	MasaKerjaGolongan        string  `json:"masa_kerja_golongan"`
	JenisJabatanID           *int16  `json:"jenis_jabatan_id"`
	JabatanID                string  `json:"jabatan_id"`
	UnitOrganisasiID         string  `json:"unit_organisasi_id"`
	LokasiKerjaID            string  `json:"lokasi_kerja_id"`
	GolonganRuangAwalID      *int16  `json:"golongan_ruang_awal_id"`
	GolonganRuangAkhirID     *int16  `json:"golongan_ruang_akhir_id"`
	TMTGolongan              db.Date `json:"tmt_golongan"`
	NomorSKASN               string  `json:"nomor_sk_asn"`
	TMTASN                   db.Date `json:"tmt_asn"`
	StatusPNS                string  `json:"status_pns"`
	EmailDikbud              string  `json:"email_dikbud"`
	EmailPribadi             string  `json:"email_pribadi"`
	Alamat                   string  `json:"alamat"`
	NoHP                     string  `json:"no_hp"`
	NoKontakDarurat          string  `json:"no_kontak_darurat"`
	NomorSuratDokter         string  `json:"nomor_surat_dokter"`
	TanggalSuratDokter       db.Date `json:"tanggal_surat_dokter"`
	NomorSuratBebasNarkoba   string  `json:"nomor_surat_bebas_narkoba"`
	TanggalSuratBebasNarkoba db.Date `json:"tanggal_surat_bebas_narkoba"`
	NomorCatatanPolisi       string  `json:"nomor_catatan_polisi"`
	TanggalCatatanPolisi     db.Date `json:"tanggal_catatan_polisi"`
	AkteKelahiran            string  `json:"akte_kelahiran"`
	NomorBPJS                string  `json:"nomor_bpjs"`
	NPWP                     string  `json:"npwp"`
	TanggalNPWP              db.Date `json:"tanggal_npwp"`
	NomorTaspen              string  `json:"nomor_taspen"`
	MkBulan                  *int16  `json:"mk_bulan"`
	MkTahun                  *int16  `json:"mk_tahun"`
	MkBulanSwasta            *int16  `json:"mk_bulan_swasta"`
	MkTahunSwasta            *int16  `json:"mk_tahun_swasta"`
}

type postAdminRequest struct {
	NIP string `param:"nip"`
	updatePegawaiParams
}

func (h *handler) putAdmin(c echo.Context) error {
	var req postAdminRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	ctx := c.Request().Context()
	found, err := h.service.updatePegawai(ctx, req.NIP, req.updatePegawaiParams)
	if err != nil {
		var multiErr *api.MultiError
		if errors.As(err, &multiErr) {
			return echo.NewHTTPError(http.StatusBadRequest, multiErr.Error())
		}

		slog.ErrorContext(ctx, "Error updating data pegawai.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	if !found {
		return echo.NewHTTPError(http.StatusNotFound, "data tidak ditemukan")
	}

	return c.NoContent(http.StatusNoContent)
}
