package lokasi

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

func Test_handler_ListRefLokasi(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_lokasi ("id", "nama", "created_at", "updated_at", "deleted_at") values
		(1, 'Jakarta', '2024-01-01', '2024-01-01', null),
		(2, 'Bandung', '2024-01-01', '2024-01-01', null),
		(3, 'Semarang', '2024-01-01', '2024-01-01', null),
		(4, 'Yogyakarta', '2024-01-01', '2024-01-01', now())
		;
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgx)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: admin get all lokasi with default pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "1", "nama": "Jakarta"},
					{"id": "2", "nama": "Bandung"},
					{"id": "3", "nama": "Semarang"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name: "ok: query search",
			requestQuery: url.Values{
				"nama": []string{"Jak"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "1", "nama": "Jakarta"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name: "ok: pagination limit 3 offset 2",
			requestQuery: url.Values{
				"limit":  []string{"3"},
				"offset": []string{"2"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "3", "nama": "Semarang"}
				],
				"meta": {"limit": 3, "offset": 2, "total": 3}
			}`,
		},
		{
			name:             "error: invalid token",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/lokasi", nil)
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
