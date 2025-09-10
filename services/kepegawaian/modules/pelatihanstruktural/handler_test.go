package pelatihanstruktural

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
// 	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
// )

// func Test_handler_list(t *testing.T) {
// 	t.Parallel()

// 	dbData := `
// 		insert into pegawai
// 			("ID", "PNS_ID", "NIP_LAMA", "NIP_BARU", "NAMA", "GELAR_DEPAN", "GELAR_BELAKANG", "TEMPAT_LAHIR_ID", "TGL_LAHIR", "JENIS_KELAMIN", "AGAMA_ID", "JENIS_KAWIN_ID", "NIK", "NOMOR_DARURAT", "NOMOR_HP", "EMAIL", "ALAMAT", "NPWP", "BPJS", "JENIS_PEGAWAI_ID", "KEDUDUKAN_HUKUM_ID", "STATUS_CPNS_PNS", "KARTU_PEGAWAI", "NOMOR_SK_CPNS", "TGL_SK_CPNS", "TMT_CPNS",  "TMT_PNS",   "GOL_AWAL_ID", "GOL_ID", "TMT_GOLONGAN", "MK_TAHUN", "MK_BULAN", "JENIS_JABATAN_IDx", "JABATAN_ID", "TMT_JABATAN", "PENDIDIKAN_ID", "TAHUN_LULUS", "KPKN_ID", "LOKASI_KERJA_ID", "UNOR_ID", "UNOR_INDUK_ID", "INSTANSI_INDUK_ID", "INSTANSI_KERJA_ID", "SATUAN_KERJA_INDUK_ID", "SATUAN_KERJA_KERJA_ID", "GOLONGAN_DARAH", "PHOTO", "TMT_PENSIUN", "LOKASI_KERJA", "JML_ISTRI", "JML_ANAK", "NO_SURAT_DOKTER", "TGL_SURAT_DOKTER", "NO_BEBAS_NARKOBA", "TGL_BEBAS_NARKOBA", "NO_CATATAN_POLISI", "TGL_CATATAN_POLISI", "AKTE_KELAHIRAN", "STATUS_HIDUP", "AKTE_MENINGGAL", "TGL_MENINGGAL", "NO_ASKES", "NO_TASPEN", "TGL_NPWP",  "TEMPAT_LAHIR", "PENDIDIKAN", "TK_PENDIDIKAN", "TEMPAT_LAHIR_NAMA", "JENIS_JABATAN_NAMA", "JABATAN_NAMA", "KPKN_NAMA", "INSTANSI_INDUK_NAMA", "INSTANSI_KERJA_NAMA", "SATUAN_KERJA_INDUK_NAMA", "SATUAN_KERJA_NAMA", "JABATAN_INSTANSI_ID", "BUP", "JABATAN_INSTANSI_NAMA", "JENIS_JABATAN_ID", terminated_date, status_pegawai, "JABATAN_PPNPN", "JABATAN_INSTANSI_REAL_ID", "CREATED_DATE", "CREATED_BY", "UPDATED_DATE", "UPDATED_BY", "EMAIL_DIKBUD_BAK", "EMAIL_DIKBUD", "KODECEPAT", "IS_DOSEN", "MK_TAHUN_SWASTA", "MK_BULAN_SWASTA", "KK", "NIDN", "KET", "NO_SK_PEMBERHENTIAN", status_pegawai_backup, "MASA_KERJA", "KARTU_ASN") values
// 			(11,   '1a',     '1b',       '1c',       '1d',   '1e',          '1f',             '1g',              '2000-01-02','1h',            21,         '5',              '1k',  '1l',            '1m',       '1n',    '1o',     '1p',   '1q',   '31',               '1s',                 '1t',              '1u',            '1v',            '2000-01-03',  '2000-01-04','2000-01-04','1z',          1,        '2000-01-05',   '1ac',      '1ad',      '1ae',               '1af',        '2000-01-06',  '1ah',           '1ai',         '1aj',     '1ak',             '1al',     '1am',           '1an',               '1ao',               '1ap',                   '1aq',                   '1ar',            '1as',   '2000-01-07',  '1au',          '1',         '1',        '1ax',             '2000-01-08',       '1az',              '2000-01-09',        '1bb',               '2000-01-10',         '1bd',            '1be',          '1bf',            '2000-01-11',    '1bh',      'bi',        '2000-01-12','1bk',          '1bl',        '1bm',           '1bn',               '1bo',                '1bp',          '1bq',       '1br',                 '1bs',                 '1bt',                     '1bu',               '1bv',                 1,     '1bx',                   1,                  '2000-01-13',    '1',            '1cb',           '1cc',                      '2000-01-14',   1,            '2000-01-15',   1,            '1ch',              '1ci',          '1cj',       1,          1,                 1,                 '1cn','1co',  '1cp', '1cq',                 1,                     '1cs',        '1ct');
// 		insert into users
// 			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
// 			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
// 		insert into rwt_diklat_struktural
// 			("ID", "PNS_ID", "PNS_NIP", "PNS_NAMA", "ID_DIKLAT", "NAMA_DIKLAT", "NOMOR", "TANGGAL",   "TAHUN", "STATUS_DATA", "FILE_BASE64", "KETERANGAN_BERKAS", "LAMA", "CREATED_DATE", "SIASN_ID") values
// 			('21', '21a',    '1c',      '21b',      '21c',       '21d',         '21e',   '2000-01-01',1000,    '21f',         '21g',         '21h',               1,      '2000-01-02',   '21i'),
// 			('22', '22a',    '1c',      '22b',      '22c',       '22d',         null,    '2001-01-01',1000,    '22f',         '22g',         '22h',               1,      '2001-01-02',   '22i'),
// 			('23', '23a',    '1c',      '23b',      '23c',       '23d',         '23e',   '2002-01-01',1000,    '23f',         '23g',         '23h',               1,      '2002-01-02',   '23i'),
// 			('24', '24a',    '2c',      '24b',      '24c',       '24d',         '24e',   '2003-01-01',1000,    '24f',         '24g',         '24h',               1,      '2003-01-02',   '24i');
// 	`

