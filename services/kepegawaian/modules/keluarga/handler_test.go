package keluarga

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
	dbData := `
		insert into kepegawaian.users
			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
		insert into kepegawaian.orang_tua
			("ID", "HUBUNGAN", "ALAMAT", "NO_TLP", "NO_HP", "STATUS_PERKAWINAN", "AKTE_KELAHIRAN", "STATUS_HIDUP", "AKTE_MENINGGAL", "TGL_MENINGGAL", "NO_NPWP", "TANGGAL_NPWP", "NAMA", "GELAR_DEPAN", "GELAR_BELAKANG", "TEMPAT_LAHIR", "TANGGAL_LAHIR", "JENIS_KELAMIN", "AGAMA", "EMAIL", "JENIS_DOKUMEN_ID", "NO_DOKUMEN_ID", "FOTO", "KODE", "NIP", "PNS_ID") values
			(21,   1,          '21a',    '21b',    '21c',   '21d',               '21e',            1,              '21f',            '2000-01-01',    '21g',     '2000-01-02',   '21h',  '21i',         '21j',            '21k',          '2000-01-03',    '21l',           '21',    '21m',   '21n',              '21o',           '21p',  1,      '1c',  '21i'),
			(22,   2,          '22a',    '22b',    '22c',   '22d',               '22e',            2,              '22f',            '2001-01-01',    '22g',     '2001-01-02',   '22h',  '22i',         '22j',            '22k',          '2001-01-03',    '22l',           '22',    '22m',   '22n',              '22o',           '22p',  2,      '1c',  '22i'),
			(23,   3,          '23a',    '23b',    '23c',   '23d',               '23e',            3,              '23f',            '2002-01-01',    '23g',     '2002-01-02',   '23h',  '23i',         '23j',            '23k',          '2002-01-03',    '23l',           '23',    '23m',   '23n',              '23o',           '23p',  1,      '2c',  '23i');
		insert into kepegawaian.istri
			("ID", "PNS", "NAMA", "TANGGAL_MENIKAH", "AKTE_NIKAH", "TANGGAL_MENINGGAL", "AKTE_MENINGGAL", "TANGGAL_CERAI", "AKTE_CERAI", "KARSUS", "STATUS", "HUBUNGAN", "PNS_ID", "NIP") values
			(31,   1,     '31a',  '2000-01-01',      '31c',        '2000-01-01',        '31e',            '2000-01-02',    '31f',        '31g',    1,        1,          '31i',    '1c'),
			(32,   2,     '32a',  '2001-01-01',      '32c',        '2001-01-01',        '32e',            '2001-01-02',    '32f',        '32g',    2,        2,          '32i',    '1c'),
			(33,   3,     '33a',  '2002-01-01',      '33c',        '2002-01-01',        '33e',            '2002-01-02',    '33f',        '33g',    3,        1,          '33i',    '2c');
		insert into kepegawaian.anak
			("ID", "PASANGAN", "NAMA", "JENIS_KELAMIN", "TANGGAL_LAHIR", "TEMPAT_LAHIR", "STATUS_ANAK", "PNS_ID", "NIP") values
			(11,   1,          '11a',  '1',             '2000-01-01',    '11b',          '1',           '11c',    '1c'),
			(12,   2,          '12a',  '2',             '2001-01-01',    '12b',          '2',           '12c',    '1c'),
			(13,   3,          '13a',  '3',             '2002-01-01',    '13b',          '3',           '13c',    '2c');
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
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(41)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":            21,
						"nama":          "21h",
						"peran":         "AYAH",
						"tanggal_lahir": "2000-01-03"
					},
					{
						"id":            22,
						"nama":          "22h",
						"peran":         "IBU",
						"tanggal_lahir": "2001-01-03"
					},
					{
						"id":            31,
						"nama":          "31a",
						"peran":         "ISTRI"
					},
					{
						"id":            32,
						"nama":          "32a",
						"peran":         "SUAMI"
					},
					{
						"id":            11,
						"nama":          "11a",
						"peran":         "ANAK",
						"tanggal_lahir": "2000-01-01"
					},
					{
						"id":            12,
						"nama":          "12a",
						"peran":         "ANAK",
						"tanggal_lahir": "2001-01-01"
					}
				]
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(200)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": []}`,
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
			db := dbtest.New(t, "kepegawaian", dbmigrations.FS)
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/keluarga", nil)
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

func Test_handler_listOrangTua(t *testing.T) {
	dbData := `
		insert into kepegawaian.users
			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
		insert into kepegawaian.orang_tua
			("ID", "HUBUNGAN", "ALAMAT", "NO_TLP", "NO_HP", "STATUS_PERKAWINAN", "AKTE_KELAHIRAN", "STATUS_HIDUP", "AKTE_MENINGGAL", "TGL_MENINGGAL", "NO_NPWP", "TANGGAL_NPWP", "NAMA", "GELAR_DEPAN", "GELAR_BELAKANG", "TEMPAT_LAHIR", "TANGGAL_LAHIR", "JENIS_KELAMIN", "AGAMA", "EMAIL", "JENIS_DOKUMEN_ID", "NO_DOKUMEN_ID", "FOTO", "KODE", "NIP", "PNS_ID") values
			(21,   1,          '21a',    '21b',    '21c',   '21d',               '21e',            1,              '21f',            '2000-01-01',    '21g',     '2000-01-02',   '21h',  '21i',         '21j',            '21k',          '2000-01-03',    '21l',           '21',    '21m',   '21n',              '21o',           '21p',  1,      '1c',  '21i'),
			(22,   2,          '22a',    '22b',    '22c',   '22d',               '22e',            2,              '22f',            '2001-01-01',    '22g',     '2001-01-02',   '22h',  '22i',         '22j',            '22k',          '2001-01-03',    '22l',           '22',    '22m',   '22n',              '22o',           '22p',  2,      '1c',  '22i'),
			(23,   3,          '23a',    '23b',    '23c',   '23d',               '23e',            3,              '23f',            '2002-01-01',    '23g',     '2002-01-02',   '23h',  '23i',         '23j',            '23k',          '2002-01-03',    '23l',           '23',    '23m',   '23n',              '23o',           '23p',  1,      '2c',  '23i');
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
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(41)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":            21,
						"nama":          "21h",
						"peran":         "AYAH",
						"tanggal_lahir": "2000-01-03"
					},
					{
						"id":            22,
						"nama":          "22h",
						"peran":         "IBU",
						"tanggal_lahir": "2001-01-03"
					}
				]
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(200)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": []}`,
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
			db := dbtest.New(t, "kepegawaian", dbmigrations.FS)
			_, err := db.Exec(tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/keluarga/orang-tua", nil)
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
