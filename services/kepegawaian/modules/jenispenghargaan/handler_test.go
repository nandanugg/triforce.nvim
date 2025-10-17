package jenispenghargaan

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_penghargaan (id, nama, deleted_at)
		VALUES
		(1, 'Jenis Penghargaan 1', NULL),
		(2, 'Jenis Penghargaan 2', NULL),
		(3, 'Jenis Penghargaan 3', '2023-02-20');
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("198765432100001")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: tanpa parameter apapun",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 1,
						"nama": "Jenis Penghargaan 1"
					},
					{
						"id": 2,
						"nama": "Jenis Penghargaan 2"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 2
				}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 2,
						"nama": "Jenis Penghargaan 2"
					}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 2
				}
			}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/jenis-penghargaan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			sqlc := sqlc.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			RegisterRoutes(e, sqlc, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetJenisPenghargaan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_penghargaan (id, nama, created_at, updated_at, deleted_at)
		VALUES
			(1, 'Jenis Penghargaan 1', now(), now(), NULL),
			(2, 'Jenis Penghargaan 2', now(), now(), NULL),
			(3, 'Jenis Penghargaan 3', now(), now(), now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("111")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get jenis penghargaan",
			id:               "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"nama": "Jenis Penghargaan 1"
				}
			}`,
		},
		{
			name:             "ok: get another jenis penghargaan",
			id:               "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jenis Penghargaan 2"
				}
			}`,
		},
		{
			name:             "error: jenis penghargaan not found",
			id:               "999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jenis penghargaan deleted",
			id:               "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: auth header tidak valid",
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer invalid-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jenis-penghargaan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			sqlc := sqlc.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			RegisterRoutes(e, sqlc, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateJenisPenghargaan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_penghargaan (nama, created_at, updated_at, deleted_at)
		VALUES
			('Jenis Penghargaan 1', now(), now(), NULL),
			('Jenis Penghargaan 2', now(), now(), NULL);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: create jenis penghargaan with required field",
			requestBody: `{
				"nama": "Jenis Penghargaan 3"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 3,
					"nama": "Jenis Penghargaan 3"
				}
			}`,
		},
		{
			name:        "error: missing required field nama",
			requestBody: `{}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name: "error: auth header tidak valid",
			requestBody: `{
				"nama": "Penghargaan Percobaan"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jenis-penghargaan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			sqlc := sqlc.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, sqlc, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateJenisPenghargaan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_penghargaan (
			id, nama, created_at, updated_at, deleted_at
		) VALUES
		(1, 'Jenis Penghargaan 1', now(), now(), NULL),
		(2, 'Jenis Penghargaan 2', now(), now(), NULL),
		(3, 'Jenis Penghargaan 3', now(), now(), now()),
		(4, 'Jenis Penghargaan 4', now(), now(), NULL);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: update existing jenis penghargaan",
			id:   "2",
			requestBody: `{
				"nama": "Jenis Penghargaan 2 Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jenis Penghargaan 2 Diperbarui"
				}
			}`,
		},
		{
			name: "error: update not found",
			id:   "99",
			requestBody: `{
				"nama": "Tidak Ada"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: update deleted record",
			id:   "3",
			requestBody: `{
				"nama": "Tidak Boleh Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: auth header tidak valid",
			id:          "1",
			requestBody: `{"nama": "Jenis Penghargaan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer invalid-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/jenis-penghargaan/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			r := sqlc.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteJenisPenghargaan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_penghargaan (id, nama, created_at, updated_at, deleted_at) VALUES
		(1, 'Jenis Penghargaan 1', now(), now(), NULL),
		(2, 'Jenis Penghargaan 2', now(), now(), NULL),
		(3, 'Jenis Penghargaan 3', now(), now(), now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: delete jenis penghargaan",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name: "error: delete not found jenis penghargaan",
			id:   "999",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: delete already deleted jenis penghargaan",
			id:   "3",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: auth header tidak valid",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/jenis-penghargaan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			r := sqlc.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			} else {
				assert.Empty(t, rec.Body.String())
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
