package jenishukuman

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
		INSERT INTO ref_jenis_hukuman (id, nama, tingkat_hukuman, deleted_at)
		VALUES
		(1, 'Jenis Hukuman 1', 'B', NULL),
		(2, 'Jenis Hukuman 2', 'S', NULL),
		(3, 'Jenis Hukuman 3', 'R', '2023-02-20');
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
						"nama": "Jenis Hukuman 1",
						"tingkat": "B"
					},
					{
						"id": 2,
						"nama": "Jenis Hukuman 2",
						"tingkat": "S"
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
						"nama": "Jenis Hukuman 2",
						"tingkat": "S"
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

			req := httptest.NewRequest(http.MethodGet, "/v1/jenis-hukuman", nil)
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

func Test_handler_adminGetJenisHukuman(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman (id, nama, tingkat_hukuman, deleted_at)
		VALUES
		(1, 'Jenis Hukuman 1', 'B', NULL),
		(2, 'Jenis Hukuman 2', 'S', NULL),
		(3, 'Jenis Hukuman 3', 'R', '2023-02-20');
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	sqlc := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
	RegisterRoutes(e, sqlc, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("111")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get jenis hukuman",
			id:               "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"nama": "Jenis Hukuman 1",
					"tingkat": "B"
				}
			}`,
		},
		{
			name:             "ok: get another jenis hukuman",
			id:               "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jenis Hukuman 2",
					"tingkat": "S"
				}
			}`,
		},
		{
			name:             "error: jenis hukuman not found",
			id:               "999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jenis hukuman deleted",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jenis-hukuman/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateJenisHukuman(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman (nama, tingkat_hukuman, deleted_at)
		VALUES
		('Jenis Hukuman 1', 'B', NULL),
		('Jenis Hukuman 2', 'S', NULL);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	sqlc := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, sqlc, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: create jenis hukuman with required field",
			requestBody: `{
				"nama": "Jenis Hukuman 3",
				"tingkat": "R"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 3,
					"nama": "Jenis Hukuman 3",
					"tingkat": "R"
				}
			}`,
		},
		{
			name:        "error: missing required field nama",
			requestBody: `{"tingkat": "R"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name:        "error: missing required field tingkat",
			requestBody: `{"nama": "Jenis Hukuman 3"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"tingkat\" harus diisi"}`,
		},
		{
			name: "error: auth header tidak valid",
			requestBody: `{
				"nama": "Hukuman Percobaan",
				"tingkat": "R"
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jenis-hukuman", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateJenisHukuman(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman (
			id, nama, tingkat_hukuman, deleted_at
		) VALUES
		(1, 'Jenis Hukuman 1', 'B', NULL),
		(2, 'Jenis Hukuman 2', 'S', NULL),
		(3, 'Jenis Hukuman 3', 'R', now()),
		(4, 'Jenis Hukuman 4', 'R', NULL);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	r := sqlc.New(pgxconn)
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
			name: "ok: update existing jenis hukuman",
			id:   "2",
			requestBody: `{
				"nama": "Jenis Hukuman 2 Diperbarui",
				"tingkat": "R"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "Jenis Hukuman 2 Diperbarui",
					"tingkat": "R"
				}
			}`,
		},
		{
			name: "error: missing required field nama",
			id:   "2",
			requestBody: `{
				"tingkat": "R"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name: "error: missing required field tingkat",
			id:   "2",
			requestBody: `{
				"nama": "Jenis Hukuman 2 Diperbarui"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"tingkat\" harus diisi"}`,
		},
		{
			name: "error: update not found",
			id:   "99",
			requestBody: `{
				"nama": "Tidak Ada",
				"tingkat": "R"
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
				"nama": "Tidak Boleh Diperbarui",
				"tingkat": "R"
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
			requestBody: `{"nama": "Jenis Hukuman 1 Diperbarui"}`,
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/jenis-hukuman/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteJenisHukuman(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman (id, nama, created_at, updated_at, deleted_at) VALUES
		(1, 'Jenis Hukuman 1', now(), now(), NULL),
		(2, 'Jenis Hukuman 2', now(), now(), NULL),
		(3, 'Jenis Hukuman 3', now(), now(), now());

		INSERT INTO riwayat_hukuman_disiplin  (jenis_hukuman_id, deleted_at) VALUES
		(1, now()),
		(2, null),
		(3, now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	r := sqlc.New(pgxconn)
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
			name: "ok: delete jenis hukuman",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name: "error: delete jenis hukuman referenced",
			id:   "2",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "jenis hukuman masih digunakan oleh riwayat hukuman disiplin"}`,
		},
		{
			name: "error: delete not found jenis hukuman",
			id:   "999",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: delete already deleted jenis hukuman",
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

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/jenis-hukuman/"+tt.id, nil)
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
