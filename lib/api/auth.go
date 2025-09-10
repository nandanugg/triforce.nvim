package api

import (
	"fmt"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const contextKeyUser = "http-auth-user"

const RoleAdmin = "admin"

type User struct {
	ID   int64 // @yap TODO: remove later, many logic or unit test depend on this currently
	NIP  string
	Role string
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

func NewAuthMiddleware(service string, keyfunc *Keyfunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// header := c.Request().Header.Get("Authorization")
			// if !strings.HasPrefix(header, "Bearer ") {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, "token otentikasi tidak valid")
			// }
			//
			// token := strings.TrimPrefix(header, "Bearer ")
			// claims := jwt.MapClaims{}
			// _, err := jwt.ParseWithClaims(token, &claims, keyfunc.Keyfunc)
			// if err != nil {
			// 	msg := "akses ditolak"
			// 	switch {
			// 	case errors.Is(err, jwt.ErrTokenMalformed):
			// 		msg = "token otentikasi tidak valid"
			// 	case errors.Is(err, jwt.ErrTokenExpired):
			// 		msg = "token otentikasi sudah kedaluwarsa"
			// 	case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			// 		msg = "signature token otentikasi tidak valid"
			// 	}
			// 	return echo.NewHTTPError(http.StatusUnauthorized, msg)
			// }
			//
			// switch aud := claims["aud"].(type) {
			// case string:
			// 	if aud != keyfunc.Audience {
			// 		return echo.NewHTTPError(http.StatusUnauthorized, "audience tidak valid")
			// 	}
			// case []any:
			// 	if !slices.Contains(aud, any(keyfunc.Audience)) {
			// 		return echo.NewHTTPError(http.StatusUnauthorized, "audience tidak valid")
			// 	}
			// default:
			// 	return echo.NewHTTPError(http.StatusUnauthorized, "audience tidak valid")
			// }
			//
			// nip, ok := claims["nip"].(string)
			// if !ok {
			// 	return echo.NewHTTPError(http.StatusUnauthorized, "nip tidak valid")
			// }

			user := User{NIP: "198501012020011001"}
			// if roles, ok := claims["roles"].(map[string]any); ok {
			// 	user.Role, _ = roles[service].(string)
			// }
			c.Set(contextKeyUser, &user)

			return next(c)
		}
	}
}
