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

		insert into unit_kerja(id, nama_unor, deleted_at) values
		(1, 'Unit 1', null),
		(2, 'Unit 2', null),
		(3, 'Unit 3', '2000-01-01');

		insert into ref_kelas_jabatan(id, kelas_jabatan, tunjangan_kinerja) values
		(1, 'Kelas 1', 2531250),
		(2, 'Kelas 2', 2708250);

		insert into riwayat_jabatan(id, pns_nip, jenis_jabatan_id, jabatan_id, tmt_jabatan, no_sk, tanggal_sk, satuan_kerja_id, unor_id, kelas_jabatan_id, periode_jabatan_start_date, periode_jabatan_end_date, deleted_at) values
		(1, '41', 1, '11', '2025-01-01', '1234567890', '2025-01-01', 1, 1, 1, '2024-01-01', '2024-12-31', null),
		(2, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(3, '42', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(4, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', '2000-01-01'),
		(5, '41', 3, '13', '2024-02-01', '2234567890', '2024-02-01', 3, 3, 2, '2025-01-01', '2025-12-31', null);
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
						"id":         2,
						"jenis_jabatan": "Jabatan Fungsional",
						"nama_jabatan": "12h",
						"id_jabatan": "12",
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
						"id_jabatan": "11",
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
						"id_jabatan":                 "13",
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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         1,
						"jenis_jabatan": "Jabatan Struktural",
						"nama_jabatan": "11h",
						"tmt_jabatan": "2025-01-01",
						"id_jabatan": "11",
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

		insert into unit_kerja(id, nama_unor, deleted_at) values
		(1, 'Unit 1', null),
		(2, 'Unit 2', null),
		(3, 'Unit 3', '2000-01-01');

		insert into ref_kelas_jabatan(id, kelas_jabatan, tunjangan_kinerja) values
		(1, 'Kelas 1', 2531250),
		(2, 'Kelas 2', 2708250);

		insert into pegawai (pns_id, nip_baru, nama, deleted_at)
		values ('pns-1', '41', 'Pegawai Test', null),
		('pns-2', '42', 'Pegawai Test 2', now());

		insert into riwayat_jabatan(id, pns_nip, jenis_jabatan_id, jabatan_id, tmt_jabatan, no_sk, tanggal_sk, satuan_kerja_id, unor_id, kelas_jabatan_id, periode_jabatan_start_date, periode_jabatan_end_date, deleted_at) values
		(1, '41', 1, '11', '2025-01-01', '1234567890', '2025-01-01', 1, 1, 1, '2024-01-01', '2024-12-31', null),
		(2, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(3, '42', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', null),
		(4, '41', 2, '12', '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31', '2000-01-01'),
		(5, '41', 3, '13', '2024-02-01', '2234567890', '2024-02-01', 3, 3, 2, '2025-01-01', '2025-12-31', null);
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
						"id":         2,
						"jenis_jabatan": "Jabatan Fungsional",
						"nama_jabatan": "12h",
						"id_jabatan": "12",
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
						"id_jabatan": "11",
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
						"id_jabatan":                 "13",
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
			nip:              "41",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         1,
						"jenis_jabatan": "Jabatan Struktural",
						"nama_jabatan": "11h",
						"tmt_jabatan": "2025-01-01",
						"id_jabatan": "11",
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
			wantResponseBody: `{"data": [{"id": 3, "jenis_jabatan": "Jabatan Fungsional", "nama_jabatan": "12h", "id_jabatan": "12", "tmt_jabatan": "2025-02-01", "no_sk": "2234567890", "tanggal_sk": "2025-02-01", "satuan_kerja": "Unit 2", "unit_organisasi": "Unit 2", "status_plt": false, "kelas_jabatan": "Kelas 2", "periode_jabatan_start_date": "2025-01-01", "periode_jabatan_end_date": "2025-12-31"}], "meta": {"limit": 10, "offset": 0, "total": 1}}`,
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
