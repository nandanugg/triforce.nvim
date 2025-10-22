package jeniskenaikanpangkat

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
		insert into ref_jenis_kenaikan_pangkat
		("id", "nama") values
		(1, 'Kenaikan Pangkat Reguler'),
		(2, 'Kenaikan Pangkat Pilihan'),
		(3, 'Kenaikan Pangkat Luar Biasa'),
		(4, 'Kenaikan Pangkat Pengabdian'),
		(5, 'Kenaikan Pangkat Penyesuaian'),
		(6, 'Kenaikan Pangkat Penyesuaian Ijazah'),
		(7, 'Kenaikan Pangkat Penyesuaian Jabatan'),
		(8, 'Kenaikan Pangkat Penyesuaian Golongan'),
		(9, 'Kenaikan Pangkat Penyesuaian Masa Kerja'),
		(10, 'Kenaikan Pangkat Penyesuaian Kinerja'),
		(11, 'Kenaikan Pangkat Penyesuaian Diklat'),
		(12, 'Kenaikan Pangkat Penyesuaian Sertifikasi'),
		(13, 'Kenaikan Pangkat Penyesuaian Penghargaan'),
		(14, 'Kenaikan Pangkat Penyesuaian Penugasan'),
		(15, 'Kenaikan Pangkat Penyesuaian Mutasi'),
		(16, 'Kenaikan Pangkat Penyesuaian Promosi'),
		(17, 'Kenaikan Pangkat Penyesuaian Khusus');
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
					{"id": 1, "nama": "Kenaikan Pangkat Reguler"},
					{"id": 2, "nama": "Kenaikan Pangkat Pilihan"},
					{"id": 3, "nama": "Kenaikan Pangkat Luar Biasa"},
					{"id": 4, "nama": "Kenaikan Pangkat Pengabdian"},
					{"id": 5, "nama": "Kenaikan Pangkat Penyesuaian"},
					{"id": 6, "nama": "Kenaikan Pangkat Penyesuaian Ijazah"},
					{"id": 7, "nama": "Kenaikan Pangkat Penyesuaian Jabatan"},
					{"id": 8, "nama": "Kenaikan Pangkat Penyesuaian Golongan"},
					{"id": 9, "nama": "Kenaikan Pangkat Penyesuaian Masa Kerja"},
					{"id": 10, "nama": "Kenaikan Pangkat Penyesuaian Kinerja"}
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
					{"id": 1, "nama": "Kenaikan Pangkat Reguler"},
					{"id": 2, "nama": "Kenaikan Pangkat Pilihan"},
					{"id": 3, "nama": "Kenaikan Pangkat Luar Biasa"},
					{"id": 4, "nama": "Kenaikan Pangkat Pengabdian"},
					{"id": 5, "nama": "Kenaikan Pangkat Penyesuaian"}
				],
				"meta": {
					"limit": 5,
					"offset": 0,
					"total": 17
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
					{"id": 6, "nama": "Kenaikan Pangkat Penyesuaian Ijazah"},
					{"id": 7, "nama": "Kenaikan Pangkat Penyesuaian Jabatan"},
					{"id": 8, "nama": "Kenaikan Pangkat Penyesuaian Golongan"}
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

			req := httptest.NewRequest(http.MethodGet, "/v1/jenis-kenaikan-pangkat", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetJenisKenaikanPangkat(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_kenaikan_pangkat (id, nama, created_at, updated_at, deleted_at)
		VALUES
			(1, 'Kenaikan Reguler', now(), now(), NULL),
			(2, 'Kenaikan Pilihan', now(), now(), NULL),
			(3, 'Kenaikan Luar Biasa', now(), now(), now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("111")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get jenis kenaikan pangkat",
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
			name:             "ok: get another jenis kenaikan pangkat",
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
			name:             "error: jenis kenaikan pangkat not found",
			id:               "999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jenis kenaikan pangkat deleted",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jenis-kenaikan-pangkat/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateJenisKenaikanPangkat(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_kenaikan_pangkat (nama, created_at, updated_at, deleted_at)
		VALUES
			('Kenaikan Reguler', now(), now(), NULL),
			('Kenaikan Pilihan', now(), now(), NULL);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: create jenis kenaikan pangkat with required field",
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
				"nama": "Kenaikan Pangkat Percobaan"
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jenis-kenaikan-pangkat", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateJenisKenaikanPangkat(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_kenaikan_pangkat (
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

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	r := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			name: "ok: update existing jenis kenaikan pangkat",
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/jenis-kenaikan-pangkat/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteJenisKenaikanPangkat(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_kenaikan_pangkat (id, nama, created_at, updated_at, deleted_at) VALUES
		(1, 'Reguler', now(), now(), NULL),
		(2, 'Pilihan', now(), now(), NULL),
		(3, 'Luar Biasa', now(), now(), now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	r := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, r, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: delete jenis kenaikan pangkat",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name: "error: delete not found jenis kenaikan pangkat",
			id:   "999",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: delete already deleted jenis kenaikan pangkat",
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

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/jenis-kenaikan-pangkat/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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