// 	tests := []struct {
// 		name             string
// 		dbData           string
// 		requestQuery     url.Values
// 		requestHeader    http.Header
// 		wantResponseCode int
// 		wantResponseBody string
// 	}{
// 		{
// 			name:             "ok: tanpa parameter apapun",
// 			dbData:           dbData,
// 			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
// 			wantResponseCode: http.StatusOK,
// 			wantResponseBody: `{
// 				"data": [
// 					{
// 						"id":          "21",
// 						"nama_diklat": "21d",
// 						"nomor":       "21e",
// 						"tahun":       1000,
// 						"tanggal":     "2000-01-01"
// 					},
// 					{
// 						"id":          "22",
// 						"nama_diklat": "22d",
// 						"tahun":       1000,
// 						"tanggal":     "2001-01-01"
// 					},
// 					{
// 						"id":          "23",
// 						"nama_diklat": "23d",
// 						"nomor":       "23e",
// 						"tahun":       1000,
// 						"tanggal":     "2002-01-01"
// 					}
// 				], "meta": {"limit": 10, "offset": 0, "total": 3}
// 			}`,
// 		},
// 		{
// 			name:             "ok: dengan parameter pagination",
// 			dbData:           dbData,
// 			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
// 			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
// 			wantResponseCode: http.StatusOK,
// 			wantResponseBody: `{
// 				"data": [
// 					{
// 						"id":          "22",
// 						"nama_diklat": "22d",
// 						"tahun":       1000,
// 						"tanggal":     "2001-01-01"
// 					}
// 				],
// 				"meta": {"limit": 1, "offset": 1, "total": 3}
// 			}`,
// 		},
// 		{
// 			name:             "ok: tidak ada data milik user",
// 			dbData:           dbData,
// 			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
// 			wantResponseCode: http.StatusOK,
// 			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
// 		},
// 		{
// 			name:             "error: auth header tidak valid",
// 			dbData:           dbData,
// 			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
// 			wantResponseCode: http.StatusUnauthorized,
// 			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			db := dbtest.New(t, dbmigrations.FS)
// 			_, err := db.Exec(tt.dbData)
// 			require.NoError(t, err)

// 			req := httptest.NewRequest(http.MethodGet, "/v1/pelatihan-struktural", nil)
// 			req.URL.RawQuery = tt.requestQuery.Encode()
// 			req.Header = tt.requestHeader
// 			rec := httptest.NewRecorder()

// 			e, err := api.NewEchoServer(docs.OpenAPIBytes)
// 			require.NoError(t, err)
// 			RegisterRoutes(e, db, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
// 			e.ServeHTTP(rec, req)

// 			assert.Equal(t, tt.wantResponseCode, rec.Code)
// 			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
// 			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
// 		})
// 	}
// }
