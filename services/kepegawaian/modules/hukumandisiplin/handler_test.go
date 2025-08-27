package hukumandisiplin

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
		insert into kepegawaian.rwt_hukdis
			("ID", "PNS_ID", "PNS_NIP", "NAMA", "ID_GOLONGAN", "NAMA_GOLONGAN", "ID_JENIS_HUKUMAN", "NAMA_JENIS_HUKUMAN", "SK_NOMOR", "SK_TANGGAL", "TANGGAL_MULAI_HUKUMAN", "MASA_TAHUN", "MASA_BULAN", "TANGGAL_AKHIR_HUKUMAN", "NO_PP", "NO_SK_PEMBATALAN", "TANGGAL_SK_PEMBATALAN", "ID_BKN", "FILE_BASE64", "KETERANGAN_BERKAS") values
			(11,   '11a',    '1c',      '11b',  '11',          '11c',           '11',               '11d',                '11e',      '2000-01-01', '2000-01-02',            1,            1,            '2000-01-03',            '11f',   '11g',              '2000-01-04',            '11h',    '11i',         '11j'),
			(12,   '12a',    '1c',      '12b',  '12',          '12c',           '12',               '12d',                '12e',      '2001-01-01', '2001-01-02',            2,            2,            '2001-01-03',            '12f',   '12g',              '2001-01-04',            '12h',    '12i',         '12j'),
			(13,   '13a',    '1c',      '13b',  '13',          '13c',           '13',               '13d',                '13e',      '2002-01-01', '2002-01-02',            3,            3,            '2002-01-03',            '13f',   '13g',              '2002-01-04',            '13h',    '13i',         '13j'),
			(14,   '14a',    '2c',      '14b',  '14',          '14c',           '14',               '14d',                '14e',      '2003-01-01', '2003-01-02',            4,            4,            '2003-01-03',            '14f',   '14g',              '2003-01-04',            '14h',    '14i',         '14j');
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
						"id":            11,
						"jenis_hukuman": "11d",
						"masa_bulan":    1,
						"masa_tahun":    1,
						"nomor_sk":      "11e",
						"tanggal_mulai": "2000-01-02",
						"tanggal_sk":    "2000-01-01"
					},
					{
						"id":            12,
						"jenis_hukuman": "12d",
						"masa_bulan":    2,
						"masa_tahun":    2,
						"nomor_sk":      "12e",
						"tanggal_mulai": "2001-01-02",
						"tanggal_sk":    "2001-01-01"
					},
					{
						"id":            13,
						"jenis_hukuman": "13d",
						"masa_bulan":    3,
						"masa_tahun":    3,
						"nomor_sk":      "13e",
						"tanggal_mulai": "2002-01-02",
						"tanggal_sk":    "2002-01-01"
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
						"id":            12,
						"jenis_hukuman": "12d",
						"masa_bulan":    2,
						"masa_tahun":    2,
						"nomor_sk":      "12e",
						"tanggal_mulai": "2001-01-02",
						"tanggal_sk":    "2001-01-01"
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

			db := dbtest.New(t, "kepegawaian", dbmigrations.FS)
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/hukuman-disiplin", nil)
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
