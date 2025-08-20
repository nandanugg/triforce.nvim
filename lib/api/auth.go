package api

import (
	"crypto/rsa"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v4"
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

func NewAuthMiddleware(jwtPublicKey *rsa.PublicKey) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			header := c.Request().Header.Get("Authorization")
			if !strings.HasPrefix(header, "Bearer ") {
				return echo.NewHTTPError(http.StatusUnauthorized, "token otentikasi tidak valid")
			}

			token := strings.TrimPrefix(header, "Bearer ")
			claims := jwt.MapClaims{}
			_, err := jwt.ParseWithClaims(token, &claims, func(*jwt.Token) (any, error) {
				return jwtPublicKey, nil
			})
			if err != nil {
				msg := err.Error()
				switch msg {
				case "token contains an invalid number of segments":
					msg = "token otentikasi tidak valid"
				case "Token is expired":
					msg = "token otentikasi sudah kedaluwarsa"
				case "crypto/rsa: verification error":
					msg = "signature token otentikasi tidak valid"
				default:
					msg = "akses ditolak"
				}
				return echo.NewHTTPError(http.StatusUnauthorized, msg)
			}

			user := User{ID: int64(claims["user_id"].(float64))}
			if role, ok := claims["role"]; ok {
				user.Role = role.(string)
			}
			c.Set(contextKeyUser, &user)

			return next(c)
		}
	}
}
