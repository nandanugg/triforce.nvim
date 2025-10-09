package auth

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func Test_handler_login(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
		wantRedirect     string
	}{
		{
			name: "success redirect with redirect_uri",
			requestQuery: url.Values{
				"redirect_uri": []string{"http://localhost:5173/callback"},
			},
			wantResponseCode: 302,
			wantRedirect:     `https://auth.local/realms/nexus/protocol/openid-connect/auth?client_id=my-portal&prompt=login&redirect_uri=http%3A%2F%2Flocalhost%3A5173%2Fcallback&response_type=code&scope=openid`,
		},
		{
			name:             "success redirect without redirect_uri",
			wantResponseCode: 302,
			wantRedirect:     `https://auth.local/realms/nexus/protocol/openid-connect/auth?client_id=my-portal&prompt=login&redirect_uri=https%3A%2F%2Fportal.local%2Fcallback&response_type=code&scope=openid`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			keycloak := config.Keycloak{
				PublicHost:  "https://auth.local",
				Realm:       "nexus",
				ClientID:    "my-portal",
				RedirectURI: "https://portal.local/callback",
			}
			RegisterRoutes(e, nil, keycloak, nil, nil, nil)
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			if rec.Code == 302 {
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
				assert.Empty(t, rec.Body)
			} else {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
		})
	}
}

func Test_handler_logout(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
		wantRedirect     string
	}{
		{
			name:             "missing required params",
			wantResponseCode: 400,
			wantResponseBody: `{"message": "parameter \"id_token_hint\" harus diisi"}`,
		},
		{
			name: "success redirect with post_logout_redirect_uri",
			requestQuery: url.Values{
				"id_token_hint":            []string{"e9d74da4-a4e7-4865-8a4b-601ad1bca900"},
				"post_logout_redirect_uri": []string{"http://localhost:5173/"},
			},
			wantResponseCode: 302,
			wantRedirect:     `https://auth.local/realms/nexus/protocol/openid-connect/logout?id_token_hint=e9d74da4-a4e7-4865-8a4b-601ad1bca900&post_logout_redirect_uri=http%3A%2F%2Flocalhost%3A5173%2F`,
		},
		{
			name: "success redirect without post_logout_redirect_uri",
			requestQuery: url.Values{
				"id_token_hint": []string{"e9d74da4-a4e7-4865-8a4b-601ad1bca123"},
			},
			wantResponseCode: 302,
			wantRedirect:     `https://auth.local/realms/nexus/protocol/openid-connect/logout?id_token_hint=e9d74da4-a4e7-4865-8a4b-601ad1bca123&post_logout_redirect_uri=https%3A%2F%2Fportal.local%2F`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/auth/logout", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			keycloak := config.Keycloak{
				PublicHost:            "https://auth.local",
				Realm:                 "nexus",
				PostLogoutRedirectURI: "https://portal.local/",
			}
			RegisterRoutes(e, nil, keycloak, nil, nil, nil)
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			if rec.Code == 302 {
				assert.Equal(t, tt.wantRedirect, rec.Header().Get("Location"))
				assert.Empty(t, rec.Body)
			} else {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
		})
	}
}

