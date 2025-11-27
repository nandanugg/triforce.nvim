package usulanperubahandata

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_adminList(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan (kode_jabatan) values ('K1'), ('K2');
		insert into ref_golongan (id) values (1), (2);
		insert into ref_unit_kerja
			(id,   diatasan_id, nama_unor) values
			('U1', null,        'Unor 1'),
			('U2', 'U1',        'Unor 2');
		insert into pegawai
			(pns_id,  nip_baru, nama,       gelar_depan, gelar_belakang, gol_id, unor_id, jabatan_instansi_id, deleted_at) values
			('id_1c', '12c',    'John Doe', 'Dr.',       'S.Kom',        1,      'U1',    'K1',                null),
			('id_1d', '12d',    'John Doe', 'Dr.',       'S.Kom',        1,      'U1',    'K1',                '2000-01-01'),
			('id_1e', '12e',    'User 12e', '',          '',             2,      'U2',    'K2',                null),
			('id_1f', 'f12',    'Jane Doe', null,        null,           1,      'U1',    'K1',                null);
		insert into usulan_perubahan_data
			(id, nip,   jenis_data, perubahan_data, action, created_at,   deleted_at) values
			(1,  '12c', 'data-1',   '{}',           '',     '2000-01-04', null),
			(2,  '12d', 'data-2',   '{}',           '',     '2000-01-03', null),
			(3,  '12e', 'data-1',   '{}',           '',     '2000-01-01', null),
			(4,  'f12', 'data-3',   '{}',           '',     '2000-01-02', null),
			(5,  '12c', 'data-4',   '{}',           '',     '2000-01-05', '2000-01-01'),
			(6,  '12c', 'data-1',   '{}',           '',     '2000-01-05', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Verify)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	RegisterRoutes(e, db, repo, authMw)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		requestHeader    http.Header
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: success with default pagination and without filter",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         3,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12e",
							"nama":           "User 12e",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 2 - Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         4,
						"jenis_data": "data-3",
						"pegawai": {
							"nip":            "f12",
							"nama":           "Jane Doe",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         1,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         6,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 4 }
			}`,
		},
		{
			name:          "ok: success with pagination and without filter",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"limit":  []string{"2"},
				"offset": []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         4,
						"jenis_data": "data-3",
						"pegawai": {
							"nip":            "f12",
							"nama":           "Jane Doe",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         1,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 2, "offset": 1, "total": 4 }
			}`,
		},
		{
			name:          "ok: success with filter nama",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"nama": []string{"do"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         4,
						"jenis_data": "data-3",
						"pegawai": {
							"nip":            "f12",
							"nama":           "Jane Doe",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         1,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         6,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 3 }
			}`,
		},
		{
			name:          "ok: success with filter nip",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"nip": []string{"12"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         3,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12e",
							"nama":           "User 12e",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 2 - Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         1,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 4, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					},
					{
						"id":         6,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 3 }
			}`,
		},
		{
			name:          "ok: success with filter jenis_data",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"jenis_data": []string{"data-3"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         4,
						"jenis_data": "data-3",
						"pegawai": {
							"nip":            "f12",
							"nama":           "Jane Doe",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 1 }
			}`,
		},
		{
			name:          "ok: success with filter kode_jabatan",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"kode_jabatan": []string{"K2"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         3,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12e",
							"nama":           "User 12e",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 2 - Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 1 }
			}`,
		},
		{
			name:          "ok: success with filter unit_kerja",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"unit_kerja_id": []string{"U2"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         3,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12e",
							"nama":           "User 12e",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 2 - Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 1 }
			}`,
		},
		{
			name:          "ok: success with filter golongan_id",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"golongan_id": []string{"2"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         3,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12e",
							"nama":           "User 12e",
							"gelar_depan":    "",
							"gelar_belakang": "",
							"unit_kerja":     "Unor 2 - Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 1 }
			}`,
		},
		{
			name:          "ok: success with pagination and with all filter",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"limit":         []string{"1"},
				"offset":        []string{"1"},
				"nama":          []string{"doe"},
				"nip":           []string{"12"},
				"jenis_data":    []string{"data-1"},
				"kode_jabatan":  []string{"K1"},
				"unit_kerja_id": []string{"U1"},
				"golongan_id":   []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         6,
						"jenis_data": "data-1",
						"pegawai": {
							"nip":            "12c",
							"nama":           "John Doe",
							"gelar_depan":    "Dr.",
							"gelar_belakang": "S.Kom",
							"unit_kerja":     "Unor 1"
						},
						"created_at": "` + time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
					}
				],
				"meta": { "limit": 1, "offset": 1, "total": 2 }
			}`,
		},
		{
			name:          "ok: success with empty data",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"jenis_data": []string{"data-4"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": { "limit": 10, "offset": 0, "total": 0 }
			}`,
		},
		{
			name:             "error: invalid token",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/usulan-perubahan-data", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDetail(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan
			(kode_jabatan, nama_jabatan, deleted_at) values
			('K1',         'Jabatan 1',  null),
			('K2',         'Jabatan 2',  '2000-01-01');
		insert into ref_golongan
			(id, nama,    gol_pppk, deleted_at) values
			(1,  'Gol 1', 'PPPK 1', null),
			(2,  'Gol 2', 'PPPK 2', '2000-01-01');
		insert into ref_kedudukan_hukum
			(id, is_pppk, deleted_at) values
			(1,  false,   null),
			(2,  true,    '2000-01-01'),
			(3,  true,    null);
		insert into ref_unit_kerja
			(id, diatasan_id, nama_unor, deleted_at) values
			(1,  null,        'Unor 1',  null),
			(2,  null,        'Unor 2',  '2000-01-01'),
			(3,  1,           'Unor 3',  null),
			(4,  3,           'Unor 4',  null);
		insert into pegawai
			(pns_id,  nip_baru, nama,      gelar_depan, gelar_belakang, foto,   status_cpns_pns, unor_id, jabatan_instansi_id, gol_id, kedudukan_hukum_id, deleted_at) values
			('id_1c', '1c',     'User 1c', 'Dr.',       'S.Kom',        'f.1c', 'C',             '1',     'K1',                1,      1,                  null),
			('id_1d', '1d',     'User 1d', 'Dr.',       'S.Kom',        'f.1d', 'C',             '1',     'K1',                1,      1,                  '2000-01-01'),
			('id_1e', '1e',     'User 1e', '',          '',             '',     'P',             '2',     'K2',                2,      3,                  null),
			('id_1f', '1f',     'User 1f', 'Dr.',       'S.Kom',        'f.1f', 'PNS',           '4',     'K1',                1,      3,                  null),
			('id_1g', '1g',     'User 1g', '',          '',             'f.1g', 'CPNS',          '3',     'K1',                1,      2,                  null),
			('id_2c', '2c',     'User 2c', null,        null,           null,   null,            null,    null,                null,   null,               null);
		insert into usulan_perubahan_data
			(id, nip,  jenis_data, data_id, perubahan_data,                                              action,   status,      created_at,   deleted_at) values
			(1,  '1c', 'data-1',   null,    '{"val1":[null,2], "val2":[null,"2"], "val3":[null,false]}', 'CREATE', 'Disetujui', '2000-01-01', null),
			(2,  '1d', 'data-2',   '1',     '{"val1":[1,2],    "val2":["1","2"],  "val3":[false,true]}', 'UPDATE', 'Diusulkan', '2000-01-01', null),
			(3,  '1e', 'data-3',   '2',     '{"val1":[1,null], "val2":["1",null], "val3":[true,null]}',  'DELETE', 'Ditolak',   '2000-01-01', null),
			(4,  '1f', 'data-4',   '3',     '{"val1":[1,2],    "val2":["1","2"],  "val3":[false,true]}', 'UPDATE', 'Diusulkan', '2000-01-01', null),
			(5,  '1c', 'data-5',   null,    '{"val1":[null,2], "val2":[null,"2"], "val3":[null,false]}', 'CREATE', 'Diusulkan', '2000-01-01', '2000-01-01'),
			(6,  '1g', 'data-6',   '4',     '{"val1":[1,2],    "val2":["1","2"],  "val3":[false,true]}', 'UPDATE', 'Diusulkan', '2000-01-01', null),
			(7,  '2c', 'data-7',   null,    '{"val1":[null,2], "val2":[null,"2"], "val3":[null,false]}', 'CREATE', 'Diusulkan', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Verify)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	RegisterRoutes(e, db, repo, authMw)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramJenisData   string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: success get non pppk pegawai",
			paramJenisData:   "data-1",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id":         1,
					"jenis_data": "data-1",
					"action":     "CREATE",
					"status":     "Disetujui",
					"data_id":    null,
					"pegawai": {
						"nip":            "1c",
						"nama":           "User 1c",
						"gelar_depan":    "Dr.",
						"gelar_belakang": "S.Kom",
						"photo":          "f.1c",
						"status_pns":     "CPNS",
						"jabatan":        "Jabatan 1",
						"golongan":       "Gol 1",
						"unit_kerja":     "Unor 1"
					},
					"perubahan_data": {
						"val1": [ null, 2     ],
						"val2": [ null, "2"   ],
						"val3": [ null, false ]
					},
					"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
				}
			}`,
		},
		{
			name:             "error: pegawai is deleted",
			paramJenisData:   "data-2",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "ok: success get data with deleted references",
			paramJenisData:   "data-3",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id":         3,
					"jenis_data": "data-3",
					"action":     "DELETE",
					"status":     "Ditolak",
					"data_id":    "2",
					"pegawai": {
						"nip":            "1e",
						"nama":           "User 1e",
						"gelar_depan":    "",
						"gelar_belakang": "",
						"photo":          "",
						"status_pns":     "PNS",
						"jabatan":        "",
						"golongan":       "",
						"unit_kerja":     ""
					},
					"perubahan_data": {
						"val1": [ 1,    null ],
						"val2": [ "1",  null ],
						"val3": [ true, null ]
					},
					"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
				}
			}`,
		},
		{
			name:             "ok: success get pppk pegawai",
			paramJenisData:   "data-4",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id":         4,
					"jenis_data": "data-4",
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"data_id":    "3",
					"pegawai": {
						"nip":            "1f",
						"nama":           "User 1f",
						"gelar_depan":    "Dr.",
						"gelar_belakang": "S.Kom",
						"photo":          "f.1f",
						"status_pns":     "PNS",
						"jabatan":        "Jabatan 1",
						"golongan":       "PPPK 1",
						"unit_kerja":     "Unor 4 - Unor 3 - Unor 1"
					},
					"perubahan_data": {
						"val1": [ 1,     2    ],
						"val2": [ "1",   "2"  ],
						"val3": [ false, true ]
					},
					"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
				}
			}`,
		},
		{
			name:             "error: usulan perubahan data is deleted",
			paramJenisData:   "data-5",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "ok: success get record with deleted kedudukan hukum",
			paramJenisData:   "data-6",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id":         6,
					"jenis_data": "data-6",
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"data_id":    "4",
					"pegawai": {
						"nip":            "1g",
						"nama":           "User 1g",
						"gelar_depan":    "",
						"gelar_belakang": "",
						"photo":          "f.1g",
						"status_pns":     "CPNS",
						"jabatan":        "Jabatan 1",
						"golongan":       "Gol 1",
						"unit_kerja":     "Unor 3 - Unor 1"
					},
					"perubahan_data": {
						"val1": [ 1,     2    ],
						"val2": [ "1",   "2"  ],
						"val3": [ false, true ]
					},
					"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
				}
			}`,
		},
		{
			name:             "ok: success get record with null references",
			paramJenisData:   "data-7",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id":         7,
					"jenis_data": "data-7",
					"action":     "CREATE",
					"status":     "Diusulkan",
					"data_id":    null,
					"pegawai": {
						"nip":            "2c",
						"nama":           "User 2c",
						"gelar_depan":    "",
						"gelar_belakang": "",
						"photo":          null,
						"status_pns":     "",
						"jabatan":        "",
						"golongan":       "",
						"unit_kerja":     ""
					},
					"perubahan_data": {
						"val1": [ null, 2     ],
						"val2": [ null, "2"   ],
						"val3": [ null, false ]
					},
					"created_at": "` + time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `"
				}
			}`,
		},
		{
			name:             "error: usulan perubahan data is not found",
			paramJenisData:   "data-1",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: record jenis data is different",
			paramJenisData:   "data-2",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "data-1",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/usulan-perubahan-data/"+tt.paramJenisData+"/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_markAsRead(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_2c', '2c',     null);
		insert into usulan_perubahan_data
			(id, nip,  jenis_data,           data_id, perubahan_data,                       action,   status,      catatan, read_at,      created_at,   updated_at,   deleted_at) values
			(1,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Disetujui', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(2,  '1c', 'riwayat-pendidikan', '1',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Ditolak',   null,    null,         '2000-01-01', '2000-01-01', null),
			(3,  '1c', 'riwayat-pendidikan', '2',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Diusulkan', null,    '2000-01-01', '2000-01-01', '2000-01-01', null),
			(4,  '2c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Disetujui', 'notes', null,         '2000-01-01', '2000-01-01', null),
			(5,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Disetujui', null,    '2000-01-01', '2000-01-01', '2000-01-01', '2000-01-01'),
			(6,  '1d', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Disetujui', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(7,  '1c', 'riwayat-pendidikan', '1',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Disetujui', null,    null,         '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Request)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	RegisterRoutes(e, db, repo, authMw)

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
	tests := []struct {
		name             string
		paramJenisData   string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success mark as read record with status Disetujui",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(1),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    "notes",
					"read_at":    "{read_at}",
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "ok: success mark as read record with status Ditolak",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(2),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "1",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Ditolak",
					"catatan":    nil,
					"read_at":    "{read_at}",
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat ditandai sebagai telah dibaca"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(3),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "2",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data is owned by different user",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(4),
					"nip":        "2c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    "notes",
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data is not found",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: usulan perubahan data is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(5),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: pegawai is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(6),
					"nip":        "1d",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: record jenis data is different",
			paramJenisData:   "riwayat-pelatihan",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(7),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "1",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(7),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "1",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/usulan-perubahan-data/"+tt.paramJenisData+"/"+tt.paramID+"/read", nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["read_at"] == "{read_at}" {
						assert.WithinDuration(t, time.Now(), row["read_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["read_at"] = row["read_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_delete(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_2c', '2c',     null);
		insert into usulan_perubahan_data
			(id, nip,  jenis_data,           data_id, perubahan_data,                       action,   status,      catatan, read_at,      created_at,   updated_at,   deleted_at) values
			(1,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(2,  '1c', 'riwayat-pendidikan', '1',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(3,  '1c', 'riwayat-pendidikan', '2',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(4,  '2c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(5,  '1c', 'riwayat-pendidikan', '3',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Disetujui', null,    '2000-01-01', '2000-01-01', '2000-01-01', null),
			(6,  '1c', 'riwayat-pendidikan', '4',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Ditolak',   'notes', null,         '2000-01-01', '2000-01-01', null),
			(7,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', null,    '2000-01-01', '2000-01-01', '2000-01-01', '2000-01-01'),
			(8,  '1d', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(9,  '1c', 'riwayat-pendidikan', '5',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Request)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	RegisterRoutes(e, db, repo, authMw)

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
	tests := []struct {
		name             string
		paramJenisData   string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success delete for action CREATE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(1),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": "{deleted_at}",
				},
			},
		},
		{
			name:             "ok: success delete for action UPDATE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(2),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "1",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": "{deleted_at}",
				},
			},
		},
		{
			name:             "ok: success delete for action DELETE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(3),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "2",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": "{deleted_at}",
				},
			},
		},
		{
			name:             "error: usulan perubahan data is owned by different user",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(4),
					"nip":        "2c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Disetujui",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat dihapus"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(5),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "3",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Ditolak",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat dihapus"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(6),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "4",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Ditolak",
					"catatan":    "notes",
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data is not found",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: usulan perubahan data is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(7),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: pegawai is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1d",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: record jenis data is different",
			paramJenisData:   "riwayat-pelatihan",
			paramID:          "9",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(9),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "9",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(9),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/usulan-perubahan-data/"+tt.paramJenisData+"/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["deleted_at"] = row["deleted_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminReject(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_2c', '2c',     null);
		insert into usulan_perubahan_data
			(id, nip,  jenis_data,           data_id, perubahan_data,                       action,   status,      catatan, read_at,      created_at,   updated_at,   deleted_at) values
			(1,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(2,  '1c', 'riwayat-pendidikan', '1',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(3,  '1c', 'riwayat-pendidikan', '2',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(4,  '1c', 'riwayat-pendidikan', '3',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Disetujui', null,    '2000-01-01', '2000-01-01', '2000-01-01', null),
			(5,  '1c', 'riwayat-pendidikan', '4',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Ditolak',   'notes', null,         '2000-01-01', '2000-01-01', null),
			(6,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', null,    '2000-01-01', '2000-01-01', '2000-01-01', '2000-01-01'),
			(7,  '1d', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(8,  '1c', 'riwayat-pendidikan', '5',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Verify)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	RegisterRoutes(e, db, repo, authMw)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramJenisData   string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success reject for action CREATE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(1),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Ditolak",
					"catatan":    "new notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "ok: success delete for action UPDATE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(2),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "1",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Ditolak",
					"catatan":    "new notes",
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "ok: success delete for action DELETE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": ""}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(3),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "2",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Ditolak",
					"catatan":    "",
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Disetujui",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat ditolak"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(4),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "3",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Ditolak",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat ditolak"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(5),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "4",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Ditolak",
					"catatan":    "notes",
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data is not found",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: usulan perubahan data is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(6),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: pegawai is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(7),
					"nip":        "1d",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: record jenis data is different",
			paramJenisData:   "riwayat-pelatihan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: catatan exceed 200 character",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"catatan": "` + strings.Repeat(".", 201) + `"}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"catatan\" harus 200 karakter atau kurang"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: missing catatan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{"notes": "new notes"}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"notes\" tidak didukung | parameter \"catatan\" harus diisi"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: missing request body",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			requestBody:      `{"catatan": "new notes"}`,
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/usulan-perubahan-data/"+tt.paramJenisData+"/"+tt.paramID+"/reject", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_myList(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into usulan_perubahan_data
			(id, nip,   jenis_data,           data_id, perubahan_data,                    action,   status,      catatan, read_at,      updated_at,   deleted_at) values
			(1,  '12c', 'riwayat-pendidikan', null,    '{"pendidikan_id":[null,"2"]}',    'CREATE', 'Diusulkan', null,    null,         '2000-01-03', null),
			(2,  '12c', 'riwayat-pendidikan', '1',     '{"pendidikan_id":["1","2"]}',     'UPDATE', 'Disetujui', '',      null,         '2000-01-01', null),
			(3,  '12c', 'riwayat-pendidikan', '2',     '{"pendidikan_id":["1",null]}',    'DELETE', 'Ditolak',   'notes', null,         '2000-01-02', null),
			(4,  '12c', 'riwayat-pendidikan', null,    '{"pendidikan_id":[null,"1"]}',    'CREATE', 'Disetujui', null,    '2000-01-01', '2000-01-01', null),
			(5,  '12c', 'riwayat-pendidikan', null,    '{"pendidikan_id":[null,"1"]}',    'CREATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01'),
			(6,  '12d', 'riwayat-pendidikan', null,    '{"pendidikan_id":[null,"1"]}',    'CREATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01'),
			(7,  '12c', 'riwayat-pelatihan',  null,    '{"pendidikan_id":[null,"1"]}',    'CREATE', 'Disetujui', null,    null,         '2000-01-01', null),
			(8,  '12c', 'riwayat-pendidikan', '3',     '{"tingkat_pendidikan_id":[1,2]}', 'UPDATE', 'Disetujui', null,    null,         '2000-01-04', null),
			(9,  '12e', 'riwayat-pendidikan', '3',     '{"tingkat_pendidikan_id":[1,2]}', 'UPDATE', 'Disetujui', null,    null,         '2000-01-04', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Request)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	RegisterRoutes(e, db, repo, authMw)

	authHeader12c := []string{apitest.GenerateAuthHeader("12c")}
	tests := []struct {
		name             string
		paramJenisData   string
		requestHeader    http.Header
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: success get unread record with default pagination",
			paramJenisData:   "riwayat-pendidikan",
			requestHeader:    http.Header{"Authorization": authHeader12c},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         8,
						"jenis_data": "riwayat-pendidikan",
						"action":     "UPDATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "3",
						"perubahan_data": {
							"tingkat_pendidikan_id": [ 1, 2 ]
						}
					},
					{
						"id":         1,
						"jenis_data": "riwayat-pendidikan",
						"action":     "CREATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"pendidikan_id": [ null, "2" ]
						}
					},
					{
						"id":         3,
						"jenis_data": "riwayat-pendidikan",
						"action":     "DELETE",
						"status":     "Ditolak",
						"catatan":    "notes",
						"data_id":    "2",
						"perubahan_data": {
							"pendidikan_id": [ "1", null ]
						}
					},
					{
						"id":         2,
						"jenis_data": "riwayat-pendidikan",
						"action":     "UPDATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "1",
						"perubahan_data": {
							"pendidikan_id": [ "1", "2" ]
						}
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 4 }
			}`,
		},
		{
			name:           "ok: success get unread record with pagination",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": authHeader12c},
			requestQuery: url.Values{
				"limit":  []string{"2"},
				"offset": []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         1,
						"jenis_data": "riwayat-pendidikan",
						"action":     "CREATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"pendidikan_id": [ null, "2" ]
						}
					},
					{
						"id":         3,
						"jenis_data": "riwayat-pendidikan",
						"action":     "DELETE",
						"status":     "Ditolak",
						"catatan":    "notes",
						"data_id":    "2",
						"perubahan_data": {
							"pendidikan_id": [ "1", null ]
						}
					}
				],
				"meta": { "limit": 2, "offset": 1, "total": 4 }
			}`,
		},
		{
			name:             "ok: success get with another user",
			paramJenisData:   "riwayat-pendidikan",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("12e")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":         9,
						"jenis_data": "riwayat-pendidikan",
						"action":     "UPDATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "3",
						"perubahan_data": {
							"tingkat_pendidikan_id": [ 1, 2 ]
						}
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 1 }
			}`,
		},
		{
			name:             "ok: empty record",
			paramJenisData:   "riwayat-pendidikan",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("12d")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": { "limit": 10, "offset": 0, "total": 0 }
			}`,
		},
		{
			name:             "error: unregister route in openapi",
			paramJenisData:   "riwayat-pelatihan",
			requestHeader:    http.Header{"Authorization": authHeader12c},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "route tidak ditemukan pada openapi"}`,
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "riwayat-pendidikan",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/usulan-perubahan-data/"+tt.paramJenisData, nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			if !strings.Contains(rec.Body.String(), "route tidak ditemukan pada openapi") {
				assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			}
		})
	}
}

