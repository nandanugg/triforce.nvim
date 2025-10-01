package riwayatpendidikan

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan
			(id, nama,         deleted_at) values
			(6, 'Diploma III', null),
			(7, 'Sarjana',     null),
			(8, 'Magister',    null),
			(9, 'Deleted',     '2000-01-01');
		insert into ref_pendidikan
			(id,       nama,                    deleted_at) values
			('ed-003', 'Akuntansi',             null),
			('ed-004', 'Magister Manajemen',    null),
			('ed-006', 'Diploma III Akuntansi', null),
			('ed-007', 'Diploma Deleted',       '2000-01-01');
		insert into riwayat_pendidikan (id, nip, tingkat_pendidikan_id, pendidikan_id, nama_sekolah, tahun_lulus, no_ijazah, gelar_depan, gelar_belakang, tugas_belajar, negara_sekolah, deleted_at) values
			(1, '198812252013014004', 6, 'ed-006', 'Politeknik Negeri Jakarta', '2009', 'PNJ/AK/2009/004', null, 'A.Md.', 0, 'Pendidikan Regular', null),
			(2, '198812252013014004', 7, 'ed-003', 'Universitas Airlangga', '2011', 'UNAIR/AK/2011/004', 'Dr.', null, 2, 'Program Ekstensi', null),
			(3, '198812252013014004', 8, 'ed-004', 'Universitas Airlangga', '2016', 'UNAIR/MM/2016/004', 'Prof.', 'M.M.', 1, 'Beasiswa Institusi', null),
			(4, '198812252013014004', 8, 'ed-004', 'Universitas Airlangga', '2016', 'UNAIR/MM/2016/004', 'Prof.', 'M.M.', 3, 'Beasiswa Institusi', '2000-01-01'),
			(5, '198812252013014004', 9, 'ed-007', 'Universitas Pariwisata', '2001', null, null, null, null, null, null),
			(6, '19881225201', 8, 'ed-004', 'Universitas Airlangga', '2016', 'UNAIR/MM/2016/004', 'Prof.', 'M.M.', 3, 'Beasiswa Institusi', null);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198812252013014004")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    3,
						"jenjang_pendidikan":    "Magister",
						"jurusan":               "Magister Manajemen",
						"nama_sekolah":          "Universitas Airlangga",
						"tahun_lulus":           2016,
						"nomor_ijazah":          "UNAIR/MM/2016/004",
						"gelar_depan":           "Prof.",
						"gelar_belakang":        "M.M.",
						"tugas_belajar":         "Tugas Belajar",
						"keterangan_pendidikan": "Beasiswa Institusi"
					},
					{
						"id":                    2,
						"jenjang_pendidikan":    "Sarjana",
						"jurusan":               "Akuntansi",
						"nama_sekolah":          "Universitas Airlangga",
						"tahun_lulus":           2011,
						"nomor_ijazah":          "UNAIR/AK/2011/004",
						"gelar_depan":           "Dr.",
						"gelar_belakang":        "",
						"tugas_belajar":         "Izin Belajar",
						"keterangan_pendidikan": "Program Ekstensi"
					},
					{
						"id":                    1,
						"jenjang_pendidikan":    "Diploma III",
						"jurusan":               "Diploma III Akuntansi",
						"nama_sekolah":          "Politeknik Negeri Jakarta",
						"tahun_lulus":           2009,
						"nomor_ijazah":          "PNJ/AK/2009/004",
						"gelar_depan":           "",
						"gelar_belakang":        "A.Md.",
						"tugas_belajar":         "",
						"keterangan_pendidikan": "Pendidikan Regular"
					},
					{
						"id":                    5,
						"jenjang_pendidikan":    "",
						"jurusan":               "",
						"nama_sekolah":          "Universitas Pariwisata",
						"tahun_lulus":           2001,
						"nomor_ijazah":          "",
						"gelar_depan":           "",
						"gelar_belakang":        "",
						"tugas_belajar":         "",
						"keterangan_pendidikan": ""
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"2"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198812252013014004")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    1,
						"jenjang_pendidikan":    "Diploma III",
						"jurusan":               "Diploma III Akuntansi",
						"nama_sekolah":          "Politeknik Negeri Jakarta",
						"tahun_lulus":           2009,
						"nomor_ijazah":          "PNJ/AK/2009/004",
						"gelar_depan":           "",
						"gelar_belakang":        "A.Md.",
						"tugas_belajar":         "",
						"keterangan_pendidikan": "Pendidikan Regular"
					}
				],
				"meta": {"limit": 1, "offset": 2, "total": 4}
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
			dbRepository := sqlc.New(db)
			_, err := db.Exec(t.Context(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pendidikan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, dbRepository, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_getBerkas(t *testing.T) {
	t.Parallel()

	filePath := "../../../../lib/api/sample/hello.pdf"
	pdfBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0xf8, 0xff, 0xff, 0x3f,
		0x00, 0x05, 0xfe, 0x02, 0xfe, 0xa7, 0x46, 0x90,
		0x3d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
		0x44, 0xae, 0x42, 0x60, 0x82,
	}

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		insert into riwayat_pendidikan
			(id, nip, deleted_at,   file_base64) values
			(1, '1c', null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c', null,         '` + pdfBase64 + `'),
			(3, '1c', null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c', null,         'data:application/pdf;base64,invalid'),
			(5, '1c', '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c', null,         null),
			(7, '1c', null,         '');
		`

	tests := []struct {
		name              string
		dbData            string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			dbData:            dbData,
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 riwayat pendidikan tidak valid",
			dbData:            dbData,
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pendidikan sudah dihapus",
			dbData:            dbData,
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pendidikan berisi null value",
			dbData:            dbData,
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pendidikan berupa string kosong",
			dbData:            dbData,
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pendidikan bukan milik user login",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pendidikan tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pendidikan/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseCode == http.StatusOK {
				assert.Equal(t, "inline", rec.Header().Get("Content-Disposition"))
				assert.Equal(t, tt.wantContentType, rec.Header().Get("Content-Type"))
				assert.Equal(t, tt.wantResponseBytes, rec.Body.Bytes())
			}
		})
	}
}

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan
			(id, nama,         deleted_at) values
			(6, 'Diploma III', null),
			(7, 'Sarjana',     null),
			(8, 'Magister',    null),
			(9, 'Deleted',     '2000-01-01');
		insert into ref_pendidikan
			(id,       nama,                    deleted_at) values
			('ed-003', 'Akuntansi',             null),
			('ed-004', 'Magister Manajemen',    null),
			('ed-006', 'Diploma III Akuntansi', null),
			('ed-007', 'Diploma Deleted',       '2000-01-01');
		insert into riwayat_pendidikan (id, nip, tingkat_pendidikan_id, pendidikan_id, nama_sekolah, tahun_lulus, no_ijazah, gelar_depan, gelar_belakang, tugas_belajar, negara_sekolah, deleted_at) values
			(1, '1c', 6, 'ed-006', 'Politeknik Negeri Jakarta', '2009', 'PNJ/AK/2009/004', null, 'A.Md.', 0, 'Pendidikan Regular', null),
			(2, '1c', 7, 'ed-003', 'Universitas Airlangga', '2011', 'UNAIR/AK/2011/004', 'Dr.', null, 2, 'Program Ekstensi', null),
			(3, '1c', 8, 'ed-004', 'Universitas Airlangga', '2016', 'UNAIR/MM/2016/004', 'Prof.', 'M.M.', 1, 'Beasiswa Institusi', null),
			(4, '2c', 7, 'ed-003', 'Universitas Airlangga', '2011', 'UNAIR/AK/2011/004', 'Dr.', null, 2, 'Program Ekstensi', null),
			(5, '1c', 8, 'ed-004', 'Universitas Airlangga', '2016', 'UNAIR/MM/2016/004', 'Prof.', 'M.M.', 3, 'Beasiswa Institusi', '2000-01-01'),
			(6, '1d', 7, 'ed-003', 'Universitas Airlangga', '2011', 'UNAIR/AK/2011/004', 'Dr.', null, 2, 'Program Ekstensi', null);
	`

	tests := []struct {
		name             string
		dbData           string
		nip              string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: admin dapat melihat riwayat pendidikan pegawai 1c",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    3,
						"jenjang_pendidikan":    "Magister",
						"jurusan":               "Magister Manajemen",
						"nama_sekolah":          "Universitas Airlangga",
						"tahun_lulus":           2016,
						"nomor_ijazah":          "UNAIR/MM/2016/004",
						"gelar_depan":           "Prof.",
						"gelar_belakang":        "M.M.",
						"tugas_belajar":         "Tugas Belajar",
						"keterangan_pendidikan": "Beasiswa Institusi"
					},
					{
						"id":                    2,
						"jenjang_pendidikan":    "Sarjana",
						"jurusan":               "Akuntansi",
						"nama_sekolah":          "Universitas Airlangga",
						"tahun_lulus":           2011,
						"nomor_ijazah":          "UNAIR/AK/2011/004",
						"gelar_depan":           "Dr.",
						"gelar_belakang":        "",
						"tugas_belajar":         "Izin Belajar",
						"keterangan_pendidikan": "Program Ekstensi"
					},
					{
						"id":                    1,
						"jenjang_pendidikan":    "Diploma III",
						"jurusan":               "Diploma III Akuntansi",
						"nama_sekolah":          "Politeknik Negeri Jakarta",
						"tahun_lulus":           2009,
						"nomor_ijazah":          "PNJ/AK/2009/004",
						"gelar_depan":           "",
						"gelar_belakang":        "A.Md.",
						"tugas_belajar":         "",
						"keterangan_pendidikan": "Pendidikan Regular"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pendidikan pegawai 1c dengan pagination",
			dbData:           dbData,
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    2,
						"jenjang_pendidikan":    "Sarjana",
						"jurusan":               "Akuntansi",
						"nama_sekolah":          "Universitas Airlangga",
						"tahun_lulus":           2011,
						"nomor_ijazah":          "UNAIR/AK/2011/004",
						"gelar_depan":           "Dr.",
						"gelar_belakang":        "",
						"tugas_belajar":         "Izin Belajar",
						"keterangan_pendidikan": "Program Ekstensi"
					},
					{
						"id":                    1,
						"jenjang_pendidikan":    "Diploma III",
						"jurusan":               "Diploma III Akuntansi",
						"nama_sekolah":          "Politeknik Negeri Jakarta",
						"tahun_lulus":           2009,
						"nomor_ijazah":          "PNJ/AK/2009/004",
						"gelar_depan":           "",
						"gelar_belakang":        "A.Md.",
						"tugas_belajar":         "",
						"keterangan_pendidikan": "Pendidikan Regular"
					}
				],
				"meta": {"limit": 2, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pendidikan pegawai 1d",
			dbData:           dbData,
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    6,
						"jenjang_pendidikan":    "Sarjana",
						"jurusan":               "Akuntansi",
						"nama_sekolah":          "Universitas Airlangga",
						"tahun_lulus":           2011,
						"nomor_ijazah":          "UNAIR/AK/2011/004",
						"gelar_depan":           "Dr.",
						"gelar_belakang":        "",
						"tugas_belajar":         "Izin Belajar",
						"keterangan_pendidikan": "Program Ekstensi"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pendidikan pegawai yang tidak ada data",
			dbData:           dbData,
			nip:              "999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pendidikan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
