package riwayatpelatihansiasn

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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(t.Context(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	dbRepository := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, dbRepository, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("01")}
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
						"id": 1,
						"institusi_penyelenggara": "Lembaga Administrasi Negara",
						"jenis_diklat_id": 1,
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
						"jenis_diklat_id": 3,
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
						"jenis_diklat_id": 2,
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
						"jenis_diklat_id": 1,
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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 3,
						"institusi_penyelenggara": "LAN-RI",
						"jenis_diklat_id": 3,
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("20")}},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-siasn", nil)
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
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := repo.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			name:              "error: base64 pelatihan siasn tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pelatihan siasn sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan siasn berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan siasn berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan siasn bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan siasn tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan siasn tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-siasn/%s/berkas", tt.paramID), nil)
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
	INSERT INTO ref_jenis_diklat (id, jenis_diklat, kode, deleted_at) VALUES
		(1, 'Jenis 1', '01', null),
		(2, 'Jenis 2', '02', null),
		(3, 'Jenis 3', '03', '2000-01-01');

	INSERT INTO riwayat_diklat (
		id, nip_baru, jenis_diklat_id, jenis_diklat, nama_diklat, no_sertifikat, tanggal_mulai, tanggal_selesai, tahun_diklat, durasi_jam, institusi_penyelenggara, deleted_at
	) VALUES
		(1, '1c', 1, 'jenis 1', 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00123', '2023-06-20', '2023-06-21', 2023, 120, 'Lembaga Administrasi Negara', null),
		(2, '1c', 2, 'jenis 2', 'Pelatihan Kepemimpinan Pengawas (PKP)', 'LAN-PKP-2022-00456', '2022-08-15', '2022-08-16', null, null, 'Badan Diklat Provinsi Jawa Barat', null),
		(3, '1c', 3, 'jenis 3', 'Pelatihan Kepemimpinan Nasional Tingkat II', 'LAN-PKNII-2021-00089', '2021-04-10', '2023-04-11', 2021, 12, 'LAN-RI', null),
		(4, '2c', 1, 'jenis 1', 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00234', '2023-07-05', '2023-07-6', 2023, 12, 'Badan Pengembangan Sumber Daya Manusia Daerah (BPSDMD) DKI Jakarta', null),
		(5, '1c', 1, 'jenis 1', 'Pelatihan Kepemimpinan Nasional Tingkat I', 'LAN-PKNI-2020-00077', '2020-09-12', '2020-09-13', 2020, 12, 'Lembaga Administrasi Negara', '2000-01-01'),
		(6, '1c', 1, 'jenis 1', 'Pelatihan Kepemimpinan Nasional Tingkat III', 'LAN-PKNI-2020-00077', null, null, 2022, 10, 'Lembaga Administrasi Negara', null),
		(7, '1d', 1, 'jenis 1', 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00234', '2023-07-05', '2023-07-6', 2023, 12, 'Badan Pengembangan Sumber Daya Manusia Daerah (BPSDMD) DKI Jakarta', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			name:             "ok: admin dapat melihat riwayat pelatihan siasn pegawai 1c",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 1,
						"institusi_penyelenggara": "Lembaga Administrasi Negara",
						"jenis_diklat_id": 1,
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
						"jenis_diklat_id": 3,
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
						"jenis_diklat_id": 2,
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
						"jenis_diklat_id": 1,
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
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan siasn pegawai 1c dengan pagination",
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 3,
						"institusi_penyelenggara": "LAN-RI",
						"jenis_diklat_id": 3,
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
						"jenis_diklat_id": 2,
						"jenis_diklat": "Jenis 2",
						"nama_diklat": "Pelatihan Kepemimpinan Pengawas (PKP)",
						"nomor_sertifikat": "LAN-PKP-2022-00456",
						"tahun": 2022,
						"tanggal_mulai": "2022-08-15",
						"tanggal_selesai": "2022-08-16",
						"durasi": null
					}
				],
				"meta": {"limit": 2, "offset": 1, "total": 4}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan siasn pegawai 1d",
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 7,
						"institusi_penyelenggara": "Badan Pengembangan Sumber Daya Manusia Daerah (BPSDMD) DKI Jakarta",
						"jenis_diklat_id": 1,
						"jenis_diklat": "Jenis 1",
						"nama_diklat": "Pelatihan Kepemimpinan Administrator (PKA)",
						"nomor_sertifikat": "LAN-PKA-2023-00234",
						"tahun": 2023,
						"tanggal_mulai": "2023-07-05",
						"tanggal_selesai": "2023-07-06",
						"durasi": 12
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan siasn pegawai yang tidak ada data",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pelatihan-siasn", nil)
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
			('id_1f', '1f',     null);
		insert into ref_jenis_diklat
			(id,  jenis_diklat, deleted_at) values
			('1', 'SIASN 1',    null),
			('2', 'SIASN 2',    '2000-01-01'),
			('3', 'SIASN 3',    null);
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
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 1",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": 2000,
				"durasi": 5
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      "{id}",
					"nama_diklat":             "Diklat 1",
					"jenis_diklat_id":         int16(1),
					"jenis_diklat":            "SIASN 1",
					"institusi_penyelenggara": "ITB",
					"no_sertifikat":           "SK.01",
					"tanggal_mulai":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":         time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
					"tahun_diklat":            int32(2000),
					"durasi_jam":              int32(5),
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              "{created_at}",
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 3,
				"nama_diklat": "",
				"institusi_penyelenggara": "",
				"nomor_sertifikat": "",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": null,
				"durasi": 0
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      "{id}",
					"nama_diklat":             "",
					"jenis_diklat_id":         int16(3),
					"jenis_diklat":            "SIASN 3",
					"institusi_penyelenggara": "",
					"no_sertifikat":           "",
					"tanggal_mulai":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":         time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
					"tahun_diklat":            nil,
					"durasi_jam":              int32(0),
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1e",
					"nip_baru":                "1e",
					"created_at":              "{created_at}",
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"durasi": 5
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      "{id}",
					"nama_diklat":             "Diklat 2",
					"jenis_diklat_id":         int16(1),
					"jenis_diklat":            "SIASN 1",
					"institusi_penyelenggara": "ITB",
					"no_sertifikat":           "SK.01",
					"tanggal_mulai":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":         time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
					"tahun_diklat":            nil,
					"durasi_jam":              int32(5),
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1f",
					"nip_baru":                "1f",
					"created_at":              "{created_at}",
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"durasi": 5
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
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": 0,
				"durasi": 5
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: jenis diklat is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 0,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"durasi": 5
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis diklat tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: jenis diklat is deleted",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 2,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": null,
				"durasi": 5
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis diklat tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: exceed length limit, unexpected date or data type",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": "0",
				"nama_diklat": "` + strings.Repeat(".", 701) + `",
				"institusi_penyelenggara": "` + strings.Repeat(".", 601) + `",
				"nomor_sertifikat": "` + strings.Repeat(".", 601) + `",
				"tanggal_mulai": "",
				"tanggal_selesai": "",
				"tahun": "0",
				"durasi": "0"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"durasi\" harus dalam tipe integer` +
				` | parameter \"institusi_penyelenggara\" harus 600 karakter atau kurang` +
				` | parameter \"jenis_diklat_id\" harus dalam tipe integer` +
				` | parameter \"nama_diklat\" harus 700 karakter atau kurang` +
				` | parameter \"nomor_sertifikat\" harus 600 karakter atau kurang` +
				` | parameter \"tahun\" harus dalam tipe integer` +
				` | parameter \"tanggal_mulai\" harus dalam format date` +
				` | parameter \"tanggal_selesai\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null params",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": null,
				"nama_diklat": null,
				"institusi_penyelenggara": null,
				"nomor_sertifikat": null,
				"tanggal_mulai": null,
				"tanggal_selesai": null,
				"tahun": null,
				"durasi": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"durasi\" tidak boleh null` +
				` | parameter \"institusi_penyelenggara\" tidak boleh null` +
				` | parameter \"jenis_diklat_id\" tidak boleh null` +
				` | parameter \"nama_diklat\" tidak boleh null` +
				` | parameter \"nomor_sertifikat\" tidak boleh null` +
				` | parameter \"tanggal_mulai\" tidak boleh null` +
				` | parameter \"tanggal_selesai\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"jenis_diklat_id\" harus diisi` +
				` | parameter \"nama_diklat\" harus diisi` +
				` | parameter \"institusi_penyelenggara\" harus diisi` +
				` | parameter \"nomor_sertifikat\" harus diisi` +
				` | parameter \"tanggal_mulai\" harus diisi` +
				` | parameter \"tanggal_selesai\" harus diisi` +
				` | parameter \"durasi\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-siasn", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat", "where nip_baru = $1", tt.paramNIP)
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
		insert into ref_jenis_diklat
			(id,  jenis_diklat, deleted_at) values
			('1', 'SIASN 1',    null),
			('2', 'SIASN 2',    '2000-01-01'),
			('3', 'SIASN 3',    null);
		insert into riwayat_diklat
			(id,  rumpun_diklat_id, rumpun_diklat_nama, bkn_id, diklat_struktural_id, file_base64, pns_orang_id,  nip_baru, created_at,   updated_at) values
			('1', '1abc',           'rumpun',           'bkn1', '2a',                 'data:abc',  'id_1c',       '1c',     '2000-01-01', '2000-01-01'),
			('2', '1abc',           'rumpun',           'bkn1', '2a',                 'data:abc',  'id_1c',       '1c',     '2000-01-01', '2000-01-01'),
			('3', '1abc',           'rumpun',           'bkn1', '2a',                 'data:abc',  'id_1c',       '1c',     '2000-01-01', '2000-01-01');
		insert into riwayat_diklat
			(id,  nama_diklat, pns_orang_id, nip_baru, sudah_kirim_siasn, created_at,   updated_at,   deleted_at) values
			('4', 'Diklat 4',  'id_1e',      '1e',     null,              '2000-01-01', '2000-01-01', null),
			('5', 'Diklat 5',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', '2000-01-01'),
			('6', 'Diklat 6',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                      int64(6),
			"nama_diklat":             "Diklat 6",
			"jenis_diklat_id":         nil,
			"jenis_diklat":            nil,
			"institusi_penyelenggara": nil,
			"no_sertifikat":           nil,
			"tanggal_mulai":           nil,
			"tanggal_selesai":         nil,
			"tahun_diklat":            nil,
			"durasi_jam":              nil,
			"diklat_struktural_id":    nil,
			"file_base64":             nil,
			"rumpun_diklat_nama":      nil,
			"rumpun_diklat_id":        nil,
			"sudah_kirim_siasn":       nil,
			"bkn_id":                  nil,
			"pns_orang_id":            "id_1c",
			"nip_baru":                "1c",
			"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":              nil,
		},
	}

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
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 1",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": 2000,
				"durasi": 5
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(1),
					"nama_diklat":             "Diklat 1",
					"jenis_diklat_id":         int16(1),
					"jenis_diklat":            "SIASN 1",
					"institusi_penyelenggara": "ITB",
					"no_sertifikat":           "SK.01",
					"tanggal_mulai":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":         time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
					"tahun_diklat":            int32(2000),
					"durasi_jam":              int32(5),
					"diklat_struktural_id":    "2a",
					"file_base64":             "data:abc",
					"rumpun_diklat_nama":      "rumpun",
					"rumpun_diklat_id":        "1abc",
					"sudah_kirim_siasn":       "belum",
					"bkn_id":                  "bkn1",
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 3,
				"nama_diklat": "",
				"institusi_penyelenggara": "",
				"nomor_sertifikat": "",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": null,
				"durasi": 0
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(2),
					"nama_diklat":             "",
					"jenis_diklat_id":         int16(3),
					"jenis_diklat":            "SIASN 3",
					"institusi_penyelenggara": "",
					"no_sertifikat":           "",
					"tanggal_mulai":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":         time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
					"tahun_diklat":            nil,
					"durasi_jam":              int32(0),
					"diklat_struktural_id":    "2a",
					"file_base64":             "data:abc",
					"rumpun_diklat_nama":      "rumpun",
					"rumpun_diklat_id":        "1abc",
					"sudah_kirim_siasn":       "belum",
					"bkn_id":                  "bkn1",
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"durasi": 5
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(3),
					"nama_diklat":             "Diklat 2",
					"jenis_diklat_id":         int16(1),
					"jenis_diklat":            "SIASN 1",
					"institusi_penyelenggara": "ITB",
					"no_sertifikat":           "SK.01",
					"tanggal_mulai":           time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_selesai":         time.Date(2000, 1, 5, 0, 0, 0, 0, time.UTC),
					"tahun_diklat":            nil,
					"durasi_jam":              int32(5),
					"diklat_struktural_id":    "2a",
					"file_base64":             "data:abc",
					"rumpun_diklat_nama":      "rumpun",
					"rumpun_diklat_id":        "1abc",
					"sudah_kirim_siasn":       "belum",
					"bkn_id":                  "bkn1",
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "error: riwayat pelatihan siasn is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"durasi": 5
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat pelatihan siasn is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": 0,
				"durasi": 0
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(4),
					"nama_diklat":             "Diklat 4",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1e",
					"nip_baru":                "1e",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              nil,
				},
			},
		},
		{
			name:          "error: riwayat pelatihan siasn is deleted",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 1,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": 0,
				"durasi": 5
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(5),
					"nama_diklat":             "Diklat 5",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: jenis diklat is not found",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 0,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"durasi": 5
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis diklat tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: jenis diklat is deleted",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": 2,
				"nama_diklat": "Diklat 2",
				"institusi_penyelenggara": "ITB",
				"nomor_sertifikat": "SK.01",
				"tanggal_mulai": "2000-01-01",
				"tanggal_selesai": "2000-01-05",
				"tahun": null,
				"durasi": 5
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data jenis diklat tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": "0",
				"nama_diklat": "` + strings.Repeat(".", 701) + `",
				"institusi_penyelenggara": "` + strings.Repeat(".", 601) + `",
				"nomor_sertifikat": "` + strings.Repeat(".", 601) + `",
				"tanggal_mulai": "",
				"tanggal_selesai": "",
				"tahun": "0",
				"durasi": "0"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"durasi\" harus dalam tipe integer` +
				` | parameter \"institusi_penyelenggara\" harus 600 karakter atau kurang` +
				` | parameter \"jenis_diklat_id\" harus dalam tipe integer` +
				` | parameter \"nama_diklat\" harus 700 karakter atau kurang` +
				` | parameter \"nomor_sertifikat\" harus 600 karakter atau kurang` +
				` | parameter \"tahun\" harus dalam tipe integer` +
				` | parameter \"tanggal_mulai\" harus dalam format date` +
				` | parameter \"tanggal_selesai\" harus dalam format date"}`,
			wantDBRows: defaultRows,
		},
		{
			name:          "error: null params",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_diklat_id": null,
				"nama_diklat": null,
				"institusi_penyelenggara": null,
				"nomor_sertifikat": null,
				"tanggal_mulai": null,
				"tanggal_selesai": null,
				"tahun": null,
				"durasi": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"durasi\" tidak boleh null` +
				` | parameter \"institusi_penyelenggara\" tidak boleh null` +
				` | parameter \"jenis_diklat_id\" tidak boleh null` +
				` | parameter \"nama_diklat\" tidak boleh null` +
				` | parameter \"nomor_sertifikat\" tidak boleh null` +
				` | parameter \"tanggal_mulai\" tidak boleh null` +
				` | parameter \"tanggal_selesai\" tidak boleh null"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"jenis_diklat_id\" harus diisi` +
				` | parameter \"nama_diklat\" harus diisi` +
				` | parameter \"institusi_penyelenggara\" harus diisi` +
				` | parameter \"nomor_sertifikat\" harus diisi` +
				` | parameter \"tanggal_mulai\" harus diisi` +
				` | parameter \"tanggal_selesai\" harus diisi` +
				` | parameter \"durasi\" harus diisi"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "6",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-siasn/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat", "where id = $1", tt.paramID)
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
		insert into riwayat_diklat
			(id,  nama_diklat, pns_orang_id, nip_baru, sudah_kirim_siasn, created_at,   updated_at,   deleted_at) values
			('1', 'Diklat 1',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', null),
			('2', 'Diklat 2',  'id_1e',      '1e',     null,              '2000-01-01', '2000-01-01', null),
			('3', 'Diklat 3',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Diklat 4',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                      int64(4),
			"nama_diklat":             "Diklat 4",
			"jenis_diklat_id":         nil,
			"jenis_diklat":            nil,
			"institusi_penyelenggara": nil,
			"no_sertifikat":           nil,
			"tanggal_mulai":           nil,
			"tanggal_selesai":         nil,
			"tahun_diklat":            nil,
			"durasi_jam":              nil,
			"diklat_struktural_id":    nil,
			"file_base64":             nil,
			"rumpun_diklat_nama":      nil,
			"rumpun_diklat_id":        nil,
			"sudah_kirim_siasn":       nil,
			"bkn_id":                  nil,
			"pns_orang_id":            "id_1c",
			"nip_baru":                "1c",
			"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":              nil,
		},
	}

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
					"id":                      int64(1),
					"nama_diklat":             "Diklat 1",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat pelatihan siasn is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(2),
					"nama_diklat":             "Diklat 2",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1e",
					"nip_baru":                "1e",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              nil,
				},
			},
		},
		{
			name:             "error: riwayat pelatihan siasn is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat pelatihan siasn is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(3),
					"nama_diklat":             "Diklat 3",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
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

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-siasn/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat", "where id = $1", tt.paramID)
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
		insert into riwayat_diklat
			(id,  nama_diklat, pns_orang_id, nip_baru, sudah_kirim_siasn, created_at,   updated_at,   deleted_at) values
			('1', 'Diklat 1',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', null),
			('2', 'Diklat 2',  'id_1e',      '1e',     null,              '2000-01-01', '2000-01-01', null),
			('3', 'Diklat 3',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Diklat 4',  'id_1c',      '1c',     null,              '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                      int64(4),
			"nama_diklat":             "Diklat 4",
			"jenis_diklat_id":         nil,
			"jenis_diklat":            nil,
			"institusi_penyelenggara": nil,
			"no_sertifikat":           nil,
			"tanggal_mulai":           nil,
			"tanggal_selesai":         nil,
			"tahun_diklat":            nil,
			"durasi_jam":              nil,
			"diklat_struktural_id":    nil,
			"file_base64":             nil,
			"rumpun_diklat_nama":      nil,
			"rumpun_diklat_id":        nil,
			"sudah_kirim_siasn":       nil,
			"bkn_id":                  nil,
			"pns_orang_id":            "id_1c",
			"nip_baru":                "1c",
			"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":              nil,
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
					"id":                      int64(1),
					"nama_diklat":             "Diklat 1",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              "{updated_at}",
					"deleted_at":              nil,
				},
			},
		},
		{
			name:              "error: riwayat pelatihan siasn is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat pelatihan siasn is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(2),
					"nama_diklat":             "Diklat 2",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1e",
					"nip_baru":                "1e",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              nil,
				},
			},
		},
		{
			name:              "error: riwayat pelatihan siasn is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                      int64(3),
					"nama_diklat":             "Diklat 3",
					"jenis_diklat_id":         nil,
					"jenis_diklat":            nil,
					"institusi_penyelenggara": nil,
					"no_sertifikat":           nil,
					"tanggal_mulai":           nil,
					"tanggal_selesai":         nil,
					"tahun_diklat":            nil,
					"durasi_jam":              nil,
					"diklat_struktural_id":    nil,
					"file_base64":             nil,
					"rumpun_diklat_nama":      nil,
					"rumpun_diklat_id":        nil,
					"sudah_kirim_siasn":       nil,
					"bkn_id":                  nil,
					"pns_orang_id":            "id_1c",
					"nip_baru":                "1c",
					"created_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":              time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-pelatihan-siasn/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_diklat", "where id = $1", tt.paramID)
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
