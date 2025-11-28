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
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata/usulanperubahandatatest"
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

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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

	repo := dbrepo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, pgxconn, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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

	repo := dbrepo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, pgxconn, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penghargaan/"+tt.paramID+"/berkas", &buf)
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

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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
				"message": "parameter \"jenis_penghargaan\" harus salah satu dari \"Internasional\", \"Unit Kerja (eselon 2 ke bawah)\", \"Unit Utama\", \"Nasional\", \"Instansional (Kementerian/Lembaga)\""
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

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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
				"message": "parameter \"jenis_penghargaan\" harus salah satu dari \"Internasional\", \"Unit Kerja (eselon 2 ke bawah)\", \"Unit Utama\", \"Nasional\", \"Instansional (Kementerian/Lembaga)\""
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

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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

func Test_handler_usulanPerubahanData(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, created_at, updated_at, deleted_at)
			values
			(11, 'Internasional', 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '1a', '2000-01-01', '2000-01-01', NULL),
			(12, 'Internasional', 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '1a', '2000-01-01', '2000-01-01', '2000-01-01'),
			(13, 'Nasional', 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '1d', '2000-01-01', '2000-01-01', NULL),
			(14, 'Unit Utama', 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '1e', '2000-01-01', '2000-01-01', NULL),
			(15, 'Instansional (Kementerian/Lembaga)', 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '1g', '2000-01-01', '2000-01-01', NULL),
			(16, 'Unit Kerja (eselon 2 ke bawah)', 'Penghargaan 6', 'Deskripsi Penghargaan 6', '2005-01-01', '1h', '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Request)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)
	svcRoute := usulanperubahandatatest.NewServiceRoute(db)

	RegisterRoutes(e, repo, authMw, svcRoute)

	// Query actual database rows for user 1a to use in error test cases
	actualRows1a, err := dbtest.QueryWithClause(db, "riwayat_penghargaan_umum", "where nip = $1 order by id", "1a")
	require.NoError(t, err)

	authHeader1a := []string{apitest.GenerateAuthHeader("1a")}
	tests := []struct {
		name                 string
		requestHeader        http.Header
		requestBody          string
		doRollback           bool
		wantResponsePostCode int
		wantResponsePostBody string
		wantResponseGetBody  string
		wantDBSvcRows        dbtest.Rows
		wantDBUsulanRows     dbtest.Rows
	}{
		{
			name:          "ok: success create riwayat penghargaan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"jenis_penghargaan": "Internasional",
					"nama_penghargaan": "Penghargaan Baru",
					"deskripsi": "Deskripsi Baru",
					"tanggal": "2023-01-01"
				}
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-penghargaan",
						"action":     "CREATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"jenis_penghargaan": [ null, "Internasional" ],
							"nama_penghargaan":  [ null, "Penghargaan Baru" ],
							"deskripsi":         [ null, "Deskripsi Baru" ],
							"tanggal":           [ null, "2023-01-01" ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"jenis_penghargaan":     "Internasional",
					"nama_penghargaan":      "Penghargaan Baru",
					"deskripsi_penghargaan": "Deskripsi Baru",
					"tanggal_penghargaan":   time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					"nip":                   "1c",
					"file_base64":           nil,
					"s3_file_id":            nil,
					"exist":                 true,
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1c",
					"jenis_data": "riwayat-penghargaan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"jenis_penghargaan": []any{nil, "Internasional"},
						"nama_penghargaan":  []any{nil, "Penghargaan Baru"},
						"deskripsi":         []any{nil, "Deskripsi Baru"},
						"tanggal":           []any{nil, "2023-01-01"},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: success update riwayat penghargaan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "13",
				"data": {
					"jenis_penghargaan": "Internasional",
					"nama_penghargaan": "Penghargaan Updated",
					"deskripsi": "",
					"tanggal": "2000-01-01"
				}
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-penghargaan",
						"action":     "UPDATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "13",
						"perubahan_data": {
							"jenis_penghargaan": [ "Nasional", "Internasional" ],
							"nama_penghargaan":  [ "Penghargaan 3", "Penghargaan Updated" ],
							"deskripsi":         [ "Deskripsi Penghargaan 3", null ],
							"tanggal":           [ "2002-01-01", "2000-01-01" ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    int32(13),
					"jenis_penghargaan":     "Internasional",
					"nama_penghargaan":      "Penghargaan Updated",
					"deskripsi_penghargaan": nil,
					"tanggal_penghargaan":   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"nip":                   "1d",
					"file_base64":           nil,
					"s3_file_id":            nil,
					"exist":                 true,
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1d",
					"jenis_data": "riwayat-penghargaan",
					"data_id":    "13",
					"perubahan_data": map[string]any{
						"jenis_penghargaan": []any{"Nasional", "Internasional"},
						"nama_penghargaan":  []any{"Penghargaan 3", "Penghargaan Updated"},
						"deskripsi":         []any{"Deskripsi Penghargaan 3", nil},
						"tanggal":           []any{"2002-01-01", "2000-01-01"},
					},
					"action":     "UPDATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: success delete riwayat penghargaan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1e")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "14"
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-penghargaan",
						"action":     "DELETE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "14",
						"perubahan_data": {
							"jenis_penghargaan": [ "Unit Utama", null ],
							"nama_penghargaan":  [ "Penghargaan 4", null ],
							"deskripsi":         [ "Deskripsi Penghargaan 4", null ],
							"tanggal":           [ "2003-01-01", null ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    int32(14),
					"jenis_penghargaan":     "Unit Utama",
					"nama_penghargaan":      "Penghargaan 4",
					"deskripsi_penghargaan": "Deskripsi Penghargaan 4",
					"tanggal_penghargaan":   time.Date(2003, 1, 1, 0, 0, 0, 0, time.UTC),
					"nip":                   "1e",
					"file_base64":           nil,
					"s3_file_id":            nil,
					"exist":                 true,
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            "{deleted_at}",
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1e",
					"jenis_data": "riwayat-penghargaan",
					"data_id":    "14",
					"perubahan_data": map[string]any{
						"jenis_penghargaan": []any{"Unit Utama", nil},
						"nama_penghargaan":  []any{"Penghargaan 4", nil},
						"deskripsi":         []any{"Deskripsi Penghargaan 4", nil},
						"tanggal":           []any{"2003-01-01", nil},
					},
					"action":     "DELETE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: rollback on usulan perubahan data should not CREATE record",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1f")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"jenis_penghargaan": "Nasional",
					"nama_penghargaan": "Penghargaan Baru",
					"deskripsi": "",
					"tanggal": "2000-01-01"
				}
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-penghargaan",
						"action":     "CREATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"jenis_penghargaan": [ null, "Nasional" ],
							"nama_penghargaan":  [ null, "Penghargaan Baru" ],
							"deskripsi":         [ null, null ],
							"tanggal":           [ null, "2000-01-01" ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1f",
					"jenis_data": "riwayat-penghargaan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"jenis_penghargaan": []any{nil, "Nasional"},
						"nama_penghargaan":  []any{nil, "Penghargaan Baru"},
						"deskripsi":         []any{nil, nil},
						"tanggal":           []any{nil, "2000-01-01"},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: rollback on usulan perubahan data should not UPDATE record",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1g")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "15",
				"data": {
					"jenis_penghargaan": "Nasional",
					"nama_penghargaan": "Penghargaan Updated",
					"deskripsi": "",
					"tanggal": "2000-01-01"
				}
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-penghargaan",
						"action":     "UPDATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    "15",
						"perubahan_data": {
							"jenis_penghargaan": [ "Instansional (Kementerian/Lembaga)", "Nasional" ],
							"nama_penghargaan":  [ "Penghargaan 5", "Penghargaan Updated" ],
							"deskripsi":         [ "Deskripsi Penghargaan 5", null ],
							"tanggal":           [ "2004-01-01", "2000-01-01" ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    int32(15),
					"jenis_penghargaan":     "Instansional (Kementerian/Lembaga)",
					"nama_penghargaan":      "Penghargaan 5",
					"deskripsi_penghargaan": "Deskripsi Penghargaan 5",
					"tanggal_penghargaan":   time.Date(2004, 1, 1, 0, 0, 0, 0, time.UTC),
					"nip":                   "1g",
					"file_base64":           nil,
					"s3_file_id":            nil,
					"exist":                 true,
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1g",
					"jenis_data": "riwayat-penghargaan",
					"data_id":    "15",
					"perubahan_data": map[string]any{
						"jenis_penghargaan": []any{"Instansional (Kementerian/Lembaga)", "Nasional"},
						"nama_penghargaan":  []any{"Penghargaan 5", "Penghargaan Updated"},
						"deskripsi":         []any{"Deskripsi Penghargaan 5", nil},
						"tanggal":           []any{"2004-01-01", "2000-01-01"},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: rollback on usulan perubahan data should not DELETE record",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1h")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "16"
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-penghargaan",
						"action":     "DELETE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    "16",
						"perubahan_data": {
							"jenis_penghargaan": [ "Unit Kerja (eselon 2 ke bawah)", null ],
							"nama_penghargaan":  [ "Penghargaan 6", null ],
							"deskripsi":         [ "Deskripsi Penghargaan 6", null ],
							"tanggal":           [ "2005-01-01", null ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    int32(16),
					"jenis_penghargaan":     "Unit Kerja (eselon 2 ke bawah)",
					"nama_penghargaan":      "Penghargaan 6",
					"deskripsi_penghargaan": "Deskripsi Penghargaan 6",
					"tanggal_penghargaan":   time.Date(2005, 1, 1, 0, 0, 0, 0, time.UTC),
					"nip":                   "1h",
					"file_base64":           nil,
					"s3_file_id":            nil,
					"exist":                 true,
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1h",
					"jenis_data": "riwayat-penghargaan",
					"data_id":    "16",
					"perubahan_data": map[string]any{
						"jenis_penghargaan": []any{"Unit Kerja (eselon 2 ke bawah)", nil},
						"nama_penghargaan":  []any{"Penghargaan 6", nil},
						"deskripsi":         []any{"Deskripsi Penghargaan 6", nil},
						"tanggal":           []any{"2005-01-01", nil},
					},
					"action":     "DELETE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "error: riwayat penghargaan is owned by other pegawai",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1b")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "11"
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat penghargaan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbtest.Rows{},
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: riwayat penghargaan is not found",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "DELETE",
				"data_id": "0"
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat penghargaan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    actualRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: riwayat penghargaan is deleted",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "12",
				"data": {
					"jenis_penghargaan": "Internasional",
					"nama_penghargaan": "Penghargaan Updated",
					"deskripsi": "",
					"tanggal": "2000-01-01"
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat penghargaan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    actualRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: invalid jenis_penghargaan",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"jenis_penghargaan": "invalid",
					"nama_penghargaan": "Penghargaan Baru",
					"deskripsi": "",
					"tanggal": "2000-01-01"
				}
			}`,
			wantResponsePostCode: http.StatusInternalServerError,
			wantResponsePostBody: `{"message": "Internal Server Error"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    actualRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: missing required params on data",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "doesn't match schema due to: ` +
				`Error at \"/data/jenis_penghargaan\": property \"jenis_penghargaan\" is missing` +
				` | Error at \"/data/nama_penghargaan\": property \"nama_penghargaan\" is missing` +
				` | Error at \"/data/tanggal\": property \"tanggal\" is missing Or ` +
				`Error at \"/action\": value is not one of the allowed values [\"UPDATE\"]` +
				` | Error at \"/data/jenis_penghargaan\": property \"jenis_penghargaan\" is missing` +
				` | Error at \"/data/nama_penghargaan\": property \"nama_penghargaan\" is missing` +
				` | Error at \"/data/tanggal\": property \"tanggal\" is missing | Error at \"/data_id\": property \"data_id\" is missing Or ` +
				`Error at \"/action\": value is not one of the allowed values [\"DELETE\"]` +
				` | property \"data\" is unsupported | Error at \"/data_id\": property \"data_id\" is missing"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    actualRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:                 "error: body is empty",
			requestHeader:        http.Header{"Authorization": authHeader1a},
			requestBody:          `{}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "doesn't match schema due to: ` +
				`Error at \"/action\": property \"action\" is missing` +
				` | Error at \"/data\": property \"data\" is missing Or ` +
				`Error at \"/action\": property \"action\" is missing` +
				` | Error at \"/data_id\": property \"data_id\" is missing` +
				` | Error at \"/data\": property \"data\" is missing Or ` +
				`Error at \"/action\": property \"action\" is missing` +
				` | Error at \"/data_id\": property \"data_id\" is missing"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    actualRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// validate create & approve usulan
			req := httptest.NewRequest(http.MethodPost, "/v1/usulan-perubahan-data/riwayat-penghargaan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			if tt.doRollback {
				req.URL.RawQuery = "rollback=true"
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponsePostCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponsePostBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			nip := apitest.GetNIPFromAuthHeader(req.Header.Get("Authorization"))

			actualSvcRows, err := dbtest.QueryWithClause(db, "riwayat_penghargaan_umum", "where nip = $1 order by id", nip)
			require.NoError(t, err)
			if len(tt.wantDBSvcRows) == len(actualSvcRows) {
				for i, row := range actualSvcRows {
					if tt.wantDBSvcRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBSvcRows[i]["id"] = row["id"]
						tt.wantDBSvcRows[i]["created_at"] = row["created_at"]
						tt.wantDBSvcRows[i]["updated_at"] = row["updated_at"]
					}
					if tt.wantDBSvcRows[i]["created_at"] == "{created_at}" {
						tt.wantDBSvcRows[i]["created_at"] = row["created_at"]
					}
					if tt.wantDBSvcRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBSvcRows[i]["updated_at"] = row["updated_at"]
					}
					if tt.wantDBSvcRows[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBSvcRows[i]["deleted_at"] = row["deleted_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBSvcRows, actualSvcRows)

			actualUsulanRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where nip = $1 order by id", nip)
			require.NoError(t, err)
			// Replace {id} placeholder in response body with actual IDs
			for _, row := range actualUsulanRows {
				tt.wantResponseGetBody = strings.ReplaceAll(tt.wantResponseGetBody, "{id}", fmt.Sprintf("%d", row["id"]))
			}
			if len(tt.wantDBUsulanRows) == len(actualUsulanRows) {
				for i, row := range actualUsulanRows {
					assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
					assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)

					tt.wantDBUsulanRows[i]["id"] = row["id"]
					tt.wantDBUsulanRows[i]["created_at"] = row["created_at"]
					tt.wantDBUsulanRows[i]["updated_at"] = row["updated_at"]
				}
			}
			assert.Equal(t, tt.wantDBUsulanRows, actualUsulanRows)

			// validate get usulan
			req2 := httptest.NewRequest(http.MethodGet, "/v1/usulan-perubahan-data/riwayat-penghargaan", nil)
			req2.Header = tt.requestHeader
			rec2 := httptest.NewRecorder()

			e.ServeHTTP(rec2, req2)

			assert.Equal(t, http.StatusOK, rec2.Code)
			assert.JSONEq(t, tt.wantResponseGetBody, rec2.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec2, req2, e))
		})
	}
}
