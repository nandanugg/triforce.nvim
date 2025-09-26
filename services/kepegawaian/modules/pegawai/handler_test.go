package pegawai

import (
	"context"
	"encoding/base64"
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
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_getDataPribadi(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan
			(id, no, kode_jabatan, nama_jabatan, deleted_at) values
			(1,  1,  'KJ1',        'Jabatan 1',  null),
			(2,  2,  'KJ2',        'Jabatan 2',  '2000-01-01');
		insert into ref_kedudukan_hukum
			(id, nama,  is_pppk, deleted_at) values
			(1,  'P3K', true,    null),
			(2,  'PNS', false,   null),
			(3,  'TNI', true,   '2000-01-01');
		insert into ref_golongan
			(id, nama,  nama_pangkat, gol_pppk, deleted_at) values
			(1,  'I/a', 'Pangkat 1',  'I',      null),
			(2,  'I/b', 'Pangkat 2',  'II',     '2000-01-01');
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, deleted_at) values
			('0', '1',         'Unor 0',  null),
			('1', '2',         'Unor 1',  null),
			('2', '3',         'Unor 2',  null),
			('3', '4',         'Unor 3',  null),
			('4', '5',         'Unor 4',  null),
			('5', '6',         'Unor 5',  null),
			('6', '7',         'Unor 6',  null),
			('7', '8',         'Unor 7',  null),
			('8', '9',         'Unor 8',  null),
			('9', 'A',         'Unor 9',  null),
			('A', 'B',         'Unor A',  null),
			('B', null,        'Unor B',  null),
			('C', 'D',         'Unor C',  null),
			('D', 'E',         '',        null),
			('E', 'F',         'Unor E',  null),
			('F', '6',         'Unor F',  '2000-01-01');
		insert into pegawai
			(pns_id, nip_lama, nip_baru, nama, gelar_depan, gelar_belakang, unor_id, jabatan_instansi_id, gol_id, kedudukan_hukum_id, deleted_at) values
			('aa>a', 'nip_l1', 'nip_b1', 'John Doe', 'Dr.', 'S.Kom', '0', 'KJ1', 1, 2, null),
			('1c', 'nip_l2', 'nip_b2', 'Bob', null, null, 'C', 'KJ1', 1, 1, null),
			('1d', 'nip_l3', 'nip_b3', 'Jane', '', '', 'F', 'KJ2', 2, 3, null),
			('1e', 'nip_l3', 'nip_b3', 'John Doe', '', '', '0', 'KJ1', 1, 1, '2000-01-01');
	`

	tests := []struct {
		name             string
		dbData           string
		paramPNSID       string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: success find record",
			paramPNSID:       base64.RawURLEncoding.EncodeToString([]byte(`aa>a`)),
			dbData:           dbData,
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"nip_lama":        "nip_l1",
					"nip_baru":        "nip_b1",
					"nama":            "John Doe",
					"gelar_depan":     "Dr.",
					"gelar_belakang":  "S.Kom",
					"golongan":        "I/a",
					"pangkat":         "Pangkat 1",
					"jabatan":         "Jabatan 1",
					"unit_organisasi": [ "Unor 0", "Unor 1", "Unor 2", "Unor 3", "Unor 4", "Unor 5", "Unor 6", "Unor 7", "Unor 8", "Unor 9" ]
				}
			}`,
		},
		{
			name:             "ok: success find record with authenticated user",
			paramPNSID:       base64.RawURLEncoding.EncodeToString([]byte(`1c`)),
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1d")}},
			dbData:           dbData,
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"nip_lama":        "nip_l2",
					"nip_baru":        "nip_b2",
					"nama":            "Bob",
					"gelar_depan":     "",
					"gelar_belakang":  "",
					"golongan":        "I",
					"pangkat":         "Pangkat 1",
					"jabatan":         "Jabatan 1",
					"unit_organisasi": [ "Unor C", "Unor E" ]
				}
			}`,
		},
		{
			name:             "ok: success find record with deleted reference",
			paramPNSID:       base64.RawURLEncoding.EncodeToString([]byte(`1d`)),
			dbData:           dbData,
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"nip_lama":        "nip_l3",
					"nip_baru":        "nip_b3",
					"nama":            "Jane",
					"gelar_depan":     "",
					"gelar_belakang":  "",
					"golongan":        "",
					"pangkat":         "",
					"jabatan":         "",
					"unit_organisasi": []
				}
			}`,
		},
		{
			name:             "error: base64 encoded with URLEncoding",
			paramPNSID:       base64.URLEncoding.EncodeToString([]byte(`aa>a`)), // YWE-YQ==
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: base64 encoded with StdEncoding",
			paramPNSID:       base64.StdEncoding.EncodeToString([]byte(`aa>a`)), // YWE+YQ==
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: base64 encoded with RawStdEncoding",
			paramPNSID:       base64.RawStdEncoding.EncodeToString([]byte(`aa>a`)), // YWE+YQ
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: invalid base64",
			paramPNSID:       "@abc",
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: invalid base64 utf8 value",
			paramPNSID:       "1c",
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: data pegawai deleted",
			paramPNSID:       base64.RawURLEncoding.EncodeToString([]byte(`1e`)),
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: tidak ada data pegawai milik user",
			paramPNSID:       base64.RawURLEncoding.EncodeToString([]byte(`2a`)),
			dbData:           dbData,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/pegawai/profil/"+tt.paramPNSID, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, nama_jabatan,    pemimpin_pns_id, deleted_at) values
			('unor-1', null, 'Paling Atas', 'Atasan 1', null, null),
			('unor-2', 'unor-1', 'Tengah', 'Atasan 2', null, null),
			('unor-3', 'unor-2', 'Bawah', 'Atasan 3', null, null),
			('unor-4', 'unor-1', 'Tengah deleted', 'Atasan 4', null, now()),
			('unor-5', 'unor-4', 'Bawah 2', 'Atasan 5', null, null);
		INSERT INTO ref_kedudukan_hukum (id, nama, deleted_at) VALUES 
			(1, 'Aktif', NULL),
			(2, 'Masa Persiapan Pensiun', NULL),
			(3, 'Aktif deleted', NOW());

		INSERT INTO ref_jabatan (id,no,kode_jabatan, nama_jabatan, deleted_at) VALUES
			(1,1,'JBT-001', 'Analis Kebijakan Madya', NULL),
			(2,2,'JBT-002', 'Kepala Subbagian Perencanaan', NULL),
			(3,3,'JBT-003', 'Staf Administrasi Umum', NOW());

		INSERT INTO ref_golongan (id, nama, deleted_at) VALUES
			(10, 'III/a', NULL),
			(11, 'III/b', NULL),
			(12, 'III/c', NOW());

		INSERT INTO pegawai (pns_id, nip_baru, nama, gelar_depan, gelar_belakang, gol_id, jabatan_instansi_id, unor_id, kedudukan_hukum_id, status_cpns_pns, deleted_at) VALUES
			(1001, '199001012022031001', 'Budi Santoso', 'Drs.', 'M.Pd.', 10, 'JBT-001', 'unor-1', 1, 'PNS', NULL),
			(1002, '198903152022041002', 'Siti Aminah', NULL, 'S.Sos.', 11, 'JBT-002', 'unor-2', 1,'CPNS', NULL),
			(1003, '198812312020121003', 'Andi Rahman', NULL, NULL, 10, 'JBT-001', 'unor-3', 2,'PNS', NULL),
			(1004, '199505052022051004', 'Lina Pratiwi', NULL, NULL, 10, 'JBT-001', 'unor-4', 3, 'PNS', NULL),
			(1005, '199707072022061005', 'Rizky Fauzan', NULL, NULL, 10, 'JBT-003', 'unor-5', 1, 'PNS', NULL),
			(1006, '199808082022071006', 'Sari Dewi', NULL, NULL, 12, 'JBT-001', 'unor-1', 1, 'PNS', NULL),
			(1007, '199709092022081007', 'Agung Herkules', NULL, NULL, 12, 'JBT-001', 'unor-1', 1, 'PNS', NOW());
	`
	tests := []struct {
		name             string
		dbData           string
		requestHeader    http.Header
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "M.Pd.",
						"gelar_depan": "Drs.",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Budi Santoso",
						"nip": "199001012022031001",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					},
					{
						"gelar_belakang": "S.Sos.",
						"gelar_depan": "",
						"golongan": "III/b",
						"jabatan": "Kepala Subbagian Perencanaan",
						"nama": "Siti Aminah",
						"nip": "198903152022041002",
						"status": "CPNS",
						"unit_kerja": "Tengah - Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Andi Rahman",
						"nip": "198812312020121003",
						"status": "MPP",
						"unit_kerja": "Bawah - Tengah - Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "",
						"nama": "Rizky Fauzan",
						"nip": "199707072022061005",
						"status": "PNS",
						"unit_kerja": "Bawah 2"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Sari Dewi",
						"nip": "199808082022071006",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 5
				}
			}`,
		},
		{
			name:          "ok with limit offset",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"limit":  []string{"2"},
				"offset": []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "S.Sos.",
						"gelar_depan": "",
						"golongan": "III/b",
						"jabatan": "Kepala Subbagian Perencanaan",
						"nama": "Siti Aminah",
						"nip": "198903152022041002",
						"status": "CPNS",
						"unit_kerja": "Tengah - Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Andi Rahman",
						"nip": "198812312020121003",
						"status": "MPP",
						"unit_kerja": "Bawah - Tengah - Paling Atas"
					}
				],
				"meta": {
					"limit": 2,
					"offset": 1,
					"total": 5
				}
			}`,
		},
		{
			name:          "ok with keyword nama",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"keyword": []string{"Siti"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "S.Sos.",
						"gelar_depan": "",
						"golongan": "III/b",
						"jabatan": "Kepala Subbagian Perencanaan",
						"nama": "Siti Aminah",
						"nip": "198903152022041002",
						"status": "CPNS",
						"unit_kerja": "Tengah - Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:          "ok with keyword nip",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"keyword": []string{"199"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "M.Pd.",
						"gelar_depan": "Drs.",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Budi Santoso",
						"nip": "199001012022031001",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "",
						"nama": "Rizky Fauzan",
						"nip": "199707072022061005",
						"status": "PNS",
						"unit_kerja": "Bawah 2"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Sari Dewi",
						"nip": "199808082022071006",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 3
				}
			}`,
		},
		{
			name:          "ok with unor",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"unor_id": []string{"unor-1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "M.Pd.",
						"gelar_depan": "Drs.",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Budi Santoso",
						"nip": "199001012022031001",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					},
					{
						"gelar_belakang": "S.Sos.",
						"gelar_depan": "",
						"golongan": "III/b",
						"jabatan": "Kepala Subbagian Perencanaan",
						"nama": "Siti Aminah",
						"nip": "198903152022041002",
						"status": "CPNS",
						"unit_kerja": "Tengah - Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Andi Rahman",
						"nip": "198812312020121003",
						"status": "MPP",
						"unit_kerja": "Bawah - Tengah - Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "",
						"nama": "Rizky Fauzan",
						"nip": "199707072022061005",
						"status": "PNS",
						"unit_kerja": "Bawah 2"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Sari Dewi",
						"nip": "199808082022071006",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					}
					
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 5
				}
			}`,
		},
		{
			name:          "ok with golonganID",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"golongan_id": []string{"11"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "S.Sos.",
						"gelar_depan": "",
						"golongan": "III/b",
						"jabatan": "Kepala Subbagian Perencanaan",
						"nama": "Siti Aminah",
						"nip": "198903152022041002",
						"status": "CPNS",
						"unit_kerja": "Tengah - Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:          "ok with status PNS",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"status": []string{"PNS"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "M.Pd.",
						"gelar_depan": "Drs.",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Budi Santoso",
						"nip": "199001012022031001",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "",
						"nama": "Rizky Fauzan",
						"nip": "199707072022061005",
						"status": "PNS",
						"unit_kerja": "Bawah 2"
					},
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Sari Dewi",
						"nip": "199808082022071006",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 3
				}
			}`,
		},
		{
			name:          "ok with status MPP",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"status": []string{"MPP"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "",
						"gelar_depan": "",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Andi Rahman",
						"nip": "198812312020121003",
						"status": "MPP",
						"unit_kerja": "Bawah - Tengah - Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:          "ok empty with random unor",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"unit_id": []string{"random-unor"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 0
				}
			}`,
		},
		{
			name:          "ok with all filter",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"status":      []string{"PNS"},
				"unit_id":     []string{"unor-1"},
				"golongan_id": []string{"10"},
				"jabatan_id":  []string{"JBT-001"},
				"keyword":     []string{"19"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"gelar_belakang": "M.Pd.",
						"gelar_depan": "Drs.",
						"golongan": "III/a",
						"jabatan": "Analis Kebijakan Madya",
						"nama": "Budi Santoso",
						"nip": "199001012022031001",
						"status": "PNS",
						"unit_kerja": "Paling Atas"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:          "ok empty data with all filter",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			requestQuery: url.Values{
				"status":      []string{"PNS"},
				"unit_id":     []string{"unor-1"},
				"golongan_id": []string{"10"},
				"jabatan_id":  []string{"JBT-010"},
				"keyword":     []string{"19"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 0
				}
			}`,
		},
		{
			name:             "error: user bukan admin",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: invalid token",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(pgxconn), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
