package pegawai

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
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
			RegisterRoutes(e, sqlc.New(pgxconn))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

// requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
