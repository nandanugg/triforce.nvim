package riwayatjabatan

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
	dbrepository "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan(id, nama_jabatan, kode_jabatan, deleted_at) values
		(11, '11h', '11', null),
		(12, '12h', '12', null),
		(13, '13h', '13', '2000-01-01');

		insert into ref_jenis_jabatan(id, nama, deleted_at) values
		(1, 'Jabatan Struktural', null),
		(2, 'Jabatan Fungsional', null),
		(3, 'Jabatan Deleted', '2000-01-01');

		insert into ref_unit_kerja(id, nama_unor, deleted_at) values
		(1, 'Unit 1', null),
		(2, 'Unit 2', null),
		(3, 'Unit 3', '2000-01-01');

		insert into ref_kelas_jabatan(id, kelas_jabatan, tunjangan_kinerja, deleted_at) values
		(1, 'Kelas 1', 2531250, null),
		(2, 'Kelas 2', 2708250, null),
		(3, 'Kelas 3', 1231312, '2000-01-01');

		insert into riwayat_jabatan(id, pns_nip, jenis_jabatan_id, jabatan_id, tmt_jabatan, no_sk, tanggal_sk, satuan_kerja_id, unor_id, kelas_jabatan_id, periode_jabatan_start_date, periode_jabatan_end_date, deleted_at) values
		(1, '41', 1, '11', '2025-01-01', '1234567890', '2025-01-01', 1, 1, 1, '2024-01-01', '2024-12-31', null),
		(2, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(3, '42', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(4, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', '2000-01-01'),
		(5, '41', 3, '13', '2024-02-01', '2234567890', '2024-02-01', 3, 3, 3, '2025-01-01', '2025-12-31', null),
		(6, '41', null, null, '2023-02-01', '2234567890', '2024-02-01', null, null, null, '2025-01-01', '2025-12-31', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepository.New(db)
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
						"id":                         2,
						"jenis_jabatan_id":           2,
						"jenis_jabatan":              "Jabatan Fungsional",
						"nama_jabatan":               "12h",
						"id_jabatan":                 "12",
						"tmt_jabatan":                "2025-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2025-02-01",
						"satuan_kerja_id":            "2",
						"satuan_kerja":               "Unit 2",
						"status_plt":                 false,
						"kelas_jabatan_id":           2,
						"kelas_jabatan":              "Kelas 2",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi_id":         "2",
						"unit_organisasi":            "Unit 2"
					},
					{
						"id":                         1,
						"jenis_jabatan_id":           1,
						"jenis_jabatan":              "Jabatan Struktural",
						"nama_jabatan":               "11h",
						"id_jabatan":                 "11",
						"tmt_jabatan":                "2025-01-01",
						"no_sk":                      "1234567890",
						"tanggal_sk":                 "2025-01-01",
						"satuan_kerja_id":            "1",
						"satuan_kerja":               "Unit 1",
						"status_plt":                 false,
						"kelas_jabatan_id":           1,
						"kelas_jabatan":              "Kelas 1",
						"periode_jabatan_start_date": "2024-01-01",
						"periode_jabatan_end_date":   "2024-12-31",
						"unit_organisasi_id":         "1",
						"unit_organisasi":            "Unit 1"
					},
					{
						"id":                         5,
						"jenis_jabatan_id":           3,
						"jenis_jabatan":              "",
						"nama_jabatan":               "",
						"id_jabatan":                 "13",
						"tmt_jabatan":                "2024-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2024-02-01",
						"satuan_kerja_id":            "3",
						"satuan_kerja":               "",
						"status_plt":                 false,
						"kelas_jabatan_id":           3,
						"kelas_jabatan":              "",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi_id":         "3",
						"unit_organisasi":            ""
					},
					{
						"id":                         6,
						"jenis_jabatan_id":           null,
						"jenis_jabatan":              "",
						"nama_jabatan":               "",
						"id_jabatan":                 null,
						"tmt_jabatan":                "2023-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2024-02-01",
						"satuan_kerja_id":            null,
						"satuan_kerja":               "",
						"status_plt":                 false,
						"kelas_jabatan_id":           null,
						"kelas_jabatan":              "",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi_id":         null,
						"unit_organisasi":            ""
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
						"id":                         1,
						"jenis_jabatan_id":           1,
						"jenis_jabatan":              "Jabatan Struktural",
						"nama_jabatan":               "11h",
						"id_jabatan":                 "11",
						"tmt_jabatan":                "2025-01-01",
						"no_sk":                      "1234567890",
						"tanggal_sk":                 "2025-01-01",
						"satuan_kerja_id":            "1",
						"satuan_kerja":               "Unit 1",
						"status_plt":                 false,
						"kelas_jabatan_id":           1,
						"kelas_jabatan":              "Kelas 1",
						"periode_jabatan_start_date": "2024-01-01",
						"periode_jabatan_end_date":   "2024-12-31",
						"unit_organisasi_id":         "1",
						"unit_organisasi":            "Unit 1"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}`,
		},
		{
			name:             "ok: tidak ada data riwayat jabatan",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-jabatan", nil)
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
		insert into riwayat_jabatan
			(id, pns_nip, deleted_at,   file_base64) values
			(1, '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c',     null,         '` + pdfBase64 + `'),
			(3, '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c',     null,         'data:application/pdf;base64,invalid'),
			(5, '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c',     null,         null),
			(7, '1c',     null,         '');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepository.New(pgxconn)
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
			name:              "error: base64 riwayat jabatan tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat jabatan sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat jabatan berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat jabatan berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat jabatan bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat jabatan tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-jabatan/%s/berkas", tt.paramID), nil)
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
		insert into ref_jabatan(id, nama_jabatan, kode_jabatan, deleted_at) values
		(11, '11h', '11', null),
		(12, '12h', '12', null),
		(13, '13h', '13', '2000-01-01');

		insert into ref_jenis_jabatan(id, nama, deleted_at) values
		(1, 'Jabatan Struktural', null),
		(2, 'Jabatan Fungsional', null),
		(3, 'Jabatan Deleted', '2000-01-01');

		insert into ref_unit_kerja(id, nama_unor, deleted_at) values
		(1, 'Unit 1', null),
		(2, 'Unit 2', null),
		(3, 'Unit 3', '2000-01-01');

		insert into ref_kelas_jabatan(id, kelas_jabatan, tunjangan_kinerja, deleted_at) values
		(1, 'Kelas 1', 2531250, null),
		(2, 'Kelas 2', 2708250, null),
		(3, 'Kelas 3', 2999912, '2000-01-01');

		insert into pegawai (pns_id, nip_baru, nama, deleted_at)
		values ('pns-1', '41', 'Pegawai Test', null),
		('pns-2', '42', 'Pegawai Test 2', now());

		insert into riwayat_jabatan(id, pns_nip, jenis_jabatan_id, jabatan_id, tmt_jabatan, no_sk, tanggal_sk, satuan_kerja_id, unor_id, kelas_jabatan_id, periode_jabatan_start_date, periode_jabatan_end_date, deleted_at) values
		(1, '41', 1, '11', '2025-01-01', '1234567890', '2025-01-01', 1, 1, 1, '2024-01-01', '2024-12-31', null),
		(2, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(3, '42', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(4, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', '2000-01-01'),
		(5, '41', 3, '13', '2024-02-01', '2234567890', '2024-02-01', 3, 3, 3, '2025-01-01', '2025-12-31', null),
		(6, '41', null, null, '2023-02-01', '2234567890', '2024-02-01', null, null, null, '2025-01-01', '2025-12-31', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepository.New(db)
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
			name:             "ok: nip 41 data returned",
			nip:              "41",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                         2,
						"jenis_jabatan_id":           2,
						"jenis_jabatan":              "Jabatan Fungsional",
						"nama_jabatan":               "12h",
						"id_jabatan":                 "12",
						"tmt_jabatan":                "2025-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2025-02-01",
						"satuan_kerja_id":            "2",
						"satuan_kerja":               "Unit 2",
						"status_plt":                 false,
						"kelas_jabatan_id":           2,
						"kelas_jabatan":              "Kelas 2",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi_id":         "2",
						"unit_organisasi":            "Unit 2"
					},
					{
						"id":                         1,
						"jenis_jabatan_id":           1,
						"jenis_jabatan":              "Jabatan Struktural",
						"nama_jabatan":               "11h",
						"id_jabatan":                 "11",
						"tmt_jabatan":                "2025-01-01",
						"no_sk":                      "1234567890",
						"tanggal_sk":                 "2025-01-01",
						"satuan_kerja_id":            "1",
						"satuan_kerja":               "Unit 1",
						"status_plt":                 false,
						"kelas_jabatan_id":           1,
						"kelas_jabatan":              "Kelas 1",
						"periode_jabatan_start_date": "2024-01-01",
						"periode_jabatan_end_date":   "2024-12-31",
						"unit_organisasi_id":         "1",
						"unit_organisasi":            "Unit 1"
					},
					{
						"id":                         5,
						"jenis_jabatan_id":           3,
						"jenis_jabatan":              "",
						"nama_jabatan":               "",
						"id_jabatan":                 "13",
						"tmt_jabatan":                "2024-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2024-02-01",
						"satuan_kerja_id":            "3",
						"satuan_kerja":               "",
						"status_plt":                 false,
						"kelas_jabatan_id":           3,
						"kelas_jabatan":              "",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi_id":         "3",
						"unit_organisasi":            ""
					},
					{
						"id":                         6,
						"jenis_jabatan_id":           null,
						"jenis_jabatan":              "",
						"nama_jabatan":               "",
						"id_jabatan":                 null,
						"tmt_jabatan":                "2023-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2024-02-01",
						"satuan_kerja_id":            null,
						"satuan_kerja":               "",
						"status_plt":                 false,
						"kelas_jabatan_id":           null,
						"kelas_jabatan":              "",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi_id":         null,
						"unit_organisasi":            ""
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
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
						"id":                         1,
						"jenis_jabatan_id":           1,
						"jenis_jabatan":              "Jabatan Struktural",
						"nama_jabatan":               "11h",
						"id_jabatan":                 "11",
						"tmt_jabatan":                "2025-01-01",
						"no_sk":                      "1234567890",
						"tanggal_sk":                 "2025-01-01",
						"satuan_kerja_id":            "1",
						"satuan_kerja":               "Unit 1",
						"status_plt":                 false,
						"kelas_jabatan_id":           1,
						"kelas_jabatan":              "Kelas 1",
						"periode_jabatan_start_date": "2024-01-01",
						"periode_jabatan_end_date":   "2024-12-31",
						"unit_organisasi_id":         "1",
						"unit_organisasi":            "Unit 1"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
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
			name:             "ok: nip 42 gets empty data (deleted pegawai)",
			nip:              "42",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [{"id": 3, "jenis_jabatan_id": 2, "jenis_jabatan": "Jabatan Fungsional", "nama_jabatan": "12h", "id_jabatan": "12", "tmt_jabatan": "2025-02-01", "no_sk": "2234567890", "tanggal_sk": "2025-02-01", "satuan_kerja_id": "2", "satuan_kerja": "Unit 2", "unit_organisasi_id": "2", "unit_organisasi": "Unit 2", "status_plt": false, "kelas_jabatan_id": 2, "kelas_jabatan": "Kelas 2", "periode_jabatan_start_date": "2025-01-01", "periode_jabatan_end_date": "2025-12-31"}], "meta": {"limit": 10, "offset": 0, "total": 1}}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-jabatan", nil)
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
		insert into riwayat_jabatan
			(id, pns_nip, deleted_at,   file_base64) values
			(1, '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c',     null,         '` + pdfBase64 + `'),
			(3, '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c',     null,         'data:application/pdf;base64,invalid'),
			(5, '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c',     null,         null),
			(7, '1c',     null,         ''),
			(8, '2a',     null,         'data:application/pdf;base64,` + pdfBase64 + `');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepository.New(pgxconn)
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
			name:              "error: base64 riwayat jabatan tidak valid",
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat jabatan sudah dihapus",
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat jabatan berisi null value",
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat jabatan berupa string kosong",
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat jabatan with wrong nip",
			nip:               "wrong-nip",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat jabatan tidak ditemukan",
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-jabatan/%s/berkas", tt.nip, tt.paramID), nil)
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
		insert into ref_jenis_jabatan
			(id,  nama,      deleted_at) values
			('1', 'Jenis 1', null),
			('2', 'Jenis 2', '2000-01-01');
		insert into ref_jabatan
			(id,  no, kode_jabatan, nama_jabatan, kode_bkn, deleted_at) values
			('2', 1,  'K1',         'Jabatan 1',  'BKN.1',  null),
			('3', 2,  'K2',         'Jabatan 2',  'BKN.2',  '2000-01-01');
		insert into ref_unit_kerja
			(id,  nama_unor, deleted_at) values
			('3', 'Unor 1',  null),
			('4', 'Unor 2',  '2000-01-01'),
			('5', 'Unor 3',  null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, dbrepository.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
				"jenis_jabatan_id": 1,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-02",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         "{id}",
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            "3",
					"unor_id":                    "5",
					"unor_id_bkn":                "5",
					"unor":                       "Unor 3",
					"jenis_jabatan_id":           int32(1),
					"jenis_jabatan":              "Jenis 1",
					"jabatan_id":                 "K1",
					"jabatan_id_bkn":             "BKN.1",
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_plt":                 false,
					"periode_jabatan_start_date": time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					"periode_jabatan_end_date":   time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
					"no_sk":                      "SK.01",
					"tanggal_sk":                 time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   "User 1c",
					"created_at":                 "{created_at}",
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": null,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "",
				"tanggal_sk": "2000-01-02",
				"status_plt": null,
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         "{id}",
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            "3",
					"unor_id":                    "5",
					"unor_id_bkn":                "5",
					"unor":                       "Unor 3",
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 "K1",
					"jabatan_id_bkn":             "BKN.1",
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_plt":                 nil,
					"periode_jabatan_start_date": time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					"periode_jabatan_end_date":   time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
					"no_sk":                      "",
					"tanggal_sk":                 time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1e",
					"pns_nip":                    "1e",
					"pns_nama":                   "User 1e",
					"created_at":                 "{created_at}",
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-02",
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         "{id}",
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            "3",
					"unor_id":                    "5",
					"unor_id_bkn":                "5",
					"unor":                       "Unor 3",
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 "K1",
					"jabatan_id_bkn":             "BKN.1",
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_plt":                 nil,
					"periode_jabatan_start_date": time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					"periode_jabatan_end_date":   time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
					"no_sk":                      "SK.01",
					"tanggal_sk":                 time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1f",
					"pns_nip":                    "1f",
					"pns_nama":                   "User 1f",
					"created_at":                 "{created_at}",
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 1,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
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
				"jenis_jabatan_id": 1,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: jenis jabatan or jabatan or satuan kerja or unit organisasi is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 0,
				"id_jabatan": "K0",
				"satuan_kerja_id": "0",
				"unit_organisasi_id": "0",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis jabatan tidak ditemukan | data jabatan tidak ditemukan | data satuan kerja tidak ditemukan | data unit organisasi tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: jenis jabatan or jabatan or satuan kerja or unit organisasi is deleted",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 2,
				"id_jabatan": "K2",
				"satuan_kerja_id": "4",
				"unit_organisasi_id": "4",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis jabatan tidak ditemukan | data jabatan tidak ditemukan | data satuan kerja tidak ditemukan | data unit organisasi tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: exceed length limit, unexpected date or data type",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": "0",
				"id_jabatan": 1,
				"satuan_kerja_id": 1,
				"unit_organisasi_id": 1,
				"tmt_jabatan": "",
				"no_sk": "` + strings.Repeat(".", 101) + `",
				"tanggal_sk": "",
				"status_plt": "false",
				"periode_jabatan_start_date": "",
				"periode_jabatan_end_date": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id_jabatan\" harus dalam tipe string` +
				` | parameter \"jenis_jabatan_id\" harus dalam tipe integer` +
				` | parameter \"no_sk\" harus 100 karakter atau kurang` +
				` | parameter \"periode_jabatan_end_date\" harus dalam format date` +
				` | parameter \"periode_jabatan_start_date\" harus dalam format date` +
				` | parameter \"satuan_kerja_id\" harus dalam tipe string` +
				` | parameter \"status_plt\" harus dalam tipe boolean` +
				` | parameter \"tanggal_sk\" harus dalam format date` +
				` | parameter \"tmt_jabatan\" harus dalam format date` +
				` | parameter \"unit_organisasi_id\" harus dalam tipe string"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null params",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": null,
				"id_jabatan": null,
				"satuan_kerja_id": null,
				"unit_organisasi_id": null,
				"tmt_jabatan": null,
				"no_sk": null,
				"tanggal_sk": null,
				"status_plt": null,
				"periode_jabatan_start_date": null,
				"periode_jabatan_end_date": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id_jabatan\" tidak boleh null` +
				` | parameter \"no_sk\" tidak boleh null` +
				` | parameter \"periode_jabatan_end_date\" tidak boleh null` +
				` | parameter \"periode_jabatan_start_date\" tidak boleh null` +
				` | parameter \"satuan_kerja_id\" tidak boleh null` +
				` | parameter \"tanggal_sk\" tidak boleh null` +
				` | parameter \"tmt_jabatan\" tidak boleh null` +
				` | parameter \"unit_organisasi_id\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"id_jabatan\" harus diisi` +
				` | parameter \"satuan_kerja_id\" harus diisi` +
				` | parameter \"unit_organisasi_id\" harus diisi` +
				` | parameter \"periode_jabatan_start_date\" harus diisi` +
				` | parameter \"periode_jabatan_end_date\" harus diisi` +
				` | parameter \"no_sk\" harus diisi` +
				` | parameter \"tanggal_sk\" harus diisi` +
				` | parameter \"tmt_jabatan\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-jabatan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_jabatan", "where pns_nip = $1", tt.paramNIP)
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
		insert into pegawai
			(pns_id,  nip_baru, nama,      deleted_at) values
			('id_1c', '1c',     'User 1c', null),
			('id_1d', '1d',     'User 1d', '2000-01-01'),
			('id_1e', '1e',     'User 1e', null);
		insert into ref_jenis_jabatan
			(id,  nama,      deleted_at) values
			('1', 'Jenis 1', null),
			('2', 'Jenis 2', '2000-01-01');
		insert into ref_jabatan
			(id,  no, kode_jabatan, nama_jabatan, kode_bkn, deleted_at) values
			('2', 1,  'K1',         'Jabatan 1',  'BKN.1',  null),
			('3', 2,  'K2',         'Jabatan 2',  'BKN.2',  '2000-01-01');
		insert into ref_unit_kerja
			(id,  nama_unor, deleted_at) values
			('3', 'Unor 1',  null),
			('4', 'Unor 2',  '2000-01-01'),
			('5', 'Unor 3',  null);
		insert into ref_kelas_jabatan (id) values (1);
		insert into riwayat_jabatan
			(id,  tmt_pelantikan, is_active, eselon_id, eselon, eselon1, eselon2, eselon3, eselon4, catatan, jenis_sk, status_satker, status_biro, kelas_jabatan_id, bkn_id, file_base64, keterangan_berkas, tabel_mutasi_id, pns_id,  pns_nip, created_at,   updated_at) values
			('1', '1999-01-01',   0,         'ide0',    'e0',   'e1',    'e2',    'e3',    'e4',    'ket',   '1',      0,             0,           1,                'bkn1', 'data:abc',  'abc',             1,               'id_1c', '1c',    '2000-01-01', '2000-01-01'),
			('2', '1999-01-01',   0,         'ide0',    'e0',   'e1',    'e2',    'e3',    'e4',    'ket',   '1',      0,             0,           1,                'bkn1', 'data:abc',  'abc',             1,               'id_1c', '1c',    '2000-01-01', '2000-01-01'),
			('3', '1999-01-01',   0,         'ide0',    'e0',   'e1',    'e2',    'e3',    'e4',    'ket',   '1',      0,             0,           1,                'bkn1', 'data:abc',  'abc',             1,               'id_1c', '1c',    '2000-01-01', '2000-01-01');
		insert into riwayat_jabatan
			(id,  nama_jabatan, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('4', 'Jabatan 4',  'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('5', 'Jabatan 5',  'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('6', 'Jabatan 6',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                         int64(6),
			"kelas_jabatan_id":           nil,
			"satuan_kerja_id":            nil,
			"unor_id":                    nil,
			"unor_id_bkn":                nil,
			"unor":                       nil,
			"jenis_jabatan_id":           nil,
			"jenis_jabatan":              nil,
			"jabatan_id":                 nil,
			"jabatan_id_bkn":             nil,
			"nama_jabatan":               "Jabatan 6",
			"tmt_jabatan":                nil,
			"status_plt":                 nil,
			"periode_jabatan_start_date": nil,
			"periode_jabatan_end_date":   nil,
			"no_sk":                      nil,
			"tanggal_sk":                 nil,
			"bkn_id":                     nil,
			"tmt_pelantikan":             nil,
			"is_active":                  nil,
			"eselon_id":                  nil,
			"eselon":                     nil,
			"eselon1":                    nil,
			"eselon2":                    nil,
			"eselon3":                    nil,
			"eselon4":                    nil,
			"catatan":                    nil,
			"jenis_sk":                   nil,
			"status_satker":              nil,
			"status_biro":                nil,
			"tabel_mutasi_id":            nil,
			"file_base64":                nil,
			"keterangan_berkas":          nil,
			"pns_id":                     "id_1c",
			"pns_nip":                    "1c",
			"pns_nama":                   nil,
			"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                 nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, dbrepository.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
				"jenis_jabatan_id": 1,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-02",
				"status_plt": true,
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(1),
					"kelas_jabatan_id":           int32(1),
					"satuan_kerja_id":            "3",
					"unor_id":                    "5",
					"unor_id_bkn":                "5",
					"unor":                       "Unor 3",
					"jenis_jabatan_id":           int32(1),
					"jenis_jabatan":              "Jenis 1",
					"jabatan_id":                 "K1",
					"jabatan_id_bkn":             "BKN.1",
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_plt":                 true,
					"periode_jabatan_start_date": time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					"periode_jabatan_end_date":   time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
					"no_sk":                      "SK.01",
					"tanggal_sk":                 time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					"bkn_id":                     "bkn1",
					"tmt_pelantikan":             time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
					"is_active":                  int16(0),
					"eselon_id":                  "ide0",
					"eselon":                     "e0",
					"eselon1":                    "e1",
					"eselon2":                    "e2",
					"eselon3":                    "e3",
					"eselon4":                    "e4",
					"catatan":                    "ket",
					"jenis_sk":                   "1",
					"status_satker":              int32(0),
					"status_biro":                int32(0),
					"tabel_mutasi_id":            int64(1),
					"file_base64":                "data:abc",
					"keterangan_berkas":          "abc",
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": null,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "",
				"tanggal_sk": "2000-01-02",
				"status_plt": null,
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(2),
					"kelas_jabatan_id":           int32(1),
					"satuan_kerja_id":            "3",
					"unor_id":                    "5",
					"unor_id_bkn":                "5",
					"unor":                       "Unor 3",
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 "K1",
					"jabatan_id_bkn":             "BKN.1",
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_plt":                 nil,
					"periode_jabatan_start_date": time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					"periode_jabatan_end_date":   time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
					"no_sk":                      "",
					"tanggal_sk":                 time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					"bkn_id":                     "bkn1",
					"tmt_pelantikan":             time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
					"is_active":                  int16(0),
					"eselon_id":                  "ide0",
					"eselon":                     "e0",
					"eselon1":                    "e1",
					"eselon2":                    "e2",
					"eselon3":                    "e3",
					"eselon4":                    "e4",
					"catatan":                    "ket",
					"jenis_sk":                   "1",
					"status_satker":              int32(0),
					"status_biro":                int32(0),
					"tabel_mutasi_id":            int64(1),
					"file_base64":                "data:abc",
					"keterangan_berkas":          "abc",
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-02",
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(3),
					"kelas_jabatan_id":           int32(1),
					"satuan_kerja_id":            "3",
					"unor_id":                    "5",
					"unor_id_bkn":                "5",
					"unor":                       "Unor 3",
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 "K1",
					"jabatan_id_bkn":             "BKN.1",
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_plt":                 nil,
					"periode_jabatan_start_date": time.Date(2000, 1, 3, 0, 0, 0, 0, time.UTC),
					"periode_jabatan_end_date":   time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC),
					"no_sk":                      "SK.01",
					"tanggal_sk":                 time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC),
					"bkn_id":                     "bkn1",
					"tmt_pelantikan":             time.Date(1999, 1, 1, 0, 0, 0, 0, time.UTC),
					"is_active":                  int16(0),
					"eselon_id":                  "ide0",
					"eselon":                     "e0",
					"eselon1":                    "e1",
					"eselon2":                    "e2",
					"eselon3":                    "e3",
					"eselon4":                    "e4",
					"catatan":                    "ket",
					"jenis_sk":                   "1",
					"status_satker":              int32(0),
					"status_biro":                int32(0),
					"tabel_mutasi_id":            int64(1),
					"file_base64":                "data:abc",
					"keterangan_berkas":          "abc",
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "error: riwayat jabatan is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 1,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-02",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat jabatan is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-02",
				"periode_jabatan_start_date": "2000-01-03",
				"periode_jabatan_end_date": "2000-01-04"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(4),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 4",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1e",
					"pns_nip":                    "1e",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:          "error: riwayat jabatan is deleted",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 1,
				"id_jabatan": "K1",
				"satuan_kerja_id": "3",
				"unit_organisasi_id": "5",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(5),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 5",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: jenis jabatan or jabatan or satuan kerja or unit organisasi is not found",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 0,
				"id_jabatan": "K0",
				"satuan_kerja_id": "0",
				"unit_organisasi_id": "0",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis jabatan tidak ditemukan | data jabatan tidak ditemukan | data satuan kerja tidak ditemukan | data unit organisasi tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: jenis jabatan or jabatan or satuan kerja or unit organisasi is deleted",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": 2,
				"id_jabatan": "K2",
				"satuan_kerja_id": "4",
				"unit_organisasi_id": "4",
				"tmt_jabatan": "2000-01-01",
				"no_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"status_plt": false,
				"periode_jabatan_start_date": "2000-01-01",
				"periode_jabatan_end_date": "2000-01-01"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis jabatan tidak ditemukan | data jabatan tidak ditemukan | data satuan kerja tidak ditemukan | data unit organisasi tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": "0",
				"id_jabatan": 1,
				"satuan_kerja_id": 1,
				"unit_organisasi_id": 1,
				"tmt_jabatan": "",
				"no_sk": "` + strings.Repeat(".", 101) + `",
				"tanggal_sk": "",
				"status_plt": "false",
				"periode_jabatan_start_date": "",
				"periode_jabatan_end_date": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id_jabatan\" harus dalam tipe string` +
				` | parameter \"jenis_jabatan_id\" harus dalam tipe integer` +
				` | parameter \"no_sk\" harus 100 karakter atau kurang` +
				` | parameter \"periode_jabatan_end_date\" harus dalam format date` +
				` | parameter \"periode_jabatan_start_date\" harus dalam format date` +
				` | parameter \"satuan_kerja_id\" harus dalam tipe string` +
				` | parameter \"status_plt\" harus dalam tipe boolean` +
				` | parameter \"tanggal_sk\" harus dalam format date` +
				` | parameter \"tmt_jabatan\" harus dalam format date` +
				` | parameter \"unit_organisasi_id\" harus dalam tipe string"}`,
			wantDBRows: defaultRows,
		},
		{
			name:          "error: null params",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_jabatan_id": null,
				"id_jabatan": null,
				"satuan_kerja_id": null,
				"unit_organisasi_id": null,
				"tmt_jabatan": null,
				"no_sk": null,
				"tanggal_sk": null,
				"status_plt": null,
				"periode_jabatan_start_date": null,
				"periode_jabatan_end_date": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id_jabatan\" tidak boleh null` +
				` | parameter \"no_sk\" tidak boleh null` +
				` | parameter \"periode_jabatan_end_date\" tidak boleh null` +
				` | parameter \"periode_jabatan_start_date\" tidak boleh null` +
				` | parameter \"satuan_kerja_id\" tidak boleh null` +
				` | parameter \"tanggal_sk\" tidak boleh null` +
				` | parameter \"tmt_jabatan\" tidak boleh null` +
				` | parameter \"unit_organisasi_id\" tidak boleh null"}`,
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
				` | parameter \"id_jabatan\" harus diisi` +
				` | parameter \"satuan_kerja_id\" harus diisi` +
				` | parameter \"unit_organisasi_id\" harus diisi` +
				` | parameter \"periode_jabatan_start_date\" harus diisi` +
				` | parameter \"periode_jabatan_end_date\" harus diisi` +
				` | parameter \"no_sk\" harus diisi` +
				` | parameter \"tanggal_sk\" harus diisi` +
				` | parameter \"tmt_jabatan\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-jabatan/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_jabatan", "where id = $1", tt.paramID)
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
		insert into riwayat_jabatan
			(id,  nama_jabatan, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('1', 'Jabatan 1',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null),
			('2', 'Jabatan 2',  'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('3', 'Jabatan 3',  'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Jabatan 4',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                         int64(4),
			"kelas_jabatan_id":           nil,
			"satuan_kerja_id":            nil,
			"unor_id":                    nil,
			"unor_id_bkn":                nil,
			"unor":                       nil,
			"jenis_jabatan_id":           nil,
			"jenis_jabatan":              nil,
			"jabatan_id":                 nil,
			"jabatan_id_bkn":             nil,
			"nama_jabatan":               "Jabatan 4",
			"tmt_jabatan":                nil,
			"status_plt":                 nil,
			"periode_jabatan_start_date": nil,
			"periode_jabatan_end_date":   nil,
			"no_sk":                      nil,
			"tanggal_sk":                 nil,
			"bkn_id":                     nil,
			"tmt_pelantikan":             nil,
			"is_active":                  nil,
			"eselon_id":                  nil,
			"eselon":                     nil,
			"eselon1":                    nil,
			"eselon2":                    nil,
			"eselon3":                    nil,
			"eselon4":                    nil,
			"catatan":                    nil,
			"jenis_sk":                   nil,
			"status_satker":              nil,
			"status_biro":                nil,
			"tabel_mutasi_id":            nil,
			"file_base64":                nil,
			"keterangan_berkas":          nil,
			"pns_id":                     "id_1c",
			"pns_nip":                    "1c",
			"pns_nama":                   nil,
			"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                 nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, dbrepository.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
					"id":                         int64(1),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat jabatan is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(2),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 2",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1e",
					"pns_nip":                    "1e",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:             "error: riwayat jabatan is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat jabatan is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(3),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 3",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
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

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-jabatan/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_jabatan", "where id = $1", tt.paramID)
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
		insert into riwayat_jabatan
			(id,  nama_jabatan, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('1', 'Jabatan 1',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null),
			('2', 'Jabatan 2',  'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('3', 'Jabatan 3',  'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Jabatan 4',  'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                         int64(4),
			"kelas_jabatan_id":           nil,
			"satuan_kerja_id":            nil,
			"unor_id":                    nil,
			"unor_id_bkn":                nil,
			"unor":                       nil,
			"jenis_jabatan_id":           nil,
			"jenis_jabatan":              nil,
			"jabatan_id":                 nil,
			"jabatan_id_bkn":             nil,
			"nama_jabatan":               "Jabatan 4",
			"tmt_jabatan":                nil,
			"status_plt":                 nil,
			"periode_jabatan_start_date": nil,
			"periode_jabatan_end_date":   nil,
			"no_sk":                      nil,
			"tanggal_sk":                 nil,
			"bkn_id":                     nil,
			"tmt_pelantikan":             nil,
			"is_active":                  nil,
			"eselon_id":                  nil,
			"eselon":                     nil,
			"eselon1":                    nil,
			"eselon2":                    nil,
			"eselon3":                    nil,
			"eselon4":                    nil,
			"catatan":                    nil,
			"jenis_sk":                   nil,
			"status_satker":              nil,
			"status_biro":                nil,
			"tabel_mutasi_id":            nil,
			"file_base64":                nil,
			"keterangan_berkas":          nil,
			"pns_id":                     "id_1c",
			"pns_nip":                    "1c",
			"pns_nama":                   nil,
			"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                 nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, dbrepository.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
					"id":                         int64(1),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 1",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 "{updated_at}",
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:              "error: riwayat jabatan is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat jabatan is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(2),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 2",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1e",
					"pns_nip":                    "1e",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 nil,
				},
			},
		},
		{
			name:              "error: riwayat jabatan is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                         int64(3),
					"kelas_jabatan_id":           nil,
					"satuan_kerja_id":            nil,
					"unor_id":                    nil,
					"unor_id_bkn":                nil,
					"unor":                       nil,
					"jenis_jabatan_id":           nil,
					"jenis_jabatan":              nil,
					"jabatan_id":                 nil,
					"jabatan_id_bkn":             nil,
					"nama_jabatan":               "Jabatan 3",
					"tmt_jabatan":                nil,
					"status_plt":                 nil,
					"periode_jabatan_start_date": nil,
					"periode_jabatan_end_date":   nil,
					"no_sk":                      nil,
					"tanggal_sk":                 nil,
					"bkn_id":                     nil,
					"tmt_pelantikan":             nil,
					"is_active":                  nil,
					"eselon_id":                  nil,
					"eselon":                     nil,
					"eselon1":                    nil,
					"eselon2":                    nil,
					"eselon3":                    nil,
					"eselon4":                    nil,
					"catatan":                    nil,
					"jenis_sk":                   nil,
					"status_satker":              nil,
					"status_biro":                nil,
					"tabel_mutasi_id":            nil,
					"file_base64":                nil,
					"keterangan_berkas":          nil,
					"pns_id":                     "id_1c",
					"pns_nip":                    "1c",
					"pns_nama":                   nil,
					"created_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-jabatan/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_jabatan", "where id = $1", tt.paramID)
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
