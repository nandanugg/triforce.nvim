package golongan_test

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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/golongan"
)

func Test_handler_ListRefGolongan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_golongan
		("id", "nama", "nama_pangkat") values
		(1, 'I/a', 'Juru Muda'),
		(2, 'I/b', 'Juru Muda Tingkat I'),
		(3, 'I/c', 'Juru'),
		(4, 'I/d', 'Juru Tingkat I'),
		(5, 'II/a', 'Pengatur Muda'),
		(6, 'II/b', 'Pengatur Muda Tingkat I'),
		(7, 'II/c', 'Pengatur'),
		(8, 'II/d', 'Pengatur Tingkat I'),
		(9, 'III/a', 'Penata Muda'),
		(10, 'III/b', 'Penata Muda Tingkat I'),
		(11, 'III/c', 'Penata'),
		(12, 'III/d', 'Penata Tingkat I'),
		(13, 'IV/a', 'Pembina'),
		(14, 'IV/b', 'Pembina Tingkat I'),
		(15, 'IV/c', 'Pembina Utama Muda'),
		(16, 'IV/d', 'Pembina Utama Madya'),
		(17, 'IV/e', 'Pembina Utama');
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
			name:             "ok: get data with default pagination",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "I/a", "nama_pangkat": "Juru Muda"},
					{"id": 2, "nama": "I/b", "nama_pangkat": "Juru Muda Tingkat I"},
					{"id": 3, "nama": "I/c", "nama_pangkat": "Juru"},
					{"id": 4, "nama": "I/d", "nama_pangkat": "Juru Tingkat I"},
					{"id": 5, "nama": "II/a", "nama_pangkat": "Pengatur Muda"},
					{"id": 6, "nama": "II/b", "nama_pangkat": "Pengatur Muda Tingkat I"},
					{"id": 7, "nama": "II/c", "nama_pangkat": "Pengatur"},
					{"id": 8, "nama": "II/d", "nama_pangkat": "Pengatur Tingkat I"},
					{"id": 9, "nama": "III/a", "nama_pangkat": "Penata Muda"},
					{"id": 10, "nama": "III/b", "nama_pangkat": "Penata Muda Tingkat I"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 17
				}
			}`,
		},
		{
			name:   "ok: with pagination limit 5",
			dbData: dbData,
			requestQuery: url.Values{
				"limit": []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "I/a", "nama_pangkat": "Juru Muda"},
					{"id": 2, "nama": "I/b", "nama_pangkat": "Juru Muda Tingkat I"},
					{"id": 3, "nama": "I/c", "nama_pangkat": "Juru"},
					{"id": 4, "nama": "I/d", "nama_pangkat": "Juru Tingkat I"},
					{"id": 5, "nama": "II/a", "nama_pangkat": "Pengatur Muda"}
				],
				"meta": {
					"limit": 5,
					"offset": 0,
					"total": 17
				}
			}`,
		},
		{
			name:   "ok: with pagination limit 3 offset 5",
			dbData: dbData,
			requestQuery: url.Values{
				"limit":  []string{"3"},
				"offset": []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 6, "nama": "II/b", "nama_pangkat": "Pengatur Muda Tingkat I"},
					{"id": 7, "nama": "II/c", "nama_pangkat": "Pengatur"},
					{"id": 8, "nama": "II/d", "nama_pangkat": "Pengatur Tingkat I"}
				],
				"meta": {
					"limit": 3,
					"offset": 5,
					"total": 17
				}
			}`,
		},
		{
			name:             "ok: empty data",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: missing auth header",
			dbData:           dbData,
			requestHeader:    http.Header{},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/golongan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			golongan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminListRefGolongan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_golongan
		("id", "nama", "nama_pangkat", "nama_2", "gol", "gol_pppk", "created_at", "updated_at", "deleted_at") values
		(1, 'I/a', 'Juru Muda', '1a', '1', '1a', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(2, 'I/b', 'Juru Muda Tingkat I', '1b', '1', '1b', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(3, 'I/c', 'Juru', '1c', '1', '1c', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(4, 'I/d', 'Juru Tingkat I', '1d', '1', '1d', '2025-01-01 00:00:00', '2025-01-01T00:00:00', now());
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
			name:             "ok: get data with default pagination",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "I/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"},
					{"id": 2, "nama": "I/b", "nama_pangkat": "Juru Muda Tingkat I", "nama_2": "1b", "gol": 1, "gol_pppk": "1b"},
					{"id": 3, "nama": "I/c", "nama_pangkat": "Juru", "nama_2": "1c", "gol": 1, "gol_pppk": "1c"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 3
				}
			}`,
		},
		{
			name:   "ok: with pagination limit 1 and offset 1",
			dbData: dbData,
			requestQuery: url.Values{
				"limit":  []string{"1"},
				"offset": []string{"1"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 2, "nama": "I/b", "nama_pangkat": "Juru Muda Tingkat I", "nama_2": "1b", "gol": 1, "gol_pppk": "1b"}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 3
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/golongan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			golongan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateRefGolongan(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name             string
		requestQuery     url.Values
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: create golongan",
			requestBody: `{"nama": "I/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id": 1, "nama": "I/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}}`,
		},
		{
			name:        "error: auth header tidak valid",
			requestBody: `{"nama": "I/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
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

			_, err := pgxconn.Exec(context.Background(), "")
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/golongan", strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			golongan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateRefGolongan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_golongan
		("id", "nama", "nama_pangkat", "nama_2", "gol", "gol_pppk", "created_at", "updated_at", "deleted_at") values
		(1, 'I/a', 'Juru Muda', '1a', '1', '1a', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(2, 'I/b', 'Juru Muda Tingkat I', '1b', '1', '1b', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(3, 'I/c', 'Juru', '1c', '1', '1c', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(4, 'I/d', 'Juru Tingkat I', '1d', '1', '1d', '2025-01-01 00:00:00', '2025-01-01T00:00:00', now());
	`

	tests := []struct {
		name             string
		id               string
		requestQuery     url.Values
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: update golongan",
			id:          "1",
			requestBody: `{"nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"id": 1, "nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}}`,
		},
		{
			name:        "error: data sudah dihapus",
			id:          "4",
			requestBody: `{"nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: id tidak valid",
			id:          "100",
			requestBody: `{"nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: auth header tidak valid",
			id:          "1",
			requestBody: `{"nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
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

			_, err := pgxconn.Exec(context.Background(), dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/golongan/"+tt.id, strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			golongan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteRefGolongan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_golongan
		("id", "nama", "nama_pangkat", "nama_2", "gol", "gol_pppk", "created_at", "updated_at", "deleted_at") values
		(1, 'I/a', 'Juru Muda', '1a', '1', '1a', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(2, 'I/b', 'Juru Muda Tingkat I', '1b', '1', '1b', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(3, 'I/c', 'Juru', '1c', '1', '1c', '2025-01-01T00:00:00', '2025-01-01T00:00:00', null),
		(4, 'I/d', 'Juru Tingkat I', '1d', '1', '1d', '2025-01-01 00:00:00', '2025-01-01T00:00:00', now());
	`

	tests := []struct {
		name             string
		id               string
		requestQuery     url.Values
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: delete golongan",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `{}`,
		},
		{
			name: "error: data sudah dihapus",
			id:   "4",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: id tidak valid",
			id:          "100",
			requestBody: `{"nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: auth header tidak valid",
			id:          "1",
			requestBody: `{"nama": "II/a", "nama_pangkat": "Juru Muda", "nama_2": "1a", "gol": 1, "gol_pppk": "1a"}`,
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

			_, err := pgxconn.Exec(context.Background(), dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/golongan/"+tt.id, strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			golongan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if rec.Body.String() != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			} else {
				assert.Empty(t, rec.Body.String())
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
