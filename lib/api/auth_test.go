package api

import (
	"bytes"
	"context"
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
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/migrations"
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

	dbData := `
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page1', 'Page 1', null),
			(2,  'portal', 'page2', 'Page 2', null),
			(3,  'portal', 'page3', 'Page 3', '2000-01-01');
		insert into permission
			(id, kode,    nama,    deleted_at) values
			(1,  'write', 'Write', null),
			(2,  'read',  'Read',  null),
			(3,  'del',   'Del',   '2000-01-01'),
			(4,  'exp',   'Exp',   null);
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           1,             null),
			(2,  1,           2,             null),
			(3,  2,           1,             null),
			(4,  2,           2,             null),
			(5,  1,           1,             '2000-01-01'),
			(6,  1,           3,             null),
			(7,  3,           1,             null),
			(8,  1,           4,             null);
		insert into role
			(id, nama,       is_default, is_aktif, deleted_at) values
			(1,  'admin',    false,      true,     null),
			(2,  'pegawai',  true,       true,     null),
			(3,  'delete',   true,       true,     '2000-01-01'),
			(4,  'inactive', true,       false,    null);
		insert into role_resource_permission
			(role_id, resource_permission_id, deleted_at) values
			(1,       1,                      null),
			(1,       4,                      '2000-01-01'),
			(1,       6,                      null),
			(1,       8,                      null),
			(2,       2,                      null),
			(2,       5,                      null),
			(2,       7,                      null),
			(2,       8,                      null),
			(3,       3,                      null),
			(3,       8,                      '2000-01-01'),
			(4,       3,                      null);
		insert into user_role
			(nip,   role_id, deleted_at) values
			('100', 1,       null),
			('101', 1,       '2000-01-01'),
			('102', 3,       null),
			('102', 4,       null),
			('103', 1,       null),
			('103', 2,       null),
			('103', 3,       null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	tests := []struct {
		name                   string
		resourcePermissionKode string
		requestHeader          http.Header
		wantResponseCode       int
		wantResponseBody       string
		wantUser               *User
		wantIDs                any
	}{
		{
			name:                   "ok: allow special permission",
			resourcePermissionKode: "*",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "2",
					"nip":       "99",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "99"},
			wantIDs:          map[string]string{"keycloakID": "1", "zimbraID": "2"},
		},
		{
			name:                   "ok: valid auth header with string audience and nip with role",
			resourcePermissionKode: "portal.page1.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "2",
					"nip":       "100",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
			wantIDs:          map[string]string{"keycloakID": "1", "zimbraID": "2"},
		},
		{
			name:                   "error: valid auth header with string audience and nip with deleted user_role",
			resourcePermissionKode: "portal.page1.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "2",
					"nip":       "101",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"keycloakID": "1", "zimbraID": "2"},
		},
		{
			name:                   "error: valid auth header with string audience and nip with deleted role_resource_permission",
			resourcePermissionKode: "portal.page2.read",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub": "1",
					"nip": "100",
					"aud": "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"keycloakID": "1"},
		},
		{
			name:                   "error: valid auth header with string audience and nip with deleted permission",
			resourcePermissionKode: "portal.page1.del",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip": "100",
					"aud": "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name: "ok: valid auth header with list audience and nip with default roles",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "",
					"zimbra_id": "",
					"nip":       "99",
					"aud":       []string{"nexus", "testing"},
				})},
			},
			resourcePermissionKode: "portal.page1.read",
			wantResponseCode:       http.StatusOK,
			wantUser:               &User{NIP: "99"},
		},
		{
			name: "error: valid auth header with list audience and nip with default roles accessing kode in specific role and deleted resource_permission",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"zimbra_id": "1",
					"nip":       "99",
					"aud":       []string{"testing", "nexus"},
				})},
			},
			resourcePermissionKode: "portal.page1.write",
			wantResponseCode:       http.StatusForbidden,
			wantResponseBody:       `{"message": "akses ditolak"}`,
			wantIDs:                map[string]string{"zimbraID": "1"},
		},
		{
			name:                   "error: valid auth header with list audience and nip with deleted resource",
			resourcePermissionKode: "portal.page3.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "",
					"nip":       "99",
					"aud":       []string{"nexus", "testing", "portal"},
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"keycloakID": "1"},
		},
		{
			name:                   "error: valid auth header with string audience and nip with deleted default role and default inactive role",
			resourcePermissionKode: "portal.page2.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "",
					"zimbra_id": "1",
					"nip":       "99",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"zimbraID": "1"},
		},
		{
			name:                   "error: valid auth header with string audience and nip with deleted role and inactive role",
			resourcePermissionKode: "portal.page2.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip": "102",
					"aud": "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:                   "ok: valid auth header with string audience, multiple roles have access to kode, and user with multiple roles",
			resourcePermissionKode: "portal.page1.exp",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       1,
					"zimbra_id": 1,
					"nip":       "103",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "103"},
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
					header := generateHeader(jwt.MapClaims{"sub": "12", "zimbra_id": "12", "nip": "100", "aud": "testing"})
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
				"Authorization": []string{generateHeader(jwt.MapClaims{"sub": "12", "aud": "testing"})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "nip tidak valid" }`,
		},
		{
			name: "error: audience is nil",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"sub": "12", "nip": "100", "aud": nil})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
		},
		{
			name: "error: missing audience",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"sub": "12", "nip": "100"})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
		},
		{
			name: "error: different string audience",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"sub": "12", "nip": "100", "aud": "nexus"})},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{ "message": "audience tidak valid" }`,
		},
		{
			name: "error: audience not in the array list",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{"sub": "12", "nip": "100", "aud": []string{"nexus", "portal"}})},
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

			handler := func(echo.Context) error { return nil }
			logMw := func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					err := next(c)
					assert.Equal(t, tt.wantUser, CurrentUser(c))
					assert.Equal(t, tt.wantIDs, c.Get(contextKeyUserIDs))
					return err
				}
			}

			svc := NewAuthService(db)
			middleware := NewAuthMiddleware(svc, keyfunc)

			e := echo.New()
			e.Add(http.MethodGet, "/", handler, logMw, middleware(tt.resourcePermissionKode))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody == "" {
				assert.Empty(t, rec.Body)
			} else {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
		})
	}
}
