package jenissatker

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
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_satker
		("id", "nama") values
		(1, 'Satker Reguler'),
		(2, 'Satker Pilihan'),
		(3, 'Satker Luar Biasa'),
		(4, 'Satker Pengabdian'),
		(5, 'Satker Penyesuaian'),
		(6, 'Satker Penyesuaian Ijazah'),
		(7, 'Satker Penyesuaian Jabatan'),
		(8, 'Satker Penyesuaian Golongan'),
		(9, 'Satker Penyesuaian Masa Kerja'),
		(10, 'Satker Penyesuaian Kinerja'),
		(11, 'Satker Penyesuaian Diklat'),
		(12, 'Satker Penyesuaian Sertifikasi'),
		(13, 'Satker Penyesuaian Penghargaan'),
		(14, 'Satker Penyesuaian Penugasan'),
		(15, 'Satker Penyesuaian Mutasi'),
		(16, 'Satker Penyesuaian Promosi'),
		(17, 'Satker Penyesuaian Khusus');
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("41")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get data with default pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "Satker Reguler"},
					{"id": 2, "nama": "Satker Pilihan"},
					{"id": 3, "nama": "Satker Luar Biasa"},
					{"id": 4, "nama": "Satker Pengabdian"},
					{"id": 5, "nama": "Satker Penyesuaian"},
					{"id": 6, "nama": "Satker Penyesuaian Ijazah"},
					{"id": 7, "nama": "Satker Penyesuaian Jabatan"},
					{"id": 8, "nama": "Satker Penyesuaian Golongan"},
					{"id": 9, "nama": "Satker Penyesuaian Masa Kerja"},
					{"id": 10, "nama": "Satker Penyesuaian Kinerja"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 17
				}
			}`,
		},
		{
			name: "ok: with pagination limit 5",
			requestQuery: url.Values{
				"limit": []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "Satker Reguler"},
					{"id": 2, "nama": "Satker Pilihan"},
					{"id": 3, "nama": "Satker Luar Biasa"},
					{"id": 4, "nama": "Satker Pengabdian"},
					{"id": 5, "nama": "Satker Penyesuaian"}
				],
				"meta": {
					"limit": 5,
					"offset": 0,
					"total": 17
				}
			}`,
		},
		{
			name: "ok: with nama filter",
			requestQuery: url.Values{
				"nama": []string{"Pengabdian"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 4, "nama": "Satker Pengabdian"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name: "ok: with pagination limit 3 offset 5",
			requestQuery: url.Values{
				"limit":  []string{"3"},
				"offset": []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 6, "nama": "Satker Penyesuaian Ijazah"},
					{"id": 7, "nama": "Satker Penyesuaian Jabatan"},
					{"id": 8, "nama": "Satker Penyesuaian Golongan"}
				],
				"meta": {
					"limit": 3,
					"offset": 5,
					"total": 17
				}
			}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: missing auth header",
			requestHeader:    http.Header{},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/jenis-satker", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repo.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetJenisSatker(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_satker (id, nama, created_at, updated_at, deleted_at)
		VALUES
			(1, 'Kenaikan Reguler', now(), now(), NULL),
			(2, 'Kenaikan Pilihan', now(), now(), NULL),
			(3, 'Kenaikan Luar Biasa', now(), now(), now());
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
			name:             "ok: get jenis satker",
			id:               "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"nama": "Kenaikan Reguler"
				}
			}`,
		},
		{
			name:             "ok: get another jenis satker",
			id:               "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Kenaikan Pilihan"
				}
			}`,
		},
		{
			name:             "error: jenis satker not found",
			id:               "999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jenis satker deleted",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jenis-satker/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repo.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateJenisSatker(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_satker (nama, created_at, updated_at, deleted_at)
		VALUES
			('Kenaikan Reguler', now(), now(), NULL),
			('Kenaikan Pilihan', now(), now(), NULL);
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
			name: "ok: create jenis satker with required field",
			requestBody: `{
				"nama": "Kenaikan Luar Biasa"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 3,
					"nama": "Kenaikan Luar Biasa"
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
				"nama": "Satker Percobaan"
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jenis-satker", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repo.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateJenisSatker(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_satker (
			id, nama, created_at, updated_at, deleted_at
		) VALUES
		(1, 'Reguler', now(), now(), NULL),
		(2, 'Jabatan Struktural', now(), now(), NULL),
		(3, 'Jabatan Fungsional', now(), now(), now()),
		(4, 'Penyesuaian Ijazah', now(), now(), NULL);
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
			name: "ok: update existing jenis satker",
			id:   "2",
			requestBody: `{
				"nama": "Jabatan Struktural Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jabatan Struktural Diperbarui"
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
			requestBody: `{"nama": "Reguler Diperbarui"}`,
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/jenis-satker/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			r := repo.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteJenisSatker(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_satker (id, nama, created_at, updated_at, deleted_at) VALUES
		(1, 'Reguler', now(), now(), NULL),
		(2, 'Pilihan', now(), now(), NULL),
		(3, 'Luar Biasa', now(), now(), now());
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
			name: "ok: delete jenis satker",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name: "error: delete not found jenis satker",
			id:   "999",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: delete already deleted jenis satker",
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

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/jenis-satker/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			r := repo.New(pgxconn)
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
