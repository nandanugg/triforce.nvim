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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198706102020121001")}},
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
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198706102020121001")}},
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
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": { "limit": 10, "offset": 0, "total": 0 } }`,
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

			pgxConn := dbtest.New(t, dbmigrations.FS)
			dbRepository := repo.New(pgxConn)
			_, err := pgxConn.Exec(t.Context(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-fungsional", nil)
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
			paramID:           "a",
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
			name:              "error: base64 pelatihan fungsional tidak valid",
			dbData:            dbData,
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pelatihan fungsional sudah dihapus",
			dbData:            dbData,
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan fungsional berisi null value",
			dbData:            dbData,
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan fungsional berupa string kosong",
			dbData:            dbData,
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan fungsional bukan milik user login",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan fungsional tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan fungsional tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-fungsional/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repo.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
