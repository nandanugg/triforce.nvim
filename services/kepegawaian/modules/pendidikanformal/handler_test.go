package pendidikanformal

import (
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
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO tingkat_pendidikan (id, nama, abbreviation, tingkat) VALUES
			(1, 'Sekolah Dasar', 'SD', 1),
			(2, 'Sekolah Menengah Pertama', 'SMP', 2),
			(3, 'Sekolah Menengah Atas', 'SMA', 3),
			(4, 'Diploma I', 'D1', 4),
			(5, 'Diploma II', 'D2', 5),
			(6, 'Diploma III', 'D3', 6),
			(7, 'Sarjana', 'S1', 7),
			(8, 'Magister', 'S2', 8),
			(9, 'Doktor', 'S3', 9);
		INSERT INTO pendidikan (id, tingkat_pendidikan_id, nama) VALUES
		('ed-003', 7, 'Akuntansi'),
		('ed-004', 8, 'Magister Manajemen'),
		('ed-006', 6, 'Diploma III Akuntansi');

		INSERT INTO pegawai (
		    id, pns_id, nip_baru, nama, gelar_depan, gelar_belakang, 
		    tgl_lahir, jenis_kelamin, tingkat_pendidikan_id
		) VALUES
		(1, 'pns-004', '198812252013014004', 'Maya Sari', NULL, 'S.E., M.M.', '1988-12-25', 'P', 8);

		INSERT INTO riwayat_pendidikan (
		    id, pns_id_3, pns_id, tingkat_pendidikan_id, pendidikan_id,
		    nama_sekolah, tahun_lulus, no_ijazah, gelar_depan, gelar_belakang,
		    tugas_belajar, negara_sekolah
		) VALUES
		(1, 'pns-004', 'pns-004', 6, 'ed-006', 'Politeknik Negeri Jakarta',
		 '2009', 'PNJ/AK/2009/004', NULL, 'A.Md.', 0, 'Pendidikan Regular'),
		(2, 'pns-004', 'pns-004', 7, 'ed-003', 'Universitas Airlangga',
		 '2011', 'UNAIR/AK/2011/004', NULL, 'S.E.', 0, 'Program Ekstensi'),
		(3, 'pns-004', 'pns-004', 8, 'ed-004', 'Universitas Airlangga',
		 '2016', 'UNAIR/MM/2016/004', NULL, 'M.M.', 1, 'Beasiswa Institusi');

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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198812252013014004")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 1,
					"jenjang_pendidikan": "Diploma III",
					"jurusan": "Diploma III Akuntansi",
					"nama_sekolah": "Politeknik Negeri Jakarta",
					"tahun_lulus": "2009",
					"nomor_ijazah": "PNJ/AK/2009/004",
					"gelar_depan": "",
					"gelar_belakang": "A.Md.",
					"tugas_belajar": "",
					"keterangan_pendidikan": "Pendidikan Regular"
				},
				{
					"id": 2,
					"jenjang_pendidikan": "Sarjana",
					"jurusan": "Akuntansi",
					"nama_sekolah": "Universitas Airlangga",
					"tahun_lulus": "2011",
					"nomor_ijazah": "UNAIR/AK/2011/004",
					"gelar_depan": "",
					"gelar_belakang": "S.E.",
					"tugas_belajar": "",
					"keterangan_pendidikan": "Program Ekstensi"
				},
				{
					"id": 3,
					"jenjang_pendidikan": "Magister",
					"jurusan": "Magister Manajemen",
					"nama_sekolah": "Universitas Airlangga",
					"tahun_lulus": "2016",
					"nomor_ijazah": "UNAIR/MM/2016/004",
					"gelar_depan": "",
					"gelar_belakang": "M.M.",
					"tugas_belajar": "Tugas Belajar",
					"keterangan_pendidikan": "Beasiswa Institusi"
				}			]
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": []}`,
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

			db := dbtest.NewPgxPool(t, dbmigrations.FS)
			dbRepository := repo.New(db)
			_, err := db.Exec(t.Context(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/pendidikan-formal", nil)
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
