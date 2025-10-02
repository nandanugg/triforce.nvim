package riwayatkepangkatan

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
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_kenaikan_pangkat ("id", "nama", "deleted_at") values
			('21', 'jenis-kp-1', null),
			('22', 'jenis-kp-2', null),
			('23', 'jenis-kp-3', null),
			('24', 'jenis-kp-4', null),
			('25', 'jenis-kp-5', now());

		insert into ref_golongan ("id", "nama", "nama_pangkat", "deleted_at") values
			('21', 'diamond 1', 'petik 1', null),
			('22', 'diamond 2', 'petik 2', null),
			('23', 'diamond 3', 'petik 3', null),
			('24', 'diamond 4', 'petik 4', null),
			('25', 'diamond 5', 'petik 5', now());

		insert into riwayat_golongan ("id", "pns_nip", "jenis_kp_id", "golongan_id", "tmt_golongan", "sk_nomor", "sk_tanggal", "mk_golongan_tahun", "mk_golongan_bulan", "no_bkn", "tanggal_bkn", "jumlah_angka_kredit_tambahan", "jumlah_angka_kredit_utama", "deleted_at") values
			('21', '41', '21', '21', '2000-01-03', 'nomor-sk-1', '2000-01-01', 1, 2, 'no-bkn-1', '2000-01-02', 1, 2, null),
			('22', '41', '22', '22', '2001-01-03', 'nomor-sk-2', '2001-01-01', 1, 2, 'no-bkn-2', '2001-01-02', 1, 2, null),
			('23', '41', '23', '23', '2002-01-03', 'nomor-sk-3', '2002-01-01', 1, 2, 'no-bkn-3', '2002-01-02', 1, 2, null),
			('24', '42', '24', '24', '2003-01-03', 'nomor-sk-4', '2003-01-01', 1, 2, 'no-bkn-4', '2003-01-02', 1, 2, null),
			('25', '41', '25', '25', '2004-01-03', 'nomor-sk-5', '2004-01-01', 1, 2, 'no-bkn-5', '2004-01-02', 1, 2, now()),
			('26', '41', '25', '25', '2005-01-03', 'nomor-sk-6', '2005-01-01', 1, 2, 'no-bkn-6', '2005-01-02', 1, 2, null);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                26,
						"id_jenis_kp":       null,
						"nama_jenis_kp":     "",
						"id_golongan":       null,
						"nama_golongan":     "",
						"nama_golongan_pangkat": "",
						"tmt_golongan":      "2005-01-03",
						"sk_nomor":          "nomor-sk-6",
						"sk_tanggal":        "2005-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-6",
						"tanggal_bkn":       "2005-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					},
					{
						"id":                23,
						"id_jenis_kp":       23,
						"nama_jenis_kp":     "jenis-kp-3",
						"id_golongan":       23,
						"nama_golongan":     "diamond 3",
						"nama_golongan_pangkat": "petik 3",
						"tmt_golongan":      "2002-01-03",
						"sk_nomor":          "nomor-sk-3",
						"sk_tanggal":        "2002-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-3",
						"tanggal_bkn":       "2002-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					},
					{
						"id":                22,
						"id_jenis_kp":       22,
						"nama_jenis_kp":     "jenis-kp-2",
						"id_golongan":       22,
						"nama_golongan":     "diamond 2",
						"nama_golongan_pangkat": "petik 2",
						"tmt_golongan":      "2001-01-03",
						"sk_nomor":          "nomor-sk-2",
						"sk_tanggal":        "2001-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-2",
						"tanggal_bkn":       "2001-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					},
					{
						"id":                21,
						"id_jenis_kp":       21,
						"nama_jenis_kp":     "jenis-kp-1",
						"id_golongan":       21,
						"nama_golongan":     "diamond 1",
						"nama_golongan_pangkat": "petik 1",
						"tmt_golongan":      "2000-01-03",
						"sk_nomor":          "nomor-sk-1",
						"sk_tanggal":        "2000-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-1",
						"tanggal_bkn":       "2000-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                           23,
						"id_jenis_kp":                  23,
						"nama_jenis_kp":                "jenis-kp-3",
						"id_golongan":                  23,
						"nama_golongan":                "diamond 3",
						"nama_golongan_pangkat":        "petik 3",
						"tmt_golongan":                 "2002-01-03",
						"sk_nomor":                     "nomor-sk-3",
						"sk_tanggal":                   "2002-01-01",
						"mk_golongan_tahun":            1,
						"mk_golongan_bulan":            2,
						"no_bkn":                       "no-bkn-3",
						"tanggal_bkn":                  "2002-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
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
			_, err := db.Exec(context.Background(), tt.dbData)
			repo := dbrepo.New(db)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-kepangkatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
		insert into riwayat_golongan
			(id, pns_nip, deleted_at,   file_base64) values
			(1, '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c',     null,         '` + pdfBase64 + `'),
			(3, '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c',     null,         'data:application/pdf;base64,invalid'),
			(5, '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c',     null,         null),
			(7, '1c',     null,         '');
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
			name:              "error: base64 riwayat kepangkatan tidak valid",
			dbData:            dbData,
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat kepangkatan sudah dihapus",
			dbData:            dbData,
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berisi null value",
			dbData:            dbData,
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berupa string kosong",
			dbData:            dbData,
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kepangkatan bukan milik user login",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kepangkatan tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-kepangkatan/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepo.New(pgxconn)
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

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_kenaikan_pangkat ("id", "nama", "deleted_at") values
			('21', 'jenis-kp-1', null),
			('22', 'jenis-kp-2', null),
			('23', 'jenis-kp-3', null),
			('24', 'jenis-kp-4', null),
			('25', 'jenis-kp-5', now());

		insert into ref_golongan ("id", "nama", "nama_pangkat", "deleted_at") values
			('21', 'diamond 1', 'petik 1', null),
			('22', 'diamond 2', 'petik 2', null),
			('23', 'diamond 3', 'petik 3', null),
			('24', 'diamond 4', 'petik 4', null),
			('25', 'diamond 5', 'petik 5', now());

		insert into riwayat_golongan ("id", "pns_nip", "jenis_kp_id", "golongan_id", "tmt_golongan", "sk_nomor", "sk_tanggal", "mk_golongan_tahun", "mk_golongan_bulan", "no_bkn", "tanggal_bkn", "jumlah_angka_kredit_tambahan", "jumlah_angka_kredit_utama", "deleted_at") values
			('21', '41', '21', '21', '2000-01-03', 'nomor-sk-1', '2000-01-01', 1, 2, 'no-bkn-1', '2000-01-02', 1, 2, null),
			('22', '41', '22', '22', '2001-01-03', 'nomor-sk-2', '2001-01-01', 1, 2, 'no-bkn-2', '2001-01-02', 1, 2, null),
			('23', '41', '23', '23', '2002-01-03', 'nomor-sk-3', '2002-01-01', 1, 2, 'no-bkn-3', '2002-01-02', 1, 2, null),
			('24', '42', '24', '24', '2003-01-03', 'nomor-sk-4', '2003-01-01', 1, 2, 'no-bkn-4', '2003-01-02', 1, 2, null),
			('25', '41', '25', '25', '2004-01-03', 'nomor-sk-5', '2004-01-01', 1, 2, 'no-bkn-5', '2004-01-02', 1, 2, now()),
			('26', '41', '25', '25', '2005-01-03', 'nomor-sk-6', '2005-01-01', 1, 2, 'no-bkn-6', '2005-01-02', 1, 2, null);
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
			name:             "ok: tanpa parameter apapun",
			dbData:           dbData,
			nip:              "41",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                26,
						"id_jenis_kp":       null,
						"nama_jenis_kp":     "",
						"id_golongan":       null,
						"nama_golongan":     "",
						"nama_golongan_pangkat": "",
						"tmt_golongan":      "2005-01-03",
						"sk_nomor":          "nomor-sk-6",
						"sk_tanggal":        "2005-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-6",
						"tanggal_bkn":       "2005-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					},
					{
						"id":                23,
						"id_jenis_kp":       23,
						"nama_jenis_kp":     "jenis-kp-3",
						"id_golongan":       23,
						"nama_golongan":     "diamond 3",
						"nama_golongan_pangkat": "petik 3",
						"tmt_golongan":      "2002-01-03",
						"sk_nomor":          "nomor-sk-3",
						"sk_tanggal":        "2002-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-3",
						"tanggal_bkn":       "2002-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					},
					{
						"id":                22,
						"id_jenis_kp":       22,
						"nama_jenis_kp":     "jenis-kp-2",
						"id_golongan":       22,
						"nama_golongan":     "diamond 2",
						"nama_golongan_pangkat": "petik 2",
						"tmt_golongan":      "2001-01-03",
						"sk_nomor":          "nomor-sk-2",
						"sk_tanggal":        "2001-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-2",
						"tanggal_bkn":       "2001-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					},
					{
						"id":                21,
						"id_jenis_kp":       21,
						"nama_jenis_kp":     "jenis-kp-1",
						"id_golongan":       21,
						"nama_golongan":     "diamond 1",
						"nama_golongan_pangkat": "petik 1",
						"tmt_golongan":      "2000-01-03",
						"sk_nomor":          "nomor-sk-1",
						"sk_tanggal":        "2000-01-01",
						"mk_golongan_tahun": 1,
						"mk_golongan_bulan": 2,
						"no_bkn":            "no-bkn-1",
						"tanggal_bkn":       "2000-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			nip:              "41",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                           23,
						"id_jenis_kp":                  23,
						"nama_jenis_kp":                "jenis-kp-3",
						"id_golongan":                  23,
						"nama_golongan":                "diamond 3",
						"nama_golongan_pangkat":        "petik 3",
						"tmt_golongan":                 "2002-01-03",
						"sk_nomor":                     "nomor-sk-3",
						"sk_tanggal":                   "2002-01-01",
						"mk_golongan_tahun":            1,
						"mk_golongan_bulan":            2,
						"no_bkn":                       "no-bkn-3",
						"tanggal_bkn":                  "2002-01-02",
						"jumlah_angka_kredit_tambahan": 1,
						"jumlah_angka_kredit_utama":    2
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			nip:              "200",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
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
			nip:              "123456789",
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
			repo := dbrepo.New(db)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-kepangkatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
		insert into riwayat_golongan
			(id, pns_nip, deleted_at,   file_base64) values
			(1, '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c',     null,         '` + pdfBase64 + `'),
			(3, '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c',     null,         'data:application/pdf;base64,invalid'),
			(5, '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c',     null,         null),
			(7, '1c',     null,         '');
		`

	tests := []struct {
		name              string
		dbData            string
		nip               string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 riwayat kepangkatan tidak valid",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat kepangkatan sudah dihapus",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berisi null value",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berupa string kosong",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kepangkatan tidak ditemukan",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: user is not an admin",
			dbData:            dbData,
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode:  http.StatusForbidden,
			wantResponseBytes: []byte(`{"message": "akses ditolak"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
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

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-kepangkatan/%s/berkas", tt.nip, tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepo.New(pgxconn)
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
