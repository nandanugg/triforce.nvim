package auth

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func Test_handler_login(t *testing.T) {
	t.Parallel()

	t.Run("success redirect", func(t *testing.T) {
		t.Parallel()

		req := httptest.NewRequest(http.MethodGet, "/v1/auth/login", nil)
		rec := httptest.NewRecorder()

		e, err := api.NewEchoServer(docs.OpenAPIBytes)
		require.NoError(t, err)

		keycloak := config.Keycloak{
			Host:        "https://auth.portal.local",
			Realm:       "nexus",
			ClientID:    "my-portal",
			RedirectURI: "https://portal.local/callback",
		}
		RegisterRoutes(e, keycloak, nil)
		e.ServeHTTP(rec, req)

		assert.Equal(t, 302, rec.Code)
		assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		assert.Equal(t, "https://auth.portal.local/realms/nexus/protocol/openid-connect/auth?client_id=my-portal&prompt=login&redirect_uri=https%3A%2F%2Fportal.local%2Fcallback&response_type=code&scope=openid", rec.Header().Get("Location"))
		assert.Empty(t, rec.Body)
	})
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
			name:             "missing id_token_hint",
			wantResponseCode: 400,
			wantResponseBody: `{"message": "parameter \"id_token_hint\" harus diisi"}`,
		},
		{
			name: "success redirect",
			requestQuery: url.Values{
				"id_token_hint": []string{"e9d74da4-a4e7-4865-8a4b-601ad1bca900"},
			},
			wantResponseCode: 302,
			wantRedirect:     `https://auth.portal.local/realms/nexus/protocol/openid-connect/logout?id_token_hint=e9d74da4-a4e7-4865-8a4b-601ad1bca900&post_logout_redirect_uri=https%3A%2F%2Fportal.local%2F`,
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
				Host:                  "https://auth.portal.local",
				Realm:                 "nexus",
				PostLogoutRedirectURI: "https://portal.local/",
			}
			RegisterRoutes(e, keycloak, nil)
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

	tests := []struct {
		name             string
		requestBody      string
		keycloakRespCode int
		keycloakRespBody []byte
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "success exchange token",
			requestBody:      `{"code": "my-code"}`,
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
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "foo",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
		},
		{
			name:             "error exchange token, keycloak return status 4xx",
			requestBody:      `{"code": "my-code"}`,
			keycloakRespCode: 422,
			keycloakRespBody: []byte(`{"error": "invalid code"}`),
			wantResponseCode: 422,
			wantResponseBody: `{"error": "invalid code"}`,
		},
		{
			name:             "error exchange token, keycloak return status 5xx",
			requestBody:      `{"code": "my-code"}`,
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
			name:             "missing `code` in request body, and request body have additional parameter",
			requestBody:      `{"token": "my-code"}`,
			wantResponseCode: 400,
			wantResponseBody: `{"message": "parameter \"token\" tidak didukung | parameter \"code\" harus diisi"}`,
		},
		{
			name:             "keycloak is unreachable",
			requestBody:      `{"code": "my-code"}`,
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
							"grant_type":    []string{"authorization_code"},
							"code":          []string{"my-code"},
							"redirect_uri":  []string{"https://portal.local/callback"},
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

			keycloak := config.Keycloak{
				Host:         keycloakHost,
				Realm:        "nexus",
				RedirectURI:  "https://portal.local/callback",
				ClientID:     "my-portal",
				ClientSecret: "my-secret",
			}
			RegisterRoutes(e, keycloak, &http.Client{})
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_refreshToken(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		requestBody      string
		keycloakRespCode int
		keycloakRespBody []byte
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "success refresh token",
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
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": {
					"access_token": "foo",
					"expires_in": 60,
					"id_token": "bar",
					"refresh_token": "baz",
					"refresh_expires_in": 300
				}
			}`,
		},
		{
			name:             "error refresh token, keycloak return status 4xx",
			requestBody:      `{"refresh_token": "my-code"}`,
			keycloakRespCode: 422,
			keycloakRespBody: []byte(`{"error": "invalid code"}`),
			wantResponseCode: 422,
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

			keycloak := config.Keycloak{
				Host:         keycloakHost,
				Realm:        "nexus",
				ClientID:     "my-portal",
				ClientSecret: "my-secret",
			}
			RegisterRoutes(e, keycloak, &http.Client{})
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}
