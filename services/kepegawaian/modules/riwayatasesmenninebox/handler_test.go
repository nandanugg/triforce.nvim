package riwayatasesmenninebox

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
		insert into riwayat_asesmen_nine_box
			(id, pns_nip, kesimpulan, tahun, deleted_at) values
			(11, '1c',    '11c',      2021,  null),
			(12, '1c',    '12c',      2022,  null),
			(13, '1c',    '13c',      null,  null),
			(14, '2c',    '14c',      2023,  null),
			(15, '1c',    '15c',      2023,  '2020-01-01');
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
						"id":         12,
						"tahun":      2022,
						"kesimpulan": "12c"
					},
					{
						"id":         11,
						"tahun":      2021,
						"kesimpulan": "11c"
					},
					{
						"id":         13,
						"tahun":      null,
						"kesimpulan": "13c"
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
						"id":         11,
						"tahun":      2021,
						"kesimpulan": "11c"
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-asesmen-nine-box", nil)
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
		insert into riwayat_asesmen_nine_box
			(id, pns_nip, kesimpulan, tahun, deleted_at) values
			(11, '1c',    '11c',      2021,  null),
			(12, '1c',    '12c',      2022,  null),
			(13, '1c',    '13c',      null,  null),
			(14, '2c',    '14c',      2023,  null),
			(15, '1c',    '15c',      2023,  '2020-01-01'),
			(16, '1d',    '16d',      2021,  null);
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
			name:             "ok: admin dapat melihat riwayat asesmen nine box pegawai 1c",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         12,
						"tahun":      2022,
						"kesimpulan": "12c"
					},
					{
						"id":         11,
						"tahun":      2021,
						"kesimpulan": "11c"
					},
					{
						"id":         13,
						"tahun":      null,
						"kesimpulan": "13c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat asesmen nine box pegawai 1d",
			dbData:           dbData,
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         16,
						"tahun":      2021,
						"kesimpulan": "16d"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat asesmen nine box pegawai 2c",
			dbData:           dbData,
			nip:              "2c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         14,
						"tahun":      2023,
						"kesimpulan": "14c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat asesmen nine box pegawai dengan pagination",
			dbData:           dbData,
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         11,
						"tahun":      2021,
						"kesimpulan": "11c"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: admin melihat pegawai yang tidak memiliki riwayat asesmen nine box",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-asesmen-nine-box", nil)
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
