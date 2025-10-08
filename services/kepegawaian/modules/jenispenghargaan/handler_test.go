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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
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

	tests := []struct {
		name             string
		dbData           string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: tanpa parameter apapun",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
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
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
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
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/jenis-penghargaan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			sqlc := sqlc.New(pgxconn)
			RegisterRoutes(e, sqlc, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get jenis penghargaan",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
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
			dbData:           dbData,
			id:               "2",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
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
			dbData:           dbData,
			id:               "999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jenis penghargaan deleted",
			dbData:           dbData,
			id:               "3",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "999")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer invalid-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jenis-penghargaan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			sqlc := sqlc.New(pgxconn)
			RegisterRoutes(e, sqlc, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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

	tests := []struct {
		name             string
		dbData           string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:   "ok: create jenis penghargaan with required field",
			dbData: dbData,
			requestBody: `{
				"nama": "Jenis Penghargaan 3"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
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
			dbData:      dbData,
			requestBody: `{}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name:   "error: auth header tidak valid",
			dbData: dbData,
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
		{
			name:   "error: user is not an admin",
			dbData: dbData,
			requestBody: `{
				"nama": "Penghargaan Percobaan"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jenis-penghargaan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			sqlc := sqlc.New(pgxconn)
			RegisterRoutes(e, sqlc, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:   "ok: update existing jenis penghargaan",
			dbData: dbData,
			id:     "2",
			requestBody: `{
				"nama": "Jenis Penghargaan 2 Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{
					apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin),
				},
				"Content-Type": []string{"application/json"},
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
			name:   "error: update not found",
			dbData: dbData,
			id:     "99",
			requestBody: `{
				"nama": "Tidak Ada"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{
					apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin),
				},
				"Content-Type": []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: update deleted record",
			dbData: dbData,
			id:     "3",
			requestBody: `{
				"nama": "Tidak Boleh Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{
					apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin),
				},
				"Content-Type": []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: auth header tidak valid",
			dbData:      dbData,
			id:          "1",
			requestBody: `{"nama": "Jenis Penghargaan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer invalid-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:        "error: user is not an admin",
			dbData:      dbData,
			id:          "1",
			requestBody: `{"nama": "Jenis Penghargaan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{
					apitest.GenerateAuthHeader(config.Service, "987654321"),
				},
				"Content-Type": []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/jenis-penghargaan/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			r := sqlc.New(pgxconn)
			RegisterRoutes(e, r, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:   "ok: delete jenis penghargaan",
			dbData: dbData,
			id:     "1",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:   "error: delete not found jenis penghargaan",
			dbData: dbData,
			id:     "999",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: delete already deleted jenis penghargaan",
			dbData: dbData,
			id:     "3",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: auth header tidak valid",
			dbData: dbData,
			id:     "1",
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:   "error: user is not an admin",
			dbData: dbData,
			id:     "1",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)

			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/jenis-penghargaan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			r := sqlc.New(pgxconn)
			RegisterRoutes(e, r, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