func Test_handler_exchangeToken(t *testing.T) {
	t.Parallel()

	generateToken := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, _ := token.SignedString(apitest.JwtPrivateKey)
		return tokenString
	}

	generateTokenWithKID := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = "my-kid"
		tokenString, _ := token.SignedString(apitest.JwtPrivateKey)
		return tokenString
	}

	tests := []struct {
		name             string
		dbData           string
		requestBody      string
		keycloakRespCode int
		keycloakRespBody []byte
		wantRedirectURI  string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "success exchange token user without roles using keycloak_id without forwarded header",
			dbData: `
				insert into "user" (id, source, nip, created_at, updated_at) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1c', '2000-01-01', '2000-01-01'),
					('00000000-0000-0000-0000-000000000002', 'keycloak', '1a', '2000-01-01', '2000-01-01'),
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1b', '2000-01-01', '2000-01-01');
			`,
			requestBody:      `{"code": "my-code", "redirect_uri": "http://localhost:5173/callback"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantRedirectURI:  "http://localhost:5173/callback",
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "` + generateTokenWithKID(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "nip": "1c"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					"source":          "keycloak",
					"nip":             "1c",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   "{last_login_at}",
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					"source":          "zimbra",
					"nip":             "1b",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					"source":          "keycloak",
					"nip":             "1a",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
			},
		},
		{
			name: "success exchange token user without deleted roles using zimbra_id as priority with forwarded header without redirect_uri",
			dbData: `
				insert into "user" (id, source, nip, created_at, updated_at) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1c', '2000-01-01', '2000-01-01'),
					('00000000-0000-0000-0000-000000000002', 'zimbra', '1b', '2000-01-01', '2000-01-01');
				insert into role (id, service, nama, deleted_at) values
					(1, 'portal', 'admin', '2000-01-01'),
					(2, 'portal', 'admin', null);
				insert into user_role (nip, role_id, deleted_at) values
					('1c', 2, null),
					('1b', 1, null),
					('1b', 2, '2000-01-01');
			`,
			requestBody:      `{"code": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantRedirectURI:  "https://portal.local/callback",
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "` + generateTokenWithKID(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002", "nip": "1b"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					"source":          "keycloak",
					"nip":             "1c",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					"source":          "zimbra",
					"nip":             "1b",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   "{last_login_at}",
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
			},
		},
		{
			name: "success exchange token user with multiple roles order by updated_at fallback using keycloak_id when user zimbra is not found with forwarded host only and with redirect_uri",
			dbData: `
				insert into "user" (id, source, nip, created_at, updated_at, last_login_at) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
				insert into role (id, service, nama) values
					(1, 'portal', 'admin'),
					(2, 'portal', 'pegawai'),
					(3, 'kepegawaian', 'admin'),
					(4, 'kepegawaian', 'pegawai'),
					(5, 'kepegawaian', 'guest');
				insert into user_role (nip, role_id, updated_at, deleted_at) values
					('1c', 1, '2000-01-02', null),
					('1c', 2, '2000-01-01', null),
					('1c', 3, '2000-01-01', null),
					('1c', 4, '2000-01-02', null),
					('1c', 5, '2000-01-03', '2000-01-03');
			`,
			requestBody:      `{"code": "my-code", "redirect_uri": "http://localhost:5173/callback"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantRedirectURI:  "http://localhost:5173/callback",
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "` + generateTokenWithKID(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002", "nip": "1c", "roles": map[string]string{"portal": "admin", "kepegawaian": "pegawai"}}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					"source":          "keycloak",
					"nip":             "1c",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   "{last_login_at}",
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
			},
		},
		{
			name: "error exchange token, user keycloak or zimbra not found",
			dbData: `
				insert into "user" (id, source, nip, created_at, updated_at, deleted_at) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1a', '2000-01-01', '2000-01-01', '2000-02-02'),
					('00000000-0000-0000-0000-000000000002', 'zimbra', '1b', '2000-01-01', '2000-01-01', '2000-01-01'),
					('00000000-0000-0000-0000-000000000002', 'keycloak', '1a', '2000-01-01', '2000-01-01', null),
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1b', '2000-01-01', '2000-01-01', null);
			`,
			requestBody:      `{"code": "my-code", "redirect_uri": "http://localhost:5173/callback"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantRedirectURI:  "http://localhost:5173/callback",
			wantResponseCode: 422,
			wantResponseBody: `{"message": "user tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					"source":          "keycloak",
					"nip":             "1a",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      time.Date(2000, 2, 2, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000001")),
					"source":          "zimbra",
					"nip":             "1b",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					"source":          "keycloak",
					"nip":             "1a",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      nil,
				},
				{
					"id":              [16]byte(uuid.MustParse("00000000-0000-0000-0000-000000000002")),
					"source":          "zimbra",
					"nip":             "1b",
					"nama":            nil,
					"email":           nil,
					"unit_organisasi": nil,
					"last_login_at":   nil,
					"created_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":      time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error exchange token, user keycloak or zimbra is not uuid",
			requestBody:      `{"code": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "1", "zimbra_id": "2"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantRedirectURI:  "https://portal.local/callback",
			wantResponseCode: 422,
			wantResponseBody: `{"message": "user tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error exchange token, invalid access_token from keycloak",
			requestBody:      `{"code": "my-code", "redirect_uri": "http://localhost:5173/callback"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "foo",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantRedirectURI:  "http://localhost:5173/callback",
			wantResponseCode: 500,
			wantResponseBody: `{"message": "Internal Server Error"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error exchange token, keycloak return status 4xx",
			requestBody:      `{"code": "my-code", "redirect_uri": "http://localhost:5173/callback"}`,
			keycloakRespCode: 409,
			keycloakRespBody: []byte(`{"error": "invalid code"}`),
			wantRedirectURI:  "http://localhost:5173/callback",
			wantResponseCode: 409,
			wantResponseBody: `{"error": "invalid code"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error exchange token, keycloak return status 5xx",
			requestBody:      `{"code": "my-code"}`,
			keycloakRespCode: 502,
			keycloakRespBody: []byte(`{"error": "server down"}`),
			wantRedirectURI:  "https://portal.local/callback",
			wantResponseCode: 500,
			wantResponseBody: `{"message": "Internal Server Error"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "missing request body",
			wantResponseCode: 400,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "missing `code` in request body, and request body have additional parameter",
			requestBody:      `{"token": "my-code"}`,
			wantResponseCode: 400,
			wantResponseBody: `{"message": "parameter \"token\" tidak didukung | parameter \"code\" harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "keycloak is unreachable",
			requestBody:      `{"code": "my-code"}`,
			wantResponseCode: 500,
			wantResponseBody: `{"message": "Internal Server Error"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keycloakHost := "http://localhost:8000"

			if tt.keycloakRespCode != 0 {
				keycloakMux := http.NewServeMux()
				keycloakMux.HandleFunc("POST /realms/nexus/protocol/openid-connect/token", func(w http.ResponseWriter, r *http.Request) {
					_ = r.ParseForm()

					assert.Equal(t,
						url.Values{
							"grant_type":    []string{"authorization_code"},
							"code":          []string{"my-code"},
							"redirect_uri":  []string{tt.wantRedirectURI},
							"client_id":     []string{"my-portal"},
							"client_secret": []string{"my-secret"},
						}, r.PostForm)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.keycloakRespCode)
					_, _ = w.Write(tt.keycloakRespBody)
				})
				keycloakSrv := httptest.NewServer(keycloakMux)
				defer keycloakSrv.Close()

				keycloakHost = keycloakSrv.URL
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/exchange-token", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			db := dbtest.New(t, dbmigrations.FS)
			_, err = db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			keycloak := config.Keycloak{
				Host:         keycloakHost,
				Realm:        "nexus",
				ClientID:     "my-portal",
				ClientSecret: "my-secret",
				KID:          "my-kid",
				RedirectURI:  "https://portal.local/callback",
			}
			repo := sqlc.New(db)
			RegisterRoutes(e, repo, keycloak, &http.Client{}, apitest.JwtPrivateKey, apitest.Keyfunc.Keyfunc)
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())

			actualRows, err := dbtest.QueryAll(db, `"user"`, "id, source")
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["last_login_at"] == "{last_login_at}" {
						assert.WithinDuration(t, time.Now(), row["last_login_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["last_login_at"] = row["last_login_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_refreshToken(t *testing.T) {
	t.Parallel()

	generateToken := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		tokenString, _ := token.SignedString(apitest.JwtPrivateKey)
		return tokenString
	}

	generateTokenWithKID := func(claims jwt.MapClaims) string {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		token.Header["kid"] = "my-kid"
		tokenString, _ := token.SignedString(apitest.JwtPrivateKey)
		return tokenString
	}

	tests := []struct {
		name             string
		dbData           string
		requestBody      string
		keycloakRespCode int
		keycloakRespBody []byte
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "success refresh token user without roles using keycloak_id",
			dbData: `
				insert into "user" (id, source, nip) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1c');
			`,
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "` + generateTokenWithKID(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "nip": "1c"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
		},
		{
			name: "success refresh token user without deleted roles using zimbra_id as priority",
			dbData: `
				insert into "user" (id, source, nip) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1c'),
					('00000000-0000-0000-0000-000000000002', 'zimbra', '1b');
				insert into role (id, service, nama, deleted_at) values
					(1, 'portal', 'admin', '2000-01-01'),
					(2, 'portal', 'admin', null);
				insert into user_role (nip, role_id, deleted_at) values
					('1c', 2, null),
					('1b', 1, null),
					('1b', 2, '2000-01-01');
			`,
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "` + generateTokenWithKID(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002", "nip": "1b"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
		},
		{
			name: "success refresh token user with multiple roles order by updated_at fallback using keycloak_id when user zimbra is not found",
			dbData: `
				insert into "user" (id, source, nip) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1c');
				insert into role (id, service, nama) values
					(1, 'portal', 'admin'),
					(2, 'portal', 'pegawai'),
					(3, 'kepegawaian', 'admin'),
					(4, 'kepegawaian', 'pegawai'),
					(5, 'kepegawaian', 'guest');
				insert into user_role (nip, role_id, updated_at, deleted_at) values
					('1c', 1, '2000-01-02', null),
					('1c', 2, '2000-01-01', null),
					('1c', 3, '2000-01-01', null),
					('1c', 4, '2000-01-02', null),
					('1c', 5, '2000-01-03', '2000-01-03');
			`,
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "` + generateTokenWithKID(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002", "nip": "1c", "roles": map[string]string{"portal": "admin", "kepegawaian": "pegawai"}}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
		},
		{
			name: "error refresh token, user keycloak or zimbra not found",
			dbData: `
				insert into "user" (id, source, nip, deleted_at) values
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1a', '2000-02-02'),
					('00000000-0000-0000-0000-000000000002', 'zimbra', '1b', '2000-01-01'),
					('00000000-0000-0000-0000-000000000002', 'keycloak', '1a', null),
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1b', null);
			`,
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "00000000-0000-0000-0000-000000000001", "zimbra_id": "00000000-0000-0000-0000-000000000002"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantResponseCode: 422,
			wantResponseBody: `{"message": "user tidak ditemukan"}`,
		},
		{
			name:             "error refresh token, user keycloak or zimbra is not uuid",
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "` + generateToken(jwt.MapClaims{"sub": "1", "zimbra_id": "2"}) + `",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantResponseCode: 422,
			wantResponseBody: `{"message": "user tidak ditemukan"}`,
		},
		{
			name:             "error refresh token, invalid access_token from keycloak",
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 200,
			keycloakRespBody: []byte(`{
					"access_token": "foo",
					"expires_in": 60,
					"id_token": "bar",
					"not-before-policy": 0,
					"refresh_expires_in": 300,
					"refresh_token": "baz",
					"scope": "openid",
					"session_state": "state",
					"token_type": "Bearer"
				}`),
			wantResponseCode: 500,
			wantResponseBody: `{"message": "Internal Server Error"}`,
		},
		{
			name:             "error refresh token, keycloak return status 4xx",
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 409,
			keycloakRespBody: []byte(`{"error": "invalid code"}`),
			wantResponseCode: 409,
			wantResponseBody: `{"error": "invalid code"}`,
		},
		{
			name:             "error refresh token, keycloak return status 5xx",
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 502,
			keycloakRespBody: []byte(`{"error": "server down"}`),
			wantResponseCode: 500,
			wantResponseBody: `{"message": "Internal Server Error"}`,
		},
		{
			name:             "missing request body",
			wantResponseCode: 400,
			wantResponseBody: `{"message": "request body harus diisi"}`,
		},
		{
			name:             "missing `refresh_token` in request body, and request body have additional parameter",
			requestBody:      `{"token": "my-code"}`,
			wantResponseCode: 400,
			wantResponseBody: `{"message": "parameter \"token\" tidak didukung | parameter \"refresh_token\" harus diisi"}`,
		},
		{
			name:             "keycloak is unreachable",
			requestBody:      `{"refresh_token": "my-code"}`,
			wantResponseCode: 500,
			wantResponseBody: `{"message": "Internal Server Error"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			keycloakHost := "http://localhost:8000"

			if tt.keycloakRespCode != 0 {
				keycloakMux := http.NewServeMux()
				keycloakMux.HandleFunc("POST /realms/nexus/protocol/openid-connect/token", func(w http.ResponseWriter, r *http.Request) {
					_ = r.ParseForm()

					assert.Equal(t,
						url.Values{
							"grant_type":    []string{"refresh_token"},
							"refresh_token": []string{"my-code"},
							"client_id":     []string{"my-portal"},
							"client_secret": []string{"my-secret"},
						}, r.PostForm)

					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(tt.keycloakRespCode)
					_, _ = w.Write(tt.keycloakRespBody)
				})
				keycloakSrv := httptest.NewServer(keycloakMux)
				defer keycloakSrv.Close()

				keycloakHost = keycloakSrv.URL
			}

			req := httptest.NewRequest(http.MethodPost, "/v1/auth/refresh-token", strings.NewReader(tt.requestBody))
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			db := dbtest.New(t, dbmigrations.FS)
			_, err = db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			keycloak := config.Keycloak{
				Host:         keycloakHost,
				Realm:        "nexus",
				ClientID:     "my-portal",
				ClientSecret: "my-secret",
				KID:          "my-kid",
			}
			repo := sqlc.New(db)
			RegisterRoutes(e, repo, keycloak, &http.Client{}, apitest.JwtPrivateKey, apitest.Keyfunc.Keyfunc)
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())

			// should not update last_login_at
			actualRows, err := dbtest.QueryAll(db, `"user"`, "id, source")
			require.NoError(t, err)
			for _, row := range actualRows {
				lastLoginAt, ok := row["last_login_at"]
				assert.True(t, ok)
				assert.Nil(t, lastLoginAt)
			}
		})
	}
}
