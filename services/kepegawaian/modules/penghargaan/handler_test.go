package penghargaan

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_penghargaan
			(id, nama, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', NULL),
			(12, 'Jenis Penghargaan 2', NULL),
			(14, 'Jenis Penghargaan 3', now());

		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan_id, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 11, 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 12, 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 14, 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 11, 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 11, 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 13,
						"deskripsi": "Deskripsi Penghargaan 3",
						"jenis_penghargaan": "",
						"nama_penghargaan": "Penghargaan 3",
						"tanggal": "2002-01-01"
					},
					{
						"id":                12,
						"deskripsi":         "Deskripsi Penghargaan 2",
						"jenis_penghargaan": "Jenis Penghargaan 2",
						"nama_penghargaan":  "Penghargaan 2",
						"tanggal":           "2001-01-01"
					},
					{
						"id":                11,
						"deskripsi":         "Deskripsi Penghargaan 1",
						"jenis_penghargaan": "Jenis Penghargaan 1",
						"nama_penghargaan":  "Penghargaan 1",
						"tanggal":           "2000-01-01"
					}

				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                12,
						"deskripsi":         "Deskripsi Penghargaan 2",
						"jenis_penghargaan": "Jenis Penghargaan 2",
						"nama_penghargaan":  "Penghargaan 2",
						"tanggal":           "2001-01-01"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			repo := repo.New(db)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-penghargaan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_getBerkas(t *testing.T) {
	t.Parallel()

	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0xf8, 0xff, 0xff, 0x3f,
		0x00, 0x05, 0xfe, 0x02, 0xfe, 0xa7, 0x46, 0x90,
		0x3d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
		0x44, 0xae, 0x42, 0x60, 0x82,
	}

	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		insert into ref_jenis_penghargaan
			(id, nama, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', NULL);

		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan_id, nama_penghargaan, deskripsi_penghargaan, file_base64, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 11, 'Penghargaan 1', 'Deskripsi Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2000-01-01', '41', NULL),
			(12, 11, 'Penghargaan 2', 'Deskripsi Penghargaan 2', 'data:image/png;base64,invalid', '2001-01-01', '41', NULL),
			(13, 11, 'Penghargaan 3', 'Deskripsi Penghargaan 3', 'data:image/png;base64,invalid', '2002-01-01', '41', now()),
			(14, 11, 'Penghargaan 4', 'Deskripsi Penghargaan 4', NULL, '2003-01-01', '41', NULL),
			(15, 11, 'Penghargaan 5', 'Deskripsi Penghargaan 5', '', '2004-01-01', '41', NULL);
		`

	tests := []struct {
		name              string
		dbData            string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			dbData:            dbData,
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 berkas riwayat penghargaan tidak valid",
			dbData:            dbData,
			paramID:           "12",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan sudah dihapus",
			dbData:            dbData,
			paramID:           "13",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat penghargaan berisi null value",
			dbData:            dbData,
			paramID:           "14",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat penghargaan berupa string kosong",
			dbData:            dbData,
			paramID:           "15",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan bukan milik user login",
			dbData:            dbData,
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "42")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-penghargaan/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repo.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseCode == http.StatusOK {
				assert.Equal(t, "inline", rec.Header().Get("Content-Disposition"))
				assert.Equal(t, tt.wantContentType, rec.Header().Get("Content-Type"))
				assert.Equal(t, tt.wantResponseBytes, rec.Body.Bytes())
			} else {
				assert.JSONEq(t, string(tt.wantResponseBytes), rec.Body.String())
			}
		})
	}
}
