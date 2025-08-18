package jabatan

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
		insert into kepegawaian.rwt_jabatan
			("ID_BKN", "PNS_ID", "PNS_NIP", "PNS_NAMA", "ID_UNOR", "UNOR", "ID_JENIS_JABATAN", "JENIS_JABATAN", "ID_JABATAN", "NAMA_JABATAN", "ID_ESELON", "ESELON", "TMT_JABATAN", "NOMOR_SK", "TANGGAL_SK", "ID_SATUAN_KERJA", "TMT_PELANTIKAN", "IS_ACTIVE", "ESELON1", "ESELON2", "ESELON3", "ESELON4", "ID", "CATATAN", "JENIS_SK", "LAST_UPDATED", "STATUS_SATKER", "STATUS_BIRO", "ID_JABATAN_BKN", "ID_UNOR_BKN", "JABATAN_TERAKHIR", "FILE_BASE64", "KETERANGAN_BERKAS", "ID_TABEL_MUTASI", "TERMINATED_DATE") values
			('11',     '11a',    '1c',      '11b',      '11c',     '11d',  '11e',              '11f',           '11g',        '11h',          '11i',       '11j',    '2000-01-01',  '11l',      '2000-01-02', '11n',             '2000-01-03',     '1',         '11q',     '11r',     '11s',     '11t',     11,   '11u',     '11v',      '2000-01-04',   1,               1,             '11z',            '11aa',        1,                  '11ac',        '11ad',              1,                 '2000-01-05'),
			('12',     '12a',    '1c',      '12b',      '12c',     '12d',  '12e',              '12f',           '12g',        '12h',          '12i',       '12j',    '2001-01-01',  '12l',      '2001-01-02', '12n',             '2001-01-03',     '2',         '12q',     '12r',     '12s',     '12t',     12,   '12u',     '12v',      '2001-01-04',   2,               2,             '12z',            '12aa',        2,                  '12ac',        '12ad',              2,                 '2001-01-05'),
			('13',     '13a',    '1c',      '13b',      '13c',     '13d',  '13e',              '13f',           '13g',        '13h',          '13i',       '13j',    '2002-01-01',  '13l',      '2002-01-02', '13n',             '2002-01-03',     '3',         '13q',     '13r',     '13s',     '13t',     13,   '13u',     '13v',      '2002-01-04',   3,               3,             '13z',            '13aa',        3,                  '13ac',        '13ad',              3,                 '2002-01-05'),
			('14',     '14a',    '2c',      '14b',      '14c',     '14d',  '14e',              '14f',           '14g',        '14h',          '14i',       '14j',    '2003-01-01',  '14l',      '2003-01-02', '14n',             '2003-01-03',     '3',         '14q',     '14r',     '14s',     '14t',     14,   '14u',     '14v',      '2003-01-04',   4,               4,             '14z',            '14aa',        4,                  '14ac',        '14ad',              4,                 '2003-01-05');
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
						"id":         "11",
						"jabatan":    "11h",
						"tmt":        "2000-01-01",
						"unit_kerja": "TODO: Unit Kerja"
					},
					{
						"id":         "12",
						"jabatan":    "12h",
						"tmt":        "2001-01-01",
						"unit_kerja": "TODO: Unit Kerja"
					},
					{
						"id":         "13",
						"jabatan":    "13h",
						"tmt":        "2002-01-01",
						"unit_kerja": "TODO: Unit Kerja"
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
						"id":         "12",
						"jabatan":    "12h",
						"tmt":        "2001-01-01",
						"unit_kerja": "TODO: Unit Kerja"
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

			req := httptest.NewRequest(http.MethodGet, "/jabatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, db, api.NewAuthMiddleware(apitest.JWTPublicKey))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
