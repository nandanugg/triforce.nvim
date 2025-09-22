package riwayatjabatan

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	dbrepository "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan(id, no, nama_jabatan, kode_jabatan, deleted_at) values
		(11, 1, '11h', '11h', null),
		(12, 2, '12h', '12h', null),
		(13, 3, '13h', '13h', '2000-01-01');

		insert into ref_jenis_jabatan(id, nama, deleted_at) values
		(1, 'Jabatan Struktural', null),
		(2, 'Jabatan Fungsional', null),
		(3, 'Jabatan Deleted', '2000-01-01');

		insert into unit_kerja(id, nama_unor, deleted_at) values
		(1, 'Unit 1', null),
		(2, 'Unit 2', null),
		(3, 'Unit 3', '2000-01-01');

		insert into ref_kelas_jabatan(id, kelas_jabatan, tunjangan_kinerja) values
		(1, 'Kelas 1', 2531250),
		(2, 'Kelas 2', 2708250);

		insert into riwayat_jabatan(id, pns_nip, jenis_jabatan_id, jabatan_id, tmt_jabatan, no_sk, tanggal_sk, satuan_kerja_id, unor_id, kelas_jabatan_id, periode_jabatan_start_date, periode_jabatan_end_date, deleted_at) values
		(1, '41', 1, 11, '2025-01-01', '1234567890', '2025-01-01', 1, 1, 1, '2024-01-01', '2024-12-31', null),
		(2, '41', 2, 12, '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(3, '42', 2, 12, '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(4, '41', 2, 12, '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', '2000-01-01'),
		(5, '41', 3, 13, '2024-02-01', '2234567890', '2024-02-01', 3, 3, 2, '2025-01-01', '2025-12-31', null);
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
						"id":         2,
						"jenis_jabatan": "Jabatan Fungsional",
						"nama_jabatan": "12h",
						"kode_jabatan": "12h",
						"tmt_jabatan": "2025-02-01",
						"no_sk": "2234567890",
						"tanggal_sk": "2025-02-01",
						"satuan_kerja": "Unit 2",
						"status_plt": false,
						"kelas_jabatan": "Kelas 2",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date": "2025-12-31",
						"unit_organisasi": "Unit 2"
					},
					{
						"id":         1,
						"jenis_jabatan": "Jabatan Struktural",
						"nama_jabatan": "11h",
						"kode_jabatan": "11h",
						"tmt_jabatan": "2025-01-01",
						"no_sk": "1234567890",
						"tanggal_sk": "2025-01-01",
						"satuan_kerja": "Unit 1",
						"status_plt": false,
						"kelas_jabatan": "Kelas 1",
						"periode_jabatan_start_date": "2024-01-01",
						"periode_jabatan_end_date": "2024-12-31",
						"unit_organisasi": "Unit 1"
					},
					{
						"id":                         5,
						"jenis_jabatan":              "",
						"nama_jabatan":               "",
						"kode_jabatan": 			  "",
						"tmt_jabatan":                "2024-02-01",
						"no_sk":                      "2234567890",
						"tanggal_sk":                 "2024-02-01",
						"satuan_kerja":               "",
						"status_plt":                 false,
						"kelas_jabatan":              "Kelas 2",
						"periode_jabatan_start_date": "2025-01-01",
						"periode_jabatan_end_date":   "2025-12-31",
						"unit_organisasi":            ""
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
						"id":         1,
						"jenis_jabatan": "Jabatan Struktural",
						"nama_jabatan": "11h",
						"tmt_jabatan": "2025-01-01",
						"kode_jabatan": "11h",
						"no_sk": "1234567890",
						"tanggal_sk": "2025-01-01",
						"satuan_kerja": "Unit 1",
						"status_plt": false,
						"kelas_jabatan": "Kelas 1",
						"periode_jabatan_start_date": "2024-01-01",
						"periode_jabatan_end_date": "2024-12-31",
						"unit_organisasi": "Unit 1"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data riwayat jabatan",
			dbData:           ``,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-jabatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := dbrepository.New(db)
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
			name:              "ok: valid pdf with data: prefix",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			dbData:            dbData,
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 riwayat jabatan tidak valid",
			dbData:            dbData,
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat jabatan sudah dihapus",
			dbData:            dbData,
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat jabatan berisi null value",
			dbData:            dbData,
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat jabatan berupa string kosong",
			dbData:            dbData,
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat jabatan bukan milik user login",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat jabatan tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat jabatan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			paramID:           "1",
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-jabatan/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepository.New(pgxconn)
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
