package pegawai

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
		insert into kepegawaian.golongan
			("ID", "NAMA", "NAMA_PANGKAT", "NAMA2", "GOL", "GOL_PPPK") values
			(11,   '11a',  '11b',          '11c',   11,    '11d'),
			(12,   '12a',  '12b',          '12c',   12,    '12d');
		insert into kepegawaian.unitkerja
			("ID", "NO",  "KODE_INTERNAL", "NAMA_UNOR", "ESELON_ID", "CEPAT_KODE", "NAMA_JABATAN", "NAMA_PEJABAT", "DIATASAN_ID", "INSTANSI_ID", "PEMIMPIN_NON_PNS_ID", "PEMIMPIN_PNS_ID", "JENIS_UNOR_ID", "UNOR_INDUK", "JUMLAH_IDEAL_STAFF", "ORDER", "deleted", "IS_SATKER", "ESELON_1", "ESELON_2", "ESELON_3", "ESELON_4", "EXPIRED_DATE", "KETERANGAN", "JENIS_SATKER", "ABBREVIATION", "UNOR_INDUK_PENYETARAAN", "JABATAN_ID", "WAKTU", "PERATURAN") values
			('31', '31a', '31b',           '31c',       '31d',       '31e',        '31f',          '31g',          '31h',         '31i',         '31j',                 '31k',             '31l',           '31m',        '31n',                31,      31,        31,          '31o',      '31p',      '31q',      '31r',      '2000-01-01',   '31s',        '31t',          '31u',          '31v',                    '31w',        '31x',   '31y'),
			('32', '32a', '32b',           '32c',       '32d',       '32e',        '32f',          '32g',          '32h',         '32i',         '32j',                 '32k',             '32l',           '32m',        '32n',                32,      32,        32,          '32o',      '32p',      '32q',      '32r',      '2000-01-01',   '32s',        '32t',          '32u',          '32v',                    '32w',        '32x',   '32y');
		insert into kepegawaian.pegawai
			("ID", "PNS_ID", "NIP_LAMA", "NIP_BARU", "NAMA", "GELAR_DEPAN", "GELAR_BELAKANG", "TEMPAT_LAHIR_ID", "TGL_LAHIR",  "JENIS_KELAMIN", "AGAMA_ID", "JENIS_KAWIN_ID", "NIK", "NOMOR_DARURAT", "NOMOR_HP", "EMAIL", "ALAMAT", "NPWP", "BPJS", "JENIS_PEGAWAI_ID", "KEDUDUKAN_HUKUM_ID", "STATUS_CPNS_PNS", "KARTU_PEGAWAI", "NOMOR_SK_CPNS", "TGL_SK_CPNS", "TMT_CPNS",   "TMT_PNS",    "GOL_AWAL_ID", "GOL_ID", "TMT_GOLONGAN", "MK_TAHUN", "MK_BULAN", "JENIS_JABATAN_IDx", "JABATAN_ID", "TMT_JABATAN", "PENDIDIKAN_ID", "TAHUN_LULUS", "KPKN_ID", "LOKASI_KERJA_ID", "UNOR_ID", "UNOR_INDUK_ID", "INSTANSI_INDUK_ID", "INSTANSI_KERJA_ID", "SATUAN_KERJA_INDUK_ID", "SATUAN_KERJA_KERJA_ID", "GOLONGAN_DARAH", "PHOTO", "TMT_PENSIUN", "LOKASI_KERJA", "JML_ISTRI", "JML_ANAK", "NO_SURAT_DOKTER", "TGL_SURAT_DOKTER", "NO_BEBAS_NARKOBA", "TGL_BEBAS_NARKOBA", "NO_CATATAN_POLISI", "TGL_CATATAN_POLISI", "AKTE_KELAHIRAN", "STATUS_HIDUP", "AKTE_MENINGGAL", "TGL_MENINGGAL", "NO_ASKES", "NO_TASPEN", "TGL_NPWP",   "TEMPAT_LAHIR", "PENDIDIKAN", "TK_PENDIDIKAN", "TEMPAT_LAHIR_NAMA", "JENIS_JABATAN_NAMA", "JABATAN_NAMA", "KPKN_NAMA", "INSTANSI_INDUK_NAMA", "INSTANSI_KERJA_NAMA", "SATUAN_KERJA_INDUK_NAMA", "SATUAN_KERJA_NAMA", "JABATAN_INSTANSI_ID", "BUP",  "JABATAN_INSTANSI_NAMA", "JENIS_JABATAN_ID", terminated_date, status_pegawai, "JABATAN_PPNPN", "JABATAN_INSTANSI_REAL_ID", "CREATED_DATE", "CREATED_BY", "UPDATED_DATE", "UPDATED_BY", "EMAIL_DIKBUD_BAK", "EMAIL_DIKBUD", "KODECEPAT", "IS_DOSEN", "MK_TAHUN_SWASTA", "MK_BULAN_SWASTA", "KK",   "NIDN", "KET",  "NO_SK_PEMBERHENTIAN", "status_pegawai_backup", "MASA_KERJA", "KARTU_ASN") values
			(41,   '41a',    '41b',      '41c',      '41d',  '41e',         '41f',            '41g',             '2000-01-01', '41h',           41,         '41i',            '41j', '41k',           '41l',      '41m',   '41n',    '41o',  '41p',  '41q',              '41r',                'PNS',             '41t',           '41u',           '2000-01-02',  '2000-01-03', '2000-01-04', '41v',         11,       '2000-01-05',   '41w',      '41x',      '41y',               '21',         '2000-01-06',  '41aa',          '41ab',        '41ac',    '41ad',            '31',      '41af',          '41ag',              '41ah',              '41ai',                  '41aj',                  '41ak',           '41al',  null,          '41am',         '1',         '1',        '41an',            '2000-01-08',       '41ao',             '2000-01-09',        '41ap',              '2000-01-10',         '41aq',           '41ar',         '41as',           null,            '41at',     '41au',      '2000-01-12', '41av',         '41aw',       '41',            '41ax',              '41ay',               '41az',         '41ba',      '41bb',                '41bc',                '41bd',                    '41be',              '41bf',                41,     '41bh',                  41,                 '2000-01-13',    41,             '41bi',          '41bj',                     '2000-01-14',   41,           '2000-01-15',   41,           '41bm',             '41bn',         '41bo',      41,         41,                41,                '41bp', '41bq', '41br', null,                  41,                      '41bu',       '41bv'),
			(42,   '42a',    '42b',      '42c',      '42d',  '42e',         '42f',            '42g',             '2001-01-01', '42h',           42,         '42i',            '42j', '42k',           '42l',      '42m',   '42n',    '42o',  '42p',  '42q',              '42r',                'CPNS',             '42t',           '42u',           '2001-01-02',  '2001-01-03', '2001-01-04', '42v',         12,       '2001-01-05',   '42w',      '42x',      '42y',               '22',         '2001-01-06',  '42aa',          '42ab',        '42ac',    '42ad',            '32',      '42af',          '42ag',              '42ah',              '42ai',                  '42aj',                  '42ak',           '42al',  null,          '42am',         '2',         '2',        '42an',            '2001-01-08',       '42ao',             '2001-01-09',        '42ap',              '2001-01-10',         '42aq',           '42ar',         '42as',           null,            '42at',     '42au',      '2001-01-12', '42av',         '42aw',       '42',            '42ax',              '42ay',               '42az',         '42ba',      '42bb',                '42bc',                '42bd',                    '42be',              '42bf',                42,     '42bh',                  42,                 '2001-01-13',    42,             '42bi',          '42bj',                     '2001-01-14',   42,           '2001-01-15',   42,           '42bm',             '42bn',         '42bo',      42,         42,                42,                '42bp', '42bq', '42br', null,                  42,                      '42bu',       '42bv'),
			(43,   '43a',    '43b',      '43c',      '43d',  '43e',         '43f',            '43g',             '2002-01-01', '43h',           43,         '43i',            '43j', '43k',           '43l',      '43m',   '43n',    '43o',  '43p',  '43q',              '43r',                'PNS',             '43t',           '43u',           '2002-01-02',  '2002-01-03', '2002-01-04', '43v',         11,       '2002-01-05',   '43w',      '43x',      '43y',               '21',         '2002-01-06',  '43aa',          '43ab',        '43ac',    '43ad',            '31',      '43af',          '43ag',              '43ah',              '43ai',                  '43aj',                  '43ak',           '43al',  null,          '43am',         '3',         '3',        '43an',            '2002-01-08',       '43ao',             '2002-01-09',        '43ap',              '2002-01-10',         '43aq',           '43ar',         '43as',           null,            '43at',     '43au',      '2002-01-12', '43av',         '43aw',       '43',            '43ax',              '43ay',               '43az',         '43ba',      '43bb',                '43bc',                '43bd',                    '43be',              '43bf',                43,     '43bh',                  43,                 '2002-01-13',    43,             '43bi',          '43bj',                     '2002-01-14',   43,           '2002-01-15',   43,           '43bm',             '43bn',         '43bo',      43,         43,                43,                '43bp', '43bq', '43br', null,                  43,                      '43bu',       '43bv'),
			(44,   '44a',    '44b',      '44c',      '44d',  '44e',         '44f',            '44g',             '2003-01-01', '44h',           44,         '44i',            '44j', '44k',           '44l',      '44m',   '44n',    '44o',  '44p',  '44q',              '44r',                'CPNS',             '44t',           '44u',           '2003-01-02',  '2003-01-03', '2003-01-04', '44v',         12,       '2003-01-05',   '44w',      '44x',      '44y',               '22',         '2003-01-06',  '44aa',          '44ab',        '44ac',    '44ad',            '32',      '44af',          '44ag',              '44ah',              '44ai',                  '44aj',                  '44ak',           '44al',  null,          '44am',         '4',         '4',        '44an',            '2003-01-08',       '44ao',             '2003-01-09',        '44ap',              '2003-01-10',         '44aq',           '44ar',         '44as',           null,            '44at',     '44au',      '2003-01-12', '44av',         '44aw',       '44',            '44ax',              '44ay',               '44az',         '44ba',      '44bb',                '44bc',                '44bd',                    '44be',              '44bf',                44,     '44bh',                  44,                 '2003-01-13',    44,             '44bi',          '44bj',                     '2003-01-14',   44,           '2003-01-15',   44,           '44bm',             '44bn',         '44bo',      44,         44,                44,                '44bp', '44bq', '44br', null,                  44,                      '44bu',       '44bv'),
			(45,   '45a',    '45b',      '45c',      '45d',  '45e',         '45f',            '45g',             '2004-01-01', '45h',           45,         '45i',            '45j', '45k',           '45l',      '45m',   '45n',    '45o',  '45p',  '45q',              '45r',                'PNS',             '45t',           '45u',           '2004-01-02',  '2004-01-03', '2004-01-04', '45v',         11,       '2004-01-05',   '45w',      '45x',      '45y',               '21',         '2004-01-06',  '45aa',          '45ab',        '45ac',    '45ad',            '31',      '45af',          '45ag',              '45ah',              '45ai',                  '45aj',                  '45ak',           '45al',  null,          '45am',         '5',         '5',        '45an',            '2004-01-08',       '45ao',             '2004-01-09',        '45ap',              '2004-01-10',         '45aq',           '45ar',         '45as',           null,            '45at',     '45au',      '2004-01-12', '45av',         '45aw',       '45',            '45ax',              '45ay',               '45az',         '45ba',      '45bb',                '45bc',                '45bd',                    '45be',              '45bf',                45,     '45bh',                  45,                 '2004-01-13',    45,             '45bi',          '45bj',                     '2004-01-14',   45,           '2004-01-15',   45,           '45bm',             '45bn',         '45bo',      45,         45,                45,                '45bp', '45bq', '45br', '45bs',                45,                      '45bu',       '45bv'),
			(46,   '46a',    '46b',      '46c',      '46d',  '46e',         '46f',            '46g',             '2005-01-01', '46h',           46,         '46i',            '46j', '46k',           '46l',      '46m',   '46n',    '46o',  '46p',  '46q',              '46r',                'CPNS',             '46t',           '46u',           '2005-01-02',  '2005-01-03', '2005-01-04', '46v',         12,       '2005-01-05',   '46w',      '46x',      '46y',               '22',         '2005-01-06',  '46aa',          '46ab',        '46ac',    '46ad',            '32',      '46af',          '46ag',              '46ah',              '46ai',                  '46aj',                  '46ak',           '46al',  null,          '46am',         '6',         '6',        '46an',            '2005-01-08',       '46ao',             '2005-01-09',        '46ap',              '2005-01-10',         '46aq',           '46ar',         '46as',           '2005-01-11',    '46at',     '46au',      '2005-01-12', '46av',         '46aw',       '46',            '46ax',              '46ay',               '46az',         '46ba',      '46bb',                '46bc',                '46bd',                    '46be',              '46bf',                46,     '46bh',                  46,                 '2005-01-13',    46,             '46bi',          '46bj',                     '2005-01-14',   46,           '2005-01-15',   46,           '46bm',             '46bn',         '46bo',      46,         46,                46,                '46bp', '46bq', '46br', null,                  46,                      '46bu',       '46bv'),
			(47,   '47a',    '47b',      '47c',      '47d',  '47e',         '47f',            '47g',             '2006-01-01', '47h',           47,         '47i',            '47j', '47k',           '47l',      '47m',   '47n',    '47o',  '47p',  '47q',              '47r',                'PNS',             '47t',           '47u',           '2006-01-02',  '2006-01-03', '2006-01-04', '47v',         11,       '2006-01-05',   '47w',      '47x',      '47y',               '21',         '2006-01-06',  '47aa',          '47ab',        '47ac',    '47ad',            '31',      '47af',          '47ag',              '47ah',              '47ai',                  '47aj',                  '47ak',           '47al',  '2006-01-07',  '47am',         '7',         '7',        '47an',            '2006-01-08',       '47ao',             '2006-01-09',        '47ap',              '2006-01-10',         '47aq',           '47ar',         '47as',           null,            '47at',     '47au',      '2006-01-12', '47av',         '47aw',       '47',            '47ax',              '47ay',               '47az',         '47ba',      '47bb',                '47bc',                '47bd',                    '47be',              '47bf',                47,     '47bh',                  47,                 '2006-01-13',    47,             '47bi',          '47bj',                     '2006-01-14',   47,           '2006-01-15',   47,           '47bm',             '47bn',         '47bo',      47,         47,                47,                '47bp', '47bq', '47br', null,                  47,                      '47bu',       '47bv');
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":             41,
						"gelar_belakang": "41f",
						"gelar_depan":    "41e",
						"golongan":       "11a",
						"jabatan":        "41az",
						"nama_pegawai":   "41d",
						"nip":            "41c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					},
					{
						"id":             42,
						"gelar_belakang": "42f",
						"gelar_depan":    "42e",
						"golongan":       "12a",
						"jabatan":        "42az",
						"nama_pegawai":   "42d",
						"nip":            "42c",
						"status_pegawai": "CPNS",
						"unit_kerja":     "32c"
					},
					{
						"id":             43,
						"gelar_belakang": "43f",
						"gelar_depan":    "43e",
						"golongan":       "11a",
						"jabatan":        "43az",
						"nama_pegawai":   "43d",
						"nip":            "43c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					},
					{
						"id":             44,
						"gelar_belakang": "44f",
						"gelar_depan":    "44e",
						"golongan":       "12a",
						"jabatan":        "44az",
						"nama_pegawai":   "44d",
						"nip":            "44c",
						"status_pegawai": "CPNS",
						"unit_kerja":     "32c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":             42,
						"gelar_belakang": "42f",
						"gelar_depan":    "42e",
						"golongan":       "12a",
						"jabatan":        "42az",
						"nama_pegawai":   "42d",
						"nip":            "42c",
						"status_pegawai": "CPNS",
						"unit_kerja":     "32c"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}`,
		},
		{
			name:             "ok: filter cari nama",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"41d"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":             41,
						"gelar_belakang": "41f",
						"gelar_depan":    "41e",
						"golongan":       "11a",
						"jabatan":        "41az",
						"nama_pegawai":   "41d",
						"nip":            "41c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: filter cari nip",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"42c"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":             42,
						"gelar_belakang": "42f",
						"gelar_depan":    "42e",
						"golongan":       "12a",
						"jabatan":        "42az",
						"nama_pegawai":   "42d",
						"nip":            "42c",
						"status_pegawai": "CPNS",
						"unit_kerja":     "32c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: filter cari jabatan",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"43az"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":             43,
						"gelar_belakang": "43f",
						"gelar_depan":    "43e",
						"golongan":       "11a",
						"jabatan":        "43az",
						"nama_pegawai":   "43d",
						"nip":            "43c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: filter unit",
			dbData:           dbData,
			requestQuery:     url.Values{"unit_id": []string{"31"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"gelar_belakang": "41f",
						"gelar_depan":    "41e",
						"golongan":       "11a",
						"id":             41,
						"jabatan":        "41az",
						"nama_pegawai":   "41d",
						"nip":            "41c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					},
					{
						"gelar_belakang": "43f",
						"gelar_depan":    "43e",
						"golongan":       "11a",
						"id":             43,
						"jabatan":        "43az",
						"nama_pegawai":   "43d",
						"nip":            "43c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok: filter golongan",
			dbData:           dbData,
			requestQuery:     url.Values{"golongan_id": []string{"11"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"gelar_belakang": "41f",
						"gelar_depan":    "41e",
						"golongan":       "11a",
						"id":             41,
						"jabatan":        "41az",
						"nama_pegawai":   "41d",
						"nip":            "41c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					},
					{
						"gelar_belakang": "43f",
						"gelar_depan":    "43e",
						"golongan":       "11a",
						"id":             43,
						"jabatan":        "43az",
						"nama_pegawai":   "43d",
						"nip":            "43c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok: filter jabatan",
			dbData:           dbData,
			requestQuery:     url.Values{"jabatan_id": []string{"21"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"gelar_belakang": "41f",
						"gelar_depan":    "41e",
						"golongan":       "11a",
						"id":             41,
						"jabatan":        "41az",
						"nama_pegawai":   "41d",
						"nip":            "41c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					},
					{
						"gelar_belakang": "43f",
						"gelar_depan":    "43e",
						"golongan":       "11a",
						"id":             43,
						"jabatan":        "43az",
						"nama_pegawai":   "43d",
						"nip":            "43c",
						"status_pegawai": "PNS",
						"unit_kerja":     "31c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok: filter status",
			dbData:           dbData,
			requestQuery:     url.Values{"status": []string{"CPNS"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"gelar_belakang": "42f",
						"gelar_depan":    "42e",
						"golongan":       "12a",
						"id":             42,
						"jabatan":        "42az",
						"nama_pegawai":   "42d",
						"nip":            "42c",
						"status_pegawai": "CPNS",
						"unit_kerja":     "32c"
					},
					{
						"gelar_belakang": "44f",
						"gelar_depan":    "44e",
						"golongan":       "12a",
						"id":             44,
						"jabatan":        "44az",
						"nama_pegawai":   "44d",
						"nip":            "44c",
						"status_pegawai": "CPNS",
						"unit_kerja":     "32c"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok: tidak ada data ditemukan",
			dbData:           dbData,
			requestQuery:     url.Values{"cari": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: parameter 'status' invalid",
			dbData:           dbData,
			requestQuery:     url.Values{"status": []string{"KONTRAK"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "admin")}},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"status\" harus salah satu dari \"PNS\", \"CPNS\""}`,
		},
		{
			name:             "error: role user bukan admin",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(1, "bukan_admin")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "Forbidden"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{"Bearer admin"}},
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

			req := httptest.NewRequest(http.MethodGet, "/pegawai", nil)
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

func Test_handler_listStatusPegawai(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into kepegawaian.jenis_pegawai
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

			req := httptest.NewRequest(http.MethodGet, "/status-pegawai", nil)
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
