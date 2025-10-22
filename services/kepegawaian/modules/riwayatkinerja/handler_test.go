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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-kinerja", nil)
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
			name:             "ok: admin dapat melihat riwayat kinerja pegawai 1c",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
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
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
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
			nip:              "2c",
			requestHeader:    http.Header{"Authorization": authHeader},
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
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
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
			nip:              "999",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-kinerja", nil)
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
