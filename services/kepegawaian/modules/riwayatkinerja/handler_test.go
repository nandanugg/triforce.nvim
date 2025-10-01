package riwayatkinerja

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
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_kinerja
			(id, nip,  tahun, rating_hasil_kerja, rating_perilaku_kerja, predikat_kinerja, deleted_at) values
			(1,  '1c', 2020,  'Sangat Baik',      'Sangat Baik',         'Sangat Baik',    null),
			(2,  '1c', 2023,  'Baik',             'Baik',                'Baik',           null),
			(3,  '1c', null,  'Cukup',            'Cukup',               'Cukup',          null),
			(4,  '2c', 2020,  'Sangat Baik',      'Sangat Baik',         'Sangat Baik',    null),
			(5,  '1c', 2020,  'Sangat Baik',      'Sangat Baik',         'Sangat Baik',    '2020-01-01');
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
						"id":              2,
						"hasil_kinerja":   "Baik",
						"kuadran_kinerja": "Baik",
						"perilaku_kerja":  "Baik",
						"tahun":           2023
					},
					{
						"id":              1,
						"hasil_kinerja":   "Sangat Baik",
						"kuadran_kinerja": "Sangat Baik",
						"perilaku_kerja":  "Sangat Baik",
						"tahun":           2020
					},
					{
						"id":              3,
						"hasil_kinerja":   "Cukup",
						"kuadran_kinerja": "Cukup",
						"perilaku_kerja":  "Cukup",
						"tahun":           null
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
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
						"id":              1,
						"hasil_kinerja":   "Sangat Baik",
						"kuadran_kinerja": "Sangat Baik",
						"perilaku_kerja":  "Sangat Baik",
						"tahun":           2020
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
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
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-kinerja", nil)
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

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_kinerja
			(id, nip,  tahun, rating_hasil_kerja, rating_perilaku_kerja, predikat_kinerja, deleted_at) values
			(1,  '1c', 2020,  'Sangat Baik',      'Sangat Baik',         'Sangat Baik',    null),
			(2,  '1c', 2023,  'Baik',             'Baik',                'Baik',           null),
			(3,  '1c', null,  'Cukup',            'Cukup',               'Cukup',          null),
			(4,  '2c', 2020,  'Sangat Baik',      'Sangat Baik',         'Sangat Baik',    null),
			(5,  '1c', 2020,  'Sangat Baik',      'Sangat Baik',         'Sangat Baik',    '2020-01-01'),
			(6,  '1d', 2021,  'Baik',             'Baik',                'Baik',           null);
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
			name:             "ok: admin dapat melihat riwayat kinerja pegawai 1c",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":              2,
						"hasil_kinerja":   "Baik",
						"kuadran_kinerja": "Baik",
						"perilaku_kerja":  "Baik",
						"tahun":           2023
					},
					{
						"id":              1,
						"hasil_kinerja":   "Sangat Baik",
						"kuadran_kinerja": "Sangat Baik",
						"perilaku_kerja":  "Sangat Baik",
						"tahun":           2020
					},
					{
						"id":              3,
						"hasil_kinerja":   "Cukup",
						"kuadran_kinerja": "Cukup",
						"perilaku_kerja":  "Cukup",
						"tahun":           null
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat kinerja pegawai 1d",
			dbData:           dbData,
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":              6,
						"hasil_kinerja":   "Baik",
						"kuadran_kinerja": "Baik",
						"perilaku_kerja":  "Baik",
						"tahun":           2021
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat kinerja pegawai 2c",
			dbData:           dbData,
			nip:              "2c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":              4,
						"hasil_kinerja":   "Sangat Baik",
						"kuadran_kinerja": "Sangat Baik",
						"perilaku_kerja":  "Sangat Baik",
						"tahun":           2020
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat kinerja pegawai dengan pagination",
			dbData:           dbData,
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":              1,
						"hasil_kinerja":   "Sangat Baik",
						"kuadran_kinerja": "Sangat Baik",
						"perilaku_kerja":  "Sangat Baik",
						"tahun":           2020
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: admin melihat pegawai yang tidak memiliki riwayat kinerja",
			dbData:           dbData,
			nip:              "999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: user bukan admin",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-kinerja", nil)
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
