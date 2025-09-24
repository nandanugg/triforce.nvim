package riwayatpelatihansiasn

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
	INSERT INTO ref_jenis_diklat (id, jenis_diklat, kode, deleted_at) VALUES
		(1, 'Jenis 1', '01', null),
		(2, 'Jenis 2', '02', null),
		(3, 'Jenis 3', '03', '2000-01-01');

	INSERT INTO riwayat_diklat (
		id, nip_baru, jenis_diklat_id, jenis_diklat, nama_diklat, no_sertifikat, tanggal_mulai, tanggal_selesai, tahun_diklat, durasi_jam, institusi_penyelenggara, deleted_at
	) VALUES
		(1, '01', 1, 'jenis 1', 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00123', '2023-06-20', '2023-06-21', 2023, 120, 'Lembaga Administrasi Negara', null),
		(2, '01', 2, 'jenis 2', 'Pelatihan Kepemimpinan Pengawas (PKP)', 'LAN-PKP-2022-00456', '2022-08-15', '2022-08-16', null, null, 'Badan Diklat Provinsi Jawa Barat', null),
		(3, '01', 3, 'jenis 3', 'Pelatihan Kepemimpinan Nasional Tingkat II', 'LAN-PKNII-2021-00089', '2021-04-10', '2023-04-11', 2021, 12, 'LAN-RI', null),
		(4, '02', 1, 'jenis 1', 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00234', '2023-07-05', '2023-07-6', 2023, 12, 'Badan Pengembangan Sumber Daya Manusia Daerah (BPSDMD) DKI Jakarta', null),
		(5, '01', 1, 'jenis 1', 'Pelatihan Kepemimpinan Nasional Tingkat I', 'LAN-PKNI-2020-00077', '2020-09-12', '2020-09-13', 2020, 12, 'Lembaga Administrasi Negara', '2000-01-01'),
		(6, '01', 1, 'jenis 1', 'Pelatihan Kepemimpinan Nasional Tingkat III', 'LAN-PKNI-2020-00077', null, null, 2022, 10, 'Lembaga Administrasi Negara', null);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "01")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 1,
						"institusi_penyelenggara": "Lembaga Administrasi Negara",
						"jenis_diklat": "Jenis 1",
						"nama_diklat": "Pelatihan Kepemimpinan Administrator (PKA)",
						"nomor_sertifikat": "LAN-PKA-2023-00123",
						"tahun": 2023,
						"tanggal_mulai": "2023-06-20",
						"tanggal_selesai": "2023-06-21",
						"durasi": 120
					},
					{
						"id": 3,
						"institusi_penyelenggara": "LAN-RI",
						"jenis_diklat": "",
						"nama_diklat": "Pelatihan Kepemimpinan Nasional Tingkat II",
						"nomor_sertifikat": "LAN-PKNII-2021-00089",
						"tahun": 2021,
						"tanggal_mulai": "2021-04-10",
						"tanggal_selesai": "2023-04-11",
						"durasi": 12
					},
					{
						"id": 2,
						"institusi_penyelenggara": "Badan Diklat Provinsi Jawa Barat",
						"jenis_diklat": "Jenis 2",
						"nama_diklat": "Pelatihan Kepemimpinan Pengawas (PKP)",
						"nomor_sertifikat": "LAN-PKP-2022-00456",
						"tahun": 2022,
						"tanggal_mulai": "2022-08-15",
						"tanggal_selesai": "2022-08-16",
						"durasi": null
					},
					{
						"id": 6,
						"institusi_penyelenggara": "Lembaga Administrasi Negara",
						"jenis_diklat": "Jenis 1",
						"nama_diklat": "Pelatihan Kepemimpinan Nasional Tingkat III",
						"nomor_sertifikat": "LAN-PKNI-2020-00077",
						"tahun": 2022,
						"tanggal_mulai": null,
						"tanggal_selesai": null,
						"durasi": 10
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}
			`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "01")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 3,
						"institusi_penyelenggara": "LAN-RI",
						"jenis_diklat": "",
						"nama_diklat": "Pelatihan Kepemimpinan Nasional Tingkat II",
						"nomor_sertifikat": "LAN-PKNII-2021-00089",
						"tahun": 2021,
						"tanggal_mulai": "2021-04-10",
						"tanggal_selesai": "2023-04-11",
						"durasi": 12
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}
			`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "20")}},
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
			dbRepository := repo.New(db)
			_, err := db.Exec(t.Context(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-siasn", nil)
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
		insert into riwayat_diklat
			(id, nip_baru, deleted_at,   file_base64) values
			(1,  '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2,  '1c',     null,         '` + pdfBase64 + `'),
			(3,  '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			(4,  '1c',     null,         'data:application/pdf;base64,invalid'),
			(5,  '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6,  '1c',     null,         null),
			(7,  '1c',     null,         '');
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
			name:              "error: base64 pelatihan siasn tidak valid",
			dbData:            dbData,
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pelatihan siasn sudah dihapus",
			dbData:            dbData,
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan siasn berisi null value",
			dbData:            dbData,
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan siasn berupa string kosong",
			dbData:            dbData,
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan siasn bukan milik user login",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan siasn tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-siasn/%s/berkas", tt.paramID), nil)
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
