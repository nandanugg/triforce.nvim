package datapribadi

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

func Test_handler_getDataPribadi(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into kepegawaian.pegawai
		("ID", "PNS_ID", "NIP_LAMA", "NIP_BARU", "NAMA", "GELAR_DEPAN", "GELAR_BELAKANG", "TEMPAT_LAHIR_ID", "TGL_LAHIR", "JENIS_KELAMIN", "AGAMA_ID", "JENIS_KAWIN_ID", "NIK", "NOMOR_DARURAT", "NOMOR_HP", "EMAIL", "ALAMAT", "NPWP", "BPJS", "JENIS_PEGAWAI_ID", "KEDUDUKAN_HUKUM_ID", "STATUS_CPNS_PNS", "KARTU_PEGAWAI", "NOMOR_SK_CPNS", "TGL_SK_CPNS", "TMT_CPNS",  "TMT_PNS",   "GOL_AWAL_ID", "GOL_ID", "TMT_GOLONGAN", "MK_TAHUN", "MK_BULAN", "JENIS_JABATAN_IDx", "JABATAN_ID", "TMT_JABATAN", "PENDIDIKAN_ID", "TAHUN_LULUS", "KPKN_ID", "LOKASI_KERJA_ID", "UNOR_ID", "UNOR_INDUK_ID", "INSTANSI_INDUK_ID", "INSTANSI_KERJA_ID", "SATUAN_KERJA_INDUK_ID", "SATUAN_KERJA_KERJA_ID", "GOLONGAN_DARAH", "PHOTO", "TMT_PENSIUN", "LOKASI_KERJA", "JML_ISTRI", "JML_ANAK", "NO_SURAT_DOKTER", "TGL_SURAT_DOKTER", "NO_BEBAS_NARKOBA", "TGL_BEBAS_NARKOBA", "NO_CATATAN_POLISI", "TGL_CATATAN_POLISI", "AKTE_KELAHIRAN", "STATUS_HIDUP", "AKTE_MENINGGAL", "TGL_MENINGGAL", "NO_ASKES", "NO_TASPEN", "TGL_NPWP",  "TEMPAT_LAHIR", "PENDIDIKAN", "TK_PENDIDIKAN", "TEMPAT_LAHIR_NAMA", "JENIS_JABATAN_NAMA", "JABATAN_NAMA", "KPKN_NAMA", "INSTANSI_INDUK_NAMA", "INSTANSI_KERJA_NAMA", "SATUAN_KERJA_INDUK_NAMA", "SATUAN_KERJA_NAMA", "JABATAN_INSTANSI_ID", "BUP", "JABATAN_INSTANSI_NAMA", "JENIS_JABATAN_ID", terminated_date, status_pegawai, "JABATAN_PPNPN", "JABATAN_INSTANSI_REAL_ID", "CREATED_DATE", "CREATED_BY", "UPDATED_DATE", "UPDATED_BY", "EMAIL_DIKBUD_BAK", "EMAIL_DIKBUD", "KODECEPAT", "IS_DOSEN", "MK_TAHUN_SWASTA", "MK_BULAN_SWASTA", "KK", "NIDN", "KET", "NO_SK_PEMBERHENTIAN", status_pegawai_backup, "MASA_KERJA", "KARTU_ASN") values
		(11,   '1a',     '1b',       '1c',       '1d',   '1e',          '1f',             '1g',              '2000-01-02','1h',            21,         '5',              '1k',  '1l',            '1m',       '1n',    '1o',     '1p',   '1q',   '31',               '1s',                 '1t',              '1u',            '1v',            '2000-01-03',  '2000-01-04','2000-01-04','1z',          1,        '2000-01-05',   '1ac',      '1ad',      '1ae',               '1af',        '2000-01-06',  '1ah',           '1ai',         '1aj',     '1ak',             '1al',     '1am',           '1an',               '1ao',               '1ap',                   '1aq',                   '1ar',            '1as',   '2000-01-07',  '1au',          '1',         '1',        '1ax',             '2000-01-08',       '1az',              '2000-01-09',        '1bb',               '2000-01-10',         '1bd',            '1be',          '1bf',            '2000-01-11',    '1bh',      'bi',        '2000-01-12','1bk',          '1bl',        '1bm',           '1bn',               '1bo',                '1bp',          '1bq',       '1br',                 '1bs',                 '1bt',                     '1bu',               '1bv',                 1,     '1bx',                   1,                  '2000-01-13',    '1',            '1cb',           '1cc',                      '2000-01-14',   1,            '2000-01-15',   1,            '1ch',              '1ci',          '1cj',       1,          1,                 1,                 '1cn','1co',  '1cp', '1cq',                 1,                     '1cs',        '1ct');

		insert into kepegawaian.agama
		("ID", "NAMA", "NCSISTIME") values
		(21,   '21a',  '21b');

		insert into kepegawaian.jenis_pegawai
		("ID", "NAMA") values
		('31', '31a');

		insert into kepegawaian.users
		(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
		(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);

		insert into kepegawaian.jenis_kawin
		("ID", "NAMA") values
		('5',  '5a');

		insert into kepegawaian.unitkerja
		("ID",  "NO",  "KODE_INTERNAL", "NAMA_UNOR", "ESELON_ID", "CEPAT_KODE", "NAMA_JABATAN", "NAMA_PEJABAT", "DIATASAN_ID", "INSTANSI_ID", "PEMIMPIN_NON_PNS_ID", "PEMIMPIN_PNS_ID", "JENIS_UNOR_ID", "UNOR_INDUK", "JUMLAH_IDEAL_STAFF", "ORDER", "deleted", "IS_SATKER", "ESELON_1", "ESELON_2", "ESELON_3", "ESELON_4", "EXPIRED_DATE", "KETERANGAN", "JENIS_SATKER", "ABBREVIATION", "UNOR_INDUK_PENYETARAAN", "JABATAN_ID", "WAKTU", "PERATURAN") values
		('1al', '61a', '61b',           '61c',       '61d',       '61e',        '61f',          '61g',          '61h',         '61i',         '61j',                 '61k',             '61l',           '61m',        '61n',                61,      61,        61,          '61o',      '61p',      '61q',      '61r',      '2000-01-01',   '61s',        '61t',          '61u',          '61v',                    '61w',        '61x',   '61y');
	`

	tests := []struct {
		name             string
		dbData           string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(41)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
				{
					"data": {
						"agama": "21a",
						"akte_kelahiran": "1bd",
						"alamat": "1o",
						"email_dikbud": "1ci",
						"email_lain": "1n",
						"gaji_pokok": "TODO: GajiPokok",
						"gelar_belakang": "1f",
						"gelar_depan": "1e",
						"golongan_ruang_awal": "TODO: GolonganRuangAwal",
						"golongan_ruang_terakhir": "TODO: GolonganRuangTerakhir",
						"jabatan": "1bp",
						"unit_kerja": "61c",
						"id": 11,
						"jenis_kelamin": "1h",
						"jenis_pegawai": "31a",
						"kartu_pegawai": "1u",
						"lokasi_kerja": "1au",
						"masa_kerja": "1cs",
						"nama": "1d",
						"nik": "1k",
						"nip": "1b",
						"nip_baru": "1c",
						"nomor_bpjs": "1q",
						"nomor_catatan_polisi": "1bb",
						"nomor_darurat": "1l",
						"nomor_hp": "1m",
						"nomor_surat_bebas_narkoba": "1az",
						"nomor_surat_dokter": "1ax",
						"npwp": "1p",
						"pangkat_golongan_aktif": "TODO: PangkatGolonganAktif",
						"pendidikan": "1bl",
						"photo": "1as",
						"sk_asn": "TODO: SKASN",
						"status_perkawinan": "5a",
						"status_pns": "TODO: StatusPNS",
						"tanggal_catatan_polisi": "2000-01-10",
						"tanggal_lahir": "2000-01-02",
						"tanggal_npwp": "2000-01-12",
						"tanggal_surat_bebas_narkoba": "2000-01-09",
						"tanggal_surat_dokter": "2000-01-08",
						"tempat_lahir": "1bk",
						"tingkat_pendidikan": "1bm",
						"tmt_asn": "1990-01-01",
						"tmt_golongan": "2000-01-05"
					}
				}
			`,
		},
		{
			name:             "error: tidak ada data pegawai milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(200)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
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

			req := httptest.NewRequest(http.MethodGet, "/data-pribadi", nil)
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

func Test_handler_listStatusPernikahan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into kepegawaian.jenis_kawin
			("ID", "NAMA") values
			('1',  'a'),
			('2',  'c'),
			('3',  'b');
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
					{"id": "1", "nama": "a"},
					{"id": "3", "nama": "b"},
					{"id": "2", "nama": "c"}
				]
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/status-pernikahan", nil)
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
