package pegawai

import (
	"encoding/base64"
	"log/slog"
	"net/http"
	"unicode/utf8"

	"github.com/labstack/echo/v4"
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
