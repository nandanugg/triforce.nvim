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

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func TestNewAuthMiddleware(t *testing.T) {
	t.Parallel()

	privateKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	publicKey := &privateKey.PublicKey
	keyfunc := &Keyfunc{
		Keyfunc:  func(*jwt.Token) (any, error) { return publicKey, nil },
		Audience: "testing",
	}

	generateHeader := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, _ := token.SignedString(privateKey)
		return "Bearer " + tokenString
	}

	tests := []struct {
		name             string
		service          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantUser         *User
	}{
		{
			name: "ok: valid auth header with string audience without role",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip": "100",
					"aud": "testing",
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
		},
		{
			name:    "ok: valid auth header with string audience with service role",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   "testing",
					"roles": map[string]any{"portal": "admin"},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100", Role: "admin"},
		},
		{
			name:    "ok: valid auth header with string audience with other service role",
			service: "kepegawaian",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   "testing",
					"roles": map[string]any{"portal": "admin"},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
		},
		{
			name:    "ok: valid auth header with string audience with role is array (unsupported)",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   "testing",
					"roles": map[string]any{"portal": []string{"admin", "pegawai"}},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
		},
		{
			name: "ok: valid auth header with list audience with empty role",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   []string{"nexus", "testing"},
					"roles": map[string]any{},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
		},
		{
			name:    "ok: valid auth header with list audience with multiple service role",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   []string{"nexus", "testing"},
					"roles": map[string]any{"kepegawaian": "admin", "portal": "pegawai"}},
				)},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100", Role: "pegawai"},
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
				"Authorization": []string{generateHeader(jwt.MapClaims{"exp": jwt.NewNumericDate(time.Now().Add(-1 * time.Minute))})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "token otentikasi sudah kedaluwarsa" }`,
		},
		{
			name: "error: tampered jwt payload",
			requestHeader: http.Header{
				"Authorization": []string{func() string {
					header := generateHeader(jwt.MapClaims{"nip": "100", "aud": "testing"})
					encodedClaims := strings.Split(header, ".")[1]
					claims, _ := base64.RawStdEncoding.DecodeString(encodedClaims)
					claims = bytes.ReplaceAll(claims, []byte("100"), []byte("200"))

					return strings.ReplaceAll(header, encodedClaims, base64.RawStdEncoding.EncodeToString(claims))
				}()},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "signature token otentikasi tidak valid" }`,
		},
		{
			name: "error: missing nip",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"aud": "testing"})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "nip tidak valid" }`,
		},
		{
			name: "error: audience is nil",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"nip": "100", "aud": nil})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
		},
		{
			name: "error: missing audience",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"nip": "100"})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
		},
		{
			name: "error: different string audience",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"nip": "100", "aud": "nexus"})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
		},
		{
			name: "error: audience not in the array list",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"nip": "100", "aud": []string{"nexus", "portal"}})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
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
			// middleware := NewAuthMiddleware(tt.service, keyfunc)
			middleware := NewAuthMiddleware(keyfunc)
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
