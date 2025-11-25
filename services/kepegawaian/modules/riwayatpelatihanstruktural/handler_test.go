package riwayatpelatihanstruktural

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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_diklat_struktural
			(id, pns_nip, nama_diklat, tanggal, lama, nomor, tahun, deleted_at) values
			('uuid-11', '1c', 'Diklat 11', '2000-01-01', 24, 'SK11', 2000, null),
			('uuid-12', '1c', 'Diklat 12', '2001-01-01', 16, 'SK12', 2001, null),
			('uuid-9', '1c', 'Diklat 12', '2001-01-01', 16, 'SK12', 2001, '2001-01-01'),
			('uuid-13', '1c', 'Diklat 13', '2002-01-01', 40, 'SK13', 2002, null),
			('uuid-14', '2c', 'Diklat 14', '2003-01-01', 8,  'SK14', 2003, null),
			('uuid-15', '1c', 'Diklat 15', '2004-01-01', 4,  'SK15', 2004, null);
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
						"id": "uuid-15",
						"nama_diklat": "Diklat 15",
						"tanggal": "2004-01-01",
						"tahun": 2004,
						"nomor": "SK15",
						"lama": 4
					},
					{
						"id": "uuid-13",
						"nama_diklat": "Diklat 13",
						"tanggal": "2002-01-01",
						"tahun": 2002,
						"nomor": "SK13",
						"lama": 40
					},
					{
						"id": "uuid-12",
						"nama_diklat": "Diklat 12",
						"tanggal": "2001-01-01",
						"tahun": 2001,
						"nomor": "SK12",
						"lama": 16
					},
					{
						"id": "uuid-11",
						"nama_diklat": "Diklat 11",
						"tanggal": "2000-01-01",
						"tahun": 2000,
						"nomor": "SK11",
						"lama": 24
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
						"id": "uuid-13",
						"nama_diklat": "Diklat 13",
						"tanggal": "2002-01-01",
						"tahun": 2002,
						"nomor": "SK13",
						"lama": 40
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-struktural", nil)
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
		insert into riwayat_diklat_struktural
			(id, pns_nip, file_base64, deleted_at) values
			('uuid-pdf', '1a', 'data:application/pdf;base64,` + pdfBase64 + `', null),
			('uuid-3', '1b', 'data:application/pdf;base64,` + pdfBase64 + `', null),
			('uuid-png', '1a', 'data:images/png;base64,` + pngBase64 + `', null),
			('uuid-x', '1a', 'data:application/pdf;base64,` + pdfBase64 + `', '2001-01-01'),
			('uuid-inv','1a', 'data:application/pdf;base64,invalid', null),
			('uuid-null', '1a', null, null),
			('uuid-empty', '1a', '', null);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("1a")}
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
			paramID:           "uuid-pdf",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			paramID:           "uuid-png",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 tidak valid",
			paramID:           "uuid-inv",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			paramID:           "uuid-x",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			paramID:           "uuid-null",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			paramID:           "uuid-empty",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: berkas tidak ditemukan",
			paramID:           "uuid-2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: ambil data dari user lain",
			paramID:           "uuid-3",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: auth header tidak valid",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-struktural/%s/berkas", tt.paramID), nil)
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
		insert into riwayat_diklat_struktural
			(id, pns_nip, nama_diklat, tanggal, lama, nomor, tahun, deleted_at) values
			('uuid-11', '1c', 'Diklat 11', '2000-01-01', 24, 'SK11', 2000, null),
			('uuid-12', '1c', 'Diklat 12', '2001-01-01', 16, 'SK12', 2001, null),
			('uuid-13', '1c', 'Diklat 13', '2002-01-01', 40, 'SK13', 2002, null),
			('uuid-14', '1c', 'Diklat 14', '2003-01-01', 8,  'SK14', 2003, null),
			('uuid-15', '1c', 'Diklat 15', '2004-01-01', 4,  'SK15', 2004, null),
			('uuid-x', '1c', 'Diklat 15', '2004-01-01', 4,  'SK15', 2004, '2001-01-01'),
			('uuid-16', '1d', 'Diklat X',  '2010-01-01', 12, 'SK16', 2010, null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			name:             "ok: admin dapat melihat riwayat pelatihan struktural pegawai 1c",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id":"uuid-15","nama_diklat":"Diklat 15","tanggal":"2004-01-01","tahun":2004,"nomor":"SK15","lama":4},
					{"id":"uuid-14","nama_diklat":"Diklat 14","tanggal":"2003-01-01","tahun":2003,"nomor":"SK14","lama":8},
					{"id":"uuid-13","nama_diklat":"Diklat 13","tanggal":"2002-01-01","tahun":2002,"nomor":"SK13","lama":40},
					{"id":"uuid-12","nama_diklat":"Diklat 12","tanggal":"2001-01-01","tahun":2001,"nomor":"SK12","lama":16},
					{"id":"uuid-11","nama_diklat":"Diklat 11","tanggal":"2000-01-01","tahun":2000,"nomor":"SK11","lama":24}
				],
				"meta": {"limit": 10, "offset": 0, "total": 5}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan struktural pegawai 1c dengan pagination",
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id":"uuid-14","nama_diklat":"Diklat 14","tanggal":"2003-01-01","tahun":2003,"nomor":"SK14","lama":8},
					{"id":"uuid-13","nama_diklat":"Diklat 13","tanggal":"2002-01-01","tahun":2002,"nomor":"SK13","lama":40}
				],
				"meta": {"limit": 2, "offset": 1, "total": 5}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan struktural pegawai 1d",
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id":"uuid-16","nama_diklat":"Diklat X","tanggal":"2010-01-01","tahun":2010,"nomor":"SK16","lama":12}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan struktural pegawai yang tidak ada data",
			nip:              "999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pelatihan-struktural", nil)
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
		insert into riwayat_diklat_struktural
			(id, pns_nip, file_base64, deleted_at) values
			('uuid-pdf', '1a', 'data:application/pdf;base64,` + pdfBase64 + `', null),
			('uuid-3', '1b', 'data:application/pdf;base64,` + pdfBase64 + `', null),
			('uuid-png', '1c', 'data:images/png;base64,` + pngBase64 + `', null),
			('uuid-x', '1d', 'data:application/pdf;base64,` + pdfBase64 + `', '2001-01-01'),
			('uuid-inv','1c', 'data:application/pdf;base64,invalid', null),
			('uuid-null', '1x', null, null),
			('uuid-empty', '1s', '', null);
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
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
		nip               string
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			paramID:           "uuid-pdf",
			nip:               "1a",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			paramID:           "uuid-png",
			nip:               "1c",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 tidak valid",
			paramID:           "uuid-inv",
			nip:               "1c",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			paramID:           "uuid-x",
			nip:               "1d",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			paramID:           "uuid-null",
			nip:               "1x",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			paramID:           "uuid-empty",
			nip:               "1s",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: berkas tidak ditemukan",
			paramID:           "uuid-2",
			nip:               "1t",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: auth header tidak valid",
			nip:               "1c",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-pelatihan-struktural/%s/berkas", tt.nip, tt.paramID), nil)
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
			(pns_id,  nip_baru, nama,      deleted_at) values
			('id_1a', '1a',     'User 1a', null),
			('id_1c', '1c',     'User 1c', null),
			('id_1d', '1d',     'User 1d', '2000-01-01'),
			('id_1e', '1e',     'User 1e', null),
			('id_1f', '1f',     'User 1f', null);
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
				"nama_diklat": "Diklat 1",
				"nomor": "SK.01",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "{id}",
					"nama_diklat":       "Diklat 1",
					"jenis_diklat_id":   nil,
					"nomor":             "SK.01",
					"tanggal":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tahun":             int16(2000),
					"lama":              float32(5),
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          "User 1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "",
				"nomor": "",
				"tanggal": "2000-01-01",
				"tahun": 0,
				"lama": 0
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "{id}",
					"nama_diklat":       "",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tahun":             int16(0),
					"lama":              float32(0),
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1e",
					"pns_nip":           "1e",
					"pns_nama":          "User 1e",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "Diklat 1",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "{id}",
					"nama_diklat":       "Diklat 1",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tahun":             int16(2000),
					"lama":              float32(5),
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1f",
					"pns_nip":           "1f",
					"pns_nama":          "User 1f",
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
				"nama_diklat": "",
				"tanggal": "2000-01-01",
				"tahun": 0,
				"lama": 0
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
				"nama_diklat": "",
				"nomor": "",
				"tanggal": "2000-01-01",
				"tahun": 0,
				"lama": 0
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: exceed length limit, unexpected date or data type",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "` + strings.Repeat(".", 201) + `",
				"nomor": "` + strings.Repeat(".", 301) + `",
				"tanggal": "",
				"tahun": "0",
				"lama": "0"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"lama\" harus dalam tipe integer` +
				` | parameter \"nama_diklat\" harus 200 karakter atau kurang` +
				` | parameter \"nomor\" harus 300 karakter atau kurang` +
				` | parameter \"tahun\" harus dalam tipe integer` +
				` | parameter \"tanggal\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null params",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": null,
				"nomor": null,
				"tanggal": null,
				"tahun": null,
				"lama": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"lama\" tidak boleh null` +
				` | parameter \"nama_diklat\" tidak boleh null` +
				` | parameter \"nomor\" tidak boleh null` +
				` | parameter \"tahun\" tidak boleh null` +
				` | parameter \"tanggal\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama_diklat\" harus diisi` +
				` | parameter \"tanggal\" harus diisi` +
				` | parameter \"tahun\" harus diisi` +
				` | parameter \"lama\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-struktural", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat_struktural", "where pns_nip = $1", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["id"] = row["id"]
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]

						tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", fmt.Sprintf("%v", row["id"]))
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
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into ref_jenis_diklat_struktural (id) values (1);
		insert into riwayat_diklat_struktural
			(id,  jenis_diklat_id, siasn_id,  status_data, file_base64, keterangan_berkas, pns_id,  pns_nip, pns_nama,  created_at,   updated_at) values
			('1', 1,               'siasn01', '0',         'data:abc',  'abc',             'id_1c', '1c',    'User 1c', '2000-01-01', '2000-01-01'),
			('2', 1,               'siasn01', '0',         'data:abc',  'abc',             'id_1c', '1c',    'User 1c', '2000-01-01', '2000-01-01'),
			('3', 1,               'siasn01', '0',         'data:abc',  'abc',             'id_1c', '1c',    'User 1c', '2000-01-01', '2000-01-01');
		insert into riwayat_diklat_struktural
			(id,  nama_diklat, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('4', 'Diklat 4',  'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('5', 'Diklat 5',  'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('6', 'Diklat 6',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                "6",
			"nama_diklat":       "Diklat 6",
			"jenis_diklat_id":   nil,
			"nomor":             nil,
			"tanggal":           nil,
			"tahun":             nil,
			"lama":              nil,
			"siasn_id":          nil,
			"file_base64":       nil,
			"s3_file_id":        nil,
			"keterangan_berkas": nil,
			"status_data":       nil,
			"pns_id":            "id_1c",
			"pns_nip":           "1c",
			"pns_nama":          nil,
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
				"nama_diklat": "Diklat 1",
				"nomor": "SK.01",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                "1",
					"nama_diklat":       "Diklat 1",
					"jenis_diklat_id":   int32(1),
					"nomor":             "SK.01",
					"tanggal":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tahun":             int16(2000),
					"lama":              float32(5),
					"siasn_id":          "siasn01",
					"file_base64":       "data:abc",
					"s3_file_id":        nil,
					"keterangan_berkas": "abc",
					"status_data":       "0",
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          "User 1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "",
				"nomor": "",
				"tanggal": "2000-01-01",
				"tahun": 0,
				"lama": 0
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                "2",
					"nama_diklat":       "",
					"jenis_diklat_id":   int32(1),
					"nomor":             nil,
					"tanggal":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tahun":             int16(0),
					"lama":              float32(0),
					"siasn_id":          "siasn01",
					"file_base64":       "data:abc",
					"s3_file_id":        nil,
					"keterangan_berkas": "abc",
					"status_data":       "0",
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          "User 1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "Diklat 1",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                "3",
					"nama_diklat":       "Diklat 1",
					"jenis_diklat_id":   int32(1),
					"nomor":             nil,
					"tanggal":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tahun":             int16(2000),
					"lama":              float32(5),
					"siasn_id":          "siasn01",
					"file_base64":       "data:abc",
					"s3_file_id":        nil,
					"keterangan_berkas": "abc",
					"status_data":       "0",
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          "User 1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: riwayat pelatihan struktural is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "Diklat 1",
				"nomor": "SK.01",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat pelatihan struktural is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "Diklat 1",
				"nomor": "SK.01",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "4",
					"nama_diklat":       "Diklat 4",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1e",
					"pns_nip":           "1e",
					"pns_nama":          nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: riwayat pelatihan struktural is deleted",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "Diklat 1",
				"tanggal": "2000-01-01",
				"tahun": 2000,
				"lama": 5
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "5",
					"nama_diklat":       "Diklat 5",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": "` + strings.Repeat(".", 201) + `",
				"nomor": "` + strings.Repeat(".", 301) + `",
				"tanggal": "",
				"tahun": "0",
				"lama": "0"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"lama\" harus dalam tipe integer` +
				` | parameter \"nama_diklat\" harus 200 karakter atau kurang` +
				` | parameter \"nomor\" harus 300 karakter atau kurang` +
				` | parameter \"tahun\" harus dalam tipe integer` +
				` | parameter \"tanggal\" harus dalam format date"}`,
			wantDBRows: defaultRows,
		},
		{
			name:          "error: null params",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama_diklat": null,
				"nomor": null,
				"tanggal": null,
				"tahun": null,
				"lama": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"lama\" tidak boleh null` +
				` | parameter \"nama_diklat\" tidak boleh null` +
				` | parameter \"nomor\" tidak boleh null` +
				` | parameter \"tahun\" tidak boleh null` +
				` | parameter \"tanggal\" tidak boleh null"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama_diklat\" harus diisi` +
				` | parameter \"tanggal\" harus diisi` +
				` | parameter \"tahun\" harus diisi` +
				` | parameter \"lama\" harus diisi"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-struktural/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat_struktural", "where id = $1", tt.paramID)
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
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_diklat_struktural
			(id,  nama_diklat, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('1', 'Diklat 1',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null),
			('2', 'Diklat 2',  'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('3', 'Diklat 3',  'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Diklat 4',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                "4",
			"nama_diklat":       "Diklat 4",
			"jenis_diklat_id":   nil,
			"nomor":             nil,
			"tanggal":           nil,
			"tahun":             nil,
			"lama":              nil,
			"siasn_id":          nil,
			"file_base64":       nil,
			"s3_file_id":        nil,
			"keterangan_berkas": nil,
			"status_data":       nil,
			"pns_id":            "id_1c",
			"pns_nip":           "1c",
			"pns_nama":          nil,
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
					"id":                "1",
					"nama_diklat":       "Diklat 1",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat pelatihan struktural is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "2",
					"nama_diklat":       "Diklat 2",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1e",
					"pns_nip":           "1e",
					"pns_nama":          nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:             "error: riwayat pelatihan struktural is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat pelatihan struktural is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "3",
					"nama_diklat":       "Diklat 3",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          nil,
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
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-struktural/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat_struktural", "where id = $1", tt.paramID)
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
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_diklat_struktural
			(id,  nama_diklat, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('1', 'Diklat 1',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null),
			('2', 'Diklat 2',  'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('3', 'Diklat 3',  'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Diklat 4',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                "4",
			"nama_diklat":       "Diklat 4",
			"jenis_diklat_id":   nil,
			"nomor":             nil,
			"tanggal":           nil,
			"tahun":             nil,
			"lama":              nil,
			"siasn_id":          nil,
			"file_base64":       nil,
			"s3_file_id":        nil,
			"keterangan_berkas": nil,
			"status_data":       nil,
			"pns_id":            "id_1c",
			"pns_nip":           "1c",
			"pns_nama":          nil,
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
					"id":                "1",
					"nama_diklat":       "Diklat 1",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:              "error: riwayat pelatihan struktural is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat pelatihan struktural is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "2",
					"nama_diklat":       "Diklat 2",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1e",
					"pns_nip":           "1e",
					"pns_nama":          nil,
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:              "error: riwayat pelatihan struktural is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                "3",
					"nama_diklat":       "Diklat 3",
					"jenis_diklat_id":   nil,
					"nomor":             nil,
					"tanggal":           nil,
					"tahun":             nil,
					"lama":              nil,
					"siasn_id":          nil,
					"file_base64":       nil,
					"s3_file_id":        nil,
					"keterangan_berkas": nil,
					"status_data":       nil,
					"pns_id":            "id_1c",
					"pns_nip":           "1c",
					"pns_nama":          nil,
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-struktural/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat_struktural", "where id = $1", tt.paramID)
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
