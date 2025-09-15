package jabatan

import (
	"context"
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
	dbrepository "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan(id, no, nama_jabatan, kode_jabatan) values
		(11, 1, '11h', '11h'),
		(12, 2, '12h', '12h'),
		(13, 3, '13h', '13h');
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
						"id":         11,
						"kode_jabatan": "11h",
						"nama_jabatan": "11h"
					},
					{
						"id":         12,
						"kode_jabatan": "12h",
						"nama_jabatan": "12h"
					},
					{
						"id":         13,
						"kode_jabatan": "13h",
						"nama_jabatan": "13h"
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
						"id":         12,
						"kode_jabatan": "12h",
						"nama_jabatan": "12h"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data jabatan",
			dbData:           ``,
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
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/jabatan", nil)
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

func Test_handler_listRiwayatJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan(id, no, nama_jabatan, kode_jabatan) values
		(11, 1, '11h', '11h'),
		(12, 2, '12h', '12h');

		insert into ref_jenis_jabatan(id, nama) values
		(1, 'Jabatan Struktural'),
		(2, 'Jabatan Fungsional');

		insert into unit_kerja(id, nama_unor) values
		(1, 'Unit 1'),
		(2, 'Unit 2');

		insert into ref_kelas_jabatan(id, kelas_jabatan, tunjangan_kinerja) values
		(1, 'Kelas 1', 2531250),
		(2, 'Kelas 2', 2708250);

		insert into riwayat_jabatan(id, pns_nip, jenis_jabatan_id, jabatan_id, tmt_jabatan, no_sk, tanggal_sk, satuan_kerja_id, unor_id, kelas_jabatan_id, periode_jabatan_start_date, periode_jabatan_end_date) values
		(1, '41', 1, 11, '2025-01-01', '1234567890', '2025-01-01', 1, 1, 1, '2024-01-01', '2024-12-31'),
		(2, '41', 2, 12, '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31'),
		(3, '42', 2, 12, '2025-02-01', '2234567890', '2025-02-01', 2, 2, 2, '2025-01-01', '2025-12-31')
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
						"tmt_jabatan": "2025-01-01",
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
				"meta": {"limit": 10, "offset": 0, "total": 2}
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
				"meta": {"limit": 1, "offset": 1, "total": 2}
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
