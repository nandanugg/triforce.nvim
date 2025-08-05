package api

import (
	"bytes"
	"crypto/rand"
	"crypto/rsa"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthMiddleware(t *testing.T) {
	t.Parallel()

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey

	tests := []struct {
		name             string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantUser         *User
	}{
		{
			name: "ok: valid auth header",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader("user1", privateKey)},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{ID: "user1"},
		},
		{
			name:             "error: missing auth header",
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "token otentikasi tidak valid" }`,
		},
		{
			name:             "error: invalid auth header format",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "token otentikasi tidak valid" }`,
		},
		{
			name: "error: expired auth token",
			requestHeader: http.Header{
				"Authorization": []string{func() string {
					token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{
						"exp": jwt.NewNumericDate(time.Now().Add(-1 * time.Minute)),
					})
					tokenString, _ := token.SignedString(privateKey)
					return "Bearer " + tokenString
				}()},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "token otentikasi sudah kedaluwarsa" }`,
		},
		{
			name: "error: tampered jwt payload",
			requestHeader: http.Header{
				"Authorization": []string{func() string {
					header := generateHeader("user1", privateKey)
					encodedClaims := strings.Split(header, ".")[1]
					claims, _ := base64.RawStdEncoding.DecodeString(encodedClaims)
					claims = bytes.ReplaceAll(claims, []byte("user1"), []byte("admin"))

					return strings.ReplaceAll(header, encodedClaims, base64.RawStdEncoding.EncodeToString(claims))
				}()},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "signature token otentikasi tidak valid" }`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/", nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			var actualUser *User

			e := echo.New()
			handler := func(c echo.Context) error {
				actualUser = CurrentUser(c)
				return nil
			}
			middleware := NewAuthMiddleware(publicKey)
			e.Add(http.MethodGet, "/", handler, middleware)
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody == "" {
				assert.Empty(t, rec.Body)
			} else {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}

			assert.Equal(t, tt.wantUser, actualUser)
		})
	}
}

func generateHeader(userID string, signingKey *rsa.PrivateKey) string {
	token := jwt.NewWithClaims(jwt.SigningMethodRS256, jwt.MapClaims{"user_id": userID})
	tokenString, _ := token.SignedString(signingKey)
	return "Bearer " + tokenString
}
