package riwayatkepangkatan

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
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/usulanperubahandata/usulanperubahandatatest"
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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

	authHeader := []string{apitest.GenerateAuthHeader("41")}
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
						"id":                "26",
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
						"id":                "23",
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
						"id":                "22",
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
						"id":                "21",
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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                           "23",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-kepangkatan", nil)
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
		insert into riwayat_golongan
			(id,   pns_nip, deleted_at,   file_base64) values
			('1a', '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			('2',  '1c',     null,         '` + pdfBase64 + `'),
			('3',  '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			('4',  '1c',     null,         'data:application/pdf;base64,invalid'),
			('5',  '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			('6',  '1c',     null,         null),
			('7',  '1c',     null,         '');
		`
	db := dbtest.New(t, dbmigrations.FS)
	_, err = db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
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
			paramID:           "1a",
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
			name:              "error: base64 riwayat kepangkatan tidak valid",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat kepangkatan sudah dihapus",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berisi null value",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berupa string kosong",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kepangkatan bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kepangkatan tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-kepangkatan/%s/berkas", tt.paramID), nil)
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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
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
			name:             "ok: tanpa parameter apapun",
			nip:              "41",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                "26",
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
						"id":                "23",
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
						"id":                "22",
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
						"id":                "21",
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
			nip:              "41",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id":                           "23",
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
			nip:              "200",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			nip:              "123456789",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-kepangkatan", nil)
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
		insert into riwayat_golongan
			(id,   pns_nip, deleted_at,   file_base64) values
			('1a', '1c',     null,         'data:application/pdf;base64,` + pdfBase64 + `'),
			('2',  '1c',     null,         '` + pdfBase64 + `'),
			('3',  '1c',     null,         'data:images/png;base64,` + pngBase64 + `'),
			('4',  '1c',     null,         'data:application/pdf;base64,invalid'),
			('5',  '1c',     '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			('6',  '1c',     null,         null),
			('7',  '1c',     null,         '');
		`
	db := dbtest.New(t, dbmigrations.FS)
	_, err = db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	authMw := api.NewAuthMiddleware(authSvc, apitest.Keyfunc)
	svcRoute := usulanperubahandata.RegisterRoutes(e, db, repo, authMw)
	RegisterRoutes(e, repo, authMw, svcRoute)

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
			paramID:           "1a",
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
			name:              "error: base64 riwayat kepangkatan tidak valid",
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat kepangkatan sudah dihapus",
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berisi null value",
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: base64 riwayat kepangkatan berupa string kosong",
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat kepangkatan tidak ditemukan",
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat kepangkatan tidak ditemukan"}`),
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-kepangkatan/%s/berkas", tt.nip, tt.paramID), nil)
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
		insert into pegawai
			(pns_id,  nip_baru, nama,      deleted_at) values
			('id_1a', '1a',     'User 1a', null),
			('id_1c', '1c',     'John',    null),
			('id_1d', '1d',     'Jane',    '2000-01-01'),
			('id_1e', '1e',     'User 1e', null),
			('id_1f', '1f',     'User 1f', null);
		insert into ref_golongan
			(id,  nama,    nama_pangkat, deleted_at) values
			('1', 'Gol 1', 'I',          null),
			('2', 'Gol 2', 'II',         '2000-01-01'),
			('3', 'Gol 3', 'III',        null);
		insert into ref_jenis_kenaikan_pangkat
			(id, nama,  deleted_at) values
			(1,  'KP1', null),
			(2,  'KP2', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
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
				"jenis_kp_id": 1,
				"golongan_id": 1,
				"tmt_golongan": "2001-01-01",
				"nomor_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "BKN.01",
				"tanggal_bkn": "2002-01-01",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 2,
				"jumlah_angka_kredit_utama": 3,
				"jumlah_angka_kredit_tambahan": 0
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "{id}",
					"jenis_kp_id":                  int32(1),
					"kode_jenis_kp":                "1",
					"jenis_kp":                     "KP1",
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "SK.01",
					"no_bkn":                       "BKN.01",
					"jumlah_angka_kredit_utama":    int32(3),
					"jumlah_angka_kredit_tambahan": int32(0),
					"mk_golongan_tahun":            int16(1),
					"mk_golongan_bulan":            int16(2),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                 time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     "John",
					"created_at":                   "{created_at}",
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": null,
				"golongan_id": 3,
				"tmt_golongan": null,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "",
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 1,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "{id}",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  int16(3),
					"golongan_nama":                "Gol 3",
					"pangkat_nama":                 "III",
					"sk_nomor":                     "",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            int16(1),
					"mk_golongan_bulan":            int16(1),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1e",
					"pns_nip":                      "1e",
					"pns_nama":                     "User 1e",
					"created_at":                   "{created_at}",
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": "{id}" }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "{id}",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            int16(0),
					"mk_golongan_bulan":            int16(0),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1f",
					"pns_nip":                      "1f",
					"pns_nama":                     "User 1f",
					"created_at":                   "{created_at}",
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
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
				"jenis_kp_id": 1,
				"golongan_id": 1,
				"tmt_golongan": null,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "",
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 1,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: golongan or jenis kenaikan pangkat is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": 0,
				"golongan_id": 0,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis kenaikan pangkat tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: golongan or jenis kenaikan pangkat is deleted",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": 2,
				"golongan_id": 2,
				"tmt_golongan": null,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "",
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 1,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis kenaikan pangkat tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: exceed length limit, unexpected date or data type",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": "1",
				"golongan_id": "1",
				"tmt_golongan": "",
				"nomor_sk": "` + strings.Repeat(".", 101) + `",
				"tanggal_sk": "",
				"nomor_bkn": "` + strings.Repeat(".", 101) + `",
				"tanggal_bkn": "",
				"masa_kerja_golongan_tahun": "1",
				"masa_kerja_golongan_bulan": "1",
				"jumlah_angka_kredit_utama": "1",
				"jumlah_angka_kredit_tambahan": "1"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"golongan_id\" harus dalam tipe integer` +
				` | parameter \"jenis_kp_id\" harus dalam tipe integer` +
				` | parameter \"jumlah_angka_kredit_tambahan\" harus dalam tipe integer` +
				` | parameter \"jumlah_angka_kredit_utama\" harus dalam tipe integer` +
				` | parameter \"masa_kerja_golongan_bulan\" harus dalam tipe integer` +
				` | parameter \"masa_kerja_golongan_tahun\" harus dalam tipe integer` +
				` | parameter \"nomor_bkn\" harus 100 karakter atau kurang` +
				` | parameter \"nomor_sk\" harus 100 karakter atau kurang` +
				` | parameter \"tanggal_bkn\" harus dalam format date` +
				` | parameter \"tanggal_sk\" harus dalam format date` +
				` | parameter \"tmt_golongan\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null params",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": null,
				"golongan_id": null,
				"tmt_golongan": null,
				"nomor_sk": null,
				"tanggal_sk": null,
				"nomor_bkn": null,
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": null,
				"masa_kerja_golongan_bulan": null,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"golongan_id\" tidak boleh null` +
				` | parameter \"masa_kerja_golongan_bulan\" tidak boleh null` +
				` | parameter \"masa_kerja_golongan_tahun\" tidak boleh null` +
				` | parameter \"nomor_bkn\" tidak boleh null` +
				` | parameter \"nomor_sk\" tidak boleh null` +
				` | parameter \"tanggal_sk\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"golongan_id\" harus diisi` +
				` | parameter \"nomor_sk\" harus diisi` +
				` | parameter \"tanggal_sk\" harus diisi` +
				` | parameter \"masa_kerja_golongan_tahun\" harus diisi` +
				` | parameter \"masa_kerja_golongan_bulan\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kepangkatan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_golongan", "where pns_nip = $1", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["id"] = row["id"]
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]

						tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", row["id"].(string))
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
			(pns_id,  nip_baru, nama,   deleted_at) values
			('id_1c', '1c',     'John', null),
			('id_1d', '1d',     'Jane', '2000-01-01'),
			('id_1e', '1e',     'Doe',  null);
		insert into ref_golongan
			(id,  nama,    nama_pangkat, deleted_at) values
			('1', 'Gol 1', 'I',          null),
			('2', 'Gol 2', 'II',         '2000-01-01');
		insert into ref_jenis_kenaikan_pangkat
			(id, nama,  deleted_at) values
			(1,  'KP1', null),
			(2,  'KP2', '2000-01-01');
		insert into riwayat_golongan
			(id,  golongan_nama, status_satker, status_biro, pangkat_terakhir, bkn_id, file_base64, keterangan_berkas, arsip_id, golongan_asal, basic, sk_type, kanreg, kpkn, keterangan, lpnk, jenis_riwayat, pns_id,  pns_nip, created_at,   updated_at) values
			('1', 'Gol 1',       1,             1,           1,                '1',    'data:abc',  'abc',             1,        'f',           '1',   1,       '1',    '1',  'ket',      '1',  '1',           'id_1c', '1c',    '2000-01-01', '2000-01-01'),
			('2', 'Gol 1',       1,             1,           1,                '1',    'data:abc',  'abc',             1,        'f',           '1',   1,       '1',    '1',  'ket',      '1',  '1',           'id_1c', '1c',    '2000-01-01', '2000-01-01'),
			('3', 'Gol 1',       1,             1,           1,                '1',    'data:abc',  'abc',             1,        'f',           '1',   1,       '1',    '1',  'ket',      '1',  '1',           'id_1c', '1c',    '2000-01-01', '2000-01-01');
		insert into riwayat_golongan
			(id,  golongan_nama, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('4', 'Gol 1',       'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('5', 'Gol 1',       'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('6', 'Gol 1',       'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                           "6",
			"jenis_kp_id":                  nil,
			"kode_jenis_kp":                nil,
			"jenis_kp":                     nil,
			"golongan_id":                  nil,
			"golongan_nama":                "Gol 1",
			"pangkat_nama":                 nil,
			"sk_nomor":                     nil,
			"no_bkn":                       nil,
			"jumlah_angka_kredit_utama":    nil,
			"jumlah_angka_kredit_tambahan": nil,
			"mk_golongan_tahun":            nil,
			"mk_golongan_bulan":            nil,
			"sk_tanggal":                   nil,
			"tanggal_bkn":                  nil,
			"tmt_golongan":                 nil,
			"status_satker":                nil,
			"status_biro":                  nil,
			"pangkat_terakhir":             nil,
			"bkn_id":                       nil,
			"file_base64":                  nil,
			"s3_file_id":                   nil,
			"keterangan_berkas":            nil,
			"arsip_id":                     nil,
			"golongan_asal":                nil,
			"basic":                        nil,
			"sk_type":                      nil,
			"kanreg":                       nil,
			"kpkn":                         nil,
			"keterangan":                   nil,
			"lpnk":                         nil,
			"jenis_riwayat":                nil,
			"pns_id":                       "id_1c",
			"pns_nip":                      "1c",
			"pns_nama":                     nil,
			"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                   nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
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
				"jenis_kp_id": 1,
				"golongan_id": 1,
				"tmt_golongan": "2001-01-01",
				"nomor_sk": "SK.01",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "BKN.01",
				"tanggal_bkn": "2002-01-01",
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 2,
				"jumlah_angka_kredit_utama": 3,
				"jumlah_angka_kredit_tambahan": 0
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "1",
					"jenis_kp_id":                  int32(1),
					"kode_jenis_kp":                "1",
					"jenis_kp":                     "KP1",
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "SK.01",
					"no_bkn":                       "BKN.01",
					"jumlah_angka_kredit_utama":    int32(3),
					"jumlah_angka_kredit_tambahan": int32(0),
					"mk_golongan_tahun":            int16(1),
					"mk_golongan_bulan":            int16(2),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                 time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_satker":                int32(1),
					"status_biro":                  int32(1),
					"pangkat_terakhir":             int32(1),
					"bkn_id":                       "1",
					"file_base64":                  "data:abc",
					"s3_file_id":                   nil,
					"keterangan_berkas":            "abc",
					"arsip_id":                     int64(1),
					"golongan_asal":                "f",
					"basic":                        "1",
					"sk_type":                      int16(1),
					"kanreg":                       "1",
					"kpkn":                         "1",
					"keterangan":                   "ket",
					"lpnk":                         "1",
					"jenis_riwayat":                "1",
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": null,
				"golongan_id": 1,
				"tmt_golongan": null,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "",
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 1,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "2",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            int16(1),
					"mk_golongan_bulan":            int16(1),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                int32(1),
					"status_biro":                  int32(1),
					"pangkat_terakhir":             int32(1),
					"bkn_id":                       "1",
					"file_base64":                  "data:abc",
					"s3_file_id":                   nil,
					"keterangan_berkas":            "abc",
					"arsip_id":                     int64(1),
					"golongan_asal":                "f",
					"basic":                        "1",
					"sk_type":                      int16(1),
					"kanreg":                       "1",
					"kpkn":                         "1",
					"keterangan":                   "ket",
					"lpnk":                         "1",
					"jenis_riwayat":                "1",
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			paramNIP:      "1c",
			paramID:       "3",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "3",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            int16(0),
					"mk_golongan_bulan":            int16(0),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                int32(1),
					"status_biro":                  int32(1),
					"pangkat_terakhir":             int32(1),
					"bkn_id":                       "1",
					"file_base64":                  "data:abc",
					"s3_file_id":                   nil,
					"keterangan_berkas":            "abc",
					"arsip_id":                     int64(1),
					"golongan_asal":                "f",
					"basic":                        "1",
					"sk_type":                      int16(1),
					"kanreg":                       "1",
					"kpkn":                         "1",
					"keterangan":                   "ket",
					"lpnk":                         "1",
					"jenis_riwayat":                "1",
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "error: riwayat golongan is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat golongan is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "4",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"golongan_id": 1,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "4",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1e",
					"pns_nip":                      "1e",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:          "error: riwayat golongan is deleted",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": 1,
				"golongan_id": 1,
				"tmt_golongan": null,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "",
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 1,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "5",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: golongan or jenis kenaikan pangkat is not found",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": 0,
				"golongan_id": 0,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"masa_kerja_golongan_tahun": 0,
				"masa_kerja_golongan_bulan": 0
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis kenaikan pangkat tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: golongan or jenis kenaikan pangkat is deleted",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": 2,
				"golongan_id": 2,
				"tmt_golongan": null,
				"nomor_sk": "",
				"tanggal_sk": "2000-01-01",
				"nomor_bkn": "",
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": 1,
				"masa_kerja_golongan_bulan": 1,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis kenaikan pangkat tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": "1",
				"golongan_id": "1",
				"tmt_golongan": "",
				"nomor_sk": "` + strings.Repeat(".", 101) + `",
				"tanggal_sk": "",
				"nomor_bkn": "` + strings.Repeat(".", 101) + `",
				"tanggal_bkn": "",
				"masa_kerja_golongan_tahun": "1",
				"masa_kerja_golongan_bulan": "1",
				"jumlah_angka_kredit_utama": "1",
				"jumlah_angka_kredit_tambahan": "1"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"golongan_id\" harus dalam tipe integer` +
				` | parameter \"jenis_kp_id\" harus dalam tipe integer` +
				` | parameter \"jumlah_angka_kredit_tambahan\" harus dalam tipe integer` +
				` | parameter \"jumlah_angka_kredit_utama\" harus dalam tipe integer` +
				` | parameter \"masa_kerja_golongan_bulan\" harus dalam tipe integer` +
				` | parameter \"masa_kerja_golongan_tahun\" harus dalam tipe integer` +
				` | parameter \"nomor_bkn\" harus 100 karakter atau kurang` +
				` | parameter \"nomor_sk\" harus 100 karakter atau kurang` +
				` | parameter \"tanggal_bkn\" harus dalam format date` +
				` | parameter \"tanggal_sk\" harus dalam format date` +
				` | parameter \"tmt_golongan\" harus dalam format date"}`,
			wantDBRows: defaultRows,
		},
		{
			name:          "error: null params",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_kp_id": null,
				"golongan_id": null,
				"tmt_golongan": null,
				"nomor_sk": null,
				"tanggal_sk": null,
				"nomor_bkn": null,
				"tanggal_bkn": null,
				"masa_kerja_golongan_tahun": null,
				"masa_kerja_golongan_bulan": null,
				"jumlah_angka_kredit_utama": null,
				"jumlah_angka_kredit_tambahan": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"golongan_id\" tidak boleh null` +
				` | parameter \"masa_kerja_golongan_bulan\" tidak boleh null` +
				` | parameter \"masa_kerja_golongan_tahun\" tidak boleh null` +
				` | parameter \"nomor_bkn\" tidak boleh null` +
				` | parameter \"nomor_sk\" tidak boleh null` +
				` | parameter \"tanggal_sk\" tidak boleh null"}`,
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
				` | parameter \"golongan_id\" harus diisi` +
				` | parameter \"nomor_sk\" harus diisi` +
				` | parameter \"tanggal_sk\" harus diisi` +
				` | parameter \"masa_kerja_golongan_tahun\" harus diisi` +
				` | parameter \"masa_kerja_golongan_bulan\" harus diisi"}`,
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kepangkatan/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_golongan", "where id = $1", tt.paramID)
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
		insert into riwayat_golongan
			(id,  golongan_nama, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('1', 'Gol 1',       'id_1c', '1c',    '2000-01-01', '2000-01-01', null),
			('2', null,          'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('3', null,          'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Gol 1',       'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                           "4",
			"jenis_kp_id":                  nil,
			"kode_jenis_kp":                nil,
			"jenis_kp":                     nil,
			"golongan_id":                  nil,
			"golongan_nama":                "Gol 1",
			"pangkat_nama":                 nil,
			"sk_nomor":                     nil,
			"no_bkn":                       nil,
			"jumlah_angka_kredit_utama":    nil,
			"jumlah_angka_kredit_tambahan": nil,
			"mk_golongan_tahun":            nil,
			"mk_golongan_bulan":            nil,
			"sk_tanggal":                   nil,
			"tanggal_bkn":                  nil,
			"tmt_golongan":                 nil,
			"status_satker":                nil,
			"status_biro":                  nil,
			"pangkat_terakhir":             nil,
			"bkn_id":                       nil,
			"file_base64":                  nil,
			"s3_file_id":                   nil,
			"keterangan_berkas":            nil,
			"arsip_id":                     nil,
			"golongan_asal":                nil,
			"basic":                        nil,
			"sk_type":                      nil,
			"kanreg":                       nil,
			"kpkn":                         nil,
			"keterangan":                   nil,
			"lpnk":                         nil,
			"jenis_riwayat":                nil,
			"pns_id":                       "id_1c",
			"pns_nip":                      "1c",
			"pns_nama":                     nil,
			"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                   nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
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
					"id":                           "1",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat golongan is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "2",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                nil,
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1e",
					"pns_nip":                      "1e",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:             "error: riwayat golongan is not found",
			paramNIP:         "1c",
			paramID:          "abc",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat golongan is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "3",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                nil,
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
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

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kepangkatan/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_golongan", "where id = $1", tt.paramID)
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
		insert into riwayat_golongan
			(id,  golongan_nama, status_satker, status_biro, pangkat_terakhir, bkn_id, file_base64, keterangan_berkas, arsip_id, golongan_asal, basic, sk_type, kanreg, kpkn, keterangan, lpnk, jenis_riwayat, pns_id,  pns_nip, created_at,   updated_at) values
			('1', 'Gol 1',       1,             1,           1,                '1',    'data:abc',  'abc',             1,        'f',           '1',   1,       '1',    '1',  'ket',      '1',  '1',           'id_1c', '1c',    '2000-01-01', '2000-01-01');
		insert into riwayat_golongan
			(id,  golongan_nama, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('2', 'Gol 1',       'id_1e', '1e',    '2000-01-01', '2000-01-01', null),
			('3', 'Gol 1',       'id_1c', '1c',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('4', 'Gol 1',       'id_1c', '1c',    '2000-01-01', '2000-01-01', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                           "4",
			"jenis_kp_id":                  nil,
			"kode_jenis_kp":                nil,
			"jenis_kp":                     nil,
			"golongan_id":                  nil,
			"golongan_nama":                "Gol 1",
			"pangkat_nama":                 nil,
			"sk_nomor":                     nil,
			"no_bkn":                       nil,
			"jumlah_angka_kredit_utama":    nil,
			"jumlah_angka_kredit_tambahan": nil,
			"mk_golongan_tahun":            nil,
			"mk_golongan_bulan":            nil,
			"sk_tanggal":                   nil,
			"tanggal_bkn":                  nil,
			"tmt_golongan":                 nil,
			"status_satker":                nil,
			"status_biro":                  nil,
			"pangkat_terakhir":             nil,
			"bkn_id":                       nil,
			"file_base64":                  nil,
			"s3_file_id":                   nil,
			"keterangan_berkas":            nil,
			"arsip_id":                     nil,
			"golongan_asal":                nil,
			"basic":                        nil,
			"sk_type":                      nil,
			"kanreg":                       nil,
			"kpkn":                         nil,
			"keterangan":                   nil,
			"lpnk":                         nil,
			"jenis_riwayat":                nil,
			"pns_id":                       "id_1c",
			"pns_nip":                      "1c",
			"pns_nama":                     nil,
			"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                   nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
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
					"id":                           "1",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                int32(1),
					"status_biro":                  int32(1),
					"pangkat_terakhir":             int32(1),
					"bkn_id":                       "1",
					"file_base64":                  "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"s3_file_id":                   nil,
					"keterangan_berkas":            "abc",
					"arsip_id":                     int64(1),
					"golongan_asal":                "f",
					"basic":                        "1",
					"sk_type":                      int16(1),
					"kanreg":                       "1",
					"kpkn":                         "1",
					"keterangan":                   "ket",
					"lpnk":                         "1",
					"jenis_riwayat":                "1",
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:              "error: riwayat golongan is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat golongan is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "2",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1e",
					"pns_nip":                      "1e",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   nil,
				},
			},
		},
		{
			name:              "error: riwayat golongan is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                           "3",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 nil,
					"sk_nomor":                     nil,
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
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

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-kepangkatan/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_golongan", "where id = $1", tt.paramID)
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
			(pns_id,  nip_baru, nama) values
			('id_1a', '1a',     'User 1a'),
			('id_1b', '1b',     'User 1b'),
			('id_1c', '1c',     'User 1c'),
			('id_1d', '1d',     'User 1d'),
			('id_1e', '1e',     'User 1e'),
			('id_1f', '1f',     'User 1f'),
			('id_1g', '1g',     'User 1g'),
			('id_1h', '1h',     'User 1h');
		insert into ref_golongan
			(id,  nama,    nama_pangkat, deleted_at) values
			('1', 'Gol 1', 'I',          null),
			('2', 'Gol 2', 'II',         '2000-01-01');
		insert into ref_jenis_kenaikan_pangkat
			(id, nama,  deleted_at) values
			(1,  'KP1', null),
			(2,  'KP2', '2000-01-01');
		insert into riwayat_golongan
			(id,    sk_nomor, pns_id,  pns_nip, created_at,   updated_at,   deleted_at) values
			('001', 'SK1',    'id_1a', '1a',    '2000-01-01', '2000-01-01', null),
			('002', 'SK1',    'id_1a', '1a',    '2000-01-01', '2000-01-01', '2000-01-01'),
			('003', 'SK1',    'id_1d', '1d',    '2000-01-01', '2000-01-01', null),
			('004', 'SK1',    'id_1g', '1g',    '2000-01-01', '2000-01-01', null),
			('005', 'SK1',    'id_1h', '1h',    '2000-01-01', '2000-01-01', null);
		insert into riwayat_golongan
			(id,    jenis_kp_id, kode_jenis_kp, jenis_kp, golongan_id, golongan_nama, pangkat_nama, sk_nomor, no_bkn, jumlah_angka_kredit_utama, jumlah_angka_kredit_tambahan, mk_golongan_tahun, mk_golongan_bulan, sk_tanggal,   tanggal_bkn,  tmt_golongan, status_satker, status_biro, pangkat_terakhir, bkn_id, file_base64, keterangan_berkas, arsip_id, golongan_asal, basic, sk_type, kanreg, kpkn, keterangan, lpnk, jenis_riwayat, pns_id,  pns_nip, pns_nama,  created_at,   updated_at) values
			('006', 2,           2,             'KP',     2,           'Gol 1',       'I',          'SK',     'BKN',  100,                       5,                            2,                 2,                 '2000-01-01', '2000-01-01', '2000-01-01', 1,             1,           1,                '1',    'data:abc',  'abc',             1,        'f',           '1',   1,       '1',    '1',  'ket',      '1',  '1',           'id_1d', '1d',    'User.1d', '2000-01-01', '2000-01-01'),
			('007', 1,           1,             'KP',     1,           'Gol 1',       'I',          'SK',     'BKN',  100,                       5,                            2,                 2,                 '2000-01-01', '2000-01-01', '2000-01-01', 1,             1,           1,                '1',    'data:abc',  'abc',             1,        'f',           '1',   1,       '1',    '1',  'ket',      '1',  '1',           'id_1e', '1e',    'User.1e', '2000-01-01', '2000-01-01');
	`

	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	dbRows1a := dbtest.Rows{
		{
			"id":                           "001",
			"jenis_kp_id":                  nil,
			"kode_jenis_kp":                nil,
			"jenis_kp":                     nil,
			"golongan_id":                  nil,
			"golongan_nama":                nil,
			"pangkat_nama":                 nil,
			"sk_nomor":                     "SK1",
			"no_bkn":                       nil,
			"jumlah_angka_kredit_utama":    nil,
			"jumlah_angka_kredit_tambahan": nil,
			"mk_golongan_tahun":            nil,
			"mk_golongan_bulan":            nil,
			"sk_tanggal":                   nil,
			"tanggal_bkn":                  nil,
			"tmt_golongan":                 nil,
			"status_satker":                nil,
			"status_biro":                  nil,
			"pangkat_terakhir":             nil,
			"bkn_id":                       nil,
			"file_base64":                  nil,
			"s3_file_id":                   nil,
			"keterangan_berkas":            nil,
			"arsip_id":                     nil,
			"golongan_asal":                nil,
			"basic":                        nil,
			"sk_type":                      nil,
			"kanreg":                       nil,
			"kpkn":                         nil,
			"keterangan":                   nil,
			"lpnk":                         nil,
			"jenis_riwayat":                nil,
			"pns_id":                       "id_1a",
			"pns_nip":                      "1a",
			"pns_nama":                     nil,
			"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                   nil,
		},
		{
			"id":                           "002",
			"jenis_kp_id":                  nil,
			"kode_jenis_kp":                nil,
			"jenis_kp":                     nil,
			"golongan_id":                  nil,
			"golongan_nama":                nil,
			"pangkat_nama":                 nil,
			"sk_nomor":                     "SK1",
			"no_bkn":                       nil,
			"jumlah_angka_kredit_utama":    nil,
			"jumlah_angka_kredit_tambahan": nil,
			"mk_golongan_tahun":            nil,
			"mk_golongan_bulan":            nil,
			"sk_tanggal":                   nil,
			"tanggal_bkn":                  nil,
			"tmt_golongan":                 nil,
			"status_satker":                nil,
			"status_biro":                  nil,
			"pangkat_terakhir":             nil,
			"bkn_id":                       nil,
			"file_base64":                  nil,
			"s3_file_id":                   nil,
			"keterangan_berkas":            nil,
			"arsip_id":                     nil,
			"golongan_asal":                nil,
			"basic":                        nil,
			"sk_type":                      nil,
			"kanreg":                       nil,
			"kpkn":                         nil,
			"keterangan":                   nil,
			"lpnk":                         nil,
			"jenis_riwayat":                nil,
			"pns_id":                       "id_1a",
			"pns_nip":                      "1a",
			"pns_nama":                     nil,
			"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
		},
	}

	repo := dbrepo.New(db)
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
			name:          "ok: success create riwayat kepangkatan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"jenis_kp_id": 1,
					"golongan_id": 1,
					"tmt_golongan": "2001-01-01",
					"nomor_sk": "SK.01",
					"tanggal_sk": "2000-01-01",
					"nomor_bkn": "BKN.01",
					"tanggal_bkn": "2002-01-01",
					"masa_kerja_golongan_tahun": 1,
					"masa_kerja_golongan_bulan": 2,
					"jumlah_angka_kredit_utama": 3,
					"jumlah_angka_kredit_tambahan": 0
				}
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-kepangkatan",
						"action":     "CREATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"jenis_kp_id":                  [ null, 1            ],
							"nama_jenis_kp":                [ null, "KP1"        ],
							"golongan_id":                  [ null, 1            ],
							"nama_golongan":                [ null, "Gol 1"      ],
							"nama_golongan_pangkat":        [ null, "I"          ],
							"tmt_golongan":                 [ null, "2001-01-01" ],
							"nomor_sk":                     [ null, "SK.01"      ],
							"tanggal_sk":                   [ null, "2000-01-01" ],
							"nomor_bkn":                    [ null, "BKN.01"     ],
							"tanggal_bkn":                  [ null, "2002-01-01" ],
							"masa_kerja_golongan_tahun":    [ null, 1            ],
							"masa_kerja_golongan_bulan":    [ null, 2            ],
							"jumlah_angka_kredit_utama":    [ null, 3            ],
							"jumlah_angka_kredit_tambahan": [ null, 0            ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                           "{id}",
					"jenis_kp_id":                  int32(1),
					"kode_jenis_kp":                "1",
					"jenis_kp":                     "KP1",
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "SK.01",
					"no_bkn":                       "BKN.01",
					"jumlah_angka_kredit_utama":    int32(3),
					"jumlah_angka_kredit_tambahan": int32(0),
					"mk_golongan_tahun":            int16(1),
					"mk_golongan_bulan":            int16(2),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                 time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1c",
					"pns_nip":                      "1c",
					"pns_nama":                     "User 1c",
					"created_at":                   "{created_at}",
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1c",
					"jenis_data": "riwayat-kepangkatan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"jenis_kp_id":                  []any{nil, float64(1)},
						"nama_jenis_kp":                []any{nil, "KP1"},
						"golongan_id":                  []any{nil, float64(1)},
						"nama_golongan":                []any{nil, "Gol 1"},
						"nama_golongan_pangkat":        []any{nil, "I"},
						"tmt_golongan":                 []any{nil, "2001-01-01"},
						"nomor_sk":                     []any{nil, "SK.01"},
						"tanggal_sk":                   []any{nil, "2000-01-01"},
						"nomor_bkn":                    []any{nil, "BKN.01"},
						"tanggal_bkn":                  []any{nil, "2002-01-01"},
						"masa_kerja_golongan_tahun":    []any{nil, float64(1)},
						"masa_kerja_golongan_bulan":    []any{nil, float64(2)},
						"jumlah_angka_kredit_utama":    []any{nil, float64(3)},
						"jumlah_angka_kredit_tambahan": []any{nil, float64(0)},
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
			name:          "ok: success update riwayat kepangkatan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "006",
				"data": {
					"jenis_kp_id": null,
					"golongan_id": 1,
					"tmt_golongan": null,
					"nomor_sk": "",
					"tanggal_sk": "2000-01-01",
					"nomor_bkn": "",
					"tanggal_bkn": null,
					"masa_kerja_golongan_tahun": 1,
					"masa_kerja_golongan_bulan": 1,
					"jumlah_angka_kredit_utama": null,
					"jumlah_angka_kredit_tambahan": null
				}
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-kepangkatan",
						"action":     "UPDATE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "006",
						"perubahan_data": {
							"jenis_kp_id":                  [ 2,            null         ],
							"nama_jenis_kp":                [ null,         null         ],
							"golongan_id":                  [ 2,            1            ],
							"nama_golongan":                [ null,         "Gol 1"      ],
							"nama_golongan_pangkat":        [ null,         "I"          ],
							"tmt_golongan":                 [ "2000-01-01", null         ],
							"nomor_sk":                     [ "SK",         ""           ],
							"tanggal_sk":                   [ "2000-01-01", "2000-01-01" ],
							"nomor_bkn":                    [ "BKN",        null         ],
							"tanggal_bkn":                  [ "2000-01-01", null         ],
							"masa_kerja_golongan_tahun":    [ 2,            1            ],
							"masa_kerja_golongan_bulan":    [ 2,            1            ],
							"jumlah_angka_kredit_utama":    [ 100,          null         ],
							"jumlah_angka_kredit_tambahan": [ 5,            null         ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                           "003",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                nil,
					"pangkat_nama":                 nil,
					"sk_nomor":                     "SK1",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1d",
					"pns_nip":                      "1d",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   nil,
				},
				{
					"id":                           "006",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            int16(1),
					"mk_golongan_bulan":            int16(1),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                int32(1),
					"status_biro":                  int32(1),
					"pangkat_terakhir":             int32(1),
					"bkn_id":                       "1",
					"file_base64":                  "data:abc",
					"s3_file_id":                   nil,
					"keterangan_berkas":            "abc",
					"arsip_id":                     int64(1),
					"golongan_asal":                "f",
					"basic":                        "1",
					"sk_type":                      int16(1),
					"kanreg":                       "1",
					"kpkn":                         "1",
					"keterangan":                   "ket",
					"lpnk":                         "1",
					"jenis_riwayat":                "1",
					"pns_id":                       "id_1d",
					"pns_nip":                      "1d",
					"pns_nama":                     "User.1d",
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   "{updated_at}",
					"deleted_at":                   nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1d",
					"jenis_data": "riwayat-kepangkatan",
					"data_id":    "006",
					"perubahan_data": map[string]any{
						"jenis_kp_id":                  []any{float64(2), nil},
						"nama_jenis_kp":                []any{nil, nil},
						"golongan_id":                  []any{float64(2), float64(1)},
						"nama_golongan":                []any{nil, "Gol 1"},
						"nama_golongan_pangkat":        []any{nil, "I"},
						"tmt_golongan":                 []any{"2000-01-01", nil},
						"nomor_sk":                     []any{"SK", ""},
						"tanggal_sk":                   []any{"2000-01-01", "2000-01-01"},
						"nomor_bkn":                    []any{"BKN", nil},
						"tanggal_bkn":                  []any{"2000-01-01", nil},
						"masa_kerja_golongan_tahun":    []any{float64(2), float64(1)},
						"masa_kerja_golongan_bulan":    []any{float64(2), float64(1)},
						"jumlah_angka_kredit_utama":    []any{float64(100), nil},
						"jumlah_angka_kredit_tambahan": []any{float64(5), nil},
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
			name:          "ok: success delete riwayat kepangkatan",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1e")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "007"
			}`,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-kepangkatan",
						"action":     "DELETE",
						"status":     "Disetujui",
						"catatan":    "",
						"data_id":    "007",
						"perubahan_data": {
							"jenis_kp_id":                  [ 1,            null ],
							"nama_jenis_kp":                [ "KP1",        null ],
							"golongan_id":                  [ 1,            null ],
							"nama_golongan":                [ "Gol 1",      null ],
							"nama_golongan_pangkat":        [ "I",          null ],
							"tmt_golongan":                 [ "2000-01-01", null ],
							"nomor_sk":                     [ "SK",         null ],
							"tanggal_sk":                   [ "2000-01-01", null ],
							"nomor_bkn":                    [ "BKN",        null ],
							"tanggal_bkn":                  [ "2000-01-01", null ],
							"masa_kerja_golongan_tahun":    [ 2,            null ],
							"masa_kerja_golongan_bulan":    [ 2,            null ],
							"jumlah_angka_kredit_utama":    [ 100,          null ],
							"jumlah_angka_kredit_tambahan": [ 5,            null ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                           "007",
					"jenis_kp_id":                  int32(1),
					"kode_jenis_kp":                "1",
					"jenis_kp":                     "KP",
					"golongan_id":                  int16(1),
					"golongan_nama":                "Gol 1",
					"pangkat_nama":                 "I",
					"sk_nomor":                     "SK",
					"no_bkn":                       "BKN",
					"jumlah_angka_kredit_utama":    int32(100),
					"jumlah_angka_kredit_tambahan": int32(5),
					"mk_golongan_tahun":            int16(2),
					"mk_golongan_bulan":            int16(2),
					"sk_tanggal":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_bkn":                  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tmt_golongan":                 time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_satker":                int32(1),
					"status_biro":                  int32(1),
					"pangkat_terakhir":             int32(1),
					"bkn_id":                       "1",
					"file_base64":                  "data:abc",
					"s3_file_id":                   nil,
					"keterangan_berkas":            "abc",
					"arsip_id":                     int64(1),
					"golongan_asal":                "f",
					"basic":                        "1",
					"sk_type":                      int16(1),
					"kanreg":                       "1",
					"kpkn":                         "1",
					"keterangan":                   "ket",
					"lpnk":                         "1",
					"jenis_riwayat":                "1",
					"pns_id":                       "id_1e",
					"pns_nip":                      "1e",
					"pns_nama":                     "User.1e",
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   "{deleted_at}",
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1e",
					"jenis_data": "riwayat-kepangkatan",
					"data_id":    "007",
					"perubahan_data": map[string]any{
						"jenis_kp_id":                  []any{float64(1), nil},
						"nama_jenis_kp":                []any{"KP1", nil},
						"golongan_id":                  []any{float64(1), nil},
						"nama_golongan":                []any{"Gol 1", nil},
						"nama_golongan_pangkat":        []any{"I", nil},
						"tmt_golongan":                 []any{"2000-01-01", nil},
						"nomor_sk":                     []any{"SK", nil},
						"tanggal_sk":                   []any{"2000-01-01", nil},
						"nomor_bkn":                    []any{"BKN", nil},
						"tanggal_bkn":                  []any{"2000-01-01", nil},
						"masa_kerja_golongan_tahun":    []any{float64(2), nil},
						"masa_kerja_golongan_bulan":    []any{float64(2), nil},
						"jumlah_angka_kredit_utama":    []any{float64(100), nil},
						"jumlah_angka_kredit_tambahan": []any{float64(5), nil},
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
			name:          "ok: rollback on usulan perubahan data should not CREATE record",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1f")}},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"golongan_id": 1,
					"nomor_sk": "",
					"tanggal_sk": "2000-01-01",
					"masa_kerja_golongan_tahun": 0,
					"masa_kerja_golongan_bulan": 0
				}
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-kepangkatan",
						"action":     "CREATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    null,
						"perubahan_data": {
							"jenis_kp_id":                  [ null, null         ],
							"nama_jenis_kp":                [ null, null         ],
							"golongan_id":                  [ null, 1            ],
							"nama_golongan":                [ null, "Gol 1"      ],
							"nama_golongan_pangkat":        [ null, "I"          ],
							"tmt_golongan":                 [ null, null         ],
							"nomor_sk":                     [ null, ""           ],
							"tanggal_sk":                   [ null, "2000-01-01" ],
							"nomor_bkn":                    [ null, null         ],
							"tanggal_bkn":                  [ null, null         ],
							"masa_kerja_golongan_tahun":    [ null, 0            ],
							"masa_kerja_golongan_bulan":    [ null, 0            ],
							"jumlah_angka_kredit_utama":    [ null, null         ],
							"jumlah_angka_kredit_tambahan": [ null, null         ]
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
					"jenis_data": "riwayat-kepangkatan",
					"data_id":    nil,
					"perubahan_data": map[string]any{
						"jenis_kp_id":                  []any{nil, nil},
						"nama_jenis_kp":                []any{nil, nil},
						"golongan_id":                  []any{nil, float64(1)},
						"nama_golongan":                []any{nil, "Gol 1"},
						"nama_golongan_pangkat":        []any{nil, "I"},
						"tmt_golongan":                 []any{nil, nil},
						"nomor_sk":                     []any{nil, ""},
						"tanggal_sk":                   []any{nil, "2000-01-01"},
						"nomor_bkn":                    []any{nil, nil},
						"tanggal_bkn":                  []any{nil, nil},
						"masa_kerja_golongan_tahun":    []any{nil, float64(0)},
						"masa_kerja_golongan_bulan":    []any{nil, float64(0)},
						"jumlah_angka_kredit_utama":    []any{nil, nil},
						"jumlah_angka_kredit_tambahan": []any{nil, nil},
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
			name:          "ok: rollback on usulan perubahan data should not UPDATE record",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1g")}},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "004",
				"data": {
					"golongan_id": 1,
					"nomor_sk": "",
					"tanggal_sk": "2000-01-01",
					"masa_kerja_golongan_tahun": 0,
					"masa_kerja_golongan_bulan": 0
				}
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-kepangkatan",
						"action":     "UPDATE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    "004",
						"perubahan_data": {
							"jenis_kp_id":                  [ null,  null         ],
							"nama_jenis_kp":                [ null,  null         ],
							"golongan_id":                  [ null,  1            ],
							"nama_golongan":                [ null,  "Gol 1"      ],
							"nama_golongan_pangkat":        [ null,  "I"          ],
							"tmt_golongan":                 [ null,  null         ],
							"nomor_sk":                     [ "SK1", ""           ],
							"tanggal_sk":                   [ null,  "2000-01-01" ],
							"nomor_bkn":                    [ null,  null         ],
							"tanggal_bkn":                  [ null,  null         ],
							"masa_kerja_golongan_tahun":    [ null,  0            ],
							"masa_kerja_golongan_bulan":    [ null,  0            ],
							"jumlah_angka_kredit_utama":    [ null,  null         ],
							"jumlah_angka_kredit_tambahan": [ null,  null         ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                           "004",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                nil,
					"pangkat_nama":                 nil,
					"sk_nomor":                     "SK1",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1g",
					"pns_nip":                      "1g",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1g",
					"jenis_data": "riwayat-kepangkatan",
					"data_id":    "004",
					"perubahan_data": map[string]any{
						"jenis_kp_id":                  []any{nil, nil},
						"nama_jenis_kp":                []any{nil, nil},
						"golongan_id":                  []any{nil, float64(1)},
						"nama_golongan":                []any{nil, "Gol 1"},
						"nama_golongan_pangkat":        []any{nil, "I"},
						"tmt_golongan":                 []any{nil, nil},
						"nomor_sk":                     []any{"SK1", ""},
						"tanggal_sk":                   []any{nil, "2000-01-01"},
						"nomor_bkn":                    []any{nil, nil},
						"tanggal_bkn":                  []any{nil, nil},
						"masa_kerja_golongan_tahun":    []any{nil, float64(0)},
						"masa_kerja_golongan_bulan":    []any{nil, float64(0)},
						"jumlah_angka_kredit_utama":    []any{nil, nil},
						"jumlah_angka_kredit_tambahan": []any{nil, nil},
					},
					"action":     "UPDATE",
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
			name:          "ok: rollback on usulan perubahan data should not DELETE record",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1h")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "005"
			}`,
			doRollback:           true,
			wantResponsePostCode: http.StatusNoContent,
			wantResponseGetBody: `{
				"data": [
					{
						"id":         {id},
						"jenis_data": "riwayat-kepangkatan",
						"action":     "DELETE",
						"status":     "Diusulkan",
						"catatan":    "",
						"data_id":    "005",
						"perubahan_data": {
							"jenis_kp_id":                  [ null,  null ],
							"nama_jenis_kp":                [ null,  null ],
							"golongan_id":                  [ null,  null ],
							"nama_golongan":                [ null,  null ],
							"nama_golongan_pangkat":        [ null,  null ],
							"tmt_golongan":                 [ null,  null ],
							"nomor_sk":                     [ "SK1", null ],
							"tanggal_sk":                   [ null,  null ],
							"nomor_bkn":                    [ null,  null ],
							"tanggal_bkn":                  [ null,  null ],
							"masa_kerja_golongan_tahun":    [ null,  null ],
							"masa_kerja_golongan_bulan":    [ null,  null ],
							"jumlah_angka_kredit_utama":    [ null,  null ],
							"jumlah_angka_kredit_tambahan": [ null,  null ]
						}
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
			wantDBSvcRows: dbtest.Rows{
				{
					"id":                           "005",
					"jenis_kp_id":                  nil,
					"kode_jenis_kp":                nil,
					"jenis_kp":                     nil,
					"golongan_id":                  nil,
					"golongan_nama":                nil,
					"pangkat_nama":                 nil,
					"sk_nomor":                     "SK1",
					"no_bkn":                       nil,
					"jumlah_angka_kredit_utama":    nil,
					"jumlah_angka_kredit_tambahan": nil,
					"mk_golongan_tahun":            nil,
					"mk_golongan_bulan":            nil,
					"sk_tanggal":                   nil,
					"tanggal_bkn":                  nil,
					"tmt_golongan":                 nil,
					"status_satker":                nil,
					"status_biro":                  nil,
					"pangkat_terakhir":             nil,
					"bkn_id":                       nil,
					"file_base64":                  nil,
					"s3_file_id":                   nil,
					"keterangan_berkas":            nil,
					"arsip_id":                     nil,
					"golongan_asal":                nil,
					"basic":                        nil,
					"sk_type":                      nil,
					"kanreg":                       nil,
					"kpkn":                         nil,
					"keterangan":                   nil,
					"lpnk":                         nil,
					"jenis_riwayat":                nil,
					"pns_id":                       "id_1h",
					"pns_nip":                      "1h",
					"pns_nama":                     nil,
					"created_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":                   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":                   nil,
				},
			},
			wantDBUsulanRows: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1h",
					"jenis_data": "riwayat-kepangkatan",
					"data_id":    "005",
					"perubahan_data": map[string]any{
						"jenis_kp_id":                  []any{nil, nil},
						"nama_jenis_kp":                []any{nil, nil},
						"golongan_id":                  []any{nil, nil},
						"nama_golongan":                []any{nil, nil},
						"nama_golongan_pangkat":        []any{nil, nil},
						"tmt_golongan":                 []any{nil, nil},
						"nomor_sk":                     []any{"SK1", nil},
						"tanggal_sk":                   []any{nil, nil},
						"nomor_bkn":                    []any{nil, nil},
						"tanggal_bkn":                  []any{nil, nil},
						"masa_kerja_golongan_tahun":    []any{nil, nil},
						"masa_kerja_golongan_bulan":    []any{nil, nil},
						"jumlah_angka_kredit_utama":    []any{nil, nil},
						"jumlah_angka_kredit_tambahan": []any{nil, nil},
					},
					"action":     "DELETE",
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
			name:          "error: riwayat kepangkatan is owned by other pegawai",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1b")}},
			requestBody: `{
				"action": "DELETE",
				"data_id": "001"
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat kepangkatan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbtest.Rows{},
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: riwayat kepangkatan is not found",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "DELETE",
				"data_id": "000"
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat kepangkatan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: riwayat kepangkatan is deleted",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "UPDATE",
				"data_id": "002",
				"data": {
					"jenis_kp_id": 1,
					"golongan_id": 1,
					"tmt_golongan": null,
					"nomor_sk": "",
					"tanggal_sk": "2000-01-01",
					"nomor_bkn": "",
					"tanggal_bkn": null,
					"masa_kerja_golongan_tahun": 1,
					"masa_kerja_golongan_bulan": 1,
					"jumlah_angka_kredit_utama": null,
					"jumlah_angka_kredit_tambahan": null
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data riwayat kepangkatan tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: golongan or jenis kenaikan pangkat is not found",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"jenis_kp_id": 0,
					"golongan_id": 0,
					"nomor_sk": "",
					"tanggal_sk": "2000-01-01",
					"masa_kerja_golongan_tahun": 0,
					"masa_kerja_golongan_bulan": 0
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data golongan tidak ditemukan | data jenis kenaikan pangkat tidak ditemukan"}`,
			wantResponseGetBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
			wantDBSvcRows:    dbRows1a,
			wantDBUsulanRows: dbtest.Rows{},
		},
		{
			name:          "error: golongan or jenis kenaikan pangkat is deleted",
			requestHeader: http.Header{"Authorization": authHeader1a},
			requestBody: `{
				"action": "CREATE",
				"data": {
					"jenis_kp_id": 2,
					"golongan_id": 2,
					"tmt_golongan": null,
					"nomor_sk": "",
					"tanggal_sk": "2000-01-01",
					"nomor_bkn": "",
					"tanggal_bkn": null,
					"masa_kerja_golongan_tahun": 1,
					"masa_kerja_golongan_bulan": 1,
					"jumlah_angka_kredit_utama": null,
					"jumlah_angka_kredit_tambahan": null
				}
			}`,
			wantResponsePostCode: http.StatusBadRequest,
			wantResponsePostBody: `{"message": "data golongan tidak ditemukan | data jenis kenaikan pangkat tidak ditemukan"}`,
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
				`Error at \"/data/golongan_id\": property \"golongan_id\" is missing` +
				` | Error at \"/data/nomor_sk\": property \"nomor_sk\" is missing` +
				` | Error at \"/data/tanggal_sk\": property \"tanggal_sk\" is missing` +
				` | Error at \"/data/masa_kerja_golongan_tahun\": property \"masa_kerja_golongan_tahun\" is missing` +
				` | Error at \"/data/masa_kerja_golongan_bulan\": property \"masa_kerja_golongan_bulan\" is missing Or ` +
				`Error at \"/action\": value is not one of the allowed values [\"UPDATE\"]` +
				` | Error at \"/data/golongan_id\": property \"golongan_id\" is missing` +
				` | Error at \"/data/nomor_sk\": property \"nomor_sk\" is missing` +
				` | Error at \"/data/tanggal_sk\": property \"tanggal_sk\" is missing` +
				` | Error at \"/data/masa_kerja_golongan_tahun\": property \"masa_kerja_golongan_tahun\" is missing` +
				` | Error at \"/data/masa_kerja_golongan_bulan\": property \"masa_kerja_golongan_bulan\" is missing` +
				` | Error at \"/data_id\": property \"data_id\" is missing Or ` +
				`Error at \"/action\": value is not one of the allowed values [\"DELETE\"]` +
				` | property \"data\" is unsupported | Error at \"/data_id\": property \"data_id\" is missing"}`,
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
			req := httptest.NewRequest(http.MethodPost, "/v1/usulan-perubahan-data/riwayat-kepangkatan", strings.NewReader(tt.requestBody))
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

			actualSvcRows, err := dbtest.QueryWithClause(db, "riwayat_golongan", "where pns_nip = $1 order by id", nip)
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
			req2 := httptest.NewRequest(http.MethodGet, "/v1/usulan-perubahan-data/riwayat-kepangkatan", nil)
			req2.Header = tt.requestHeader
			rec2 := httptest.NewRecorder()

			e.ServeHTTP(rec2, req2)

			assert.Equal(t, http.StatusOK, rec2.Code)
			assert.JSONEq(t, tt.wantResponseGetBody, rec2.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec2, req2, e))
		})
	}
}
