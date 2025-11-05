package api

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"slices"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

const (
	contextKeyUser    = "http-auth-user"
	contextKeyUserIDs = "http-auth-user-ids"
)

type User struct {
	NIP string
}

func CurrentUser(c echo.Context) *User {
	u, _ := c.Get(contextKeyUser).(*User)
	return u
}

type Keyfunc struct {
	jwt.Keyfunc
	Audience string
}

func NewAuthKeyfunc(host, realm, audience string) (*Keyfunc, error) {
	jwksURL := fmt.Sprintf("%s/realms/%s/protocol/openid-connect/certs", host, realm)

	keyfunc, err := keyfunc.NewDefault([]string{jwksURL})
	if err != nil {
		return nil, fmt.Errorf("new keyfunc: %w", err)
	}

	return &Keyfunc{keyfunc.Keyfunc, audience}, nil
}

type AuthService struct {
	db *pgxpool.Pool
}

func NewAuthService(db *pgxpool.Pool) *AuthService {
	return &AuthService{db: db}
}

func (s *AuthService) IsUserHasAccess(ctx context.Context, nip, kode string) (bool, error) {
	if kode == Kode_Allow {
		return true, nil
	}

	var ok bool
	err := s.db.QueryRow(ctx, "select public.is_user_has_access($1, $2)", nip, kode).Scan(&ok)
	return ok, err
}

type AuthInterface interface {
	IsUserHasAccess(ctx context.Context, nip, kode string) (bool, error)
}

// AuthMiddlewareFunc returns Echo middleware for auth and permission checks.
type AuthMiddlewareFunc func(kodeResourcePermission string) echo.MiddlewareFunc

// NewAuthMiddleware creates middleware that allows requests only if the user have permission.
func NewAuthMiddleware(svc AuthInterface, keyfunc *Keyfunc) AuthMiddlewareFunc {
	return func(kode string) echo.MiddlewareFunc {
		return func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				claims, err := validateRequestToken(c, keyfunc)
				if err != nil {
					return err
				}

				nip := claims["nip"].(string)
				ok, err := svc.IsUserHasAccess(c.Request().Context(), nip, kode)
				if err != nil {
					slog.ErrorContext(c.Request().Context(), "Error checking user access", "error", err)
					return echo.NewHTTPError(http.StatusInternalServerError)
				}

				if ok {
					c.Set(contextKeyUser, &User{NIP: nip})
					return next(c)
				}
				return echo.NewHTTPError(http.StatusForbidden, "akses ditolak")
			}
		}
	}
}

func validateRequestToken(c echo.Context, keyfunc *Keyfunc) (jwt.MapClaims, error) {
	header := c.Request().Header.Get("Authorization")
	if !strings.HasPrefix(header, "Bearer ") {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "token otentikasi tidak valid")
	}

	token := strings.TrimPrefix(header, "Bearer ")
	claims := jwt.MapClaims{}
	_, err := jwt.ParseWithClaims(token, &claims, keyfunc.Keyfunc)
	if err != nil {
		msg := "akses ditolak"
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			msg = "token otentikasi tidak valid"
		case errors.Is(err, jwt.ErrTokenExpired):
			msg = "token otentikasi sudah kedaluwarsa"
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			msg = "signature token otentikasi tidak valid"
		}
		return nil, echo.NewHTTPError(http.StatusUnauthorized, msg)
	}

	aud, err := claims.GetAudience()
	if err != nil || !slices.Contains(aud, keyfunc.Audience) {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "audience tidak valid")
	}

	if _, ok := claims["nip"].(string); !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "nip tidak valid")
	}

	// set keycloak_id & zimbra_id for logging purpose
	ids := make(map[string]string)
	if sub, _ := claims.GetSubject(); sub != "" {
		ids["keycloakID"] = sub
	}
	if zimbraID, _ := claims["zimbra_id"].(string); zimbraID != "" {
		ids["zimbraID"] = zimbraID
	}
	if len(ids) > 0 {
		c.Set(contextKeyUserIDs, ids)
	}

	return claims, nil
}
