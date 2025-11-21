package riwayatpenghargaan

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
	"net/url"
	"strconv"
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
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 'Jenis Penghargaan 2', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 'Jenis Penghargaan 3', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 'Jenis Penghargaan 1', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 'Jenis Penghargaan 1', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
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
			name:             "ok: tanpa parameter apapun",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                13,
						"deskripsi":         "Deskripsi Penghargaan 3",
						"jenis_penghargaan": "Jenis Penghargaan 3",
						"nama_penghargaan":  "Penghargaan 3",
						"tanggal":           "2002-01-01"
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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-penghargaan", nil)
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
		insert into riwayat_penghargaan_umum
			(id, nama_penghargaan, deskripsi_penghargaan, file_base64, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Penghargaan 1', 'Deskripsi Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2000-01-01', '41', NULL),
			(12, 'Penghargaan 2', 'Deskripsi Penghargaan 2', 'data:image/png;base64,invalid', '2001-01-01', '41', NULL),
			(13, 'Penghargaan 3', 'Deskripsi Penghargaan 3', 'data:image/png;base64,invalid', '2002-01-01', '41', now()),
			(14, 'Penghargaan 4', 'Deskripsi Penghargaan 4', NULL, '2003-01-01', '41', NULL),
			(15, 'Penghargaan 5', 'Deskripsi Penghargaan 5', '', '2004-01-01', '41', NULL);
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("41")}
	tests := []struct {
		name              string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 berkas riwayat penghargaan tidak valid",
			paramID:           "12",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan sudah dihapus",
			paramID:           "13",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat penghargaan berisi null value",
			paramID:           "14",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat penghargaan berupa string kosong",
			paramID:           "15",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan bukan milik user login",
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("42")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-penghargaan/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 'Jenis Penghargaan 2', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 'Jenis Penghargaan 3', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 'Jenis Penghargaan 1', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 'Jenis Penghargaan 1', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		nip              string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: tanpa parameter apapun",
			nip:              "41",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                13,
						"deskripsi":         "Deskripsi Penghargaan 3",
						"jenis_penghargaan": "Jenis Penghargaan 3",
						"nama_penghargaan":  "Penghargaan 3",
						"tanggal":           "2002-01-01"
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
			nip:              "41",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
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
			name:             "error: auth header tidak valid",
			nip:              "41",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-penghargaan", nil)
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

func Test_handler_getBerkasAdmin(t *testing.T) {
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
		insert into riwayat_penghargaan_umum
			(id, nama_penghargaan, deskripsi_penghargaan, file_base64, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Penghargaan 1', 'Deskripsi Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2000-01-01', '41', NULL),
			(12, 'Penghargaan 2', 'Deskripsi Penghargaan 2', 'data:image/png;base64,invalid', '2001-01-01', '41', NULL),
			(13, 'Penghargaan 3', 'Deskripsi Penghargaan 3', 'data:image/png;base64,invalid', '2002-01-01', '41', now()),
			(14, 'Penghargaan 4', 'Deskripsi Penghargaan 4', NULL, '2003-01-01', '41', NULL),
			(15, 'Penghargaan 5', 'Deskripsi Penghargaan 5', '', '2004-01-01', '41', NULL);
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name              string
		paramID           string
		nip               string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			paramID:           "11",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 berkas riwayat penghargaan tidak valid",
			paramID:           "12",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan sudah dihapus",
			paramID:           "13",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat penghargaan berisi null value",
			paramID:           "14",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat penghargaan berupa string kosong",
			paramID:           "15",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat penghargaan tidak ditemukan",
			paramID:           "0",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penghargaan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			paramID:           "abc",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "invalid request"}`),
		},
		{
			name:              "error: auth header tidak valid",
			paramID:           "11",
			nip:               "41",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-penghargaan/%s/berkas", tt.nip, tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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

func Test_handler_adminUploadBerkas(t *testing.T) {
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

	dbData := `
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 'Jenis Penghargaan 2', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 'Jenis Penghargaan 3', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 'Jenis Penghargaan 1', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 'Jenis Penghargaan 1', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		files            map[string][]byte
		requestHeader    http.Header
		fileContentType  string
		wantResponseCode int
		wantResponseBody string
		paramNIP         string
		paramID          string
	}{
		{
			name:     "ok: create template with file",
			paramNIP: "41",
			paramID:  "11",
			files: map[string][]byte{
				"file": pngBytes,
			},
			fileContentType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:     "error: id not exists",
			paramNIP: "1c",
			paramID:  "99",
			files: map[string][]byte{
				"file": pngBytes,
			},
			fileContentType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{
				"message": "data tidak ditemukan"
			}`,
		},
		{
			name:     "error: id is deleted",
			paramNIP: "1c",
			paramID:  "14",
			files: map[string][]byte{
				"file": pngBytes,
			},
			fileContentType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{
				"message": "data tidak ditemukan"
			}`,
		},
		{
			name:     "error: file with invalid type",
			paramNIP: "41",
			paramID:  "11",
			files: map[string][]byte{
				"file": pngBytes,
			},
			fileContentType:  "image/x-xpixmap",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"file\" harus dalam format yang sesuai"}`,
		},
		{
			name:             "error: missing file upload",
			paramNIP:         "41",
			paramID:          "11",
			fileContentType:  "application/pdf",
			files:            nil,
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"file\" harus diisi"}`,
		},
		{
			name:            "error: invalid auth header",
			paramNIP:        "41",
			paramID:         "11",
			fileContentType: "application/pdf",
			files: map[string][]byte{
				"file": pngBytes,
			},
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			for fieldName, content := range tt.files {
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, "example.bin"))
				h.Set("Content-Type", tt.fileContentType)

				part, err := writer.CreatePart(h)
				require.NoError(t, err)
				_, err = part.Write(content)
				require.NoError(t, err)
			}

			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penghargaan/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
		})
	}
}

