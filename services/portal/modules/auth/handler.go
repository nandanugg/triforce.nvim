package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/labstack/echo/v4"
)

type handler struct {
	svc *service
}

func newHandler(s *service) *handler {
	return &handler{svc: s}
}

func (h *handler) login(c echo.Context) error {
	authURL, err := h.svc.generateAuthURL()
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error generating keycloak auth URL.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, authURL)
}

type logoutRequest struct {
	IDTokenHint string `query:"id_token_hint"`
}

func (h *handler) logout(c echo.Context) error {
	var req logoutRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	logoutURL, err := h.svc.generateLogoutURL(req.IDTokenHint)
	if err != nil {
		slog.ErrorContext(c.Request().Context(), "Error generating keycloak logout URL.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.Redirect(http.StatusFound, logoutURL)
}

type exchangeTokenRequest struct {
	Code string `json:"code"`
}

type exchangeTokenResponse struct {
	Data *token `json:"data"`
}

func (h *handler) exchangeToken(c echo.Context) error {
	var req exchangeTokenRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	token, err := h.svc.exchangeToken(req.Code)
	if err != nil {
		var httpErr *httpStatusError
		if errors.As(err, &httpErr) && httpErr.code < 500 {
			return c.JSONBlob(httpErr.code, httpErr.message)
		}

		slog.ErrorContext(c.Request().Context(), "Error exchanging token.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, exchangeTokenResponse{
		Data: token,
	})
}

type refreshTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type refreshTokenResponse struct {
	Data *token `json:"data"`
}

func (h *handler) refreshToken(c echo.Context) error {
	var req refreshTokenRequest
	if err := c.Bind(&req); err != nil {
		return err
	}

	token, err := h.svc.refreshToken(req.RefreshToken)
	if err != nil {
		var httpErr *httpStatusError
		if errors.As(err, &httpErr) && httpErr.code < 500 {
			return c.JSONBlob(httpErr.code, httpErr.message)
		}

		slog.ErrorContext(c.Request().Context(), "Error refreshing token.", "error", err)
		return echo.NewHTTPError(http.StatusInternalServerError)
	}

	return c.JSON(http.StatusOK, refreshTokenResponse{
		Data: token,
	})
}
