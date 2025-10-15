package jenisjabatan_test

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
	dbrepository "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jenisjabatan"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_jabatan
			("id", "nama") values
			(1,  'a'),
			(2,  'c'),
			(3,  'b');
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
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "a"},
					{"id": 2, "nama": "c"},
					{"id": 3, "nama": "b"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok with limit 2",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"2"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "a"},
					{"id": 2, "nama": "c"}
				],
				"meta": {"limit": 2, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok with limit 2 and offset 1",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 2, "nama": "c"},
					{"id": 3, "nama": "b"}
				],
				"meta": {"limit": 2, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok with empty data",
			dbData:           ``,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
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

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/jenis-jabatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepository.New(db)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			jenisjabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetJenisJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_jabatan (id, nama, created_at, updated_at, deleted_at)
		VALUES
			(1, 'Jenis Jabatan 1', now(), now(), NULL),
			(2, 'Jenis Jabatan 2', now(), now(), NULL),
			(3, 'Jenis Jabatan 3', now(), now(), now());
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
			name:             "ok: get jenis jabatan",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("111")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"nama": "Jenis Jabatan 1"
				}
			}`,
		},
		{
			name:             "ok: get another jenis jabatan",
			dbData:           dbData,
			id:               "2",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("111")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jenis Jabatan 2"
				}
			}`,
		},
		{
			name:             "error: jenis jabatan not found",
			dbData:           dbData,
			id:               "999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("111")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jenis jabatan deleted",
			dbData:           dbData,
			id:               "3",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("111")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jenis-jabatan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			jenisjabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateJenisJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_jabatan (nama, created_at, updated_at, deleted_at)
		VALUES
			('Jenis Jabatan 1', now(), now(), NULL),
			('Jenis Jabatan 2', now(), now(), NULL);
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
			name:   "ok: create jenis jabatan with required field",
			dbData: dbData,
			requestBody: `{
				"nama": "Jenis Jabatan 3"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 3,
					"nama": "Jenis Jabatan 3"
				}
			}`,
		},
		{
			name:        "error: missing required field nama",
			dbData:      dbData,
			requestBody: `{}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name:   "error: auth header tidak valid",
			dbData: dbData,
			requestBody: `{
				"nama": "Jabatan Percobaan"
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

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jenis-jabatan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			jenisjabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateJenisJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_jabatan (
			id, nama, created_at, updated_at, deleted_at
		) VALUES
		(1, 'Jenis Jabatan 1', now(), now(), NULL),
		(2, 'Jenis Jabatan 2', now(), now(), NULL),
		(3, 'Jenis Jabatan 3', now(), now(), now()),
		(4, 'Jenis Jabatan 4', now(), now(), NULL);
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
			name:   "ok: update existing jenis jabatan",
			dbData: dbData,
			id:     "2",
			requestBody: `{
				"nama": "Jenis Jabatan 2 Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{
					apitest.GenerateAuthHeader("123456789"),
				},
				"Content-Type": []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jenis Jabatan 2 Diperbarui"
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
					apitest.GenerateAuthHeader("123456789"),
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
					apitest.GenerateAuthHeader("123456789"),
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
			requestBody: `{"nama": "Jenis Jabatan 1 Diperbarui"}`,
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
			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/jenis-jabatan/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			r := dbrepository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			jenisjabatan.RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteJenisJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_jabatan (id, nama, created_at, updated_at, deleted_at) VALUES
		(1, 'Jenis Jabatan 1', now(), now(), NULL),
		(2, 'Jenis Jabatan 2', now(), now(), NULL),
		(3, 'Jenis Jabatan 3', now(), now(), now());
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
			name:   "ok: delete jenis jabatan",
			dbData: dbData,
			id:     "1",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:   "error: delete not found jenis jabatan",
			dbData: dbData,
			id:     "999",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: delete already deleted jenis jabatan",
			dbData: dbData,
			id:     "3",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)

			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/jenis-jabatan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			r := dbrepository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			jenisjabatan.RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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
