package riwayatkenaikangajiberkala

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
