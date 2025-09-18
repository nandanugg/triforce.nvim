package riwayatkepangkatan

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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
		insert into ref_jenis_kp ("id", "nama", "deleted_at") values
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
						"id_jenis_kp":       0,
						"nama_jenis_kp":     "",
						"id_golongan":       0,
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
