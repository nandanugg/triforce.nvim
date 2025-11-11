package riwayatsertifikasi

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, file_base64, created_at,   deskripsi, deleted_at) values
			(11, '1c', 1,     '11a',            '11b',       '2000-01-01', '11c',     null),
			(12, '1c', 3,     '12a',            '12b',       '2001-01-01', null,      null),
			(13, '1c', 2,     '13a',            '13b',       '2002-01-01', '13c',     null),
			(14, '2c', 4,     '14a',            '14b',       '2003-01-01', '14c',     null),
			(15, '1c', 5,     '15a',            '15b',       '2003-01-01', '15c',     '2020-01-01');
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
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
						"id":               12,
						"nama_sertifikasi": "12a",
						"tahun":            3,
						"deskripsi":        ""
					},
					{
						"id":               13,
						"nama_sertifikasi": "13a",
						"tahun":            2,
						"deskripsi":        "13c"
					},
					{
						"id":               11,
						"nama_sertifikasi": "11a",
						"tahun":            1,
						"deskripsi":        "11c"
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
						"id":               13,
						"nama_sertifikasi": "13a",
						"tahun":            2,
						"deskripsi":        "13c"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-sertifikasi", nil)
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
		insert into riwayat_sertifikasi
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
			name:              "error: base64 riwayat sertifikasi tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sertifikasi sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat sertifikasi berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat sertifikasi berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat sertifikasi bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat sertifikasi tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-sertifikasi/%s/berkas", tt.paramID), nil)
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
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, file_base64, created_at,   deskripsi, deleted_at) values
			(11, '1c', 1,     '11a',            '11b',       '2000-01-01', '11c',     null),
			(12, '1c', 3,     '12a',            '12b',       '2001-01-01', null,      null),
			(13, '1c', 2,     '13a',            '13b',       '2002-01-01', '13c',     null),
			(14, '2c', 4,     '14a',            '14b',       '2003-01-01', '14c',     null),
			(15, '1c', 5,     '15a',            '15b',       '2003-01-01', '15c',     '2020-01-01');
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
			name:             "ok: nip 1c data returned",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":               12,
						"nama_sertifikasi": "12a",
						"tahun":            3,
						"deskripsi":        ""
					},
					{
						"id":               13,
						"nama_sertifikasi": "13a",
						"tahun":            2,
						"deskripsi":        "13c"
					},
					{
						"id":               11,
						"nama_sertifikasi": "11a",
						"tahun":            1,
						"deskripsi":        "11c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
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
						"id":               13,
						"nama_sertifikasi": "13a",
						"tahun":            2,
						"deskripsi":        "13c"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: nip 200 gets empty data",
			nip:              "200",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-sertifikasi", nil)
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
		insert into riwayat_sertifikasi
			(id, nip, deleted_at, file_base64) values
			(1, '1c', null, 'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c', null, '` + pdfBase64 + `'),
			(3, '1c', null, 'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c', null, 'data:application/pdf;base64,invalid'),
			(5, '1c', '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c', null, null),
			(7, '1c', null, ''),
			(8, '2a', null, 'data:application/pdf;base64,` + pdfBase64 + `');
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
		nip               string
		paramID           string
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
			name:              "ok: admin can access other user's berkas",
			nip:               "2a",
			paramID:           "8",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "error: base64 tidak valid",
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat with wrong nip",
			nip:               "wrong-nip",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat tidak ditemukan",
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat sertifikasi tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-sertifikasi/%s/berkas", tt.nip, tt.paramID), nil)
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
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null),
			('id_1f', '1f',     null),
			('id_1g', '1g',     null);
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, deskripsi, created_at,   updated_at) values
			(1,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01');
		SELECT setval('riwayat_sertifikasi_id_seq', 1, true);
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
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Baru",
				"deskripsi": "Deskripsi sertifikasi baru"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(1),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        "Deskripsi 1",
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
				{
					"id":               "{id}",
					"nip":              "1c",
					"tahun":            int64(2025),
					"nama_sertifikasi": "Sertifikasi Baru",
					"deskripsi":        "Deskripsi sertifikasi baru",
					"file_base64":      nil,
					"created_at":       "{created_at}",
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "ok: with null/empty values",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2020,
				"nama_sertifikasi": "Sertifikasi",
				"deskripsi": ""
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               "{id}",
					"nip":              "1e",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       "{created_at}",
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2020,
				"nama_sertifikasi": "Sertifikasi"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               "{id}",
					"nip":              "1f",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       "{created_at}",
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2020,
				"nama_sertifikasi": "Sertifikasi"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2020,
				"nama_sertifikasi": "Sertifikasi",
				"deskripsi": "Deskripsi"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: missing required params",
			paramNIP:         "1g",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "deskripsi": "Deskripsi" }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"tahun\" harus diisi` +
				` | parameter \"nama_sertifikasi\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: wrong data type",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": "2020",
				"nama_sertifikasi": 123
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama_sertifikasi\" harus dalam tipe string` +
				` | parameter \"tahun\" harus dalam tipe integer"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null required params",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": null,
				"nama_sertifikasi": null,
				"deskripsi": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"deskripsi\" tidak boleh null` +
				` | parameter \"nama_sertifikasi\" tidak boleh null` +
				` | parameter \"tahun\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"tahun\" harus diisi` +
				` | parameter \"nama_sertifikasi\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-sertifikasi", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if strings.Contains(tt.wantResponseBody, "{id}") {
				var resp struct {
					Data struct {
						ID int64 `json:"id"`
					} `json:"data"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				id := resp.Data.ID
				tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", strconv.FormatInt(id, 10))
				for i, row := range tt.wantDBRows {
					if row["id"] == "{id}" {
						tt.wantDBRows[i]["id"] = id
					}
				}
			}
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_sertifikasi", "where nip = $1 order by id", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_sertifikasi (id, nip, tahun, nama_sertifikasi, deskripsi, created_at, updated_at, deleted_at) values
			(1,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null),
			(2,  '1c', 2021,  'Sertifikasi 2',   'Deskripsi 2', '2000-01-01', '2000-01-01', null),
			(3,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null),
			(4,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null),
			(5,  '1c', 2020,  'Sertifikasi 1',   null,           '2000-01-01', '2000-01-01', null),
			(6,  '1e', 2020,  'Sertifikasi 1',   null,           '2000-01-01', '2000-01-01', null),
			(7,  '1c', 2020,  'Sertifikasi 1',   null,           '2000-01-01', '2000-01-01', '2000-01-01'),
			(8,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null),
			(9,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null),
			(10, '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null),
			(11, '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', '2000-01-01', '2000-01-01', null);
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
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Updated",
				"deskripsi": "Deskripsi Updated"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(1),
					"nip":              "1c",
					"tahun":            int64(2025),
					"nama_sertifikasi": "Sertifikasi Updated",
					"deskripsi":        "Deskripsi Updated",
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
				{
					"id":               int64(2),
					"nip":              "1c",
					"tahun":            int64(2021),
					"nama_sertifikasi": "Sertifikasi 2",
					"deskripsi":        "Deskripsi 2",
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "ok: with null/empty values",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Updated",
				"deskripsi": ""
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(3),
					"nip":              "1c",
					"tahun":            int64(2025),
					"nama_sertifikasi": "Sertifikasi Updated",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Updated"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(4),
					"nip":              "1c",
					"tahun":            int64(2025),
					"nama_sertifikasi": "Sertifikasi Updated",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "error: riwayat sertifikasi is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Updated"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat sertifikasi is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Updated"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(6),
					"nip":              "1e",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
			},
		},
		{
			name:          "error: riwayat sertifikasi is deleted",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": 2025,
				"nama_sertifikasi": "Sertifikasi Updated",
				"deskripsi": "Deskripsi Updated"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(7),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: missing required params",
			paramNIP:         "1c",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "deskripsi": "Deskripsi" }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"tahun\" harus diisi` +
				` | parameter \"nama_sertifikasi\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: wrong data type",
			paramNIP:      "1c",
			paramID:       "9",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": "2020",
				"nama_sertifikasi": 123
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama_sertifikasi\" harus dalam tipe string` +
				` | parameter \"tahun\" harus dalam tipe integer"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null required params",
			paramNIP:      "1c",
			paramID:       "10",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tahun": null,
				"nama_sertifikasi": null,
				"deskripsi": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"deskripsi\" tidak boleh null` +
				` | parameter \"nama_sertifikasi\" tidak boleh null` +
				` | parameter \"tahun\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid id",
			paramNIP:         "1c",
			paramID:          "abc",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"tahun": 2025, "nama_sertifikasi": "Test"}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			requestBody:      `{"tahun": 2025, "nama_sertifikasi": "Test"}`,
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-sertifikasi/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if len(tt.wantDBRows) > 0 {
				expectedIDs := make([]int64, len(tt.wantDBRows))
				for i, wantRow := range tt.wantDBRows {
					expectedIDs[i] = wantRow["id"].(int64)
				}
				actualRows, err := dbtest.QueryWithClause(db, "riwayat_sertifikasi", "where id = any($1::int8[]) order by id", expectedIDs)
				require.NoError(t, err)
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
				assert.Equal(t, tt.wantDBRows, actualRows)
			}
		})
	}
}

func Test_handler_adminDelete(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_sertifikasi
			(id, nip, tahun, nama_sertifikasi, created_at, updated_at, deleted_at) values
			(1, '1c', 2020, 'Sertifikasi 1', '2000-01-01', '2000-01-01', null),
			(2, '1c', 2021, 'Sertifikasi 2', '2000-01-01', '2000-01-01', null),
			(3, '1e', 2020, 'Sertifikasi 1', '2000-01-01', '2000-01-01', null),
			(4, '1c', 2020, 'Sertifikasi 1', '2000-01-01', '2000-01-01', null),
			(5, '1c', 2020, 'Sertifikasi 1', '2000-01-01', '2000-01-01', '2000-01-01');
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
					"id":               int64(1),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       "{deleted_at}",
				},
				{
					"id":               int64(2),
					"nip":              "1c",
					"tahun":            int64(2021),
					"nama_sertifikasi": "Sertifikasi 2",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
			},
		},
		{
			name:             "error: riwayat sertifikasi is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(3),
					"nip":              "1e",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
			},
		},
		{
			name:             "error: riwayat sertifikasi is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(4),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
			},
		},
		{
			name:             "error: riwayat sertifikasi is deleted",
			paramNIP:         "1c",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(5),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: unexpected data type",
			paramNIP:         "1c",
			paramID:          "abc",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-sertifikasi/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if len(tt.wantDBRows) > 0 {
				expectedIDs := make([]int64, len(tt.wantDBRows))
				for i, wantRow := range tt.wantDBRows {
					expectedIDs[i] = wantRow["id"].(int64)
				}
				actualRows, err := dbtest.QueryWithClause(db, "riwayat_sertifikasi", "where id = any($1::int8[]) order by id", expectedIDs)
				require.NoError(t, err)
				for i, row := range actualRows {
					if tt.wantDBRows[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["deleted_at"] = row["deleted_at"]
					}
				}
				assert.Equal(t, tt.wantDBRows, actualRows)
			}
		})
	}
}

func Test_handler_adminUploadBerkas(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, deskripsi, file_base64, created_at,   updated_at) values
			(1,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', 'data:abc', '2000-01-01', '2000-01-01'),
			(2,  '1c', 2021,  'Sertifikasi 2',   'Deskripsi 2', 'data:abc', '2000-01-01', '2000-01-01');
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, created_at,   updated_at) values
			(3,  '1c', 2020,  'Sertifikasi 1',   '2000-01-01', '2000-01-01');
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, created_at,   updated_at,   deleted_at) values
			(4,  '1c', 2020,  'Sertifikasi 1',   '2000-01-01', '2000-01-01', '2000-01-01');
		insert into riwayat_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, deskripsi, file_base64, created_at,   updated_at) values
			(5,  '1c', 2020,  'Sertifikasi 1',   'Deskripsi 1', 'data:abc', '2000-01-01', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":               int64(5),
			"nip":              "1c",
			"tahun":            int64(2020),
			"nama_sertifikasi": "Sertifikasi 1",
			"deskripsi":        "Deskripsi 1",
			"file_base64":      "data:abc",
			"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":       nil,
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
					"id":               int64(1),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        "Deskripsi 1",
					"file_base64":      "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       "{updated_at}",
					"deleted_at":       nil,
				},
				{
					"id":               int64(2),
					"nip":              "1c",
					"tahun":            int64(2021),
					"nama_sertifikasi": "Sertifikasi 2",
					"deskripsi":        "Deskripsi 2",
					"file_base64":      "data:abc",
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       nil,
				},
			},
		},
		{
			name:              "error: riwayat sertifikasi is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat sertifikasi is deleted",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":               int64(4),
					"nip":              "1c",
					"tahun":            int64(2020),
					"nama_sertifikasi": "Sertifikasi 1",
					"deskripsi":        nil,
					"file_base64":      nil,
					"created_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":       time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:              "error: missing file",
			paramNIP:          "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"file\" harus diisi"}`,
			wantDBRows:        defaultRows,
		},
		{
			name:              "error: invalid id",
			paramNIP:          "1c",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: invalid token",
			paramNIP:          "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBody:  `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:        dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			require.NoError(t, tt.appendRequestBody(writer))
			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-sertifikasi/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if len(tt.wantDBRows) > 0 {
				expectedIDs := make([]int64, len(tt.wantDBRows))
				for i, wantRow := range tt.wantDBRows {
					expectedIDs[i] = wantRow["id"].(int64)
				}
				actualRows, err := dbtest.QueryWithClause(db, "riwayat_sertifikasi", "where id = any($1::int8[]) order by id", expectedIDs)
				require.NoError(t, err)
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
				assert.Equal(t, tt.wantDBRows, actualRows)
			}
		})
	}
}
