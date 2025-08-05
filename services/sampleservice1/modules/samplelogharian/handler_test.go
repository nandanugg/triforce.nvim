package samplelogharian

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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/sampleservice1/dbmigrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/sampleservice1/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into sample_log_harian(tanggal, aktivitas, user_id) values
			('2000-01-02', 'Melakukan task A', 'user1'),
			('2000-01-03', 'Melakukan task B', 'user2'),
			('2000-01-04', 'Melakukan task C', 'user1'),
			('2000-01-05', 'Melakukan task D', 'user1');
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
			name:             "ok: no params",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("user1")}},
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": [
					{"tanggal":"2000-01-02", "aktivitas":"Melakukan task A"},
					{"tanggal":"2000-01-04", "aktivitas":"Melakukan task C"},
					{"tanggal":"2000-01-05", "aktivitas":"Melakukan task D"}
				],
				"meta":{"limit":10,"offset":0,"total":3}
			}`,
		},
		{
			name:             "ok: with pagination params",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("user1")}},
			wantResponseCode: 200,
			wantResponseBody: `{
				"data": [
					{"tanggal":"2000-01-04", "aktivitas":"Melakukan task C"}
				],
				"meta":{"limit":1,"offset":1,"total":3}
			}`,
		},
		{
			name:             "error: invalid parameter",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"satu"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("user1")}},
			wantResponseCode: 400,
			wantResponseBody: `{"message": "parameter \"limit\" harus dalam format yang sesuai"}`,
		},
		{
			name:             "error: invalid auth token",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: 401,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest("GET", "/sample-log-harian", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenapiBytes)
			require.NoError(t, err)
			RegisterRoutes(e, db, api.NewAuthMiddleware(apitest.JWTPublicKey))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
