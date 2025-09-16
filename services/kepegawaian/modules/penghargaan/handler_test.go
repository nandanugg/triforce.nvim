package penghargaan

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
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_penghargaan
			(id, nama, deleted_at)
			values
			(11, 'Jenis Penghargaan 1', NULL),
			(12, 'Jenis Penghargaan 2', NULL),
			(14, 'Jenis Penghargaan 3', now());

		insert into riwayat_penghargaan_umum
			(id, jenis_penghargaan_id, nama_penghargaan, deskripsi_penghargaan, tanggal_penghargaan, nip, deleted_at)
			values
			(11, 11, 'Penghargaan 1', 'Deskripsi Penghargaan 1', '2000-01-01', '41', NULL),
			(12, 12, 'Penghargaan 2', 'Deskripsi Penghargaan 2', '2001-01-01', '41', NULL),
			(13, 14, 'Penghargaan 3', 'Deskripsi Penghargaan 3', '2002-01-01', '41', NULL),
			(14, 11, 'Penghargaan 4', 'Deskripsi Penghargaan 4', '2003-01-01', '41', now()),
			(15, 11, 'Penghargaan 5', 'Deskripsi Penghargaan 5', '2004-01-01', '42', NULL);
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
						"id": 13,
						"deskripsi": "Deskripsi Penghargaan 3",
						"jenis_penghargaan": "",
						"nama_penghargaan": "Penghargaan 3",
						"tanggal": "2002-01-01"
					},
					{
						"id":                12,
						"deskripsi":         "Deskripsi Penghargaan 2",
						"jenis_penghargaan": "Jenis Penghargaan 2",
						"nama_penghargaan":  "Penghargaan 2",
						"tanggal":           "2001-01-01"
					},
					{
						"id":                11,
						"deskripsi":         "Deskripsi Penghargaan 1",
						"jenis_penghargaan": "Jenis Penghargaan 1",
						"nama_penghargaan":  "Penghargaan 1",
						"tanggal":           "2000-01-01"
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
						"id":                12,
						"deskripsi":         "Deskripsi Penghargaan 2",
						"jenis_penghargaan": "Jenis Penghargaan 2",
						"nama_penghargaan":  "Penghargaan 2",
						"tanggal":           "2001-01-01"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
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
			repo := repo.New(db)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-penghargaan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
