package datapribadi

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_getDataPribadi(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_jabatan
			(id, nama,             deleted_at) values
			(1, 'Jenis Jabatan 1', null),
			(2, 'Jenis Jabatan 2', null),
			(3, 'Jenis Jabatan 3', null),
			(4, 'Jenis Jabatan 4', '2000-01-01');
		insert into ref_kelas_jabatan
			(id, kelas_jabatan) values
			(1, 'Kelas 1');
		insert into ref_jabatan
			(id, no, kode_jabatan, jenis_jabatan, kelas, nama_jabatan, deleted_at) values
			(1,  1,  'KJ1',        1,             1,     'Jabatan 1',  null),
			(2,  2,  'KJ2',        2,             1,     'Jabatan 2',  null),
			(3,  3,  'KJ3',        3,             1,     'Jabatan 3',  null),
			(4,  4,  'KJ4',        null,          null,  'Jabatan 4',  null),
			(5,  5,  'KJ5',        2,             1,     'Jabatan 5',  '2000-01-01'),
			(6,  6,  'KJ6',        4,             1,     'Jabatan 6',  null);
		insert into ref_jenis_kawin
			(id, nama,      deleted_at) values
			(1,  'Menikah', null),
			(2,  'Cerai',   '2000-01-01');
		insert into ref_kedudukan_hukum
			(id, nama,  is_pppk, deleted_at) values
			(1,  'P3K', true,    null),
			(2,  'PNS', false,   null),
			(3,  'TNI', false,   '2000-01-01'),
			(4,  'POL', null,    null);
		insert into ref_golongan
			(id, nama,  nama_pangkat, gol_pppk, deleted_at) values
			(1,  'I/a', 'Pangkat 1',  'I',      null),
			(2,  'I/b', 'Pangkat 2',  'II',     null),
			(3,  'I/c', 'Pangkat 3',  'III',    '2000-01-01');
		insert into ref_jenis_pegawai
			(id, nama,              deleted_at) values
			(1,  'Jenis Pegawai 1', null),
			(2,  'Jenis Pegawai 2', '2000-01-01');
		insert into ref_lokasi
			(id, nama,      deleted_at) values
			(1,  'Jakarta', null),
			(2,  'Medan',   null),
			(3,  'Jogja',   '2000-01-01');
		insert into ref_pendidikan
			(id,  nama,           deleted_at) values
			('1', 'Pendidikan 1', null),
			('2', 'Pendidikan 2', '2000-01-01');
		insert into ref_tingkat_pendidikan
			(id, nama,        deleted_at) values
			(1,  'Tingkat 1', null),
			(2,  'Tingkat 2', '2000-01-01');
		insert into ref_agama
			(id, nama,      deleted_at) values
			(1,  'Kristen', null),
			(2,  'Katolik', '2000-01-01');
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, nama_jabatan,    pemimpin_pns_id, deleted_at) values
			('0', '1',         'Unor 0',  'Unit Kerja 0',  null,            null),
			('1', '2',         'Unor 1',  'Unit Kerja 1',  null,            null),
			('2', '3',         'Unor 2',  'Unit Kerja 2',  null,            null),
			('3', '4',         'Unor 3',  'Unit Kerja 3',  null,            null),
			('4', '5',         'Unor 4',  'Unit Kerja 4',  null,            null),
			('5', '6',         'Unor 5',  'Unit Kerja 5',  null,            null),
			('6', '7',         'Unor 6',  'Unit Kerja 6',  null,            null),
			('7', '8',         'Unor 7',  'Unit Kerja 7',  null,            null),
			('8', '9',         'Unor 8',  'Unit Kerja 8',  null,            null),
			('9', 'A',         'Unor 9',  'Unit Kerja 9',  null,            null),
			('A', 'B',         'Unor A',  'Unit Kerja A',  null,            null),
			('B', 'C',         'Unor B',  'Unit Kerja B',  null,            null),
			('C', 'D',         'Unor C',  'Unit Kerja C',  null,            null),
			('D', null,        'Unor D',  'Unit Kerja D',  null,            null),
			('E', 'F',         'Unor E',  'Unit Kerja E',  null,            null),
			('F', 'G',         '',        '',              null,            null),
			('G', 'H',         'Unor G',  'Unit Kerja G',  null,            null),
			('H', '6',         'Unor H',  'Unit Kerja H',  null,            '2000-01-01');
		insert into pegawai
			(nip_baru, pns_id, nama, gelar_depan, gelar_belakang, jenis_jabatan_id, tmt_jabatan, unor_id, nik, kk, jenis_kelamin, tempat_lahir, tempat_lahir_id, tanggal_lahir, tingkat_pendidikan_id, pendidikan_id, jenis_kawin_id, agama_id, email_dikbud, email, alamat, no_hp, no_darurat, jenis_pegawai_id, status_pegawai, tmt_cpns, mk_tahun_swasta, mk_bulan_swasta, masa_kerja, jabatan_instansi_id, jabatan_instansi_real_id, lokasi_kerja, lokasi_kerja_id, gol_awal_id, gol_id, tmt_golongan, no_sk_cpns, status_cpns_pns, tmt_pns, kartu_pegawai, no_surat_dokter, tanggal_surat_dokter, no_bebas_narkoba, tanggal_bebas_narkoba, no_catatan_polisi, tanggal_catatan_polisi, akte_kelahiran, bpjs, npwp, tanggal_npwp, no_taspen, kedudukan_hukum_id, terminated_date, deleted_at) values
			('1c', 'PNS_1c', 'Budi Santoso', 'Dr.', 'M.Sc', 2, '2020-01-01', '0', '3173000000000001', '3173000000000002', 'L', 'DKI Jakarta', 1, '1990-05-20', 1, '1', 1, 1, 'budi@dikbud.go.id', 'budi@gmail.com', 'Jl. Merdeka No. 123', '08123456789', '08198765432', 1, 1, '2015-06-01', 2, 12, '5 Tahun 6 Bulan', 'KJ2', 'KJ3', 'Jakarta HQ', 2, 1, 2, '2018-01-01', 'SK-CPNS-2015', 'PNS', '2017-01-01', 'KARPEG001', 'DOC-HEALTH-001', '2015-05-01', 'BN-001', '2015-05-02', 'SKCK-001', '2015-05-03', 'AKTE-001', 'BPJS-001', 'NPWP-001', '2016-01-01', 'TASPEN-001', 2, null, null),
			('1d', 'PNS_1d', 'John Doe', null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null, null),
			('1e', 'PNS_1e', 'John Santoso', 'Dr.', 'M.Sc', 4, '2020-01-01', 'E', '3173000000000001', '3173000000000002', 'L', 'DKI Jakarta', 3, '1990-05-20', 2, '2', 2, 2, 'budi@dikbud.go.id', 'budi@gmail.com', 'Jl. Merdeka No. 123', '08123456789', '08198765432', 2, 1, null, 2, 12, '5 Tahun 6 Bulan', 'KJ5', 'KJ5', 'Jakarta HQ', 3, 3, 3, '2018-01-01', 'SK-CPNS-2015', 'CPNS', '2017-01-01', 'KARPEG001', 'DOC-HEALTH-001', '2015-05-01', 'BN-001', '2015-05-02', 'SKCK-001', '2015-05-03', 'AKTE-001', 'BPJS-001', 'NPWP-001', '2016-01-01', 'TASPEN-001', 3, null, null),
			('1f', 'PNS001', 'Budi John', '', '', 1, '2020-01-01', '8', '3173000000000001', '3173000000000002', 'L', 'DKI Jakarta', 1, '1990-05-20', 1, '1', 1, 1, 'budi@dikbud.go.id', 'budi@gmail.com', 'Jl. Merdeka No. 123', '08123456789', '08198765432', 1, 1, '2015-06-01', null, null, '5 Tahun 6 Bulan', 'KJ4', 'KJ6', 'Jakarta HQ', 1, 1, 2, '2018-01-01', 'SK-CPNS-2015', 'C', '2017-01-01', 'KARPEG001', 'DOC-HEALTH-001', '2015-05-01', 'BN-001', '2015-05-02', 'SKCK-001', '2015-05-03', 'AKTE-001', 'BPJS-001', 'NPWP-001', '2016-01-01', 'TASPEN-001', 1, '3000-01-01', null),
			('1g', 'PNS_1g', 'Budi Santoso', 'Dr.', 'M.Sc', 1, '2020-01-01', 'C', '3173000000000001', '3173000000000002', 'L', 'DKI Jakarta', 1, '1990-05-20', 1, '1', 1, 1, 'budi@dikbud.go.id', 'budi@gmail.com', 'Jl. Merdeka No. 123', '08123456789', '08198765432', 1, 1, '2015-06-01', 2, 4, '5 Tahun 6 Bulan', 'KJ1', 'KJ6', 'Jakarta HQ', 2, 1, 2, '2018-01-01', 'SK-CPNS-2015', 'P', null, 'KARPEG001', 'DOC-HEALTH-001', '2015-05-01', 'BN-001', '2015-05-02', 'SKCK-001', '2015-05-03', 'AKTE-001', 'BPJS-001', 'NPWP-001', '2016-01-01', 'TASPEN-001', 4, '2000-01-01', null),
			('2c', 'PNS002', 'Budi Santoso', 'Dr.', 'M.Sc', 2, '2020-01-01', '0', '3173000000000001', '3173000000000002', 'L', 'DKI Jakarta', 1, '1990-05-20', 1, '1', 1, 1, 'budi@dikbud.go.id', 'budi@gmail.com', 'Jl. Merdeka No. 123', '08123456789', '08198765432', 1, 1, '2015-06-01', 2, 12, '5 Tahun 6 Bulan', 'KJ2', 'KJ3', 'Jakarta HQ', 1, 1, 2, '2018-01-01', 'SK-CPNS-2015', 'CPNS', '2017-01-01', 'KARPEG001', 'DOC-HEALTH-001', '2015-05-01', 'BN-001', '2015-05-02', 'SKCK-001', '2015-05-03', 'AKTE-001', 'BPJS-001', 'NPWP-001', '2016-01-01', 'TASPEN-001', 2, null, '2000-01-01');
		update unit_kerja set pemimpin_pns_id = 'PNS001' where id = '8';
	`

	masaKerjaKeseluruhan := func(date time.Time, tahun, bulan int) string {
		tahun += time.Now().Year() - date.Year()
		bulan += int(time.Now().Month()) - int(date.Month())
		if time.Now().Day() < date.Day() {
			bulan--
		}
		if bulan < 0 {
			tahun, bulan = tahun-1, bulan+12
		} else if bulan >= 12 {
			tahun, bulan = tahun+bulan/12, bulan%12
		}
		return fmt.Sprintf("%d Tahun %d Bulan", tahun, bulan)
	}

	tests := []struct {
		name             string
		dbData           string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: non pppk with status_pns & tmt_pns",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
				{
					"data": {
						"nama":                        "Budi Santoso",
						"gelar_depan":                 "Dr.",
						"gelar_belakang":              "M.Sc",
						"jabatan_aktual":              "Jabatan 3",
						"jenis_jabatan_aktual":        "Jenis Jabatan 3",
						"tmt_jabatan":                 "2020-01-01",
						"nip":                         "1c",
						"nik":                         "3173000000000001",
						"nomor_kk":                    "3173000000000002",
						"jenis_kelamin":               "L",
						"tempat_lahir":                "Jakarta",
						"tanggal_lahir":               "1990-05-20",
						"tingkat_pendidikan":          "Tingkat 1",
						"pendidikan":                  "Pendidikan 1",
						"status_perkawinan":           "Menikah",
						"agama":                       "Kristen",
						"email_dikbud":                "budi@dikbud.go.id",
						"email_pribadi":               "budi@gmail.com",
						"alamat":                      "Jl. Merdeka No. 123",
						"nomor_hp":                    "08123456789",
						"nomor_kontak_darurat":        "08198765432",
						"jenis_pegawai":               "Jenis Pegawai 1",
						"masa_kerja_keseluruhan":      "` + masaKerjaKeseluruhan(time.Date(2015, 6, 1, 0, 0, 0, 0, time.Local), 2, 12) + `",
						"masa_kerja_golongan":         "5 Tahun 6 Bulan",
						"jabatan":                     "Jabatan 2",
						"jenis_jabatan":               "Jenis Jabatan 2",
						"kelas_jabatan":               "Kelas 1",
						"lokasi_kerja":                "Medan",
						"golongan_ruang_awal":         "I/a",
						"golongan_ruang_akhir":        "I/b",
						"pangkat_akhir":               "Pangkat 2",
						"tmt_golongan":                "2018-01-01",
						"tmt_asn":                     "2015-06-01",
						"nomor_sk_asn":                "SK-CPNS-2015",
						"is_pppk":                     false,
						"status_asn":                  "PNS",
						"status_pns":                  "PNS",
						"tmt_pns":                     "2017-01-01",
						"kartu_pegawai":               "KARPEG001",
						"nomor_surat_dokter":          "DOC-HEALTH-001",
						"tanggal_surat_dokter":        "2015-05-01",
						"nomor_surat_bebas_narkoba":   "BN-001",
						"tanggal_surat_bebas_narkoba": "2015-05-02",
						"nomor_catatan_polisi":        "SKCK-001",
						"tanggal_catatan_polisi":      "2015-05-03",
						"akte_kelahiran":              "AKTE-001",
						"nomor_bpjs":                  "BPJS-001",
						"npwp":                        "NPWP-001",
						"tanggal_npwp":                "2016-01-01",
						"nomor_taspen":                "TASPEN-001",
						"unit_organisasi":             ["Unor 0", "Unor 1", "Unor 2", "Unor 3", "Unor 4", "Unor 5", "Unor 6", "Unor 7", "Unor 8", "Unor 9"]
					}
				}
			`,
		},
		{
			name:             "ok: most data is null",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1d")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
				{
					"data": {
						"nama":                        "John Doe",
						"gelar_depan":                 "",
						"gelar_belakang":              "",
						"jabatan_aktual":              "",
						"jenis_jabatan_aktual":        "",
						"tmt_jabatan":                 null,
						"nip":                         "1d",
						"nik":                         "",
						"nomor_kk":                    "",
						"jenis_kelamin":               "",
						"tempat_lahir":                "",
						"tanggal_lahir":               null,
						"tingkat_pendidikan":          "",
						"pendidikan":                  "",
						"status_perkawinan":           "",
						"agama":                       "",
						"email_dikbud":                "",
						"email_pribadi":               "",
						"alamat":                      "",
						"nomor_hp":                    "",
						"nomor_kontak_darurat":        "",
						"jenis_pegawai":               "",
						"masa_kerja_keseluruhan":      "",
						"masa_kerja_golongan":         "",
						"jabatan":                     "",
						"jenis_jabatan":               "",
						"kelas_jabatan":               "",
						"lokasi_kerja":                "",
						"golongan_ruang_awal":         "",
						"golongan_ruang_akhir":        "",
						"pangkat_akhir":               "",
						"tmt_golongan":                null,
						"tmt_asn":                     null,
						"nomor_sk_asn":                "",
						"is_pppk":                     false,
						"status_asn":                  "",
						"status_pns":                  "",
						"tmt_pns":                     null,
						"kartu_pegawai":               "",
						"nomor_surat_dokter":          "",
						"tanggal_surat_dokter":        null,
						"nomor_surat_bebas_narkoba":   "",
						"tanggal_surat_bebas_narkoba": null,
						"nomor_catatan_polisi":        "",
						"tanggal_catatan_polisi":      null,
						"akte_kelahiran":              "",
						"nomor_bpjs":                  "",
						"npwp":                        "",
						"tanggal_npwp":                null,
						"nomor_taspen":                "",
						"unit_organisasi":             []
					}
				}
			`,
		},
		{
			name:             "ok: references record is deleted with empty tmt_cpns without status_pns",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1e")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
				{
					"data": {
						"nama":                        "John Santoso",
						"gelar_depan":                 "Dr.",
						"gelar_belakang":              "M.Sc",
						"jabatan_aktual":              "",
						"jenis_jabatan_aktual":        "",
						"tmt_jabatan":                 "2020-01-01",
						"nip":                         "1e",
						"nik":                         "3173000000000001",
						"nomor_kk":                    "3173000000000002",
						"jenis_kelamin":               "L",
						"tempat_lahir":                "DKI Jakarta",
						"tanggal_lahir":               "1990-05-20",
						"tingkat_pendidikan":          "",
						"pendidikan":                  "",
						"status_perkawinan":           "",
						"agama":                       "",
						"email_dikbud":                "budi@dikbud.go.id",
						"email_pribadi":               "budi@gmail.com",
						"alamat":                      "Jl. Merdeka No. 123",
						"nomor_hp":                    "08123456789",
						"nomor_kontak_darurat":        "08198765432",
						"jenis_pegawai":               "",
						"masa_kerja_keseluruhan":      "3 Tahun 0 Bulan",
						"masa_kerja_golongan":         "5 Tahun 6 Bulan",
						"jabatan":                     "",
						"jenis_jabatan":               "",
						"kelas_jabatan":               "",
						"lokasi_kerja":                "Jakarta HQ",
						"golongan_ruang_awal":         "",
						"golongan_ruang_akhir":        "",
						"pangkat_akhir":               "",
						"tmt_golongan":                "2018-01-01",
						"tmt_asn":                     null,
						"nomor_sk_asn":                "SK-CPNS-2015",
						"is_pppk":                     false,
						"status_asn":                  "",
						"status_pns":                  "",
						"tmt_pns":                     null,
						"kartu_pegawai":               "KARPEG001",
						"nomor_surat_dokter":          "DOC-HEALTH-001",
						"tanggal_surat_dokter":        "2015-05-01",
						"nomor_surat_bebas_narkoba":   "BN-001",
						"tanggal_surat_bebas_narkoba": "2015-05-02",
						"nomor_catatan_polisi":        "SKCK-001",
						"tanggal_catatan_polisi":      "2015-05-03",
						"akte_kelahiran":              "AKTE-001",
						"nomor_bpjs":                  "BPJS-001",
						"npwp":                        "NPWP-001",
						"tanggal_npwp":                "2016-01-01",
						"nomor_taspen":                "TASPEN-001",
						"unit_organisasi":             ["Unor E", "Unor G"]
					}
				}
			`,
		},
		{
			name:             "ok: pppk with terminated date later than today",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1f")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
				{
					"data": {
						"nama":                        "Budi John",
						"gelar_depan":                 "",
						"gelar_belakang":              "",
						"jabatan_aktual":              "Unit Kerja 8",
						"jenis_jabatan_aktual":        "Struktural",
						"tmt_jabatan":                 "2020-01-01",
						"nip":                         "1f",
						"nik":                         "3173000000000001",
						"nomor_kk":                    "3173000000000002",
						"jenis_kelamin":               "L",
						"tempat_lahir":                "Jakarta",
						"tanggal_lahir":               "1990-05-20",
						"tingkat_pendidikan":          "Tingkat 1",
						"pendidikan":                  "Pendidikan 1",
						"status_perkawinan":           "Menikah",
						"agama":                       "Kristen",
						"email_dikbud":                "budi@dikbud.go.id",
						"email_pribadi":               "budi@gmail.com",
						"alamat":                      "Jl. Merdeka No. 123",
						"nomor_hp":                    "08123456789",
						"nomor_kontak_darurat":        "08198765432",
						"jenis_pegawai":               "Jenis Pegawai 1",
						"masa_kerja_keseluruhan":      "` + masaKerjaKeseluruhan(time.Date(2015, 6, 1, 0, 0, 0, 0, time.Local), 0, 0) + `",
						"masa_kerja_golongan":         "5 Tahun 6 Bulan",
						"jabatan":                     "Jabatan 4",
						"jenis_jabatan":               "",
						"kelas_jabatan":               "",
						"lokasi_kerja":                "Jakarta",
						"golongan_ruang_awal":         "I",
						"golongan_ruang_akhir":        "II",
						"pangkat_akhir":               "Pangkat 2",
						"tmt_golongan":                "2018-01-01",
						"tmt_asn":                     "2015-06-01",
						"nomor_sk_asn":                "SK-CPNS-2015",
						"is_pppk":                     true,
						"status_asn":                  "P3K",
						"status_pns":                  "",
						"tmt_pns":                     null,
						"kartu_pegawai":               "KARPEG001",
						"nomor_surat_dokter":          "DOC-HEALTH-001",
						"tanggal_surat_dokter":        "2015-05-01",
						"nomor_surat_bebas_narkoba":   "BN-001",
						"tanggal_surat_bebas_narkoba": "2015-05-02",
						"nomor_catatan_polisi":        "SKCK-001",
						"tanggal_catatan_polisi":      "2015-05-03",
						"akte_kelahiran":              "AKTE-001",
						"nomor_bpjs":                  "BPJS-001",
						"npwp":                        "NPWP-001",
						"tanggal_npwp":                "2016-01-01",
						"nomor_taspen":                "TASPEN-001",
						"unit_organisasi":             ["Unor 8", "Unor 9", "Unor A", "Unor B", "Unor C", "Unor D"]
					}
				}
			`,
		},
		{
			name:             "ok: status_pns without tmt_pns and another case with edge case",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1g")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
				{
					"data": {
						"nama":                        "Budi Santoso",
						"gelar_depan":                 "Dr.",
						"gelar_belakang":              "M.Sc",
						"jabatan_aktual":              "Jabatan 6",
						"jenis_jabatan_aktual":        "",
						"tmt_jabatan":                 "2020-01-01",
						"nip":                         "1g",
						"nik":                         "3173000000000001",
						"nomor_kk":                    "3173000000000002",
						"jenis_kelamin":               "L",
						"tempat_lahir":                "Jakarta",
						"tanggal_lahir":               "1990-05-20",
						"tingkat_pendidikan":          "Tingkat 1",
						"pendidikan":                  "Pendidikan 1",
						"status_perkawinan":           "Menikah",
						"agama":                       "Kristen",
						"email_dikbud":                "budi@dikbud.go.id",
						"email_pribadi":               "budi@gmail.com",
						"alamat":                      "Jl. Merdeka No. 123",
						"nomor_hp":                    "08123456789",
						"nomor_kontak_darurat":        "08198765432",
						"jenis_pegawai":               "Jenis Pegawai 1",
						"masa_kerja_keseluruhan":      "",
						"masa_kerja_golongan":         "5 Tahun 6 Bulan",
						"jabatan":                     "Jabatan 1",
						"jenis_jabatan":               "Jenis Jabatan 1",
						"kelas_jabatan":               "Kelas 1",
						"lokasi_kerja":                "Medan",
						"golongan_ruang_awal":         "I/a",
						"golongan_ruang_akhir":        "I/b",
						"pangkat_akhir":               "Pangkat 2",
						"tmt_golongan":                "2018-01-01",
						"tmt_asn":                     "2015-06-01",
						"nomor_sk_asn":                "SK-CPNS-2015",
						"is_pppk":                     false,
						"status_asn":                  "POL",
						"status_pns":                  "PNS",
						"tmt_pns":                     null,
						"kartu_pegawai":               "KARPEG001",
						"nomor_surat_dokter":          "DOC-HEALTH-001",
						"tanggal_surat_dokter":        "2015-05-01",
						"nomor_surat_bebas_narkoba":   "BN-001",
						"tanggal_surat_bebas_narkoba": "2015-05-02",
						"nomor_catatan_polisi":        "SKCK-001",
						"tanggal_catatan_polisi":      "2015-05-03",
						"akte_kelahiran":              "AKTE-001",
						"nomor_bpjs":                  "BPJS-001",
						"npwp":                        "NPWP-001",
						"tanggal_npwp":                "2016-01-01",
						"nomor_taspen":                "TASPEN-001",
						"unit_organisasi":             ["Unor C", "Unor D"]
					}
				}
			`,
		},
		{
			name:             "error: data pegawai deleted",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2c")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: tidak ada data pegawai milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
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

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/data-pribadi", nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
