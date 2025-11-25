package riwayatpenugasan

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

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
		insert into riwayat_penugasan
			(id, nip,  tipe_jabatan, nama_jabatan,         deskripsi_jabatan,     is_menjabat,  tanggal_mulai, tanggal_selesai,                  deleted_at) values
			(1,  '1c', 'Struktural', 'Kepala Bagian',      'Deskripsi Jabatan 1', false,        '2023-01-01',  '2023-12-31',                     null),
			(2,  '1c', 'Fungsional', 'Analis Kepegawaian', 'Deskripsi Jabatan 2', false,        '2024-01-01',  null,                             null),
			(3,  '1c', 'Struktural', 'Kepala Sub Bagian',  'Deskripsi Jabatan 3', true,         '2023-06-01',  '2024-12-31',                     null),
			(4,  '1c', 'Struktural', 'Kepala Sub Bagian',  'Deskripsi Jabatan 4', false,        '2024-01-01',  '2024-12-31',                     '2000-01-01'),
			(5,  '2a', 'Struktural', 'Kepala Bagian',      'Deskripsi Jabatan 5', false,        '2024-01-01',  '2024-12-31',                     null),
			(6,  '1c', 'Struktural', 'Kepala Bagian',      'Deskripsi Jabatan 6', false,        '2022-01-01',  current_date + interval '1 days', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
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
						"id": 2,
						"tipe_jabatan": "Fungsional",
						"nama_jabatan": "Analis Kepegawaian",
						"deskripsi_jabatan": "Deskripsi Jabatan 2",
						"tanggal_mulai": "2024-01-01",
						"tanggal_selesai": null,
						"is_menjabat": true
					},
					{
						"id": 3,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Sub Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 3",
						"tanggal_mulai": "2023-06-01",
						"tanggal_selesai": "2024-12-31",
						"is_menjabat": true
					},
					{
						"id": 1,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 1",
						"tanggal_mulai": "2023-01-01",
						"tanggal_selesai": "2023-12-31",
						"is_menjabat": false
					},
					{
						"id": 6,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 6",
						"tanggal_mulai": "2022-01-01",
						"tanggal_selesai": "` + time.Now().Add(24*time.Hour).Format(time.DateOnly) + `",
						"is_menjabat": true
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
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
						"id": 3,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Sub Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 3",
						"tanggal_mulai": "2023-06-01",
						"tanggal_selesai": "2024-12-31",
						"is_menjabat": true
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}`,
		},
		{
			name:             "ok: tidak ada data",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2c")}},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-penugasan", nil)
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

	filePath := "../../../../lib/api/sample/hello.pdf"
	pdfBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

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

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		insert into riwayat_penugasan
			(id, nip, deleted_at,   file_base64) values
			(1, '1c', null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c', null,         '` + pdfBase64 + `'),
			(3, '1c', null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c', null,         'data:application/pdf;base64,invalid'),
			(5, '1c', '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c', null,         null),
			(7, '1c', null,         '');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
	tests := []struct {
		name              string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 riwayat penugasan tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat penugasan sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat penugasan berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat penugasan berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat penugasan bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat penugasan tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
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
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-penugasan/%s/berkas", tt.paramID), nil)
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
		insert into riwayat_penugasan
			(id, nip,  tipe_jabatan, nama_jabatan,         deskripsi_jabatan,     is_menjabat,  tanggal_mulai, tanggal_selesai,                  deleted_at) values
			(1,  '1c', 'Struktural', 'Kepala Bagian',      'Deskripsi Jabatan 1', false,        '2023-01-01',  '2023-12-31',                     null),
			(2,  '1c', 'Fungsional', 'Analis Kepegawaian', 'Deskripsi Jabatan 2', false,        '2024-01-01',  null,                             null),
			(3,  '1c', 'Struktural', 'Kepala Sub Bagian',  'Deskripsi Jabatan 3', true,         '2023-06-01',  '2024-12-31',                     null),
			(4,  '1c', 'Struktural', 'Kepala Sub Bagian',  'Deskripsi Jabatan 4', false,        '2024-01-01',  '2024-12-31',                     '2000-01-01'),
			(5,  '2a', 'Struktural', 'Kepala Bagian',      'Deskripsi Jabatan 5', false,        '2024-01-01',  '2024-12-31',                     null),
			(6,  '1c', 'Struktural', 'Kepala Bagian',      'Deskripsi Jabatan 6', false,        '2022-01-01',  current_date + interval '1 days', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
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
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 2,
						"tipe_jabatan": "Fungsional",
						"nama_jabatan": "Analis Kepegawaian",
						"deskripsi_jabatan": "Deskripsi Jabatan 2",
						"tanggal_mulai": "2024-01-01",
						"tanggal_selesai": null,
						"is_menjabat": true
					},
					{
						"id": 3,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Sub Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 3",
						"tanggal_mulai": "2023-06-01",
						"tanggal_selesai": "2024-12-31",
						"is_menjabat": true
					},
					{
						"id": 1,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 1",
						"tanggal_mulai": "2023-01-01",
						"tanggal_selesai": "2023-12-31",
						"is_menjabat": false
					},
					{
						"id": 6,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 6",
						"tanggal_mulai": "2022-01-01",
						"tanggal_selesai": "` + time.Now().Add(24*time.Hour).Format(time.DateOnly) + `",
						"is_menjabat": true
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 3,
						"tipe_jabatan": "Struktural",
						"nama_jabatan": "Kepala Sub Bagian",
						"deskripsi_jabatan": "Deskripsi Jabatan 3",
						"tanggal_mulai": "2023-06-01",
						"tanggal_selesai": "2024-12-31",
						"is_menjabat": true
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}`,
		},
		{
			name:             "ok: tidak ada data",
			nip:              "2c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-penugasan", nil)
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

	filePath := "../../../../lib/api/sample/hello.pdf"
	pdfBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

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

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		insert into riwayat_penugasan
			(id, nip, deleted_at,   file_base64) values
			(1, '1c', null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c', null,         '` + pdfBase64 + `'),
			(3, '1c', null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c', null,         'data:application/pdf;base64,invalid'),
			(5, '1c', '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c', null,         null),
			(7, '1c', null,         '');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
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
			name:              "ok: valid pdf with data: prefix",
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			nip:               "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			nip:               "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 riwayat penugasan tidak valid",
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat penugasan sudah dihapus",
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat penugasan berisi null value",
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat penugasan berupa string kosong",
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat penugasan tidak ditemukan",
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat penugasan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			nip:               "1c",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-penugasan/%s/berkas", tt.nip, tt.paramID), nil)
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

func Test_handler_adminCreate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai
			(pns_id,  nip_baru, nama, deleted_at) VALUES
			('id_1a', '1a',     'Pegawai 1', NULL),
			('id_1c', '1c',     'Pegawai 2', NULL),
			('id_1d', '1d',     'Pegawai 3', '2000-01-01'),
			('id_1e', '1e',     'Pegawai 4', NULL),
			('id_1f', '1f',     'Pegawai 5', NULL),
			('id_1g', '1g',     'Pegawai 6', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"deskripsi_jabatan": "Deskripsi Jabatan 1",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "{id}",
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 1",
					"tanggal_mulai":     time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with minimal required data",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Fungsional",
				"nama_jabatan": "Analis Kepegawaian",
				"tanggal_mulai": "2023-02-15",
				"tanggal_selesai": "2023-05-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "{id}",
					"nip":               "1e",
					"tipe_jabatan":      "Fungsional",
					"nama_jabatan":      "Analis Kepegawaian",
					"deskripsi_jabatan": nil,
					"tanggal_mulai":     time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with empty string optional fields",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Sub Bagian",
				"deskripsi_jabatan": "",
				"tanggal_mulai": "2023-03-15",
				"tanggal_selesai": "2023-06-15",
				"is_menjabat": true
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "{id}",
					"nip":               "1f",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Sub Bagian",
					"deskripsi_jabatan": nil,
					"tanggal_mulai":     time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2023, 6, 15, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       true,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: missing required params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"tipe_jabatan\" harus diisi` +
				` | parameter \"nama_jabatan\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penugasan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_penugasan", "where nip = $1", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["id"] = row["id"]
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]

						tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", fmt.Sprintf("%d", row["id"]))
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai
			(pns_id,  nip_baru, deleted_at) VALUES
			('id_1c', '1c',     NULL),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     NULL);
		INSERT INTO riwayat_penugasan
			(id, nip, tipe_jabatan, nama_jabatan, deskripsi_jabatan, tanggal_mulai, tanggal_selesai, is_menjabat, created_at, updated_at, deleted_at) VALUES
			(1, '1c', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 1', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', NULL),
			(2, '1c', 'Fungsional', 'Analis Kepegawaian', 'Deskripsi Jabatan 2', '2000-01-01', NULL, false, '2000-01-01', '2000-01-01', NULL),
			(5, '1e', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 5', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', NULL),
			(6, '1c', 'Struktural', 'Kepala Sub Bagian', 'Deskripsi Jabatan 6', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', '2000-01-01'),
			(7, '1c', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 7', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                int32(7),
			"nip":               "1c",
			"tipe_jabatan":      "Struktural",
			"nama_jabatan":      "Kepala Bagian",
			"deskripsi_jabatan": "Deskripsi Jabatan 7",
			"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			"is_menjabat":       false,
			"file_base64":       nil,
			"s3_file_id":        nil,
			"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":        nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian Updated",
				"deskripsi_jabatan": "Deskripsi Jabatan Updated",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": true
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Bagian Updated",
					"deskripsi_jabatan": "Deskripsi Jabatan Updated",
					"tanggal_mulai":     time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2024, 1, 15, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       true,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with empty string optional fields",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Fungsional",
				"nama_jabatan": "Analis Kepegawaian Updated",
				"deskripsi_jabatan": "",
				"tanggal_mulai": "2023-02-15",
				"tanggal_selesai": "2023-05-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(2),
					"nip":               "1c",
					"tipe_jabatan":      "Fungsional",
					"nama_jabatan":      "Analis Kepegawaian Updated",
					"deskripsi_jabatan": nil,
					"tanggal_mulai":     time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: riwayat penugasan is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat penugasan is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(5),
					"nip":               "1e",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 5",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: riwayat penugasan is deleted",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(6),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Sub Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 6",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tipe_jabatan": "Struktural",
				"nama_jabatan": "Kepala Bagian",
				"tanggal_mulai": "2023-01-15",
				"tanggal_selesai": "2024-01-15",
				"is_menjabat": false
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "pegawai tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: missing required params",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"tipe_jabatan\" harus diisi` +
				` | parameter \"nama_jabatan\" harus diisi"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penugasan/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_penugasan", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminDelete(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai
			(pns_id,  nip_baru, deleted_at) VALUES
			('id_1c', '1c',     NULL),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     NULL);
		INSERT INTO riwayat_penugasan
			(id, nip, tipe_jabatan, nama_jabatan, deskripsi_jabatan, tanggal_mulai, tanggal_selesai, is_menjabat, created_at, updated_at, deleted_at) VALUES
			(1, '1c', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 1', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', NULL),
			(2, '1e', 'Fungsional', 'Analis Kepegawaian', 'Deskripsi Jabatan 2', '2000-01-01', NULL, false, '2000-01-01', '2000-01-01', NULL),
			(3, '1c', 'Struktural', 'Kepala Sub Bagian', 'Deskripsi Jabatan 3', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', '2000-01-01'),
			(4, '1c', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 4', '2000-01-01', '2000-12-31', false, '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success delete",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 1",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat penugasan is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(2),
					"nip":               "1e",
					"tipe_jabatan":      "Fungsional",
					"nama_jabatan":      "Analis Kepegawaian",
					"deskripsi_jabatan": "Deskripsi Jabatan 2",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   nil,
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:             "error: riwayat penugasan is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat penugasan is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(3),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Sub Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 3",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(4),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 4",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penugasan/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_penugasan", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["deleted_at"] = row["deleted_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminUploadBerkas(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai
			(pns_id,  nip_baru, deleted_at) VALUES
			('id_1c', '1c',     NULL),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     NULL);
		INSERT INTO riwayat_penugasan
			(id, nip, tipe_jabatan, nama_jabatan, deskripsi_jabatan, tanggal_mulai, tanggal_selesai, is_menjabat, file_base64, created_at, updated_at, deleted_at) VALUES
			(1, '1c', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 1', '2000-01-01', '2000-12-31', false, 'data:abc', '2000-01-01', '2000-01-01', NULL),
			(2, '1e', 'Fungsional', 'Analis Kepegawaian', 'Deskripsi Jabatan 2', '2000-01-01', NULL, false, NULL, '2000-01-01', '2000-01-01', NULL),
			(3, '1c', 'Struktural', 'Kepala Sub Bagian', 'Deskripsi Jabatan 3', '2000-01-01', '2000-12-31', false, NULL, '2000-01-01', '2000-01-01', '2000-01-01'),
			(4, '1c', 'Struktural', 'Kepala Bagian', 'Deskripsi Jabatan 4', '2000-01-01', '2000-12-31', false, NULL, '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                int32(4),
			"nip":               "1c",
			"tipe_jabatan":      "Struktural",
			"nama_jabatan":      "Kepala Bagian",
			"deskripsi_jabatan": "Deskripsi Jabatan 4",
			"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
			"is_menjabat":       false,
			"file_base64":       nil,
			"s3_file_id":        nil,
			"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":        nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	defaultRequestBody := func(writer *multipart.Writer) error {
		part, err := writer.CreateFormFile("file", "file.txt")
		if err != nil {
			return err
		}
		_, err = io.WriteString(part, "Hello World!!")
		return err
	}

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name              string
		paramNIP          string
		paramID           string
		requestHeader     http.Header
		appendRequestBody func(writer *multipart.Writer) error
		wantResponseCode  int
		wantResponseBody  string
		wantDBRows        dbtest.Rows
	}{
		{
			name:              "ok: success upload",
			paramNIP:          "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 1",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:              "error: riwayat penugasan is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat penugasan is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(2),
					"nip":               "1e",
					"tipe_jabatan":      "Fungsional",
					"nama_jabatan":      "Analis Kepegawaian",
					"deskripsi_jabatan": "Deskripsi Jabatan 2",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   nil,
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:              "error: riwayat penugasan is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(3),
					"nip":               "1c",
					"tipe_jabatan":      "Struktural",
					"nama_jabatan":      "Kepala Sub Bagian",
					"deskripsi_jabatan": "Deskripsi Jabatan 3",
					"tanggal_mulai":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":   time.Date(2000, 12, 31, 0, 0, 0, 0, time.UTC),
					"is_menjabat":       false,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:              "error: missing file",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"file\" harus diisi"}`,
			wantDBRows:        defaultRows,
		},
		{
			name:              "error: invalid token",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBody:  `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:        defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			require.NoError(t, tt.appendRequestBody(writer))
			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-penugasan/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_penugasan", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}
