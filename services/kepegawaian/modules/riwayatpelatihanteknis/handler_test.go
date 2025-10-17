package riwayatpelatihanteknis

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
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into riwayat_kursus
			(id, pns_nip, tipe_kursus, jenis_kursus, nama_kursus, tanggal_kursus, lama_kursus, institusi_penyelenggara, no_sertifikat, deleted_at) values
			(11, '1c', 'Teknis', 'Workshop', '11a', '2000-01-01', 24, 'Institution 11', 'CERT11', null),
			(12, '1c', 'Teknis', 'Seminar', '12a', '2001-01-01', 16, 'Institution 12', '', null),
			(13, '1c', 'Teknis', 'Kursus', '13a', '2002-01-01', 40, 'Institution 13', 'CERT13', null),
			(14, '2c', 'Teknis', 'Workshop', '14a', '2003-01-01', 8, 'Institution 14', 'CERT14', null),
			(15, '1c', 'Teknis', 'Seminar', '15a', '2004-01-01', 4, 'Institution 15', 'CERT15', null),
			(16, '1c', 'Teknis', 'Kursus', '16a', '2005-01-01', 20, 'Institution 16', null, null),
			(17, '1c', 'Teknis', 'Workshop', '17a', '2006-01-01', 12, 'Institution 17', 'CERT17', null),
			-- Null test cases
			(18, '1c', 'Teknis', 'Workshop', '18a', '2010-01-01', null, 'Institution 18', 'CERT18', null),
			(19, '1c', 'Teknis', 'Seminar', '19a', null, 8, 'Institution 19', 'CERT19', null),
			(20, '1c', null, null, '20a', '2012-01-01', 12, 'Institution 20', 'CERT20', null),
			(21, '1c', 'Teknis', 'Workshop', null, '2013-01-01', 16, 'Institution 21', 'CERT21', null),
			(22, '1c', 'Teknis', 'Seminar', '22a', '2014-01-01', 8, null, 'CERT22', null);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

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
						"id": 22,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Seminar",
						"nama_pelatihan": "22a",
						"tanggal_mulai": "2014-01-01",
						"tanggal_selesai": "2014-01-01",
						"tahun": 2014,
						"durasi": 8,
						"institusi_penyelenggara": "",
						"nomor_sertifikat": "CERT22"
					},
					{
						"id": 21,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "",
						"tanggal_mulai": "2013-01-01",
						"tanggal_selesai": "2013-01-01",
						"tahun": 2013,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 21",
						"nomor_sertifikat": "CERT21"
					},
					{
						"id": 20,
						"tipe_pelatihan": "",
						"jenis_pelatihan": "",
						"nama_pelatihan": "20a",
						"tanggal_mulai": "2012-01-01",
						"tanggal_selesai": "2012-01-01",
						"tahun": 2012,
						"durasi": 12,
						"institusi_penyelenggara": "Institution 20",
						"nomor_sertifikat": "CERT20"
					},
					{
						"id": 18,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "18a",
						"tanggal_mulai": "2010-01-01",
						"tanggal_selesai": "2010-01-01",
						"tahun": 2010,
						"durasi": null,
						"institusi_penyelenggara": "Institution 18",
						"nomor_sertifikat": "CERT18"
					},
					{
						"id": 17,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "17a",
						"tanggal_mulai": "2006-01-01",
						"tanggal_selesai": "2006-01-01",
						"tahun": 2006,
						"durasi": 12,
						"institusi_penyelenggara": "Institution 17",
						"nomor_sertifikat": "CERT17"
					},
					{
						"id": 16,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Kursus",
						"nama_pelatihan": "16a",
						"tanggal_mulai": "2005-01-01",
						"tanggal_selesai": "2005-01-01",
						"tahun": 2005,
						"durasi": 20,
						"institusi_penyelenggara": "Institution 16",
						"nomor_sertifikat": ""
					},
					{
						"id": 15,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Seminar",
						"nama_pelatihan": "15a",
						"tanggal_mulai": "2004-01-01",
						"tanggal_selesai": "2004-01-01",
						"tahun": 2004,
						"durasi": 4,
						"institusi_penyelenggara": "Institution 15",
						"nomor_sertifikat": "CERT15"
					},
					{
						"id": 13,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Kursus",
						"nama_pelatihan": "13a",
						"tanggal_mulai": "2002-01-01",
						"tanggal_selesai": "2002-01-02",
						"tahun": 2002,
						"durasi": 40,
						"institusi_penyelenggara": "Institution 13",
						"nomor_sertifikat": "CERT13"
					},
					{
						"id": 12,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Seminar",
						"nama_pelatihan": "12a",
						"tanggal_mulai": "2001-01-01",
						"tanggal_selesai": "2001-01-01",
						"tahun": 2001,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 12",
						"nomor_sertifikat": ""
					},
					{
						"id": 11,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "11a",
						"tanggal_mulai": "2000-01-01",
						"tanggal_selesai": "2000-01-02",
						"tahun": 2000,
						"durasi": 24,
						"institusi_penyelenggara": "Institution 11",
						"nomor_sertifikat": "CERT11"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 11}
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
						"id": 21,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "",
						"tanggal_mulai": "2013-01-01",
						"tanggal_selesai": "2013-01-01",
						"tahun": 2013,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 21",
						"nomor_sertifikat": "CERT21"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 11}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-teknis", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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
		insert into riwayat_kursus
			(id, pns_nip, deleted_at,   file_base64) values
			(1, '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c',     null,         '` + pdfBase64 + `'),
			(3, '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c',     null,         'data:application/pdf;base64,invalid'),
			(5, '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c',     null,         null),
			(7, '1c',     null,         '');
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
			name:              "error: base64 pelatihan teknis tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat pelatihan teknis sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan teknis tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan teknis berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan teknis tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat pelatihan teknis berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan teknis tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan teknis bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan teknis tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat pelatihan teknis tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat pelatihan teknis tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-pelatihan-teknis/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
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
		insert into riwayat_kursus
			(id, pns_nip, tipe_kursus, jenis_kursus, nama_kursus, tanggal_kursus, lama_kursus, institusi_penyelenggara, no_sertifikat, deleted_at) values
			(11, '1c', 'Teknis', 'Workshop', '11a', '2000-01-01', 24, 'Institution 11', 'CERT11', null),
			(12, '1c', 'Teknis', 'Seminar', '12a', '2001-01-01', 16, 'Institution 12', '', null),
			(13, '1c', 'Teknis', 'Kursus', '13a', '2002-01-01', 40, 'Institution 13', 'CERT13', null),
			(14, '2c', 'Teknis', 'Workshop', '14a', '2003-01-01', 8, 'Institution 14', 'CERT14', null),
			(15, '1c', 'Teknis', 'Seminar', '15a', '2004-01-01', 4, 'Institution 15', 'CERT15', null),
			(16, '1c', 'Teknis', 'Kursus', '16a', '2005-01-01', 20, 'Institution 16', null, null),
			(17, '1c', 'Teknis', 'Workshop', '17a', '2006-01-01', 12, 'Institution 17', 'CERT17', null),
			-- Null test cases
			(18, '1c', 'Teknis', 'Workshop', '18a', '2010-01-01', null, 'Institution 18', 'CERT18', null),
			(19, '1c', 'Teknis', 'Seminar', '19a', null, 8, 'Institution 19', 'CERT19', null),
			(20, '1c', null, null, '20a', '2012-01-01', 12, 'Institution 20', 'CERT20', null),
			(21, '1c', 'Teknis', 'Workshop', null, '2013-01-01', 16, 'Institution 21', 'CERT21', null),
			(22, '1c', 'Teknis', 'Seminar', '22a', '2014-01-01', 8, null, 'CERT22', null),
			(23, '1d', 'Teknis', 'Workshop', '23a', '2015-01-01', 16, 'Institution 23', 'CERT23', null);
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
			name:             "ok: admin dapat melihat riwayat pelatihan teknis pegawai 1c",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 22,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Seminar",
						"nama_pelatihan": "22a",
						"tanggal_mulai": "2014-01-01",
						"tanggal_selesai": "2014-01-01",
						"tahun": 2014,
						"durasi": 8,
						"institusi_penyelenggara": "",
						"nomor_sertifikat": "CERT22"
					},
					{
						"id": 21,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "",
						"tanggal_mulai": "2013-01-01",
						"tanggal_selesai": "2013-01-01",
						"tahun": 2013,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 21",
						"nomor_sertifikat": "CERT21"
					},
					{
						"id": 20,
						"tipe_pelatihan": "",
						"jenis_pelatihan": "",
						"nama_pelatihan": "20a",
						"tanggal_mulai": "2012-01-01",
						"tanggal_selesai": "2012-01-01",
						"tahun": 2012,
						"durasi": 12,
						"institusi_penyelenggara": "Institution 20",
						"nomor_sertifikat": "CERT20"
					},
					{
						"id": 18,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "18a",
						"tanggal_mulai": "2010-01-01",
						"tanggal_selesai": "2010-01-01",
						"tahun": 2010,
						"durasi": null,
						"institusi_penyelenggara": "Institution 18",
						"nomor_sertifikat": "CERT18"
					},
					{
						"id": 17,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "17a",
						"tanggal_mulai": "2006-01-01",
						"tanggal_selesai": "2006-01-01",
						"tahun": 2006,
						"durasi": 12,
						"institusi_penyelenggara": "Institution 17",
						"nomor_sertifikat": "CERT17"
					},
					{
						"id": 16,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Kursus",
						"nama_pelatihan": "16a",
						"tanggal_mulai": "2005-01-01",
						"tanggal_selesai": "2005-01-01",
						"tahun": 2005,
						"durasi": 20,
						"institusi_penyelenggara": "Institution 16",
						"nomor_sertifikat": ""
					},
					{
						"id": 15,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Seminar",
						"nama_pelatihan": "15a",
						"tanggal_mulai": "2004-01-01",
						"tanggal_selesai": "2004-01-01",
						"tahun": 2004,
						"durasi": 4,
						"institusi_penyelenggara": "Institution 15",
						"nomor_sertifikat": "CERT15"
					},
					{
						"id": 13,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Kursus",
						"nama_pelatihan": "13a",
						"tanggal_mulai": "2002-01-01",
						"tanggal_selesai": "2002-01-02",
						"tahun": 2002,
						"durasi": 40,
						"institusi_penyelenggara": "Institution 13",
						"nomor_sertifikat": "CERT13"
					},
					{
						"id": 12,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Seminar",
						"nama_pelatihan": "12a",
						"tanggal_mulai": "2001-01-01",
						"tanggal_selesai": "2001-01-01",
						"tahun": 2001,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 12",
						"nomor_sertifikat": ""
					},
					{
						"id": 11,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "11a",
						"tanggal_mulai": "2000-01-01",
						"tanggal_selesai": "2000-01-02",
						"tahun": 2000,
						"durasi": 24,
						"institusi_penyelenggara": "Institution 11",
						"nomor_sertifikat": "CERT11"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 11}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan teknis pegawai 1c dengan pagination",
			nip:              "1c",
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 21,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "",
						"tanggal_mulai": "2013-01-01",
						"tanggal_selesai": "2013-01-01",
						"tahun": 2013,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 21",
						"nomor_sertifikat": "CERT21"
					},
					{
						"id": 20,
						"tipe_pelatihan": "",
						"jenis_pelatihan": "",
						"nama_pelatihan": "20a",
						"tanggal_mulai": "2012-01-01",
						"tanggal_selesai": "2012-01-01",
						"tahun": 2012,
						"durasi": 12,
						"institusi_penyelenggara": "Institution 20",
						"nomor_sertifikat": "CERT20"
					}
				],
				"meta": {"limit": 2, "offset": 1, "total": 11}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan teknis pegawai 1d",
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 23,
						"tipe_pelatihan": "Teknis",
						"jenis_pelatihan": "Workshop",
						"nama_pelatihan": "23a",
						"tanggal_mulai": "2015-01-01",
						"tanggal_selesai": "2015-01-01",
						"tahun": 2015,
						"durasi": 16,
						"institusi_penyelenggara": "Institution 23",
						"nomor_sertifikat": "CERT23"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok: admin dapat melihat riwayat pelatihan teknis pegawai yang tidak ada data",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-pelatihan-teknis", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
