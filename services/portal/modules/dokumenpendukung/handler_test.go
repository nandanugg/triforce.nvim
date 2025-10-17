package dokumenpendukung

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
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/migrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into dokumen_pendukung
			(id, nama_tombol, nama_halaman, updated_by, updated_at,   file) values
			(11, '11a',       '11b',        '11c',      '2000-01-02', null),
			(12, '12a',       '12b',        null,       null,         '12d');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("1")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                  11,
						"nama_halaman":        "11b",
						"nama_tombol":         "11a",
						"diperbarui_oleh":     "11c",
						"terakhir_diperbarui": "2000-01-02",
						"status":              "Belum Upload"
					},
					{
						"id":           12,
						"nama_halaman": "12b",
						"nama_tombol":  "12a",
						"status":       "Sudah Upload"
					}
				]
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/dokumen-pendukung", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Informasi_Read)
			RegisterRoutes(e, db, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
