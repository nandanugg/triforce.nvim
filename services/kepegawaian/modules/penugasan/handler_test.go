package penugasan

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
		insert into kepegawaian.users
			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
		insert into kepegawaian.rwt_penugasan
			(id, tipe_jabatan, deskripsi_jabatan, tanggal_mulai, tanggal_selesai, createddate, exist) values
			(1,  'Struktural', 'Kepala Bagian',   '2023-01-01',   '2023-12-31',    now(),       true),
			(2,  'Fungsional', 'Analis Kepegawaian', '2023-06-01', '2024-05-31',   now(),       true),
			(3,  'Struktural', 'Kepala Sub Bagian', '2024-01-01',  '2024-12-31',   now(),       true);
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
						"id": 3,
						"tipe_jabatan": "Struktural",
						"deskripsi_jabatan": "Kepala Sub Bagian",
						"tanggal_mulai": "2024-01-01",
						"tanggal_selesai": "2024-12-31"
					},
					{
						"id": 2,
						"tipe_jabatan": "Fungsional",
						"deskripsi_jabatan": "Analis Kepegawaian",
						"tanggal_mulai": "2023-06-01",
						"tanggal_selesai": "2024-05-31"
					},
					{
						"id": 1,
						"tipe_jabatan": "Struktural",
						"deskripsi_jabatan": "Kepala Bagian",
						"tanggal_mulai": "2023-01-01",
						"tanggal_selesai": "2023-12-31"
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
						"id": 2,
						"tipe_jabatan": "Fungsional",
						"deskripsi_jabatan": "Analis Kepegawaian",
						"tanggal_mulai": "2023-06-01",
						"tanggal_selesai": "2024-05-31"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: tidak ada data",
			dbData:           "",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(41)}},
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

			db := dbtest.New(t, "kepegawaian", dbmigrations.FS)
			if tt.dbData != "" {
				_, err := db.Exec(tt.dbData)
				require.NoError(t, err)
			}

			req := httptest.NewRequest(http.MethodGet, "/penugasan", nil)
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
