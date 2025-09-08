package kenaikangajiberkala

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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into users
			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
		insert into rwt_kgb
			(pegawai_id, tmt_sk,       alasan, mv_kgb_id, no_sk, pejabat, id, ref,   tgl_sk,       pegawai_nama, pegawai_nip, birth_place, birth_date,   o_gol_ruang, o_gol_tmt, o_masakerja_thn, o_masakerja_bln, o_gapok, o_jabatan_text, o_tmt_jabatan, n_gol_ruang, n_gol_tmt, n_masakerja_thn, n_masakerja_bln, n_gapok, n_jabatan_text, n_tmt_jabatan, n_golongan_id, unit_kerja_text, unit_kerja_induk_text, unit_kerja_induk_id, kantor_pembayaran, last_education, last_education_date, nama_pejabat, "FILE_BASE64", "KETERANGAN_BERKAS") values
			(1,          '2000-01-01', '11a',  '1',       '11b', '11c',   11, '11d', '2000-01-02', '11e',        '1c',        '11g',       '2000-01-03', '11h',       '11i',     1,               1,               '11j',   '11k',          '2000-01-04',  '11l',       '11m',     1,               1,               '11n',   '11o',          '2000-01-05',  1,             '11q',           '11r',                 '11s',               '11t',             '11u',          '2000-01-06',        '11v',        '11w',         '11x'),
			(3,          '2001-01-01', '12a',  '2',       '12b', '12c',   12, '12d', '2001-01-02', '12e',        '1c',        '12g',       '2001-01-03', '12h',       '12i',     2,               2,               '12j',   '12k',          '2001-01-04',  '12l',       '12m',     2,               2,               '12n',   '12o',          '2001-01-05',  2,             '12q',           '12r',                 '12s',               '12t',             '12u',          '2001-01-06',        '12v',        '12w',         '12x'),
			(3,          '2002-01-01', '13a',  '3',       '13b', '13c',   13, '13d', '2002-01-02', '13e',        '1c',        '13g',       '2002-01-03', '13h',       '13i',     3,               3,               '13j',   '13k',          '2002-01-04',  '13l',       '13m',     3,               3,               '13n',   '13o',          '2002-01-05',  3,             '13q',           '13r',                 '13s',               '13t',             '13u',          '2002-01-06',        '13v',        '13w',         '13x'),
			(4,          '2003-01-01', '14a',  '4',       '14b', '14c',   14, '14d', '2003-01-02', '14e',        '2c',        '14g',       '2003-01-03', '14h',       '14i',     4,               4,               '14j',   '14k',          '2003-01-04',  '14l',       '14m',     4,               4,               '14n',   '14o',          '2003-01-05',  4,             '14q',           '14r',                 '14s',               '14t',             '14u',          '2003-01-06',        '14v',        '14w',         '14x');
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
						"id":         11,
						"gaji_pokok": "11n",
						"gol_ruang":  "11l",
						"no_sk":      "11b",
						"tgl_sk":     "2000-01-02",
						"tmt_kgb":    "2000-01-01"
					},
					{
						"id":         12,
						"gaji_pokok": "12n",
						"gol_ruang":  "12l",
						"no_sk":      "12b",
						"tgl_sk":     "2001-01-02",
						"tmt_kgb":    "2001-01-01"
					},
					{
						"id":         13,
						"gaji_pokok": "13n",
						"gol_ruang":  "13l",
						"no_sk":      "13b",
						"tgl_sk":     "2002-01-02",
						"tmt_kgb":    "2002-01-01"
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
						"id":         12,
						"gaji_pokok": "12n",
						"gol_ruang":  "12l",
						"no_sk":      "12b",
						"tgl_sk":     "2001-01-02",
						"tmt_kgb":    "2001-01-01"
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
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/kenaikan-gaji-berkala", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, db, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
