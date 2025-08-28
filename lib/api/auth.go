package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/MicahParks/keyfunc/v3"
	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const contextKeyUser = "http-auth-user"

const RoleAdmin = "admin"

type User struct {
	ID   int64
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

func NewAuthMiddleware(keyfunc *Keyfunc) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ss := strings.SplitN(c.Request().Header.Get("Authorization"), ":", 2)
			userID, err := strconv.ParseInt(ss[0], 10, 64)
			if err != nil {
				return echo.NewHTTPError(http.StatusUnauthorized, "token otentikasi tidak valid")
			}

			user := User{ID: userID}
			if len(ss) > 1 {
				user.Role = ss[1]
			}
			c.Set(contextKeyUser, &user)
			return next(c)

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
			// // TODO: @yap fix user_id & roles, role structure `realm_access.roles[]`
			// user := User{ID: int64(claims["user_id"].(float64))}
			// if role, ok := claims["role"]; ok {
			// 	user.Role = role.(string)
			// }
			// c.Set(contextKeyUser, &user)
			//
			// return next(c)
		}
	}
}