type testService struct {
	generate func() ([]byte, error)
	sync     func() error
}

func newTestService(generate func() ([]byte, error), sync func() error) *testService {
	return &testService{
		generate: generate,
		sync:     sync,
	}
}

func (s *testService) GeneratePerubahanData(context.Context, string, string, string, json.RawMessage) ([]byte, error) {
	return s.generate()
}

func (s *testService) SyncPerubahanData(context.Context, *sqlc.Queries, string, string, string, []byte) error {
	return s.sync()
}

func Test_handler_create(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into usulan_perubahan_data
			(nip,  jenis_data,           data_id, perubahan_data, action,   status,      created_at,   updated_at,   deleted_at) values
			('1c', 'riwayat-pendidikan', '1',     '{}',           'DELETE', 'Diusulkan', '2000-01-01', '2000-01-01', null),
			('1d', 'riwayat-pendidikan', '2',     '{}',           'UPDATE', 'Diusulkan', '2000-01-01', '2000-01-01', '2000-01-01'),
			('1d', 'riwayat-pendidikan', '2',     '{}',           'UPDATE', 'Disetujui', '2000-01-01', '2000-01-01', null),
			('1d', 'riwayat-pendidikan', '2',     '{}',           'DELETE', 'Ditolak',   '2000-01-01', '2000-01-01', null),
			('1d', 'riwayat-pelatihan',  '2',     '{}',           'UPDATE', 'Diusulkan', '2000-01-01', '2000-01-01', null),
			('1d', 'riwayat-pendidikan', '3',     '{}',           'UPDATE', 'Diusulkan', '2000-01-01', '2000-01-01', null),
			('1d', 'riwayat-pendidikan', null,    '{}',           'CREATE', 'Diusulkan', '2000-01-01', '2000-01-01', null),
			('1e', 'riwayat-pendidikan', null,    '{}',           'CREATE', 'Diusulkan', '2000-01-01', '2000-01-01', null),
			('1e', 'riwayat-pendidikan', '4',     '{}',           'DELETE', 'Diusulkan', '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Request)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)

	authHeader12c := []string{apitest.GenerateAuthHeader("12c")}
	tests := []struct {
		name             string
		paramJenisData   string
		requestHeader    http.Header
		requestBody      string
		generateFunc     func() ([]byte, error)
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:           "ok: success create record with data_id where user have no record",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1f")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "6"
			}`,
			generateFunc: func() ([]byte, error) {
				return json.Marshal(map[string][2]any{})
			},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             "{id}",
					"nip":            "1f",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "6",
					"perubahan_data": map[string]any{},
					"action":         "DELETE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     "{created_at}",
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name:           "ok: success create record without data_id where user have existings record",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1e")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"tingkat_pendidikan_id": 1,
					"nama_sekolah": "UI",
					"tahun_lulus": 2000,
					"nomor_ijazah": "IZ.1"
				}
			}`,
			generateFunc: func() ([]byte, error) {
				return json.Marshal(map[string][2]any{
					"nama_sekolah": {nil, nil},
				})
			},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(8),
					"nip":            "1e",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        nil,
					"perubahan_data": map[string]any{},
					"action":         "CREATE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":             int64(9),
					"nip":            "1e",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "4",
					"perubahan_data": map[string]any{},
					"action":         "DELETE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":         "{id}",
					"nip":        "1e",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"nama_sekolah": []any{nil, nil},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:           "ok: success create record with data_id where user have existings record",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "2",
				"data": {
					"tingkat_pendidikan_id": 1,
					"nama_sekolah": "UI",
					"tahun_lulus": 2000,
					"nomor_ijazah": "IZ.1"
				}
			}`,
			generateFunc: func() ([]byte, error) {
				return json.Marshal(map[string][2]any{
					"tingkat_pendidikan_id": {nil, 2},
					"nama_sekolah":          {"1", nil},
					"is_pns":                {false, true},
					"tahun_lulus":           {1999, 2000},
					"nomor_ijazah":          {"IZ.1", "IZ.2"},
				})
			},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(2),
					"nip":            "1d",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "2",
					"perubahan_data": map[string]any{},
					"action":         "UPDATE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":             int64(3),
					"nip":            "1d",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "2",
					"perubahan_data": map[string]any{},
					"action":         "UPDATE",
					"status":         "Disetujui",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":             int64(4),
					"nip":            "1d",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "2",
					"perubahan_data": map[string]any{},
					"action":         "DELETE",
					"status":         "Ditolak",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":             int64(5),
					"nip":            "1d",
					"jenis_data":     "riwayat-pelatihan",
					"data_id":        "2",
					"perubahan_data": map[string]any{},
					"action":         "UPDATE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":             int64(6),
					"nip":            "1d",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "3",
					"perubahan_data": map[string]any{},
					"action":         "UPDATE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":             int64(7),
					"nip":            "1d",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        nil,
					"perubahan_data": map[string]any{},
					"action":         "CREATE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":         "{id}",
					"nip":        "1d",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "2",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
						"nama_sekolah":          []any{"1", nil},
						"is_pns":                []any{false, true},
						"tahun_lulus":           []any{float64(1999), float64(2000)},
						"nomor_ijazah":          []any{"IZ.1", "IZ.2"},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:           "error: error create record having same data_id, jenis_data and status Diusulkan",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "1",
				"data": {
					"tingkat_pendidikan_id": 1,
					"nama_sekolah": "UI",
					"tahun_lulus": 2000,
					"nomor_ijazah": "IZ.1"
				}
			}`,
			generateFunc: func() ([]byte, error) {
				return json.Marshal(map[string][2]any{
					"val1": {nil, 1},
				})
			},
			wantResponseCode: http.StatusConflict,
			wantResponseBody: `{"message": "data dengan id ini sudah diusulkan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nip":            "1c",
					"jenis_data":     "riwayat-pendidikan",
					"data_id":        "1",
					"perubahan_data": map[string]any{},
					"action":         "DELETE",
					"status":         "Diusulkan",
					"catatan":        nil,
					"read_at":        nil,
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name:           "error: service interface return multiError",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": authHeader12c},
			requestBody: `{
				"action": "DELETE",
				"data_id": "0"
			}`,
			generateFunc: func() ([]byte, error) {
				return nil, api.NewMultiError([]error{
					errors.New("parameter tidak valid"),
					errors.New("request tidak valid"),
				})
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter tidak valid | request tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:           "error: service interface raise an error",
			paramJenisData: "riwayat-pendidikan",
			requestHeader:  http.Header{"Authorization": authHeader12c},
			requestBody: `{
				"action": "DELETE",
				"data_id": "0"
			}`,
			generateFunc: func() ([]byte, error) {
				return nil, errors.New("unexpected error")
			},
			wantResponseCode: http.StatusInternalServerError,
			wantResponseBody: `{"message": "Internal Server Error"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: unregister route in openapi",
			paramJenisData:   "riwayat-pelatihan",
			requestHeader:    http.Header{"Authorization": authHeader12c},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "route tidak ditemukan pada openapi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "riwayat-pendidikan",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/usulan-perubahan-data/"+tt.paramJenisData, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			svcRoute := RegisterRoutes(e, db, repo, authMw)

			svc := newTestService(tt.generateFunc, nil)
			svcRoute.Register(e, authMw, svc, tt.paramJenisData)

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			if !strings.Contains(rec.Body.String(), "route tidak ditemukan pada openapi") {
				assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			}

			nip := apitest.GetNIPFromAuthHeader(req.Header.Get("Authorization"))
			actualRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where nip = $1 order by id", nip)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["id"] = row["id"]
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminApprove(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_2c', '2c',     null);
		insert into usulan_perubahan_data
			(id, nip,  jenis_data,           data_id, perubahan_data,                       action,   status,      catatan, read_at,      created_at,   updated_at,   deleted_at) values
			(1,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(2,  '1c', 'riwayat-pendidikan', '1',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(3,  '1c', 'riwayat-pendidikan', '2',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(4,  '1c', 'riwayat-pendidikan', '3',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Disetujui', null,    '2000-01-01', '2000-01-01', '2000-01-01', null),
			(5,  '1c', 'riwayat-pendidikan', '4',     '{"tingkat_pendidikan_id":[1,null]}', 'DELETE', 'Ditolak',   'notes', null,         '2000-01-01', '2000-01-01', null),
			(6,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', null,    '2000-01-01', '2000-01-01', '2000-01-01', '2000-01-01'),
			(7,  '1d', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null),
			(8,  '1c', 'riwayat-pendidikan', '5',     '{"tingkat_pendidikan_id":[1,2]}',    'UPDATE', 'Diusulkan', null,    null,         '2000-01-01', '2000-01-01', null),
			(9,  '1c', 'riwayat-pendidikan', null,    '{"tingkat_pendidikan_id":[null,2]}', 'CREATE', 'Diusulkan', 'notes', '2000-01-01', '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Verify)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramJenisData   string
		paramID          string
		requestHeader    http.Header
		syncFunc         func() error
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success reject for action CREATE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			syncFunc:         func() error { return nil },
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(1),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "ok: success delete for action UPDATE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			syncFunc:         func() error { return nil },
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(2),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "1",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "ok: success delete for action DELETE with status Diusulkan",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			syncFunc:         func() error { return nil },
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(3),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "2",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: service interface raise an error",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "9",
			requestHeader:    http.Header{"Authorization": authHeader},
			syncFunc:         func() error { return errors.New("unexpected error") },
			wantResponseCode: http.StatusInternalServerError,
			wantResponseBody: `{"message": "Internal Server Error"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(9),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Disetujui",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat disetujui"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(4),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "3",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data status is Ditolak",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusUnprocessableEntity,
			wantResponseBody: `{"message": "status transisi tidak valid, data tidak dapat disetujui"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(5),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "4",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
					},
					"action":     "DELETE",
					"status":     "Ditolak",
					"catatan":    "notes",
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: usulan perubahan data is not found",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: usulan perubahan data is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(6),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: pegawai is deleted",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(7),
					"nip":        "1d",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(2)},
					},
					"action":     "CREATE",
					"status":     "Diusulkan",
					"catatan":    "notes",
					"read_at":    time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: record jenis data is different",
			paramJenisData:   "riwayat-pelatihan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: invalid token",
			paramJenisData:   "riwayat-pendidikan",
			paramID:          "8",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":         int64(8),
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), float64(2)},
					},
					"action":     "UPDATE",
					"status":     "Diusulkan",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/usulan-perubahan-data/"+tt.paramJenisData+"/"+tt.paramID+"/approve", nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			svcRoute := RegisterRoutes(e, db, repo, authMw)

			svc := newTestService(nil, tt.syncFunc)
			svcRoute.Register(e, authMw, svc, tt.paramJenisData)

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where id = $1", tt.paramID)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}
