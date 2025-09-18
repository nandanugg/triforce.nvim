package riwayatpelatihanstruktural

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
	INSERT INTO "ref_jenis_diklat_struktural" (id, nama, deleted_at) VALUES
		(1, 'Jenis 1', null),
		(2, 'Jenis 2', null),
		(3, 'Jenis 3', '2000-01-01');

	INSERT INTO "riwayat_diklat_struktural" (
	    id, pns_nip, pns_nama, jenis_diklat_id, nama_diklat, nomor, tanggal, tahun, lama, institusi_penyelenggara, deleted_at
	) VALUES
	    ('uuid-diklat-struktural-001', '199001012020121001', 'Agus Purnomo', 1, 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00123', '2023-06-20', 2023, 900, 'Lembaga Administrasi Negara', null),
	    ('uuid-diklat-struktural-002', '199001012020121001', 'Siti Rahmawati', 2, 'Pelatihan Kepemimpinan Pengawas (PKP)', 'LAN-PKP-2022-00456', '2022-08-15', null, null, 'Badan Diklat Provinsi Jawa Barat', null),
	    ('uuid-diklat-struktural-003', '199001012020121001', 'Budi Santoso', 3, 'Pelatihan Kepemimpinan Nasional Tingkat II', 'LAN-PKNII-2021-00089', '2021-04-10', 2021, 1200, 'LAN-RI', null),
	    ('uuid-diklat-struktural-004', '199305202021121002', 'Dewi Kartika', 1, 'Pelatihan Kepemimpinan Administrator (PKA)', 'LAN-PKA-2023-00234', '2023-07-05', 2023, 900, 'Badan Pengembangan Sumber Daya Manusia Daerah (BPSDMD) DKI Jakarta', null),
	    ('uuid-diklat-struktural-005', '199001012020121001', 'Ahmad Fauzi', 1, 'Pelatihan Kepemimpinan Nasional Tingkat I', 'LAN-PKNI-2020-00077', '2020-09-12', 2020, 1500, 'Lembaga Administrasi Negara', '2000-01-01'),
			('uuid-diklat-struktural-006', '199001012020121001', 'Ahmad Fauzi', 1, 'Pelatihan Kepemimpinan Nasional Tingkat III', 'LAN-PKNI-2020-00077', null, 2022, 10, 'Lembaga Administrasi Negara', null);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "199001012020121001")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "uuid-diklat-struktural-001",
						"institusi_penyelenggara": "Lembaga Administrasi Negara",
						"jenis_diklat": "Jenis 1",
						"nama_diklat": "Pelatihan Kepemimpinan Administrator (PKA)",
						"nomor_sertifikat": "LAN-PKA-2023-00123",
						"tahun": 2023,
						"tanggal_mulai": "2023-06-20",
						"tanggal_selesai": "2023-07-27",
						"durasi": 900
					},
					{
						"id": "uuid-diklat-struktural-002",
						"institusi_penyelenggara": "Badan Diklat Provinsi Jawa Barat",
						"jenis_diklat": "Jenis 2",
						"nama_diklat": "Pelatihan Kepemimpinan Pengawas (PKP)",
						"nomor_sertifikat": "LAN-PKP-2022-00456",
						"tahun": 2022,
						"tanggal_mulai": "2022-08-15",
						"tanggal_selesai": "2022-08-15",
						"durasi": 0
					},
					{
						"id": "uuid-diklat-struktural-003",
						"institusi_penyelenggara": "LAN-RI",
						"jenis_diklat": "",
						"nama_diklat": "Pelatihan Kepemimpinan Nasional Tingkat II",
						"nomor_sertifikat": "LAN-PKNII-2021-00089",
						"tahun": 2021,
						"tanggal_mulai": "2021-04-10",
						"tanggal_selesai": "2021-05-30",
						"durasi": 1200
					},
					{
						"id": "uuid-diklat-struktural-006",
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "199001012020121001")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "uuid-diklat-struktural-002",
						"institusi_penyelenggara": "Badan Diklat Provinsi Jawa Barat",
						"jenis_diklat": "Jenis 2",
						"nama_diklat": "Pelatihan Kepemimpinan Pengawas (PKP)",
						"nomor_sertifikat": "LAN-PKP-2022-00456",
						"tahun": 2022,
						"tanggal_mulai": "2022-08-15",
						"tanggal_selesai": "2022-08-15",
						"durasi": 0
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}
			`,
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
			dbRepository := repo.New(db)
			_, err := db.Exec(t.Context(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-pelatihan-struktural", nil)
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
