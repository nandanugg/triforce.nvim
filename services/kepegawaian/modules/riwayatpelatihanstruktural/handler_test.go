package riwayatpelatihanstruktural

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
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
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
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
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
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
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
			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-struktural", nil)
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
			paramID:           "uuid-pdf",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			paramID:           "uuid-png",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 tidak valid",
			dbData:            dbData,
			paramID:           "uuid-inv",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			dbData:            dbData,
			paramID:           "uuid-x",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			dbData:            dbData,
			paramID:           "uuid-null",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			dbData:            dbData,
			paramID:           "uuid-empty",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: berkas tidak ditemukan",
			dbData:            dbData,
			paramID:           "uuid-2",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: ambil data dari user lain",
			dbData:            dbData,
			paramID:           "uuid-3",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-struktural/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
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

	tests := []struct {
		name             string
		dbData           string
		nip              string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: admin dapat melihat riwayat pelatihan struktural pegawai 1c",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
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
			dbData:           dbData,
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
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
			dbData:           dbData,
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
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
			dbData:           dbData,
			nip:              "999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pelatihan-struktural", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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

	tests := []struct {
		name              string
		dbData            string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
		nip               string
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			dbData:            dbData,
			paramID:           "uuid-pdf",
			nip:               "1a",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			paramID:           "uuid-png",
			nip:               "1c",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 tidak valid",
			dbData:            dbData,
			paramID:           "uuid-inv",
			nip:               "1c",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			dbData:            dbData,
			paramID:           "uuid-x",
			nip:               "1d",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			dbData:            dbData,
			paramID:           "uuid-null",
			nip:               "1x",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			dbData:            dbData,
			paramID:           "uuid-empty",
			nip:               "1s",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: berkas tidak ditemukan",
			dbData:            dbData,
			paramID:           "uuid-2",
			nip:               "1t",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan struktural tidak ditemukan"}`),
		},
		{
			name:              "error: user bukan admin",
			dbData:            dbData,
			paramID:           "uuid-2",
			nip:               "1c",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode:  http.StatusForbidden,
			wantResponseBytes: []byte(`{"message": "akses ditolak"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			nip:               "1c",
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-pelatihan-struktural/%s/berkas", tt.nip, tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
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
