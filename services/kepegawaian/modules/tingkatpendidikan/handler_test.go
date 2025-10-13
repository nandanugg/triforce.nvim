package tingkatpendidikan

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
		INSERT INTO ref_tingkat_pendidikan (id, nama, deleted_at)
		VALUES
		(1, 'Jenis Pendidikan 1', NULL),
		(2, 'Jenis Pendidikan 2', NULL),
		(3, 'Jenis Pendidikan 3', '2023-02-20');
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
						"nama": "Jenis Pendidikan 1"
					},
					{
						"id": 2,
						"nama": "Jenis Pendidikan 2"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 2,
						"nama": "Jenis Pendidikan 2"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 2}
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

			req := httptest.NewRequest(http.MethodGet, "/v1/tingkat-pendidikan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_golongan (id, nama, nama_pangkat, deleted_at)
		VALUES
		(1, 'Golongan 1', 'Pangkat 1', NULL),
		(2, 'Golongan 2', 'Pangkat 2', NULL),
		(3, 'Golongan 3', 'Pangkat 3', '2023-02-20');

		INSERT INTO ref_tingkat_pendidikan (id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat, deleted_at)
		VALUES
		(1, 'Jenis Pendidikan 1', 'J1', 1, NULL, 1, NULL),
		(2, 'Jenis Pendidikan 2', 'J2', 1, 2, 2, NULL),
		(3, 'Jenis Pendidikan 3', 'J3', 1, 3, 3, '2023-02-20'),
		(4, 'Jenis Pendidikan 4', NULL, NULL, NULL, NULL, NULL);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 1,
						"nama": "Jenis Pendidikan 1",
						"abbreviation": "J1",
						"golongan_id": 1,
						"golongan_awal_id": null,
						"tingkat": 1,
						"nama_golongan": "Golongan 1",
						"nama_golongan_awal": null
					},
					{
						"id": 2,
						"nama": "Jenis Pendidikan 2",
						"abbreviation": "J2",
						"golongan_id": 1,
						"golongan_awal_id": 2,
						"tingkat": 2,
						"nama_golongan": "Golongan 1",
						"nama_golongan_awal": "Golongan 2"
					},
					{
						"id": 4,
						"nama": "Jenis Pendidikan 4",
						"abbreviation": null,
						"golongan_id": null,
						"golongan_awal_id": null,
						"tingkat": null,
						"nama_golongan": null,
						"nama_golongan_awal": null
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 2,
						"nama": "Jenis Pendidikan 2",
						"abbreviation": "J2",
						"golongan_id": 1,
						"golongan_awal_id": 2,
						"tingkat": 2,
						"nama_golongan": "Golongan 1",
						"nama_golongan_awal": "Golongan 2"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/tingkat-pendidikan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_getAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_tingkat_pendidikan (id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat, deleted_at)
		VALUES
		(1, 'Jenis Pendidikan 1', 'J1', 1, 1, 1, NULL),
		(2, 'Jenis Pendidikan 2', 'J2', 2, 2, 2, NULL),
		(3, 'Jenis Pendidikan 3', 'J3', 3, 3, 3, '2023-02-20'),
		(4, 'Jenis Pendidikan 4', NULL, NULL, NULL, NULL, NULL);
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
			name:             "ok: tanpa parameter apapun",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": {
					"id": 1,
					"nama": "Jenis Pendidikan 1",
					"abbreviation": "J1",
					"golongan_id": 1,
					"golongan_awal_id": 1,
					"tingkat": 1
				}
			}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			id:               "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: data tidak ditemukan",
			dbData:           dbData,
			id:               "999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: data deleted",
			dbData:           dbData,
			id:               "3",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/tingkat-pendidikan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_createAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_tingkat_pendidikan (id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat, deleted_at)
		VALUES
		(1, 'Jenis Pendidikan 1', 'J1', 1, 1, 1, NULL),
		(2, 'Jenis Pendidikan 2', 'J2', 2, 2, 2, NULL),
		(3, 'Jenis Pendidikan 3', 'J3', 3, 3, 3, '2023-02-20'),
		(4, 'Jenis Pendidikan 4', NULL, NULL, NULL, NULL, NULL);

		SELECT setval('ref_tingkat_pendidikan_id_seq', 4);
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
			name:        "ok: create tingkat pendidikan",
			dbData:      dbData,
			requestBody: `{"nama": "Jenis Pendidikan 5", "abbreviation": "J5", "golongan_id": 4, "golongan_awal_id": 4, "tingkat": 4}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id": 5, "nama": "Jenis Pendidikan 5", "abbreviation": "J5", "golongan_id": 4, "golongan_awal_id": 4, "tingkat": 4}}`,
		},
		{
			name:        "ok: create tingkat pendidikan with minimal data",
			dbData:      dbData,
			requestBody: `{"nama": "Jenis Pendidikan 5"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id": 5, "nama": "Jenis Pendidikan 5", "abbreviation": null, "golongan_id": null, "golongan_awal_id": null, "tingkat": null}}`,
		},
		{
			name:        "ok: create tingkat pendidikan with null data except nama",
			dbData:      dbData,
			requestBody: `{"nama": "Jenis Pendidikan 5", "abbreviation": null, "golongan_id": null, "golongan_awal_id": null, "tingkat": null}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id": 5, "nama": "Jenis Pendidikan 5", "abbreviation": null, "golongan_id": null, "golongan_awal_id": null, "tingkat": null}}`,
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
			name:        "error: auth header tidak valid",
			dbData:      dbData,
			requestBody: `{"nama": "Jenis Pendidikan 5"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:        "error: user is not an admin",
			dbData:      dbData,
			requestBody: `{"nama": "Jenis Pendidikan 5"}`,
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/tingkat-pendidikan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_updateAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_tingkat_pendidikan (id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat, deleted_at)
		VALUES
		(1, 'Jenis Pendidikan 1', 'J1', 1, 1, 1, NULL),
		(2, 'Jenis Pendidikan 2', 'J2', 2, 2, 2, NULL),
		(3, 'Jenis Pendidikan 3', 'J3', 3, 3, 3, '2023-02-20'),
		(4, 'Jenis Pendidikan 4', NULL, NULL, NULL, NULL, NULL);
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
			name:        "ok: update tingkat pendidikan",
			dbData:      dbData,
			id:          "1",
			requestBody: `{"nama": "Jenis Pendidikan 1 Diperbarui", "abbreviation": "J1 Diperbarui", "golongan_id": 1, "golongan_awal_id": 1, "tingkat": 1}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"id": 1, "nama": "Jenis Pendidikan 1 Diperbarui", "abbreviation": "J1 Diperbarui", "golongan_id": 1, "golongan_awal_id": 1, "tingkat": 1}}`,
		},
		{
			name:        "ok: update tingkat pendidikan with minimal data",
			dbData:      dbData,
			id:          "1",
			requestBody: `{"nama": "Jenis Pendidikan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"id": 1, "nama": "Jenis Pendidikan 1 Diperbarui", "abbreviation": null, "golongan_id": null, "golongan_awal_id": null, "tingkat": null}}`,
		},
		{
			name:        "error: missing required field nama",
			dbData:      dbData,
			id:          "1",
			requestBody: `{}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name:        "error: auth header tidak valid",
			dbData:      dbData,
			id:          "1",
			requestBody: `{"nama": "Jenis Pendidikan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:        "error: user is not an admin",
			dbData:      dbData,
			id:          "1",
			requestBody: `{"nama": "Jenis Pendidikan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:        "error: data tidak ditemukan",
			dbData:      dbData,
			id:          "999",
			requestBody: `{"nama": "Jenis Pendidikan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: data deleted",
			dbData:      dbData,
			id:          "3",
			requestBody: `{"nama": "Jenis Pendidikan 1 Diperbarui"}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/tingkat-pendidikan/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_deleteAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_tingkat_pendidikan (id, nama, abbreviation, golongan_id, golongan_awal_id, tingkat, deleted_at)
		VALUES
		(1, 'Jenis Pendidikan 1', 'J1', 1, 1, 1, NULL),
		(2, 'Jenis Pendidikan 2', 'J2', 2, 2, 2, NULL),
		(3, 'Jenis Pendidikan 3', 'J3', 3, 3, 3, '2023-02-20'),
		(4, 'Jenis Pendidikan 4', NULL, NULL, NULL, NULL, NULL);
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
			name:   "ok: delete tingkat pendidikan",
			dbData: dbData,
			id:     "1",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: "",
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
		{
			name:   "error: data tidak ditemukan",
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
			name:   "error: data deleted",
			dbData: dbData,
			id:     "3",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/tingkat-pendidikan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