func Test_handler_adminCreate(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 'Jenis Penghargaan 2', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 'Jenis Penghargaan 3', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 'Jenis Penghargaan 1', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 'Jenis Penghargaan 1', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "culpa occaecat eiusmod commodo minim veniam adipisicing pariatur reprehenderit quis",
				"tanggal": "2002-01-01"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
		},
		{
			name:          "ok: with required data",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "",
				"tanggal": null
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
		},
		{
			name:          "error: required data not included",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "",
				"deskripsi": "",
				"tanggal": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{
				"message": "parameter \"nama_penghargaan\" harus 1 karakter atau lebih"
			}`,
		},
		{
			name:          "error: invalid jenis_penghargaan",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "minim",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "culpa occaecat eiusmod commodo minim veniam adipisicing pariatur reprehenderit quis",
				"tanggal": "2002-01-01"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{
				"message": "jenis penghargaan tidak valid"
			}`,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penghargaan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if rec.Code == http.StatusCreated {
				var resp adminCreateResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)
				tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "\"{id}\"", strconv.Itoa((int(resp.Data.ID))))
			}
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 'Jenis Penghargaan 2', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 'Jenis Penghargaan 3', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 'Jenis Penghargaan 1', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 'Jenis Penghargaan 1', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			paramID:       "11",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "culpa occaecat eiusmod commodo minim veniam adipisicing pariatur reprehenderit quis",
				"tanggal": "2002-01-01"
			}`,
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:          "ok: with required data",
			paramNIP:      "1c",
			paramID:       "11",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "",
				"tanggal": null
			}`,
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:          "error: required data not included",
			paramNIP:      "1c",
			paramID:       "11",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "",
				"deskripsi": "",
				"tanggal": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{
				"message": "parameter \"nama_penghargaan\" harus 1 karakter atau lebih"
			}`,
		},
		{
			name:          "error: invalid jenis_penghargaan",
			paramNIP:      "1c",
			paramID:       "11",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "minim",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "culpa occaecat eiusmod commodo minim veniam adipisicing pariatur reprehenderit quis",
				"tanggal": "2002-01-01"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{
				"message": "jenis penghargaan tidak valid"
			}`,
		},
		{
			name:          "error: id not exists",
			paramNIP:      "1c",
			paramID:       "99",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "culpa occaecat eiusmod commodo minim veniam adipisicing pariatur reprehenderit quis",
				"tanggal": "2002-01-01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{
				"message": "data tidak ditemukan"
			}`,
		},
		{
			name:          "error: id is deleted",
			paramNIP:      "1c",
			paramID:       "14",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_penghargaan": "Internasional",
				"nama_penghargaan": "penghargaan",
				"deskripsi": "culpa occaecat eiusmod commodo minim veniam adipisicing pariatur reprehenderit quis",
				"tanggal": "2002-01-01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{
				"message": "data tidak ditemukan"
			}`,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1a",
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1a",
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penghargaan/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
		})
	}
}

func Test_handler_adminDelete(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 'Jenis Penghargaan 2', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 'Jenis Penghargaan 3', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 'Jenis Penghargaan 1', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 'Jenis Penghargaan 1', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: with all data",
			paramNIP:         "41",
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:             "error: id not exists",
			paramNIP:         "41",
			paramID:          "99",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{
				"message": "data tidak ditemukan"
			}`,
		},
		{
			name:             "error: id is deleted",
			paramNIP:         "41",
			paramID:          "14",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{
				"message": "data tidak ditemukan"
			}`,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1a",
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penghargaan/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
		})
	}
}
