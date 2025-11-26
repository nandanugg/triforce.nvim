package riwayatpendidikan

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata/usulanperubahandatatest"
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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(t.Context(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	dbRepository := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, dbRepository, authMw)
	RegisterRoutes(e, dbRepository, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("198812252013014004")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: tanpa parameter apapun",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    3,
						"tingkat_pendidikan_id": 8,
						"jenjang_pendidikan":    "Magister",
						"pendidikan_id":         "ed-004",
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
						"tingkat_pendidikan_id": 7,
						"jenjang_pendidikan":    "Sarjana",
						"pendidikan_id":         "ed-003",
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
						"tingkat_pendidikan_id": 6,
						"jenjang_pendidikan":    "Diploma III",
						"pendidikan_id":         "ed-006",
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
						"tingkat_pendidikan_id": 9,
						"jenjang_pendidikan":    "",
						"pendidikan_id":         "ed-007",
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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"2"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    1,
						"tingkat_pendidikan_id": 6,
						"jenjang_pendidikan":    "Diploma III",
						"pendidikan_id":         "ed-006",
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pendidikan", nil)
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
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, pgxconn, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
	tests := []struct {
		name              string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 riwayat pendidikan tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pendidikan sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pendidikan berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pendidikan berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pendidikan bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pendidikan tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pendidikan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pendidikan/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		nip              string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: admin dapat melihat riwayat pendidikan pegawai 1c",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    3,
						"tingkat_pendidikan_id": 8,
						"jenjang_pendidikan":    "Magister",
						"pendidikan_id":         "ed-004",
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
						"tingkat_pendidikan_id": 7,
						"jenjang_pendidikan":    "Sarjana",
						"pendidikan_id":         "ed-003",
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
						"tingkat_pendidikan_id": 6,
						"jenjang_pendidikan":    "Diploma III",
						"pendidikan_id":         "ed-006",
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
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    2,
						"tingkat_pendidikan_id": 7,
						"jenjang_pendidikan":    "Sarjana",
						"pendidikan_id":         "ed-003",
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
						"tingkat_pendidikan_id": 6,
						"jenjang_pendidikan":    "Diploma III",
						"pendidikan_id":         "ed-006",
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
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                    6,
						"tingkat_pendidikan_id": 7,
						"jenjang_pendidikan":    "Sarjana",
						"pendidikan_id":         "ed-003",
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
			nip:              "999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
		},
		{
			name:             "error: auth header tidak valid",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pendidikan", nil)
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

func Test_handler_adminCreate(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1a', '1a',     null),
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null),
			('id_1f', '1f',     null),
			('id_1g', '1g',     null);
		insert into ref_pendidikan
			(id,  deleted_at) values
			('1', null),
			('2', '2000-01-01');
		insert into ref_tingkat_pendidikan
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": "1",
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2000,
				"nomor_ijazah": "UI.01",
				"gelar_depan": "Dr.",
				"gelar_belakang": "S.Kom",
				"tugas_belajar": "Tugas Belajar",
				"negara_sekolah": "Dalam Negeri"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         "1",
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2000),
					"gelar_depan":           "Dr.",
					"gelar_belakang":        "S.Kom",
					"tugas_belajar":         int16(1),
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        "Dalam Negeri",
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with different enum data",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": "1",
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2000,
				"nomor_ijazah": "UI.01",
				"gelar_depan": "Dr.",
				"gelar_belakang": "S.Kom",
				"tugas_belajar": "Izin Belajar",
				"negara_sekolah": "Luar Negeri"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         "1",
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2000),
					"gelar_depan":           "Dr.",
					"gelar_belakang":        "S.Kom",
					"tugas_belajar":         int16(2),
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        "Luar Negeri",
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1e",
					"nip":                   "1e",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": null,
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": "",
				"gelar_depan": "",
				"gelar_belakang": "",
				"tugas_belajar": "",
				"negara_sekolah": ""
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         nil,
					"no_ijazah":             "",
					"nama_sekolah":          "",
					"tahun_lulus":           int16(0),
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1f",
					"nip":                   "1f",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1g",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2020,
				"nomor_ijazah": "UI.01"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         nil,
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2020),
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1g",
					"nip":                   "1g",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": null,
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": ""
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": "1",
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": "",
				"gelar_depan": "",
				"gelar_belakang": "",
				"tugas_belajar": "",
				"negara_sekolah": ""
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: tingkat pendidikan or pendidikan is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 0,
				"pendidikan_id": "0",
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data tingkat pendidikan tidak ditemukan | data pendidikan tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: tingkat pendidikan or pendidikan is deleted",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 2,
				"pendidikan_id": "2",
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": "",
				"gelar_depan": "",
				"gelar_belakang": "",
				"tugas_belajar": "",
				"negara_sekolah": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data tingkat pendidikan tidak ditemukan | data pendidikan tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: exceed length limit, unexpected enum or data type",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": "s1",
				"pendidikan_id": 12,
				"nama_sekolah": "` + strings.Repeat(".", 201) + `",
				"tahun_lulus": "2001",
				"nomor_ijazah": "` + strings.Repeat(".", 101) + `",
				"gelar_depan": "` + strings.Repeat(".", 51) + `",
				"gelar_belakang": "` + strings.Repeat(".", 61) + `",
				"tugas_belajar": "Aksi Belajar",
				"negara_sekolah": "Indonesia"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"gelar_belakang\" harus 60 karakter atau kurang` +
				` | parameter \"gelar_depan\" harus 50 karakter atau kurang` +
				` | parameter \"nama_sekolah\" harus 200 karakter atau kurang` +
				` | parameter \"negara_sekolah\" harus salah satu dari \"Dalam Negeri\", \"Luar Negeri\", \"\"` +
				` | parameter \"nomor_ijazah\" harus 100 karakter atau kurang` +
				` | parameter \"pendidikan_id\" harus dalam tipe string` +
				` | parameter \"tahun_lulus\" harus dalam tipe integer` +
				` | parameter \"tingkat_pendidikan_id\" harus dalam tipe integer` +
				` | parameter \"tugas_belajar\" harus salah satu dari \"Tugas Belajar\", \"Izin Belajar\", \"\""}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null params",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": null,
				"pendidikan_id": null,
				"nama_sekolah": null,
				"tahun_lulus": null,
				"nomor_ijazah": null,
				"gelar_depan": null,
				"gelar_belakang": null,
				"tugas_belajar": null,
				"negara_sekolah": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"gelar_belakang\" tidak boleh null` +
				` | parameter \"gelar_depan\" tidak boleh null` +
				` | parameter \"nama_sekolah\" tidak boleh null` +
				` | parameter \"negara_sekolah\" harus salah satu dari \"Dalam Negeri\", \"Luar Negeri\", \"\"` +
				` | parameter \"nomor_ijazah\" tidak boleh null` +
				` | parameter \"tahun_lulus\" tidak boleh null` +
				` | parameter \"tingkat_pendidikan_id\" tidak boleh null` +
				` | parameter \"tugas_belajar\" harus salah satu dari \"Tugas Belajar\", \"Izin Belajar\", \"\""}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"tingkat_pendidikan_id\" harus diisi` +
				` | parameter \"nama_sekolah\" harus diisi` +
				` | parameter \"tahun_lulus\" harus diisi` +
				` | parameter \"nomor_ijazah\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pendidikan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_pendidikan", "where nip = $1", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["id"] = row["id"]
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]

						tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", fmt.Sprintf("%d", row["id"]))
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into ref_pendidikan
			(id,  deleted_at) values
			('1', null),
			('2', '2000-01-01');
		insert into ref_tingkat_pendidikan
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into riwayat_pendidikan
			(id, nama_sekolah, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
			(1,  'UI',         'id_1c', '1c', '2000-01-01', '2000-01-01', null),
			(2,  'UI',         'id_1c', '1c', '2000-01-01', '2000-01-01', null),
			(5,  'UI',         'id_1e', '1e', '2000-01-01', '2000-01-01', null),
			(6,  'UI',         'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01'),
			(7,  'UI',         'id_1c', '1c', '2000-01-01', '2000-01-01', null);
		insert into riwayat_pendidikan
			(id, pendidikan_pertama, pendidikan_terakhir, diakui_bkn, status_satker, status_biro, tanggal_lulus, file_base64, keterangan_berkas, pns_id_3, pendidikan_id_3, pns_id,  nip,  created_at,   updated_at) values
			(3,  '1',                1,                   1,          1,             1,           '2020-01-01',  'data:abc',  'abc',             '1a',     '2',             'id_1c', '1c', '2000-01-01', '2000-01-01'),
			(4,  '1',                1,                   1,          1,             1,           '2020-01-01',  'data:abc',  'abc',             '1a',     '2',             'id_1c', '1c', '2000-01-01', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                    int32(7),
			"tingkat_pendidikan_id": nil,
			"pendidikan_id":         nil,
			"no_ijazah":             nil,
			"nama_sekolah":          "UI",
			"tahun_lulus":           nil,
			"gelar_depan":           nil,
			"gelar_belakang":        nil,
			"tugas_belajar":         nil,
			"pendidikan_pertama":    nil,
			"pendidikan_terakhir":   nil,
			"negara_sekolah":        nil,
			"diakui_bkn":            nil,
			"status_satker":         nil,
			"status_biro":           nil,
			"tanggal_lulus":         nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"pns_id_3":              nil,
			"pendidikan_id_3":       nil,
			"pns_id":                "id_1c",
			"nip":                   "1c",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": "1",
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2000,
				"nomor_ijazah": "UI.01",
				"gelar_depan": "Dr.",
				"gelar_belakang": "S.Kom",
				"tugas_belajar": "Tugas Belajar",
				"negara_sekolah": "Dalam Negeri"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(1),
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         "1",
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2000),
					"gelar_depan":           "Dr.",
					"gelar_belakang":        "S.Kom",
					"tugas_belajar":         int16(1),
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        "Dalam Negeri",
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with different enum data",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": "1",
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2000,
				"nomor_ijazah": "UI.01",
				"gelar_depan": "Dr.",
				"gelar_belakang": "S.Kom",
				"tugas_belajar": "Izin Belajar",
				"negara_sekolah": "Luar Negeri"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(2),
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         "1",
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2000),
					"gelar_depan":           "Dr.",
					"gelar_belakang":        "S.Kom",
					"tugas_belajar":         int16(2),
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        "Luar Negeri",
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": null,
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": "",
				"gelar_depan": "",
				"gelar_belakang": "",
				"tugas_belajar": "",
				"negara_sekolah": ""
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(3),
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         nil,
					"no_ijazah":             "",
					"nama_sekolah":          "",
					"tahun_lulus":           int16(0),
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    "1",
					"pendidikan_terakhir":   int32(1),
					"negara_sekolah":        nil,
					"diakui_bkn":            int32(1),
					"status_satker":         int32(1),
					"status_biro":           int32(1),
					"tanggal_lulus":         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					"file_base64":           "data:abc",
					"s3_file_id":            nil,
					"keterangan_berkas":     "abc",
					"pns_id_3":              "1a",
					"pendidikan_id_3":       "2",
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2020,
				"nomor_ijazah": "UI.01"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(4),
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         nil,
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2020),
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    "1",
					"pendidikan_terakhir":   int32(1),
					"negara_sekolah":        nil,
					"diakui_bkn":            int32(1),
					"status_satker":         int32(1),
					"status_biro":           int32(1),
					"tanggal_lulus":         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					"file_base64":           "data:abc",
					"s3_file_id":            nil,
					"keterangan_berkas":     "abc",
					"pns_id_3":              "1a",
					"pendidikan_id_3":       "2",
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "error: riwayat pendidikan is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2020,
				"nomor_ijazah": "UI.01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat pendidikan is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2020,
				"nomor_ijazah": "UI.01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(5),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "UI",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1e",
					"nip":                   "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "error: riwayat pendidikan is deleted",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 1,
				"pendidikan_id": "1",
				"nama_sekolah": "Universitas Indonesia",
				"tahun_lulus": 2000,
				"nomor_ijazah": "UI.01",
				"gelar_depan": "Dr.",
				"gelar_belakang": "S.Kom",
				"tugas_belajar": "Tugas Belajar",
				"negara_sekolah": "Dalam Negeri"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(6),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "UI",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: tingkat pendidikan or pendidikan is not found",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 0,
				"pendidikan_id": "0",
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data tingkat pendidikan tidak ditemukan | data pendidikan tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: tingkat pendidikan or pendidikan is deleted",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": 2,
				"pendidikan_id": "2",
				"nama_sekolah": "",
				"tahun_lulus": 0,
				"nomor_ijazah": "",
				"gelar_depan": "",
				"gelar_belakang": "",
				"tugas_belajar": "",
				"negara_sekolah": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data tingkat pendidikan tidak ditemukan | data pendidikan tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": "s1",
				"pendidikan_id": 12,
				"nama_sekolah": "` + strings.Repeat(".", 201) + `",
				"tahun_lulus": "2001",
				"nomor_ijazah": "` + strings.Repeat(".", 101) + `",
				"gelar_depan": "` + strings.Repeat(".", 51) + `",
				"gelar_belakang": "` + strings.Repeat(".", 61) + `",
				"tugas_belajar": "Aksi Belajar",
				"negara_sekolah": "Indonesia"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"gelar_belakang\" harus 60 karakter atau kurang` +
				` | parameter \"gelar_depan\" harus 50 karakter atau kurang` +
				` | parameter \"nama_sekolah\" harus 200 karakter atau kurang` +
				` | parameter \"negara_sekolah\" harus salah satu dari \"Dalam Negeri\", \"Luar Negeri\", \"\"` +
				` | parameter \"nomor_ijazah\" harus 100 karakter atau kurang` +
				` | parameter \"pendidikan_id\" harus dalam tipe string` +
				` | parameter \"tahun_lulus\" harus dalam tipe integer` +
				` | parameter \"tingkat_pendidikan_id\" harus dalam tipe integer` +
				` | parameter \"tugas_belajar\" harus salah satu dari \"Tugas Belajar\", \"Izin Belajar\", \"\""}`,
			wantDBRows: defaultRows,
		},
		{
			name:          "error: null params",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"tingkat_pendidikan_id": null,
				"pendidikan_id": null,
				"nama_sekolah": null,
				"tahun_lulus": null,
				"nomor_ijazah": null,
				"gelar_depan": null,
				"gelar_belakang": null,
				"tugas_belajar": null,
				"negara_sekolah": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"gelar_belakang\" tidak boleh null` +
				` | parameter \"gelar_depan\" tidak boleh null` +
				` | parameter \"nama_sekolah\" tidak boleh null` +
				` | parameter \"negara_sekolah\" harus salah satu dari \"Dalam Negeri\", \"Luar Negeri\", \"\"` +
				` | parameter \"nomor_ijazah\" tidak boleh null` +
				` | parameter \"tahun_lulus\" tidak boleh null` +
				` | parameter \"tingkat_pendidikan_id\" tidak boleh null` +
				` | parameter \"tugas_belajar\" harus salah satu dari \"Tugas Belajar\", \"Izin Belajar\", \"\""}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"tingkat_pendidikan_id\" harus diisi` +
				` | parameter \"nama_sekolah\" harus diisi` +
				` | parameter \"tahun_lulus\" harus diisi` +
				` | parameter \"nomor_ijazah\" harus diisi"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pendidikan/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_pendidikan", "where id = $1", tt.paramID)
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

func Test_handler_adminDelete(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_pendidikan
			(id, nama_sekolah,     pns_id,  nip,  created_at,   updated_at,   deleted_at) values
			(1,  'Universitas',    'id_1c', '1c', '2000-01-01', '2000-01-01', null),
			(2,  null,             'id_1e', '1e', '2000-01-01', '2000-01-01', null),
			(3,  null,             'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01'),
			(4,  'UI',             'id_1c', '1c', '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                    int32(4),
			"tingkat_pendidikan_id": nil,
			"pendidikan_id":         nil,
			"no_ijazah":             nil,
			"nama_sekolah":          "UI",
			"tahun_lulus":           nil,
			"gelar_depan":           nil,
			"gelar_belakang":        nil,
			"tugas_belajar":         nil,
			"pendidikan_pertama":    nil,
			"pendidikan_terakhir":   nil,
			"negara_sekolah":        nil,
			"diakui_bkn":            nil,
			"status_satker":         nil,
			"status_biro":           nil,
			"tanggal_lulus":         nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"pns_id_3":              nil,
			"pendidikan_id_3":       nil,
			"pns_id":                "id_1c",
			"nip":                   "1c",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success delete",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(1),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "Universitas",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat pendidikan is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(2),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          nil,
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1e",
					"nip":                   "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
		},
		{
			name:             "error: riwayat pendidikan is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat pendidikan is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(3),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          nil,
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pendidikan/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_pendidikan", "where id = $1", tt.paramID)
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

func Test_handler_adminUploadBerkas(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into riwayat_pendidikan
			(id, nama_sekolah, pendidikan_pertama, pendidikan_terakhir, diakui_bkn, status_satker, status_biro, tanggal_lulus, file_base64, keterangan_berkas, pns_id_3, pendidikan_id_3, pns_id,  nip,  created_at,   updated_at) values
			(1,  'UI',         '1',                1,                   1,          1,             1,           '2020-01-01',  'data:abc',  'abc',             '1a',     '2',             'id_1c', '1c', '2000-01-01', '2000-01-01');
		insert into riwayat_pendidikan
			(id, nama_sekolah, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
			(2,  'UI',         'id_1e', '1e', '2000-01-01', '2000-01-01', null),
			(3,  'UI',         'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01'),
			(4,  'UI',         'id_1c', '1c', '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                    int32(4),
			"tingkat_pendidikan_id": nil,
			"pendidikan_id":         nil,
			"no_ijazah":             nil,
			"nama_sekolah":          "UI",
			"tahun_lulus":           nil,
			"gelar_depan":           nil,
			"gelar_belakang":        nil,
			"tugas_belajar":         nil,
			"pendidikan_pertama":    nil,
			"pendidikan_terakhir":   nil,
			"negara_sekolah":        nil,
			"diakui_bkn":            nil,
			"status_satker":         nil,
			"status_biro":           nil,
			"tanggal_lulus":         nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"pns_id_3":              nil,
			"pendidikan_id_3":       nil,
			"pns_id":                "id_1c",
			"nip":                   "1c",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	defaultRequestBody := func(writer *multipart.Writer) error {
		part, err := writer.CreateFormFile("file", "file.txt")
		if err != nil {
			return err
		}
		_, err = io.WriteString(part, "Hello World!!")
		return err
	}

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name              string
		paramNIP          string
		paramID           string
		requestHeader     http.Header
		appendRequestBody func(writer *multipart.Writer) error
		wantResponseCode  int
		wantResponseBody  string
		wantDBRows        dbtest.Rows
	}{
		{
			name:              "ok: success upload",
			paramNIP:          "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(1),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "UI",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    "1",
					"pendidikan_terakhir":   int32(1),
					"negara_sekolah":        nil,
					"diakui_bkn":            int32(1),
					"status_satker":         int32(1),
					"status_biro":           int32(1),
					"tanggal_lulus":         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					"file_base64":           "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"s3_file_id":            nil,
					"keterangan_berkas":     "abc",
					"pns_id_3":              "1a",
					"pendidikan_id_3":       "2",
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:              "error: riwayat pendidikan is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat pendidikan is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(2),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "UI",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1e",
					"nip":                   "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
		},
		{
			name:              "error: riwayat pendidikan is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int32(3),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "UI",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:              "error: missing file",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"file\" harus diisi"}`,
			wantDBRows:        defaultRows,
		},
		{
			name:              "error: invalid token",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBody:  `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:        defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			require.NoError(t, tt.appendRequestBody(writer))
			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pendidikan/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_pendidikan", "where id = $1", tt.paramID)
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

func Test_handler_usulanPerubahanData(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(pns_id,  nip_baru) values
			('id_1a', '1a'),
			('id_1b', '1b'),
			('id_1c', '1c'),
			('id_1d', '1d'),
			('id_1e', '1e'),
			('id_1f', '1f');
		insert into ref_pendidikan
			(id,  nama,   deleted_at) values
			('1', 'PD 1', null),
			('2', 'PD 2', '2000-01-01');
		insert into ref_tingkat_pendidikan
			(id, nama,   deleted_at) values
			(1,  'TP 1', null),
			(2,  'TP 2', '2000-01-01');
		insert into riwayat_pendidikan
			(nama_sekolah, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
			('UI',         'id_1a', '1a', '2000-01-01', '2000-01-01', null),
			('UI',         'id_1a', '1a', '2000-01-01', '2000-01-01', '2000-01-01'),
			('UI',         'id_1d', '1d', '2000-01-01', '2000-01-01', null);
		insert into riwayat_pendidikan
			(tingkat_pendidikan_id, pendidikan_id, no_ijazah, nama_sekolah, tahun_lulus, gelar_depan, gelar_belakang, tugas_belajar, negara_sekolah, pendidikan_pertama, pendidikan_terakhir, diakui_bkn, status_satker, status_biro, tanggal_lulus, file_base64, keterangan_berkas, pns_id_3, pendidikan_id_3, pns_id,  nip,  created_at,   updated_at) values
			(2,                     2,             'IZ1',     'ITB',        2000,        'Dr.',       'S.Kom',        1,             'Dalam Negeri', '1',                1,                   1,          1,             1,           '2020-01-01',  'data:abc',  'abc',             '1a',     '2',             'id_1d', '1d', '2000-01-01', '2000-01-01'),
			(1,                     1,             'IZ1',     'ITB',        2000,        'Dr.',       'S.Kom',        2,             'Luar Negeri',  '1',                1,                   1,          1,             1,           '2020-01-01',  'data:abc',  'abc',             '1a',     '2',             'id_1e', '1e', '2000-01-01', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	dbRows1a := dbtest.Rows{
		{
			"id":                    int32(1),
			"tingkat_pendidikan_id": nil,
			"pendidikan_id":         nil,
			"no_ijazah":             nil,
			"nama_sekolah":          "UI",
			"tahun_lulus":           nil,
			"gelar_depan":           nil,
			"gelar_belakang":        nil,
			"tugas_belajar":         nil,
			"pendidikan_pertama":    nil,
			"pendidikan_terakhir":   nil,
			"negara_sekolah":        nil,
			"diakui_bkn":            nil,
			"status_satker":         nil,
			"status_biro":           nil,
			"tanggal_lulus":         nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"pns_id_3":              nil,
			"pendidikan_id_3":       nil,
			"pns_id":                "id_1a",
			"nip":                   "1a",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
		{
			"id":                    int32(2),
			"tingkat_pendidikan_id": nil,
			"pendidikan_id":         nil,
			"no_ijazah":             nil,
			"nama_sekolah":          "UI",
			"tahun_lulus":           nil,
			"gelar_depan":           nil,
			"gelar_belakang":        nil,
			"tugas_belajar":         nil,
			"pendidikan_pertama":    nil,
			"pendidikan_terakhir":   nil,
			"negara_sekolah":        nil,
			"diakui_bkn":            nil,
			"status_satker":         nil,
			"status_biro":           nil,
			"tanggal_lulus":         nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"pns_id_3":              nil,
			"pendidikan_id_3":       nil,
			"pns_id":                "id_1a",
			"nip":                   "1a",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
		},
	}

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_PegawaiPerubahanData_Request)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)
	svcRoute := usulanperubahandatatest.NewServiceRoute(db)

	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader1a := []string{apitest.GenerateAuthHeader("1a")}
	tests := []struct {
		name                 string
		requestHeader        http.Header
		requestBody          string
		doRollback           bool
		wantResponsePostCode int
		wantResponsePostBody string
		wantResponseGetBody  string
		wantDBSvcRows        dbtest.Rows
		wantDBUsulanRows     dbtest.Rows
	}{
		{
			name:          "ok: success create riwayat pendidikan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"tingkat_pendidikan_id": 1,
					"pendidikan_id": "1",
					"nama_sekolah": "Universitas Indonesia",
					"tahun_lulus": 2000,
					"nomor_ijazah": "UI.01",
					"gelar_depan": "Dr.",
					"gelar_belakang": "S.Kom",
					"tugas_belajar": "Tugas Belajar",
					"negara_sekolah": "Dalam Negeri"
				}
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-pendidikan",
						"action":     "CREATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"tingkat_pendidikan_id": [ null, 1                       ],
							"jenjang_pendidikan":    [ null, "TP 1"                  ],
							"pendidikan_id":         [ null, "1"                     ],
							"jurusan":               [ null, "PD 1"                  ],
							"nama_sekolah":          [ null, "Universitas Indonesia" ],
							"tahun_lulus":           [ null, 2000                    ],
							"nomor_ijazah":          [ null, "UI.01"                 ],
							"gelar_depan":           [ null, "Dr."                   ],
							"gelar_belakang":        [ null, "S.Kom"                 ],
							"tugas_belajar":         [ null, "Tugas Belajar"         ],
							"negara_sekolah":        [ null, "Dalam Negeri"          ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         "1",
					"no_ijazah":             "UI.01",
					"nama_sekolah":          "Universitas Indonesia",
					"tahun_lulus":           int16(2000),
					"gelar_depan":           "Dr.",
					"gelar_belakang":        "S.Kom",
					"tugas_belajar":         int16(1),
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        "Dalam Negeri",
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1c",
					"nip":                   "1c",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1c",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(1)},
						"jenjang_pendidikan":    []any{nil, "TP 1"},
						"pendidikan_id":         []any{nil, "1"},
						"jurusan":               []any{nil, "PD 1"},
						"nama_sekolah":          []any{nil, "Universitas Indonesia"},
						"tahun_lulus":           []any{nil, float64(2000)},
						"nomor_ijazah":          []any{nil, "UI.01"},
						"gelar_depan":           []any{nil, "Dr."},
						"gelar_belakang":        []any{nil, "S.Kom"},
						"tugas_belajar":         []any{nil, "Tugas Belajar"},
						"negara_sekolah":        []any{nil, "Dalam Negeri"},
					},
					"action":     "CREATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: success update riwayat pendidikan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "4",
				"data": {
					"tingkat_pendidikan_id": 1,
					"pendidikan_id": null,
					"nama_sekolah": "",
					"tahun_lulus": 0,
					"nomor_ijazah": "",
					"gelar_depan": "",
					"gelar_belakang": "",
					"tugas_belajar": "",
					"negara_sekolah": ""
				}
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-pendidikan",
						"action":     "UPDATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "4",
						"perubahan_data": {
							"tingkat_pendidikan_id": [ 2,               1      ],
							"jenjang_pendidikan":    [ null,            "TP 1" ],
							"pendidikan_id":         [ "2",             null   ],
							"jurusan":               [ null,            null   ],
							"nama_sekolah":          [ "ITB",           ""     ],
							"tahun_lulus":           [ 2000,            0      ],
							"nomor_ijazah":          [ "IZ1",           ""     ],
							"gelar_depan":           [ "Dr.",           null   ],
							"gelar_belakang":        [ "S.Kom",         null   ],
							"tugas_belajar":         [ "Tugas Belajar", null   ],
							"negara_sekolah":        [ "Dalam Negeri",  null   ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    int32(3),
					"tingkat_pendidikan_id": nil,
					"pendidikan_id":         nil,
					"no_ijazah":             nil,
					"nama_sekolah":          "UI",
					"tahun_lulus":           nil,
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    nil,
					"pendidikan_terakhir":   nil,
					"negara_sekolah":        nil,
					"diakui_bkn":            nil,
					"status_satker":         nil,
					"status_biro":           nil,
					"tanggal_lulus":         nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"pns_id_3":              nil,
					"pendidikan_id_3":       nil,
					"pns_id":                "id_1d",
					"nip":                   "1d",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
				{
					"id":                    int32(4),
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         nil,
					"no_ijazah":             "",
					"nama_sekolah":          "",
					"tahun_lulus":           int16(0),
					"gelar_depan":           nil,
					"gelar_belakang":        nil,
					"tugas_belajar":         nil,
					"pendidikan_pertama":    "1",
					"pendidikan_terakhir":   int32(1),
					"negara_sekolah":        nil,
					"diakui_bkn":            int32(1),
					"status_satker":         int32(1),
					"status_biro":           int32(1),
					"tanggal_lulus":         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					"file_base64":           "data:abc",
					"s3_file_id":            nil,
					"keterangan_berkas":     "abc",
					"pns_id_3":              "1a",
					"pendidikan_id_3":       "2",
					"pns_id":                "id_1d",
					"nip":                   "1d",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1d",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "4",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(2), float64(1)},
						"jenjang_pendidikan":    []any{nil, "TP 1"},
						"pendidikan_id":         []any{"2", nil},
						"jurusan":               []any{nil, nil},
						"nama_sekolah":          []any{"ITB", ""},
						"tahun_lulus":           []any{float64(2000), float64(0)},
						"nomor_ijazah":          []any{"IZ1", ""},
						"gelar_depan":           []any{"Dr.", nil},
						"gelar_belakang":        []any{"S.Kom", nil},
						"tugas_belajar":         []any{"Tugas Belajar", nil},
						"negara_sekolah":        []any{"Dalam Negeri", nil},
					},
					"action":     "UPDATE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: success delete riwayat pendidikan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1e")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "5"
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-pendidikan",
						"action":     "DELETE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "5",
						"perubahan_data": {
							"tingkat_pendidikan_id": [ 1,              null ],
							"jenjang_pendidikan":    [ "TP 1",         null ],
							"pendidikan_id":         [ "1",            null ],
							"jurusan":               [ "PD 1",         null ],
							"nama_sekolah":          [ "ITB",          null ],
							"tahun_lulus":           [ 2000,           null ],
							"nomor_ijazah":          [ "IZ1",          null ],
							"gelar_depan":           [ "Dr.",          null ],
							"gelar_belakang":        [ "S.Kom",        null ],
							"tugas_belajar":         [ "Izin Belajar", null ],
							"negara_sekolah":        [ "Luar Negeri",  null ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                    int32(5),
					"tingkat_pendidikan_id": int16(1),
					"pendidikan_id":         "1",
					"no_ijazah":             "IZ1",
					"nama_sekolah":          "ITB",
					"tahun_lulus":           int16(2000),
					"gelar_depan":           "Dr.",
					"gelar_belakang":        "S.Kom",
					"tugas_belajar":         int16(2),
					"pendidikan_pertama":    "1",
					"pendidikan_terakhir":   int32(1),
					"negara_sekolah":        "Luar Negeri",
					"diakui_bkn":            int32(1),
					"status_satker":         int32(1),
					"status_biro":           int32(1),
					"tanggal_lulus":         time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					"file_base64":           "data:abc",
					"s3_file_id":            nil,
					"keterangan_berkas":     "abc",
					"pns_id_3":              "1a",
					"pendidikan_id_3":       "2",
					"pns_id":                "id_1e",
					"nip":                   "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            "{deleted_at}",
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1e",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    "5",
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{float64(1), nil},
						"jenjang_pendidikan":    []any{"TP 1", nil},
						"pendidikan_id":         []any{"1", nil},
						"jurusan":               []any{"PD 1", nil},
						"nama_sekolah":          []any{"ITB", nil},
						"tahun_lulus":           []any{float64(2000), nil},
						"nomor_ijazah":          []any{"IZ1", nil},
						"gelar_depan":           []any{"Dr.", nil},
						"gelar_belakang":        []any{"S.Kom", nil},
						"tugas_belajar":         []any{"Izin Belajar", nil},
						"negara_sekolah":        []any{"Luar Negeri", nil},
					},
					"action":     "DELETE",
					"status":     "Disetujui",
					"catatan":    nil,
					"read_at":    nil,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name:          "ok: rollback on usulan perubahan data should not sync perubahan data",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1f")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"tingkat_pendidikan_id": 1,
					"nama_sekolah": "UI",
					"tahun_lulus": 2020,
					"nomor_ijazah": "UI.01"
				}
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-pendidikan",
						"action":     "CREATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"tingkat_pendidikan_id": [ null, 1       ],
							"jenjang_pendidikan":    [ null, "TP 1"  ],
							"pendidikan_id":         [ null, null    ],
							"jurusan":               [ null, null    ],
							"nama_sekolah":          [ null, "UI"    ],
							"tahun_lulus":           [ null, 2020    ],
							"nomor_ijazah":          [ null, "UI.01" ],
							"gelar_depan":           [ null, null    ],
							"gelar_belakang":        [ null, null    ],
							"tugas_belajar":         [ null, null    ],
							"negara_sekolah":        [ null, null    ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1f",
					"jenis_data": "riwayat-pendidikan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"tingkat_pendidikan_id": []any{nil, float64(1)},
						"jenjang_pendidikan":    []any{nil, "TP 1"},
						"pendidikan_id":         []any{nil, nil},
						"jurusan":               []any{nil, nil},
						"nama_sekolah":          []any{nil, "UI"},
						"tahun_lulus":           []any{nil, float64(2020)},
						"nomor_ijazah":          []any{nil, "UI.01"},
						"gelar_depan":           []any{nil, nil},
						"gelar_belakang":        []any{nil, nil},
						"tugas_belajar":         []any{nil, nil},
						"negara_sekolah":        []any{nil, nil},
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
			name:          "error: riwayat pendidikan is owned by other pegawai",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1b")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "1"
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat pendidikan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbtest.Rows{},
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: riwayat pendidikan is not found",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "DELETE",
				"data_id": "0"
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat pendidikan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: riwayat pendidikan is deleted",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "2",
				"data": {
					"tingkat_pendidikan_id": 1,
					"pendidikan_id": "1",
					"nama_sekolah": "Universitas Indonesia",
					"tahun_lulus": 2000,
					"nomor_ijazah": "UI.01",
					"gelar_depan": "Dr.",
					"gelar_belakang": "S.Kom",
					"tugas_belajar": "Tugas Belajar",
					"negara_sekolah": "Dalam Negeri"
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat pendidikan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: tingkat pendidikan or pendidikan is not found",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"tingkat_pendidikan_id": 0,
					"pendidikan_id": "0",
					"nama_sekolah": "",
					"tahun_lulus": 0,
					"nomor_ijazah": ""
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data tingkat pendidikan tidak ditemukan | data pendidikan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: tingkat pendidikan or pendidikan is deleted",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"tingkat_pendidikan_id": 2,
					"pendidikan_id": "2",
					"nama_sekolah": "",
					"tahun_lulus": 0,
					"nomor_ijazah": "",
					"gelar_depan": "",
					"gelar_belakang": "",
					"tugas_belajar": "",
					"negara_sekolah": ""
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data tingkat pendidikan tidak ditemukan | data pendidikan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: missing required params on data",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "doesn't match schema due to: ` +
				`Error at \"/data/tingkat_pendidikan_id\": property \"tingkat_pendidikan_id\" is missing` +
				` | Error at \"/data/nama_sekolah\": property \"nama_sekolah\" is missing` +
				` | Error at \"/data/tahun_lulus\": property \"tahun_lulus\" is missing` +
				` | Error at \"/data/nomor_ijazah\": property \"nomor_ijazah\" is missing Or ` +
				`Error at \"/action\": value is not one of the allowed values [\"UPDATE\"]` +
				` | Error at \"/data/tingkat_pendidikan_id\": property \"tingkat_pendidikan_id\" is missing` +
				` | Error at \"/data/nama_sekolah\": property \"nama_sekolah\" is missing` +
				` | Error at \"/data/tahun_lulus\": property \"tahun_lulus\" is missing` +
				` | Error at \"/data/nomor_ijazah\": property \"nomor_ijazah\" is missing` +
				` | Error at \"/data_id\": property \"data_id\" is missing Or ` +
				`Error at \"/action\": value is not one of the allowed values [\"DELETE\"]` +
				` | property \"data\" is unsupported` +
				` | Error at \"/data_id\": property \"data_id\" is missing"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:                 "error: body is empty",
			requestHeader:        http.Header{"Authorization": authHeader1a},
			requestBody:          `{}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "doesn't match schema due to: ` +
				`Error at \"/action\": property \"action\" is missing` +
				` | Error at \"/data\": property \"data\" is missing Or ` +
				`Error at \"/action\": property \"action\" is missing` +
				` | Error at \"/data_id\": property \"data_id\" is missing` +
				` | Error at \"/data\": property \"data\" is missing Or ` +
				`Error at \"/action\": property \"action\" is missing` +
				` | Error at \"/data_id\": property \"data_id\" is missing"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			// validate create & approve usulan
			req := httptest.NewRequest(http.MethodPost, "/v1/usulan-perubahan-data/riwayat-pendidikan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			if tt.doRollback {
				req.URL.RawQuery = "rollback=true"
			}
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponsePostCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponsePostBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			nip := apitest.GetNIPFromAuthHeader(req.Header.Get("Authorization"))

			actualSvcRows, err := dbtest.QueryWithClause(db, "riwayat_pendidikan", "where nip = $1 order by id", nip)
			require.NoError(t, err)
			if len(tt.wantDBSvcRows) == len(actualSvcRows) {
				for i, row := range actualSvcRows {
					if tt.wantDBSvcRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBSvcRows[i]["id"] = row["id"]
						tt.wantDBSvcRows[i]["created_at"] = row["created_at"]
						tt.wantDBSvcRows[i]["updated_at"] = row["updated_at"]
					}
					if tt.wantDBSvcRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBSvcRows[i]["updated_at"] = row["updated_at"]
					}
					if tt.wantDBSvcRows[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBSvcRows[i]["deleted_at"] = row["deleted_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBSvcRows, actualSvcRows)

			actualUsulanRows, err := dbtest.QueryWithClause(db, "usulan_perubahan_data", "where nip = $1 order by id", nip)
			require.NoError(t, err)
			if len(tt.wantDBUsulanRows) == len(actualUsulanRows) {
				for i, row := range actualUsulanRows {
					assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
					assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)

					tt.wantDBUsulanRows[i]["id"] = row["id"]
					tt.wantDBUsulanRows[i]["created_at"] = row["created_at"]
					tt.wantDBUsulanRows[i]["updated_at"] = row["updated_at"]

					tt.wantResponseGetBody = strings.ReplaceAll(tt.wantResponseGetBody, "{id}", fmt.Sprintf("%d", row["id"]))
				}
			}
			assert.Equal(t, tt.wantDBUsulanRows, actualUsulanRows)

			// validate get usulan
			req2 := httptest.NewRequest(http.MethodGet, "/v1/usulan-perubahan-data/riwayat-pendidikan", nil)
			req2.Header = tt.requestHeader
			rec2 := httptest.NewRecorder()

			e.ServeHTTP(rec2, req2)

			assert.Equal(t, http.StatusOK, rec2.Code)
			assert.JSONEq(t, tt.wantResponseGetBody, rec2.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec2, req2, e))
		})
	}
}
