package pemberitahuan

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/dbmigrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pemberitahuan
			(id, judul_berita, deskripsi_berita, status,                 updated_by, updated_at) values
			(11, '11a',        '11b',            'Aktif',                '11c',      '2000-01-02'),
			(12, '12a',        '12b',            'Menunggu Diberitakan', '12c',      '2000-01-02'),
			(13, '13a',        '13b',            'Sudah Tidak Aktif',    '13c',      '2000-01-02');
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                  11,
						"deskripsi_berita":    "11b",
						"diperbarui_oleh":     "11c",
						"judul_berita":        "11a",
						"status":              "Aktif",
						"terakhir_diperbarui": "2000-01-02"
					},
					{
						"id":                  12,
						"deskripsi_berita":    "12b",
						"diperbarui_oleh":     "12c",
						"judul_berita":        "12a",
						"status":              "Menunggu Diberitakan",
						"terakhir_diperbarui": "2000-01-02"
					},
					{
						"id":                  13,
						"deskripsi_berita":    "13b",
						"diperbarui_oleh":     "13c",
						"judul_berita":        "13a",
						"status":              "Sudah Tidak Aktif",
						"terakhir_diperbarui": "2000-01-02"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"deskripsi_berita":    "12b",
						"diperbarui_oleh":     "12c",
						"id":                  12,
						"judul_berita":        "12a",
						"status":              "Menunggu Diberitakan",
						"terakhir_diperbarui": "2000-01-02"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: cari judul",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"12a"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                  12,
						"deskripsi_berita":    "12b",
						"diperbarui_oleh":     "12c",
						"judul_berita":        "12a",
						"status":              "Menunggu Diberitakan",
						"terakhir_diperbarui": "2000-01-02"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: cari deskripsi",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"11b"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                  11,
						"deskripsi_berita":    "11b",
						"diperbarui_oleh":     "11c",
						"judul_berita":        "11a",
						"status":              "Aktif",
						"terakhir_diperbarui": "2000-01-02"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: tidak ada data ditemukan",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"22"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: user bukan admin",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "bukan_admin")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "Forbidden"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/pemberitahuan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, db, api.NewAuthMiddleware(apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
