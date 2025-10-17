package riwayatpelatihanfungsional

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
		INSERT INTO riwayat_diklat_fungsional ( id, nip_baru, nip_lama, jenis_diklat, nama_kursus, jumlah_jam, tahun, institusi_penyelenggara, jenis_kursus_sertifikat, no_sertifikat, instansi, status_data, tanggal_kursus, keterangan_berkas, lama, siasn_id, created_at, updated_at, deleted_at) VALUES
		( '9159a38e-586c-4338-8915-ab9ee5f6a937', '198706102020121001', '123456789', 'Teknis', 'Pelatihan Dasar CPNS', 120, 2020, 'LAN RI', 'A', 'LATSAR-2020-001', 'Kementerian Dalam Negeri', 'valid', '2020-03-15', 'Sertifikat Asli', 1.5, gen_random_uuid(), now(), now(), null),
		( 'add68305-1993-40e5-9e32-14de553cfd73', '198706102020121001', '123456789', 'Teknis II', 'Pelatihan Dasar CPNS II', 120, 2020, 'LAN RI', 'A', 'LATSAR-2020-002', 'Kementerian Dalam Negeri', 'valid', '2020-03-15', 'Sertifikat Asli', 1.5, gen_random_uuid(), now(), now(),null),
		( gen_random_uuid(), '198906152019031002', '987654321', 'Fungsional', 'Diklat Auditor Pertama', 80, 2021, 'BPKP', 'B', 'AUD-2021-014', 'Badan Pengawas Keuangan', 'valid', '2021-07-20', 'Softcopy', 1.0, gen_random_uuid(), now(), now(),null),
		( gen_random_uuid(), '199001052018021003', '112233445', 'Kepemimpinan', 'Diklatpim III', 200, 2022, 'LAN RI', 'C', 'PIM3-2022-045', 'Kementerian Keuangan', 'valid', '2022-09-10', 'Sudah dilegalisir', 2.0, gen_random_uuid(), now(), now(),null),
		( gen_random_uuid(), '199202202017011004', '223344556', 'Teknis', 'Pelatihan Manajemen Proyek IT', 60, 2023, 'Kemenkominfo', 'A', 'IT-2023-089', 'Kementerian Komunikasi & Informatika', 'valid', '2023-04-12', 'Discan warna', 0.5, gen_random_uuid(), now(), now(),null),
		( gen_random_uuid(), '199305302021041005', '334455667', 'Fungsional', 'Pelatihan Arsiparis Ahli', 100, 2024, 'ANRI', 'B', 'ARSIP-2024-122', 'Arsip Nasional RI', 'valid', '2024-06-18', 'Tersimpan di HRD', 1.25, gen_random_uuid(), now(), now(),null),
		( gen_random_uuid(), '199305302021041005', '334455667', 'Fungsional', 'Pelatihan Arsiparis Ahli Deleted', 100, 2024, 'ANRI', 'B', 'ARSIP-2024-122', 'Arsip Nasional RI', 'valid', '2024-06-18', 'Tersimpan di HRD', 1.25, gen_random_uuid(), now(), now(),now());
	`
	pgxConn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxConn.Exec(t.Context(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("198706102020121001")}
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
						"id": "9159a38e-586c-4338-8915-ab9ee5f6a937",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Teknis",
						"nama_diklat": "Pelatihan Dasar CPNS",
						"nomor_sertifikat": "LATSAR-2020-001",
						"tahun": 2020,
						"durasi": 120,
						"tanggal_mulai": "2020-03-15",
						"tanggal_selesai": "2020-03-20"
					},
					{
						"id": "add68305-1993-40e5-9e32-14de553cfd73",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Teknis II",
						"nama_diklat": "Pelatihan Dasar CPNS II",
						"nomor_sertifikat": "LATSAR-2020-002",
						"tahun": 2020,
						"durasi": 120,
						"tanggal_mulai": "2020-03-15",
						"tanggal_selesai": "2020-03-20"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 2 }
			}
			`,
		},
		{
			name:             "ok: dengan parameter pagination",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "add68305-1993-40e5-9e32-14de553cfd73",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Teknis II",
						"nama_diklat": "Pelatihan Dasar CPNS II",
						"nomor_sertifikat": "LATSAR-2020-002",
						"tahun": 2020,
						"durasi": 120,
						"tanggal_mulai": "2020-03-15",
						"tanggal_selesai": "2020-03-20"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 2}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": { "limit": 10, "offset": 0, "total": 0 } }`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-fungsional", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			dbRepository := repo.New(pgxConn)
			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
			RegisterRoutes(e, dbRepository, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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
		insert into riwayat_diklat_fungsional
			(id,  nip_baru, deleted_at,   file_base64) values
			('a', '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			('2', '1c',     null,         '` + pdfBase64 + `'),
			('3', '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			('4', '1c',     null,         'data:application/pdf;base64,invalid'),
			('5', '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			('6', '1c',     null,         null),
			('7', '1c',     null,         '');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

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
			paramID:           "a",
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
			name:              "error: base64 pelatihan fungsional tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pelatihan fungsional sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan fungsional berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan fungsional berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan fungsional bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan fungsional tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-fungsional/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repo.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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
		INSERT INTO riwayat_diklat_fungsional ( id, nip_baru, nip_lama, jenis_diklat, nama_kursus, jumlah_jam, tahun, institusi_penyelenggara, jenis_kursus_sertifikat, no_sertifikat, instansi, status_data, tanggal_kursus, keterangan_berkas, lama, siasn_id, created_at, updated_at, deleted_at) VALUES
		( '9159a38e-586c-4338-8915-ab9ee5f6a937', '1c', '123456789', 'Teknis', 'Pelatihan Dasar CPNS', 120, 2020, 'LAN RI', 'A', 'LATSAR-2020-001', 'Kementerian Dalam Negeri', 'valid', '2020-03-15', 'Sertifikat Asli', 1.5, '00000000-0000-0000-0000-000000000001', now(), now(), null),
		( 'add68305-1993-40e5-9e32-14de553cfd73', '1c', '123456789', 'Teknis II', 'Pelatihan Dasar CPNS II', 120, 2020, 'LAN RI', 'A', 'LATSAR-2020-002', 'Kementerian Dalam Negeri', 'valid', '2020-03-15', 'Sertifikat Asli', 1.5, '00000000-0000-0000-0000-000000000002', now(), now(),null),
		( '00000000-0000-0000-0000-000000000003', '2c', '987654321', 'Fungsional', 'Diklat Auditor Pertama', 80, 2021, 'BPKP', 'B', 'AUD-2021-014', 'Badan Pengawas Keuangan', 'valid', '2021-07-20', 'Softcopy', 1.0, '00000000-0000-0000-0000-000000000004', now(), now(),null),
		( '00000000-0000-0000-0000-000000000005', '1c', '112233445', 'Kepemimpinan', 'Diklatpim III', 200, 2022, 'LAN RI', 'C', 'PIM3-2022-045', 'Kementerian Keuangan', 'valid', '2022-09-10', 'Sudah dilegalisir', 2.0, '00000000-0000-0000-0000-000000000006', now(), now(),null),
		( '00000000-0000-0000-0000-000000000007', '1c', '223344556', 'Teknis', 'Pelatihan Manajemen Proyek IT', 60, 2023, 'Kemenkominfo', 'A', 'IT-2023-089', 'Kementerian Komunikasi & Informatika', 'valid', '2023-04-12', 'Discan warna', 0.5, '00000000-0000-0000-0000-000000000008', now(), now(),null),
		( '00000000-0000-0000-0000-000000000009', '1c', '334455667', 'Fungsional', 'Pelatihan Arsiparis Ahli', 100, 2024, 'ANRI', 'B', 'ARSIP-2024-122', 'Arsip Nasional RI', 'valid', '2024-06-18', 'Tersimpan di HRD', 1.25, '00000000-0000-0000-0000-000000000010', now(), now(),null),
		( '00000000-0000-0000-0000-000000000011', '1c', '334455667', 'Fungsional', 'Pelatihan Arsiparis Ahli Deleted', 100, 2024, 'ANRI', 'B', 'ARSIP-2024-122', 'Arsip Nasional RI', 'valid', '2024-06-18', 'Tersimpan di HRD', 1.25, '00000000-0000-0000-0000-000000000012', now(), now(),now()),
		( '00000000-0000-0000-0000-000000000000', '1d', '445566778', 'Fungsional', 'Pelatihan Arsiparis Ahli', 100, 2024, 'ANRI', 'B', 'ARSIP-2024-123', 'Arsip Nasional RI', 'valid', '2024-06-18', 'Tersimpan di HRD', 1.25, '00000000-0000-0000-0000-000000000013', now(), now(),null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

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
			name:             "ok: admin dapat melihat riwayat pelatihan fungsional pegawai 1c",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "00000000-0000-0000-0000-000000000009",
						"institusi_penyelenggara": "ANRI",
						"jenis_diklat": "Fungsional",
						"nama_diklat": "Pelatihan Arsiparis Ahli",
						"nomor_sertifikat": "ARSIP-2024-122",
						"tahun": 2024,
						"durasi": 100,
						"tanggal_mulai": "2024-06-18",
						"tanggal_selesai": "2024-06-22"
					},
					{
						"id": "00000000-0000-0000-0000-000000000007",
						"institusi_penyelenggara": "Kemenkominfo",
						"jenis_diklat": "Teknis",
						"nama_diklat": "Pelatihan Manajemen Proyek IT",
						"nomor_sertifikat": "IT-2023-089",
						"tahun": 2023,
						"durasi": 60,
						"tanggal_mulai": "2023-04-12",
						"tanggal_selesai": "2023-04-14"
					},
					{
						"id": "00000000-0000-0000-0000-000000000005",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Kepemimpinan",
						"nama_diklat": "Diklatpim III",
						"nomor_sertifikat": "PIM3-2022-045",
						"tahun": 2022,
						"durasi": 200,
						"tanggal_mulai": "2022-09-10",
						"tanggal_selesai": "2022-09-18"
					},
					{
						"id": "9159a38e-586c-4338-8915-ab9ee5f6a937",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Teknis",
						"nama_diklat": "Pelatihan Dasar CPNS",
						"nomor_sertifikat": "LATSAR-2020-001",
						"tahun": 2020,
						"durasi": 120,
						"tanggal_mulai": "2020-03-15",
						"tanggal_selesai": "2020-03-20"
					},
					{
						"id": "add68305-1993-40e5-9e32-14de553cfd73",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Teknis II",
						"nama_diklat": "Pelatihan Dasar CPNS II",
						"nomor_sertifikat": "LATSAR-2020-002",
						"tahun": 2020,
						"durasi": 120,
						"tanggal_mulai": "2020-03-15",
						"tanggal_selesai": "2020-03-20"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 5 }
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan fungsional pegawai 1c dengan pagination",
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "00000000-0000-0000-0000-000000000007",
						"institusi_penyelenggara": "Kemenkominfo",
						"jenis_diklat": "Teknis",
						"nama_diklat": "Pelatihan Manajemen Proyek IT",
						"nomor_sertifikat": "IT-2023-089",
						"tahun": 2023,
						"durasi": 60,
						"tanggal_mulai": "2023-04-12",
						"tanggal_selesai": "2023-04-14"
					},
					{
						"id": "00000000-0000-0000-0000-000000000005",
						"institusi_penyelenggara": "LAN RI",
						"jenis_diklat": "Kepemimpinan",
						"nama_diklat": "Diklatpim III",
						"nomor_sertifikat": "PIM3-2022-045",
						"tahun": 2022,
						"durasi": 200,
						"tanggal_mulai": "2022-09-10",
						"tanggal_selesai": "2022-09-18"
					}
				],
				"meta": {"limit": 2, "offset": 1, "total": 5}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan fungsional pegawai 1d",
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "00000000-0000-0000-0000-000000000000",
						"institusi_penyelenggara": "ANRI",
						"jenis_diklat": "Fungsional",
						"nama_diklat": "Pelatihan Arsiparis Ahli",
						"nomor_sertifikat": "ARSIP-2024-123",
						"tahun": 2024,
						"durasi": 100,
						"tanggal_mulai": "2024-06-18",
						"tanggal_selesai": "2024-06-22"
					}
				],
				"meta": { "limit": 10, "offset": 0, "total": 1 }
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan fungsional pegawai yang tidak ada data",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pelatihan-fungsional", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
