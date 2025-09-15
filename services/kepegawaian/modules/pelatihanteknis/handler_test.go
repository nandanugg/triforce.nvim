package pelatihanteknis

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
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	dbData := `
		insert into pegawai (pns_id, nip_baru, deleted_at) values
		(1, '1c', null),
		(2, '2c', null),
		(3, '3c', '2020-01-01');
		
		insert into riwayat_kursus
			(id, pns_id, pns_nip, tipe_kursus, jenis_kursus, nama_kursus, tanggal_kursus, lama_kursus, institusi_penyelenggara, no_sertifikat, deleted_at) values
			(11, '1', '1c', 'Teknis', 'Workshop', '11a', '2000-01-01', 24, 'Institution 11', 'CERT11', null),
			(12, '1', '1c', 'Teknis', 'Seminar', '12a', '2001-01-01', 16, 'Institution 12', '', null),
			(13, '1', '1c', 'Teknis', 'Kursus', '13a', '2002-01-01', 40, 'Institution 13', 'CERT13', null),
			(14, '2', '2c', 'Teknis', 'Workshop', '14a', '2003-01-01', 8, 'Institution 14', 'CERT14', null),
			(15, '1', '1c', 'Teknis', 'Seminar', '15a', '2004-01-01', 4, 'Institution 15', 'CERT15', null),
			(16, '1', '1c', 'Teknis', 'Kursus', '16a', '2005-01-01', 20, 'Institution 16', null, null),
			(17, '1', '1c', 'Teknis', 'Workshop', '17a', '2006-01-01', 12, 'Institution 17', 'CERT17', null),
			-- Null test cases
			(18, '1', '1c', 'Teknis', 'Workshop', '18a', '2010-01-01', null, 'Institution 18', 'CERT18', null),
			(19, '1', '1c', 'Teknis', 'Seminar', '19a', null, 8, 'Institution 19', 'CERT19', null),
			(20, '1', '1c', null, null, '20a', '2012-01-01', 12, 'Institution 20', 'CERT20', null),
			(21, '1', '1c', 'Teknis', 'Workshop', null, '2013-01-01', 16, 'Institution 21', 'CERT21', null),
			(22, '1', '1c', 'Teknis', 'Seminar', '22a', '2014-01-01', 8, null, 'CERT22', null);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
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
						"durasi": 0,
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
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
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
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "2a")}},
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
			pgxconn := dbtest.NewPgxPool(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/pelatihan-teknis", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
