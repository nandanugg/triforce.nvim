package sertifikasi

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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/dbmigrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into users
			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
		insert into rwt_sertifikasi
			(id, nip,  tahun, nama_sertifikasi, base64, createddate,  deskripsi) values
			(11, '1c', 1,     '11a',            '11b',  '2000-01-01', '11c'),
			(12, '1c', 3,     '12a',            '12b',  '2001-01-01', '12c'),
			(13, '1c', 2,     '13a',            '13b',  '2002-01-01', '13c'),
			(14, '2c', 4,     '14a',            '14b',  '2003-01-01', '14c');
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(41)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":               11,
						"nama_sertifikasi": "11a",
						"tahun":            1,
						"tanggal_upload":   "2000-01-01"
					},
					{
						"id":               13,
						"nama_sertifikasi": "13a",
						"tahun":            2,
						"tanggal_upload":   "2002-01-01"
					},
					{
						"id":               12,
						"nama_sertifikasi": "12a",
						"tahun":            3,
						"tanggal_upload":   "2001-01-01"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(41)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":               13,
						"nama_sertifikasi": "13a",
						"tahun":            2,
						"tanggal_upload":   "2002-01-01"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(200)}},
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
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/sertifikasi", nil)
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
