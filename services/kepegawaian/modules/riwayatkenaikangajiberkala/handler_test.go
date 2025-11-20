package riwayatkenaikangajiberkala

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strconv"
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
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai (id, pns_id, nip_baru, nama, created_at, deleted_at) values
			(1, 1, '1c', 'Test User 1', '2000-01-01', null),
			(2, 2, '2c', 'Test User 2', '2000-01-01', null);

		insert into ref_golongan (id, nama, nama_pangkat, created_at) values
			(1, 'III/a', 'Penata Muda', '2000-01-01'),
			(2, 'III/b', 'Penata Muda Tingkat I', '2000-01-01');

		insert into ref_unit_kerja (id, nama_unor, created_at) values
			(1, 'Unit Kerja 1', '2000-01-01'),
			(2, 'Unit Kerja 2', '2000-01-01');

		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, jabatan, tmt_jabatan, pendidikan_terakhir, tanggal_lulus_pendidikan_terakhir, kantor_pembayaran, unit_kerja_induk_id, pejabat, created_at, deleted_at) values
			(11, 1, 1, 'SK/123/2023', '2023-01-15', '2023-01-15', 2, 6, '2023-01-15', '3500000', 3500000, 'Staff A', '2023-01-15', 'S1', '2018-05-20', 'Kantor A', 1, 'Pejabat A', now(), null),
			(12, 1, 2, 'SK/124/2024', '2024-01-15', '2024-01-15', 3, 0, '2024-01-15', '4000000', 4000000, 'Staff B', '2024-01-15', 'S2', '2019-06-20', 'Kantor B', 2, 'Pejabat B', now(), null),
			(13, 2, 1, 'SK/125/2023', '2023-06-15', '2023-06-15', 1, 3, '2023-06-15', '3200000', 3200000, 'Staff C', '2023-06-15', 'S1', '2017-04-10', 'Kantor C', 1, 'Pejabat C', now(), null),
			(14, 1, 1, 'SK/126/2022', '2022-01-15', '2022-01-15', 1, 0, '2022-01-15', '3000000', 3000000, 'Staff D', '2022-01-15', 'S3', '2016-03-15', 'Kantor D', 1, 'Pejabat D', now(), '2020-01-01'),
			-- Null test cases
			(15, 1, 1, null, '2025-01-15', '2025-01-15', 4, 0, '2025-01-15', '4500000', 4500000, 'Staff E', '2025-01-15', 'S4', '2020-07-01', 'Kantor E', 2, 'Pejabat E', now(), null),
			(16, 1, 2, 'SK/127/2025', null, '2025-02-15', 5, 3, '2025-02-15', '5000000', 5000000, 'Staff F', '2025-02-15', 'S5', '2021-08-01', 'Kantor F', 1, 'Pejabat F', now(), null),
			(17, 1, 1, 'SK/128/2025', '2025-03-15', null, 6, 6, '2025-03-15', '5500000', 5500000, 'Staff G', '2025-03-15', 'S6', '2022-09-01', 'Kantor G', 2, 'Pejabat G', now(), null),
			(18, 1, 2, 'SK/129/2025', '2025-04-15', '2025-04-15', null, 9, '2025-04-15', '6000000', 6000000, 'Staff H', '2025-04-15', 'S7', '2023-10-01', 'Kantor H', 1, 'Pejabat H', now(), null),
			(19, 1, 1, 'SK/130/2025', '2025-05-15', '2025-05-15', 7, null, '2025-05-15', '6500000', 6500000, 'Staff I', '2025-05-15', 'S8', '2024-11-01', 'Kantor I', 2, 'Pejabat I', now(), null),
			(20, 1, 2, 'SK/131/2025', '2025-06-15', '2025-06-15', 8, 0, null, '7000000', 7000000, 'Staff J', '2025-06-15', 'S9', '2025-01-01', 'Kantor J', 1, 'Pejabat J', now(), null),
			(21, 1, 1, 'SK/132/2025', '2025-07-15', '2025-07-15', 9, 3, '2025-07-15', '7500000', null, null, null, null, null, null, null, null, now(), null);
	`
	pgxConn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxConn.Exec(t.Context(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	dbRepository := repo.New(pgxConn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, dbRepository, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
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
						"id": 20,
						"id_golongan": 2,
						"nama_golongan": "III/b",
						"nama_golongan_pangkat": "Penata Muda Tingkat I",
						"nomor_sk": "SK/131/2025",
						"tanggal_sk": "2025-06-15",
						"tmt_golongan": "2025-06-15",
						"masa_kerja_golongan_tahun": 8,
						"masa_kerja_golongan_bulan": 0,
						"tmt_kenaikan_gaji_berkala": null,
						"gaji_pokok": 7000000,
						"jabatan": "Staff J",
						"tmt_jabatan": "2025-06-15",
						"pendidikan": "S9",
						"tanggal_lulus": "2025-01-01",
						"kantor_pembayaran": "Kantor J",
						"unit_kerja_induk_id": "1",
						"unit_kerja_induk": "Unit Kerja 1",
						"pejabat": "Pejabat J"
					},
					{
						"id": 21,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "SK/132/2025",
						"tanggal_sk": "2025-07-15",
						"tmt_golongan": "2025-07-15",
						"masa_kerja_golongan_tahun": 9,
						"masa_kerja_golongan_bulan": 3,
						"tmt_kenaikan_gaji_berkala": "2025-07-15",
						"gaji_pokok": null,
						"jabatan": "",
						"tmt_jabatan": null,
						"pendidikan": "",
						"tanggal_lulus": null,
						"kantor_pembayaran": "",
						"unit_kerja_induk_id": "",
						"unit_kerja_induk": "",
						"pejabat": ""
					},
					{
						"id": 19,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "SK/130/2025",
						"tanggal_sk": "2025-05-15",
						"tmt_golongan": "2025-05-15",
						"masa_kerja_golongan_tahun": 7,
						"masa_kerja_golongan_bulan": null,
						"tmt_kenaikan_gaji_berkala": "2025-05-15",
						"gaji_pokok": 6500000,
						"jabatan": "Staff I",
						"tmt_jabatan": "2025-05-15",
						"pendidikan": "S8",
						"tanggal_lulus": "2024-11-01",
						"kantor_pembayaran": "Kantor I",
						"unit_kerja_induk_id": "2",
						"unit_kerja_induk": "Unit Kerja 2",
						"pejabat": "Pejabat I"
					},
					{
						"id": 18,
						"id_golongan": 2,
						"nama_golongan": "III/b",
						"nama_golongan_pangkat": "Penata Muda Tingkat I",
						"nomor_sk": "SK/129/2025",
						"tanggal_sk": "2025-04-15",
						"tmt_golongan": "2025-04-15",
						"masa_kerja_golongan_tahun": null,
						"masa_kerja_golongan_bulan": 9,
						"tmt_kenaikan_gaji_berkala": "2025-04-15",
						"gaji_pokok": 6000000,
						"jabatan": "Staff H",
						"tmt_jabatan": "2025-04-15",
						"pendidikan": "S7",
						"tanggal_lulus": "2023-10-01",
						"kantor_pembayaran": "Kantor H",
						"unit_kerja_induk_id": "1",
						"unit_kerja_induk": "Unit Kerja 1",
						"pejabat": "Pejabat H"
					},
					{
						"id": 17,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "SK/128/2025",
						"tanggal_sk": "2025-03-15",
						"tmt_golongan": null,
						"masa_kerja_golongan_tahun": 6,
						"masa_kerja_golongan_bulan": 6,
						"tmt_kenaikan_gaji_berkala": "2025-03-15",
						"gaji_pokok": 5500000,
						"jabatan": "Staff G",
						"tmt_jabatan": "2025-03-15",
						"pendidikan": "S6",
						"tanggal_lulus": "2022-09-01",
						"kantor_pembayaran": "Kantor G",
						"unit_kerja_induk_id": "2",
						"unit_kerja_induk": "Unit Kerja 2",
						"pejabat": "Pejabat G"
					},
					{
						"id": 16,
						"id_golongan": 2,
						"nama_golongan": "III/b",
						"nama_golongan_pangkat": "Penata Muda Tingkat I",
						"nomor_sk": "SK/127/2025",
						"tanggal_sk": null,
						"tmt_golongan": "2025-02-15",
						"masa_kerja_golongan_tahun": 5,
						"masa_kerja_golongan_bulan": 3,
						"tmt_kenaikan_gaji_berkala": "2025-02-15",
						"gaji_pokok": 5000000,
						"jabatan": "Staff F",
						"tmt_jabatan": "2025-02-15",
						"pendidikan": "S5",
						"tanggal_lulus": "2021-08-01",
						"kantor_pembayaran": "Kantor F",
						"unit_kerja_induk_id": "1",
						"unit_kerja_induk": "Unit Kerja 1",
						"pejabat": "Pejabat F"
					},
					{
						"id": 15,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "",
						"tanggal_sk": "2025-01-15",
						"tmt_golongan": "2025-01-15",
						"masa_kerja_golongan_tahun": 4,
						"masa_kerja_golongan_bulan": 0,
						"tmt_kenaikan_gaji_berkala": "2025-01-15",
						"gaji_pokok": 4500000,
						"jabatan": "Staff E",
						"tmt_jabatan": "2025-01-15",
						"pendidikan": "S4",
						"tanggal_lulus": "2020-07-01",
						"kantor_pembayaran": "Kantor E",
						"unit_kerja_induk_id": "2",
						"unit_kerja_induk": "Unit Kerja 2",
						"pejabat": "Pejabat E"
					},
					{
						"id": 12,
						"id_golongan": 2,
						"nama_golongan": "III/b",
						"nama_golongan_pangkat": "Penata Muda Tingkat I",
						"nomor_sk": "SK/124/2024",
						"tanggal_sk": "2024-01-15",
						"tmt_golongan": "2024-01-15",
						"masa_kerja_golongan_tahun": 3,
						"masa_kerja_golongan_bulan": 0,
						"tmt_kenaikan_gaji_berkala": "2024-01-15",
						"gaji_pokok": 4000000,
						"jabatan": "Staff B",
						"tmt_jabatan": "2024-01-15",
						"pendidikan": "S2",
						"tanggal_lulus": "2019-06-20",
						"kantor_pembayaran": "Kantor B",
						"unit_kerja_induk_id": "2",
						"unit_kerja_induk": "Unit Kerja 2",
						"pejabat": "Pejabat B"
					},
					{
						"id": 11,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "SK/123/2023",
						"tanggal_sk": "2023-01-15",
						"tmt_golongan": "2023-01-15",
						"masa_kerja_golongan_tahun": 2,
						"masa_kerja_golongan_bulan": 6,
						"tmt_kenaikan_gaji_berkala": "2023-01-15",
						"gaji_pokok": 3500000,
						"jabatan": "Staff A",
						"tmt_jabatan": "2023-01-15",
						"pendidikan": "S1",
						"tanggal_lulus": "2018-05-20",
						"kantor_pembayaran": "Kantor A",
						"unit_kerja_induk_id": "1",
						"unit_kerja_induk": "Unit Kerja 1",
						"pejabat": "Pejabat A"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 9}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":21,
						"id_golongan":1,
						"nama_golongan":"III/a",
						"nama_golongan_pangkat":"Penata Muda",
						"nomor_sk":"SK/132/2025",
						"tanggal_sk":"2025-07-15",
						"tmt_golongan":"2025-07-15",
						"masa_kerja_golongan_tahun":9,
						"masa_kerja_golongan_bulan":3,
						"tmt_kenaikan_gaji_berkala":"2025-07-15",
						"gaji_pokok":null,
						"jabatan":"",
						"tmt_jabatan":null,
						"pendidikan":"",
						"tanggal_lulus":null,
						"kantor_pembayaran":"",
						"unit_kerja_induk_id": "",
						"unit_kerja_induk":"",
						"pejabat": ""
					}
				],
				"meta": {"limit":1, "offset":1, "total":9}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("3c")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit":10, "offset":0, "total":0}
			}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestHeader:    http.Header{"Authorization": []string{"invalid"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message":"token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-kenaikan-gaji-berkala", nil)
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
		insert into pegawai (id, pns_id, nip_baru, nama, created_at, deleted_at) values
			(1, 1, '1c', 'Test User 1', '2000-01-01', null),
			(2, 2, '2c', 'Test User 2', '2000-01-01', null);

		insert into riwayat_kenaikan_gaji_berkala (id, pegawai_id, deleted_at, file_base64) values
			(1, 1, null, 'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, 1, null, '` + pdfBase64 + `'),
			(3, 1, null, 'data:images/png;base64,` + pngBase64 + `'),
			(4, 1, null, 'data:application/pdf;base64,invalid'),
			(5, 1, '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, 1, null, null),
			(7, 1, null, '');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	queries := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, queries, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("1c")}
	tests := []struct {
		name              string
		paramID           int64
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			paramID:           1,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			paramID:           2,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			paramID:           3,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 kenaikan gaji berkala tidak valid",
			paramID:           4,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat kenaikan gaji berkala sudah dihapus",
			paramID:           5,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kenaikan gaji berkala berisi null value",
			paramID:           6,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kenaikan gaji berkala berupa string kosong",
			paramID:           7,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kenaikan gaji berkala bukan milik user login",
			paramID:           1,
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kenaikan gaji berkala tidak ditemukan",
			paramID:           0,
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: auth header tidak valid",
			paramID:           1,
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-kenaikan-gaji-berkala/%d/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseCode == http.StatusOK {
				assert.Equal(t, "inline", rec.Header().Get("Content-Disposition"))
				assert.Equal(t, tt.wantContentType, rec.Header().Get("Content-Type"))
				assert.Equal(t, tt.wantResponseBytes, rec.Body.Bytes())
			} else {
				assert.JSONEq(t, string(tt.wantResponseBytes), rec.Body.String())
			}
		})
	}
}

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai (id, pns_id, nip_baru, nama, created_at, deleted_at) values
			(1, 1, '1c', 'Test User 1', '2000-01-01', null),
			(2, 2, '2c', 'Test User 2', '2000-01-01', null),
			(3, 3, '3c', 'Test User 3', '2000-01-01', '2020-01-01');

		insert into ref_golongan (id, nama, nama_pangkat, created_at) values
			(1, 'III/a', 'Penata Muda', '2000-01-01'),
			(2, 'III/b', 'Penata Muda Tingkat I', '2000-01-01');

		insert into ref_unit_kerja (id, nama_unor, created_at) values
			(1, 'Unit Kerja 1', '2000-01-01'),
			(2, 'Unit Kerja 2', '2000-01-01');

		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, jabatan, tmt_jabatan, pendidikan_terakhir, tanggal_lulus_pendidikan_terakhir, kantor_pembayaran, unit_kerja_induk_id, pejabat, created_at, deleted_at) values
			(11, 1, 1, 'SK/123/2023', '2023-01-15', '2023-01-15', 2, 6, '2023-01-15', '3500000', 3500000, 'Staff A', '2023-01-15', 'S1', '2018-05-20', 'Kantor A', 1, 'Pejabat A', now(), null),
			(12, 1, 2, 'SK/124/2024', '2024-01-15', '2024-01-15', 3, 0, '2024-01-15', '4000000', 4000000, 'Staff B', '2024-01-15', 'S2', '2019-06-20', 'Kantor B', 2, 'Pejabat B', now(), null),
			(13, 2, 1, 'SK/125/2023', '2023-06-15', '2023-06-15', 1, 3, '2023-06-15', '3200000', 3200000, 'Staff C', '2023-06-15', 'S1', '2017-04-10', 'Kantor C', 1, 'Pejabat C', now(), null),
			(14, 1, 1, 'SK/126/2022', '2022-01-15', '2022-01-15', 1, 0, '2022-01-15', '3000000', 3000000, 'Staff D', '2022-01-15', 'S3', '2016-03-15', 'Kantor D', 1, 'Pejabat D', now(), '2020-01-01');
`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	queries := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, queries, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			name:             "ok: nip 1c data returned",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 12,
						"id_golongan": 2,
						"nama_golongan": "III/b",
        				"nama_golongan_pangkat": "Penata Muda Tingkat I",
						"nomor_sk": "SK/124/2024",
						"tanggal_sk": "2024-01-15",
						"tmt_golongan": "2024-01-15",
						"masa_kerja_golongan_tahun": 3,
						"masa_kerja_golongan_bulan": 0,
						"tmt_kenaikan_gaji_berkala": "2024-01-15",
						"gaji_pokok": 4000000,
						"jabatan": "Staff B",
						"tmt_jabatan": "2024-01-15",
						"pendidikan": "S2",
						"tanggal_lulus": "2019-06-20",
						"kantor_pembayaran": "Kantor B",
						"unit_kerja_induk_id": "2",
						"unit_kerja_induk": "Unit Kerja 2",
						"pejabat": "Pejabat B"
					},
					{
						"id": 11,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "SK/123/2023",
						"tanggal_sk": "2023-01-15",
						"tmt_golongan": "2023-01-15",
						"masa_kerja_golongan_tahun": 2,
						"masa_kerja_golongan_bulan": 6,
						"tmt_kenaikan_gaji_berkala": "2023-01-15",
						"gaji_pokok": 3500000,
						"jabatan": "Staff A",
						"tmt_jabatan": "2023-01-15",
						"pendidikan": "S1",
						"tanggal_lulus": "2018-05-20",
						"kantor_pembayaran": "Kantor A",
						"unit_kerja_induk_id": "1",
						"unit_kerja_induk": "Unit Kerja 1",
						"pejabat": "Pejabat A"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 11,
						"id_golongan": 1,
						"nama_golongan": "III/a",
						"nama_golongan_pangkat": "Penata Muda",
						"nomor_sk": "SK/123/2023",
						"tanggal_sk": "2023-01-15",
						"tmt_golongan": "2023-01-15",
						"masa_kerja_golongan_tahun": 2,
						"masa_kerja_golongan_bulan": 6,
						"tmt_kenaikan_gaji_berkala": "2023-01-15",
						"gaji_pokok": 3500000,
						"jabatan": "Staff A",
						"tmt_jabatan": "2023-01-15",
						"pendidikan": "S1",
						"tanggal_lulus": "2018-05-20",
						"kantor_pembayaran": "Kantor A",
						"unit_kerja_induk_id": "1",
						"unit_kerja_induk": "Unit Kerja 1",
						"pejabat": "Pejabat A"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 2}
			}`,
		},
		{
			name:             "ok: nip 200 gets empty data",
			nip:              "200",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "ok: nip 3c gets empty data (deleted pegawai)",
			nip:              "3c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-kenaikan-gaji-berkala", nil)
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

func Test_handler_getBerkasAdmin(t *testing.T) {
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
		insert into pegawai (id, pns_id, nip_baru, nama, created_at, deleted_at) values
			(1, 1, '1c', 'Test User 1', '2000-01-01', null),
			(2, 2, '2a', 'Test User 2', '2000-01-01', null);

		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, deleted_at, file_base64) values
			(1, 1, null, 'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, 1, null, '` + pdfBase64 + `'),
			(3, 1, null, 'data:images/png;base64,` + pngBase64 + `'),
			(4, 1, null, 'data:application/pdf;base64,invalid'),
			(5, 1, '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, 1, null, null),
			(7, 1, null, ''),
			(8, 2, null, 'data:application/pdf;base64,` + pdfBase64 + `');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	queries := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, queries, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name              string
		nip               string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			nip:               "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			nip:               "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "ok: admin can access other user's berkas",
			nip:               "2a",
			paramID:           "8",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "error: base64 tidak valid",
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat with wrong nip",
			nip:               "wrong-nip",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat tidak ditemukan",
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kenaikan gaji berkala tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			nip:               "1c",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-kenaikan-gaji-berkala/%s/berkas", tt.nip, tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseCode == http.StatusOK {
				assert.Equal(t, "inline", rec.Header().Get("Content-Disposition"))
				assert.Equal(t, tt.wantContentType, rec.Header().Get("Content-Type"))
				assert.Equal(t, tt.wantResponseBytes, rec.Body.Bytes())
			} else {
				assert.JSONEq(t, string(tt.wantResponseBytes), rec.Body.String())
			}
		})
	}
}

func Test_handler_adminCreate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_lokasi (id, nama, created_at, deleted_at) VALUES
			(1, 'Jakarta deleted', '2000-01-01', now()),
			(2, 'Semarang', '2000-01-01', null),
			(3, 'Surabaya', '2000-01-01', null),
			(4, 'Bandung', '2000-01-01', null),
			(5, 'Medan', '2000-01-01', null);
		insert into pegawai
			(pns_id,  nip_baru, nama, tanggal_lahir, tempat_lahir, tempat_lahir_id, deleted_at) values
			('id_1c', '1c', 'Pegawai 1', '2000-01-01', 'Jakarta',    1, null),
			('id_1d', '1d', 'Pegawai 2', '2000-02-01', 'Semarang',    2, '2000-01-01'),
			('id_1e', '1e', 'Pegawai 3', '2000-03-01', 'Surabaya - unused',    3, null),
			('id_1f', '1f', 'Pegawai 4', '2000-04-01', 'Bandung',    4, null),
			('id_1g', '1g', 'Pegawai 5', '2000-05-01', 'Medan',    5, null);
		insert into ref_golongan (id, nama, nama_pangkat, created_at) values
			(1, 'III/a', 'Penata Muda', '2000-01-01'),
			(2, 'III/b', 'Penata Muda Tingkat I', '2000-01-01');
		insert into ref_unit_kerja (id, nama_unor, created_at) values
			(1, 'Unit Kerja 1', '2000-01-01'),
			(2, 'Unit Kerja 2', '2000-01-01');
		SELECT setval('riwayat_kgb_id_seq', 1, false);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 3,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/125/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 4500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"jabatan": "Staff Baru",
				"tmt_jabatan": "2025-01-15",
				"pendidikan": "S2",
				"tanggal_lulus": "2020-07-01",
				"kantor_pembayaran": "Kantor Baru",
				"pejabat": "Pejabat Baru",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                "{id}",
					"pegawai_id":                        int32(1),
					"pegawai_nama":                      "Pegawai 1",
					"pegawai_nip":                       "1c",
					"tempat_lahir":                      "Jakarta",
					"tanggal_lahir":                     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"n_gol_ruang":                       "III/a",
					"golongan_id":                       int32(1),
					"no_sk":                             "SK/125/2025",
					"tanggal_sk":                        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"masa_kerja_golongan_tahun":         int16(3),
					"masa_kerja_golongan_bulan":         int16(0),
					"tmt_sk":                            time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"gaji_pokok":                        int32(4500000),
					"jabatan":                           "Staff Baru",
					"tmt_jabatan":                       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"pendidikan_terakhir":               "S2",
					"tanggal_lulus_pendidikan_terakhir": time.Date(2020, 7, 1, 0, 0, 0, 0, time.UTC),
					"kantor_pembayaran":                 "Kantor Baru",
					"unit_kerja_induk_id":               "1",
					"pejabat":                           "Pejabat Baru",
					"created_at":                        "{created_at}",
					"updated_at":                        "{updated_at}",
					"deleted_at":                        nil,
					"ref":                               "{uuid}",
					"file_base64":                       nil,
					"alasan":                            nil,
					"n_gapok":                           nil,
					"keterangan_berkas":                 nil,
					"unit_kerja_induk_text":             "Unit Kerja 1",
					"mv_kgb_id":                         nil,
				},
			},
		},
		{
			name:          "ok: with null/empty values",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 2,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/126/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 4000000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"jabatan": "",
				"tmt_jabatan": null,
				"pendidikan": "",
				"tanggal_lulus": null,
				"kantor_pembayaran": "",
				"pejabat": "",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                "{id}",
					"pegawai_id":                        int32(3),
					"pegawai_nama":                      "Pegawai 3",
					"pegawai_nip":                       "1e",
					"tempat_lahir":                      "Surabaya",
					"tanggal_lahir":                     time.Date(2000, 3, 1, 0, 0, 0, 0, time.UTC),
					"n_gol_ruang":                       "III/a",
					"golongan_id":                       int32(1),
					"no_sk":                             "SK/126/2025",
					"tanggal_sk":                        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"masa_kerja_golongan_tahun":         int16(2),
					"masa_kerja_golongan_bulan":         int16(0),
					"tmt_sk":                            time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"gaji_pokok":                        int32(4000000),
					"jabatan":                           nil,
					"tmt_jabatan":                       nil,
					"pendidikan_terakhir":               nil,
					"tanggal_lulus_pendidikan_terakhir": nil,
					"kantor_pembayaran":                 nil,
					"unit_kerja_induk_id":               "1",
					"pejabat":                           nil,
					"created_at":                        "{created_at}",
					"updated_at":                        "{updated_at}",
					"deleted_at":                        nil,
					"ref":                               "{uuid}",
					"file_base64":                       nil,
					"alasan":                            nil,
					"n_gapok":                           nil,
					"keterangan_berkas":                 nil,
					"unit_kerja_induk_text":             "Unit Kerja 1",
					"mv_kgb_id":                         nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/127/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                "{id}",
					"pegawai_id":                        int32(4),
					"pegawai_nama":                      "Pegawai 4",
					"pegawai_nip":                       "1f",
					"tempat_lahir":                      "Bandung",
					"tanggal_lahir":                     time.Date(2000, 4, 1, 0, 0, 0, 0, time.UTC),
					"n_gol_ruang":                       "III/a",
					"golongan_id":                       int32(1),
					"no_sk":                             "SK/127/2025",
					"tanggal_sk":                        time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                      time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"masa_kerja_golongan_tahun":         int16(1),
					"masa_kerja_golongan_bulan":         int16(0),
					"tmt_sk":                            time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
					"gaji_pokok":                        int32(3500000),
					"jabatan":                           nil,
					"tmt_jabatan":                       nil,
					"pendidikan_terakhir":               nil,
					"tanggal_lulus_pendidikan_terakhir": nil,
					"kantor_pembayaran":                 nil,
					"unit_kerja_induk_id":               nil,
					"pejabat":                           nil,
					"created_at":                        "{created_at}",
					"updated_at":                        "{updated_at}",
					"deleted_at":                        nil,
					"ref":                               "{uuid}",
					"file_base64":                       nil,
					"alasan":                            nil,
					"n_gapok":                           nil,
					"keterangan_berkas":                 nil,
					"unit_kerja_induk_text":             nil,
					"mv_kgb_id":                         nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/128/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"unit_kerja_induk_id": "1"
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
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/129/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: unit kerja is not found",
			paramNIP:      "1g",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/129/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"unit_kerja_induk_id": "0"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data unit kerja tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kenaikan-gaji-berkala", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if strings.Contains(tt.wantResponseBody, "{id}") {
				var resp struct {
					Data struct {
						ID int64 `json:"id"`
					} `json:"data"`
				}
				require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
				id := resp.Data.ID
				tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", strconv.FormatInt(id, 10))
				for i, row := range tt.wantDBRows {
					if row["id"] == "{id}" {
						tt.wantDBRows[i]["id"] = id
					}
				}
			}
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_kenaikan_gaji_berkala", "where pegawai_id = (select id from pegawai where nip_baru = $1) order by id", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
						tt.wantDBRows[i]["ref"] = row["ref"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pegawai
			(id,pns_id,  nip_baru, nama, tanggal_lahir, tempat_lahir, deleted_at) values
			(1,'id_1c', '1c',  'Pegawai 1', '2000-01-01', 'Jakarta', null),
			(2,'id_1d', '1d',  'Pegawai 2', '2000-02-01', 'Semarang', '2000-01-01'),
			(3,'id_1e', '1e',  'Pegawai 3', '2000-03-01', 'Surabaya', null);
		insert into ref_golongan (id, nama, nama_pangkat, created_at) values
			(1, 'III/a', 'Penata Muda', '2000-01-01'),
			(2, 'III/b', 'Penata Muda Tingkat I', '2000-01-01');
		insert into ref_unit_kerja (id, nama_unor, created_at) values
			(1, 'Unit Kerja 1', '2000-01-01'),
			(2, 'Unit Kerja 2', '2000-01-01');
		insert into riwayat_kenaikan_gaji_berkala
			(id, ref, pegawai_id, pegawai_nama, pegawai_nip, tempat_lahir, tanggal_lahir, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, gaji_pokok, jabatan, tmt_jabatan, pendidikan_terakhir, tanggal_lulus_pendidikan_terakhir, kantor_pembayaran, unit_kerja_induk_id, pejabat, created_at, updated_at, deleted_at) values
			(1, '00000000-0000-0000-0000-000000000001', 1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/123/2023', '2023-01-15', '2023-01-15', 2, 6, '2023-01-15', 3500000, 'Staff A', '2023-01-15', 'S1', '2018-05-20', 'Kantor A', '1', 'Pejabat A', '2000-01-01', '2000-01-01', null),
			(2, '00000000-0000-0000-0000-000000000002', 1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/124/2024', '2024-01-15', '2024-01-15', 3, 0, '2024-01-15', 4000000, 'Staff B', '2024-01-15', 'S2', '2019-06-20', 'Kantor B', '1', 'Pejabat B', '2000-01-01', '2000-01-01', null),
			(3, '00000000-0000-0000-0000-000000000003', 1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/125/2023', '2023-06-15', '2023-06-15', 1, 3, '2023-06-15', 3200000, 'Staff C', '2023-06-15', 'S1', '2017-04-10', 'Kantor C', '1', 'Pejabat C', '2000-01-01', '2000-01-01', null),
			(4, '00000000-0000-0000-0000-000000000004', 1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/126/2022', '2022-01-15', '2022-01-15', 1, 0, '2022-01-15', 3000000, 'Staff D', '2022-01-15', 'S3', '2016-03-15', 'Kantor D', '1', 'Pejabat D', '2000-01-01', '2000-01-01', null),
			(5, '00000000-0000-0000-0000-000000000005', 1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/127/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, null, null, null, null, null, '1', null, '2000-01-01', '2000-01-01', null),
			(6, '00000000-0000-0000-0000-000000000006', 3, 'Pegawai 3', '1e', 'Surabaya', '2000-03-01', 1, 'SK/128/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, 'Staff E', '2020-01-15', 'S1', '2015-01-01', 'Kantor E', '1', 'Pejabat E', '2000-01-01', '2000-01-01', null),
			(7, '00000000-0000-0000-0000-000000000007', 1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/129/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, 'Staff F', '2020-01-15', 'S1', '2015-01-01', 'Kantor F', '1', 'Pejabat F', '2000-01-01', '2000-01-01', '2000-01-01'),
			(8, '00000000-0000-0000-0000-000000000008',  1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/130/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, 'Staff G', '2020-01-15', 'S1', '2015-01-01', 'Kantor G', '1', 'Pejabat G', '2000-01-01', '2000-01-01', null),
			(9, '00000000-0000-0000-0000-000000000009',  1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/131/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, 'Staff H', '2020-01-15', 'S1', '2015-01-01', 'Kantor H', '1', 'Pejabat H', '2000-01-01', '2000-01-01', null),
			(10, '00000000-0000-0000-0000-000000000010',  1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/132/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, 'Staff I', '2020-01-15', 'S1', '2015-01-01', 'Kantor I', '1', 'Pejabat I', '2000-01-01', '2000-01-01', null),
			(11, '00000000-0000-0000-0000-000000000011',  1, 'Pegawai 1', '1c', 'Jakarta', '2000-01-01', 1, 'SK/133/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', 3000000, 'Staff J', '2020-01-15', 'S1', '2015-01-01', 'Kantor J', '1', 'Pejabat J', '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
				"golongan_id": 2,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 4,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/125/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 5000000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"jabatan": "Staff Updated",
				"tmt_jabatan": "2025-01-15",
				"pendidikan": "S3",
				"tanggal_lulus": "2021-07-01",
				"kantor_pembayaran": "Kantor Updated",
				"pejabat": "Pejabat Updated",
				"unit_kerja_induk_id": "2"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(1),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(5000000),
					"golongan_id":                       int32(2),
					"jabatan":                           "Staff Updated",
					"kantor_pembayaran":                 "Kantor Updated",
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(4),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       "III/b",
					"no_sk":                             "SK/125/2025",
					"pegawai_id":                        int32(1),
					"pegawai_nama":                      "Pegawai 1",
					"pegawai_nip":                       "1c",
					"pejabat":                           "Pejabat Updated",
					"pendidikan_terakhir":               "S3",
					"ref":                               "00000000-0000-0000-0000-000000000001",
					"tanggal_lahir":                     time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": time.Date(2021, time.July, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_sk":                        time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Jakarta",
					"tmt_golongan":                      time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_sk":                            time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               "2",
					"unit_kerja_induk_text":             "Unit Kerja 2",
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        "{updated_at}",
					"deleted_at":                        nil,
				},
			},
		},
		{
			name:          "ok: with null/empty values",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 2,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/126/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 4000000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"jabatan": "",
				"tmt_jabatan": null,
				"pendidikan": "",
				"tanggal_lulus": null,
				"kantor_pembayaran": "",
				"pejabat": "",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(3),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(4000000),
					"golongan_id":                       int32(1),
					"jabatan":                           nil,
					"kantor_pembayaran":                 nil,
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(2),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       "III/a",
					"no_sk":                             "SK/126/2025",
					"pegawai_id":                        int32(1),
					"pegawai_nama":                      "Pegawai 1",
					"pegawai_nip":                       "1c",
					"pejabat":                           nil,
					"pendidikan_terakhir":               nil,
					"ref":                               "00000000-0000-0000-0000-000000000003",
					"tanggal_lahir":                     time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": nil,
					"tanggal_sk":                        time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Jakarta",
					"tmt_golongan":                      time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       nil,
					"tmt_sk":                            time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               "1",
					"unit_kerja_induk_text":             "Unit Kerja 1",
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        "{updated_at}",
					"deleted_at":                        nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/127/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(4),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(3500000),
					"golongan_id":                       int32(1),
					"jabatan":                           nil,
					"kantor_pembayaran":                 nil,
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(1),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       "III/a",
					"no_sk":                             "SK/127/2025",
					"pegawai_id":                        int32(1),
					"pegawai_nama":                      "Pegawai 1",
					"pegawai_nip":                       "1c",
					"pejabat":                           nil,
					"pendidikan_terakhir":               nil,
					"ref":                               "00000000-0000-0000-0000-000000000004",
					"tanggal_lahir":                     time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": nil,
					"tanggal_sk":                        time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Jakarta",
					"tmt_golongan":                      time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       nil,
					"tmt_sk":                            time.Date(2025, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               nil,
					"unit_kerja_induk_text":             nil,
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        "{updated_at}",
					"deleted_at":                        nil,
				},
			},
		},
		{
			name:          "error: riwayat kenaikan gaji berkala is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/128/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat kenaikan gaji berkala is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/129/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(6),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(3000000),
					"golongan_id":                       int32(1),
					"jabatan":                           "Staff E",
					"kantor_pembayaran":                 "Kantor E",
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(1),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       nil,
					"no_sk":                             "SK/128/2020",
					"pegawai_id":                        int32(3),
					"pegawai_nama":                      "Pegawai 3",
					"pegawai_nip":                       "1e",
					"pejabat":                           "Pejabat E",
					"pendidikan_terakhir":               "S1",
					"ref":                               "00000000-0000-0000-0000-000000000006",
					"tanggal_lahir":                     time.Date(2000, time.March, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_sk":                        time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Surabaya",
					"tmt_golongan":                      time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_sk":                            time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               "1",
					"unit_kerja_induk_text":             nil,
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                        nil,
				},
			},
		},
		{
			name:          "error: riwayat kenaikan gaji berkala is deleted",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"tmt_golongan": "2025-01-15",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 0,
				"nomor_sk": "SK/130/2025",
				"tanggal_sk": "2025-01-15",
				"gaji_pokok": 3500000,
				"tmt_kenaikan_gaji_berkala": "2025-01-15",
				"unit_kerja_induk_id": "1"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(7),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(3000000),
					"golongan_id":                       int32(1),
					"jabatan":                           "Staff F",
					"kantor_pembayaran":                 "Kantor F",
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(1),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       nil,
					"no_sk":                             "SK/129/2020",
					"pegawai_id":                        int32(1),
					"pegawai_nama":                      "Pegawai 1",
					"pegawai_nip":                       "1c",
					"pejabat":                           "Pejabat F",
					"pendidikan_terakhir":               "S1",
					"ref":                               "00000000-0000-0000-0000-000000000007",
					"tanggal_lahir":                     time.Date(2000, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_sk":                        time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Jakarta",
					"tmt_golongan":                      time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_sk":                            time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               "1",
					"unit_kerja_induk_text":             nil,
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1e",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(6),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(3000000),
					"golongan_id":                       int32(1),
					"jabatan":                           "Staff E",
					"kantor_pembayaran":                 "Kantor E",
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(1),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       nil,
					"no_sk":                             "SK/128/2020",
					"pegawai_id":                        int32(3),
					"pegawai_nama":                      "Pegawai 3",
					"pegawai_nip":                       "1e",
					"pejabat":                           "Pejabat E",
					"pendidikan_terakhir":               "S1",
					"ref":                               "00000000-0000-0000-0000-000000000006",
					"tanggal_lahir":                     time.Date(2000, time.March, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_sk":                        time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Surabaya",
					"tmt_golongan":                      time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_sk":                            time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               "1",
					"unit_kerja_induk_text":             nil,
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                        nil,
				},
			},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1e",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			requestBody:      `{"golongan_id": 1, "tmt_golongan": "2025-01-15", "masa_kerja_golongan_tahun": 1, "masa_kerja_golongan_bulan": 0, "nomor_sk": "SK/132/2025", "tanggal_sk": "2025-01-15", "gaji_pokok": 3500000, "tmt_kenaikan_gaji_berkala": "2025-01-15", "unit_kerja_induk_id": "1"}`,
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                                int64(6),
					"alasan":                            nil,
					"file_base64":                       nil,
					"gaji_pokok":                        int32(3000000),
					"golongan_id":                       int32(1),
					"jabatan":                           "Staff E",
					"kantor_pembayaran":                 "Kantor E",
					"keterangan_berkas":                 nil,
					"masa_kerja_golongan_bulan":         int16(0),
					"masa_kerja_golongan_tahun":         int16(1),
					"mv_kgb_id":                         nil,
					"n_gapok":                           nil,
					"n_gol_ruang":                       nil,
					"no_sk":                             "SK/128/2020",
					"pegawai_id":                        int32(3),
					"pegawai_nama":                      "Pegawai 3",
					"pegawai_nip":                       "1e",
					"pejabat":                           "Pejabat E",
					"pendidikan_terakhir":               "S1",
					"ref":                               "00000000-0000-0000-0000-000000000006",
					"tanggal_lahir":                     time.Date(2000, time.March, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_lulus_pendidikan_terakhir": time.Date(2015, time.January, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_sk":                        time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tempat_lahir":                      "Surabaya",
					"tmt_golongan":                      time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_jabatan":                       time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"tmt_sk":                            time.Date(2020, time.January, 15, 0, 0, 0, 0, time.UTC),
					"unit_kerja_induk_id":               "1",
					"unit_kerja_induk_text":             nil,
					"created_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                        nil,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kenaikan-gaji-berkala/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_kenaikan_gaji_berkala", "where id = $1 order by id", tt.paramID)
			require.NoError(t, err)
			for i, row := range actualRows {
				if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
					assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
					tt.wantDBRows[i]["updated_at"] = row["updated_at"]
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
		insert into ref_golongan (id, nama, nama_pangkat, created_at) values
			(1, 'III/a', 'Penata Muda', '2000-01-01');
		insert into ref_unit_kerja (id, nama_unor, created_at) values
			(1, 'Unit Kerja 1', '2000-01-01');
		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, created_at, updated_at, deleted_at) values
			(1,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/123/2023', '2023-01-15', '2023-01-15', 2, 6, '2023-01-15', '3500000', 3500000, '2000-01-01', '2000-01-01', null),
			(2,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/124/2024', '2024-01-15', '2024-01-15', 3, 0, '2024-01-15', '4000000', 4000000, '2000-01-01', '2000-01-01', null),
			(3,  (select id from pegawai where nip_baru = '1e'), 1, 'SK/125/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', '3000000', 3000000, '2000-01-01', '2000-01-01', null),
			(4,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/126/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', '3000000', 3000000, '2000-01-01', '2000-01-01', null),
			(5,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/127/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', '3000000', 3000000, '2000-01-01', '2000-01-01', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDeletedAt    string
	}{
		{
			name:             "ok: success delete",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDeletedAt:    "{deleted_at}",
		},
		{
			name:             "error: riwayat kenaikan gaji berkala is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDeletedAt:    "null",
		},
		{
			name:             "error: riwayat kenaikan gaji berkala is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDeletedAt:    "null",
		},
		{
			name:             "error: riwayat kenaikan gaji berkala is deleted",
			paramNIP:         "1c",
			paramID:          "5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDeletedAt:    "2000-01-01",
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDeletedAt:    "null",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kenaikan-gaji-berkala/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
			actualRows, err := dbtest.QueryWithClause(db, "riwayat_kenaikan_gaji_berkala", "where id = $1 order by id", tt.paramID)
			assert.NoError(t, err)
			if len(actualRows) > 0 {
				deletedAt := actualRows[0]["deleted_at"]
				switch tt.wantDeletedAt {
				case "{deleted_at}":
					assert.WithinDuration(t, time.Now(), deletedAt.(time.Time), 10*time.Second)
				case "null":
					assert.Nil(t, deletedAt)
				default:
					assert.Equal(t, tt.wantDeletedAt, deletedAt.(time.Time).Format(time.DateOnly))
				}
			}
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
		insert into ref_golongan (id, nama, nama_pangkat, created_at) values
			(1, 'III/a', 'Penata Muda', '2000-01-01');
		insert into ref_unit_kerja (id, nama_unor, created_at) values
			(1, 'Unit Kerja 1', '2000-01-01');
		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, file_base64, created_at, updated_at) values
			(1,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/123/2023', '2023-01-15', '2023-01-15', 2, 6, '2023-01-15', '3500000', 3500000, 'data:abc', '2000-01-01', '2000-01-01'),
			(2,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/124/2024', '2024-01-15', '2024-01-15', 3, 0, '2024-01-15', '4000000', 4000000, 'data:abc', '2000-01-01', '2000-01-01');
		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, created_at, updated_at) values
			(3,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/125/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', '3000000', 3000000, '2000-01-01', '2000-01-01');
		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, created_at, updated_at, deleted_at) values
			(4,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/126/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', '3000000', 3000000, '2000-01-01', '2000-01-01', '2000-01-01');
		insert into riwayat_kenaikan_gaji_berkala
			(id, pegawai_id, golongan_id, no_sk, tanggal_sk, tmt_golongan, masa_kerja_golongan_tahun, masa_kerja_golongan_bulan, tmt_sk, n_gapok, gaji_pokok, file_base64, created_at, updated_at) values
			(5,  (select id from pegawai where nip_baru = '1c'), 1, 'SK/127/2020', '2020-01-15', '2020-01-15', 1, 0, '2020-01-15', '3000000', 3000000, 'data:abc', '2000-01-01', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                        int64(5),
			"pegawai_id":                int32(1),
			"golongan_id":               int32(1),
			"no_sk":                     "SK/127/2020",
			"tanggal_sk":                time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			"tmt_golongan":              time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			"masa_kerja_golongan_tahun": int16(1),
			"masa_kerja_golongan_bulan": int16(0),
			"tmt_sk":                    time.Date(2020, 1, 15, 0, 0, 0, 0, time.UTC),
			"gaji_pokok":                int32(3000000),
			"file_base64":               "data:abc",
			"created_at":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"updated_at":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"deleted_at":                nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
					"id":                        int64(1),
					"pegawai_id":                int32(1),
					"golongan_id":               int32(1),
					"no_sk":                     "SK/123/2023",
					"tanggal_sk":                time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":              time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"masa_kerja_golongan_tahun": int16(2),
					"masa_kerja_golongan_bulan": int16(6),
					"tmt_sk":                    time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"gaji_pokok":                int32(3500000),
					"file_base64":               "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"created_at":                time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"updated_at":                "{updated_at}",
					"deleted_at":                nil,
				},
			},
		},
		{
			name:              "error: riwayat kenaikan gaji berkala is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat kenaikan gaji berkala is deleted",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: missing file",
			paramNIP:          "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"file\" harus diisi"}`,
			wantDBRows:        defaultRows,
		},
		{
			name:              "error: invalid id",
			paramNIP:          "1c",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: invalid token",
			paramNIP:          "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBody:  `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:        dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			require.NoError(t, tt.appendRequestBody(writer))
			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kenaikan-gaji-berkala/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if len(tt.wantDBRows) > 0 {
				actualRows, err := dbtest.QueryWithClause(db, "riwayat_kenaikan_gaji_berkala", "where id = $1 order by id", tt.paramID)
				require.NoError(t, err)
				for i, row := range actualRows {
					if tt.wantDBRows[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
				assert.Equal(t, tt.wantDBRows[0]["file_base64"], actualRows[0]["file_base64"])
			}
		})
	}
}
