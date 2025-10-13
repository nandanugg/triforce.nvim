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

	tests := []struct {
		name             string
		service          string
		requestHeader    http.Header
		allowedRoles     []string
		wantResponseCode int
		wantResponseBody string
		wantUser         *User
		wantIDs          any
	}{
		{
			name: "ok: valid auth header with string audience without role and none allowed roles is provided",
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
			name: "error: valid auth header with string audience without role but with allowed role",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "2",
					"nip":       "100",
					"aud":       "testing",
				})},
			},
			allowedRoles:     []string{"admin"},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{ "message": "akses ditolak" }`,
			wantIDs: map[string]string{
				"keycloakID": "1",
				"zimbraID":   "2",
			},
		},
		{
			name:    "ok: valid auth header with string audience with service role and none allowed roles is provided",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "",
					"zimbra_id": "",
					"nip":       "100",
					"aud":       "testing",
					"roles":     map[string]any{"portal": "admin"},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100", Role: "admin"},
		},
		{
			name:    "ok: valid auth header with string audience with service role and within allowed roles",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":   "1",
					"nip":   "100",
					"aud":   "testing",
					"roles": map[string]any{"portal": "admin"},
				})},
			},
			allowedRoles:     []string{"admin"},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100", Role: "admin"},
			wantIDs:          map[string]string{"keycloakID": "1"},
		},
		{
			name:    "error: valid auth header with string audience with service role and without match allowed roles",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"zimbra_id": "2",
					"nip":       "100",
					"aud":       "testing",
					"roles":     map[string]any{"portal": "admin"},
				})},
			},
			allowedRoles:     []string{"pegawai", "tester", "guest"},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{ "message": "akses ditolak" }`,
			wantIDs:          map[string]string{"zimbraID": "2"},
		},
		{
			name:    "ok: valid auth header with string audience with other service role and none allowed roles is provided",
			service: "kepegawaian",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       123,
					"zimbra_id": 123,
					"nip":       "100",
					"aud":       "testing",
					"roles":     map[string]any{"portal": "admin"},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
		},
		{
			name:    "error: valid auth header with string audience with other service role and with allowed roles",
			service: "kepegawaian",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "123",
					"zimbra_id": "",
					"nip":       "100",
					"aud":       "testing",
					"roles":     map[string]any{"portal": "admin"},
				})},
			},
			allowedRoles:     []string{"admin"},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{ "message": "akses ditolak" }`,
			wantIDs:          map[string]string{"keycloakID": "123"},
		},
		{
			name:    "ok: valid auth header with string audience with role is array (unsupported) and none allowed roles is provided",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "",
					"zimbra_id": "123",
					"nip":       "100",
					"aud":       "testing",
					"roles":     map[string]any{"portal": []string{"admin", "pegawai"}},
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
			wantIDs:          map[string]string{"zimbraID": "123"},
		},
		{
			name:    "error: valid auth header with string audience with role is array (unsupported) and with allowed roles",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   "testing",
					"roles": map[string]any{"portal": []string{"admin", "pegawai"}},
				})},
			},
			allowedRoles:     []string{"admin", "pegawai"},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{ "message": "akses ditolak" }`,
		},
		{
			name: "ok: valid auth header with list audience with empty role and none allowed roles is provided",
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
			name: "error: valid auth header with list audience with empty role and with multiple allowed roles",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   []string{"nexus", "testing"},
					"roles": map[string]any{},
				})},
			},
			allowedRoles:     []string{"admin", "tester", "pegawai"},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{ "message": "akses ditolak" }`,
		},
		{
			name:    "ok: valid auth header with list audience with multiple service role and none allowed roles is provided",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   []string{"nexus", "testing"},
					"roles": map[string]any{"kepegawaian": "admin", "portal": "pegawai"},
				},
				)},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100", Role: "pegawai"},
		},
		{
			name:    "ok: valid auth header with list audience with multiple service role and with multiple allowed roles",
			service: "portal",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"nip":   "100",
					"aud":   []string{"nexus", "testing"},
					"roles": map[string]any{"kepegawaian": "admin", "portal": "pegawai"},
				},
				)},
			},
			allowedRoles:     []string{"guest", "teacher", "pegawai", "tester"},
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
					header := generateHeader(jwt.MapClaims{"sub": "123", "zimbra_id": "123", "nip": "100", "aud": "testing"})
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

			e := echo.New()
			handler := func(echo.Context) error { return nil }
			logMw := func(next echo.HandlerFunc) echo.HandlerFunc {
				return func(c echo.Context) error {
					err := next(c)
					assert.Equal(t, tt.wantUser, CurrentUser(c))
					assert.Equal(t, tt.wantIDs, c.Get(contextKeyUserIDs))
					return err
				}
			}

			middleware := NewAuthMiddleware(tt.service, keyfunc)
			e.Add(http.MethodGet, "/", handler, logMw, middleware(tt.allowedRoles...))
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

func TestNewAuthResourcePermissionMiddleware(t *testing.T) {
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

	seedData := `
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page1', 'Page 1', null),
			(2,  'portal', 'page2', 'Page 2', null),
			(3,  'portal', 'page3', 'Page 3', '2000-01-01');
		insert into permission
			(id, kode,    nama,    deleted_at) values
			(1,  'write', 'Write', null),
			(2,  'read',  'Read',  null),
			(3,  'del',   'Del',   '2000-01-01');
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           1,             null),
			(2,  1,           2,             null),
			(3,  2,           1,             null),
			(4,  2,           2,             null),
			(5,  1,           1,             '2000-01-01'),
			(6,  1,           3,             null),
			(7,  3,           1,             null);
		insert into role
			(id, nama,      is_default, deleted_at) values
			(1,  'admin',   false,      null),
			(2,  'pegawai', true,       null),
			(3,  'del',     true,       '2000-01-01');
		insert into role_resource_permission
			(role_id, resource_permission_id, deleted_at) values
			(1,       1,                      null),
			(1,       4,                      '2000-01-01'),
			(1,       6,                      null),
			(2,       2,                      null),
			(2,       5,                      null),
			(2,       7,                      null),
			(3,       3,                      null);
	`
	tests := []struct {
		name                   string
		dbData                 string
		resourcePermissionKode string
		requestHeader          http.Header
		wantResponseCode       int
		wantResponseBody       string
		wantUser               *User
		wantIDs                any
	}{
		{
			name: "ok: valid auth header with string audience and nip with role",
			dbData: seedData + `
				insert into user_role (nip, role_id) values ('100', 1);
			`,
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
			name: "error: valid auth header with string audience and nip with deleted user_role",
			dbData: seedData + `
				insert into user_role (nip, role_id, deleted_at) values ('100', 1, '2000-01-01');
			`,
			resourcePermissionKode: "portal.page1.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "2",
					"nip":       "100",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"keycloakID": "1", "zimbraID": "2"},
		},
		{
			name: "error: valid auth header with string audience and nip with deleted role_resource_permission",
			dbData: seedData + `
				insert into user_role (nip, role_id) values ('100', 1);
			`,
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
			name: "error: valid auth header with string audience and nip with deleted permission",
			dbData: seedData + `
				insert into user_role (nip, role_id) values ('100', 1);
			`,
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
			name:   "ok: valid auth header with list audience and nip with default roles",
			dbData: seedData,
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "",
					"zimbra_id": "",
					"nip":       "100",
					"aud":       []string{"nexus", "testing"},
				})},
			},
			resourcePermissionKode: "portal.page1.read",
			wantResponseCode:       http.StatusOK,
			wantUser:               &User{NIP: "100"},
		},
		{
			name:   "error: valid auth header with list audience and nip with default roles accessing kode in specific role and deleted resource_permission",
			dbData: seedData,
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"zimbra_id": "1",
					"nip":       "100",
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
			dbData:                 seedData,
			resourcePermissionKode: "portal.page3.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "1",
					"zimbra_id": "",
					"nip":       "100",
					"aud":       []string{"nexus", "testing", "portal"},
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"keycloakID": "1"},
		},
		{
			name:                   "error: valid auth header with string audience and nip with deleted default role",
			dbData:                 seedData,
			resourcePermissionKode: "portal.page2.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       "",
					"zimbra_id": "1",
					"nip":       "100",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
			wantIDs:          map[string]string{"zimbraID": "1"},
		},
		{
			name: "error: valid auth header with string audience and nip with deleted role",
			dbData: seedData + `
				insert into user_role (nip, role_id) values ('100', 3);
			`,
			resourcePermissionKode: "portal.page2.write",
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
			name: "ok: valid auth header with string audience, multiple roles have access to kode, and user with multiple roles",
			dbData: seedData + `
				insert into role_resource_permission
					(role_id, resource_permission_id, deleted_at) values
					(2,       1,                      null),
					(3,       1,                      '2000-01-01');
				insert into user_role
					(nip,   role_id) values
					('100', 1),
					('100', 2),
					('100', 3);
			`,
			resourcePermissionKode: "portal.page1.write",
			requestHeader: http.Header{
				"Authorization": []string{generateHeader(jwt.MapClaims{
					"sub":       1,
					"zimbra_id": 1,
					"nip":       "100",
					"aud":       "testing",
				})},
			},
			wantResponseCode: http.StatusOK,
			wantUser:         &User{NIP: "100"},
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

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

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

			svc := NewAuthResourcePermissionService(db)
			middleware := NewAuthResourcePermissionMiddleware(svc, keyfunc)

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
