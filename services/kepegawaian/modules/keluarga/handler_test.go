package keluarga

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
	-- Seed ref_agama
	INSERT INTO ref_agama (id, nama) VALUES 
	(1, 'Islam'),
	(2, 'Kristen');

	-- pegawai with nip_baru = 1c
	INSERT INTO pegawai (pns_id, nip_baru, nama)
	VALUES ('pns-1', '1c', 'Pegawai Test');

	-- orang_tua linked to pegawai
	INSERT INTO orang_tua (id, hubungan, nama, no_dokumen, agama_id, pns_id)
	VALUES (21, 1, 'Ayah A', '123', 1, 'pns-1'),
	       (22, 2, 'Ibu B', '456', 2, 'pns-1');

	-- pasangan linked to pegawai
	INSERT INTO pasangan (id, pns, nama, tanggal_menikah, status, pns_id)
	VALUES (31, 1, 'Istri A', '2000-01-01', 1, 'pns-1');

	-- anak linked to pegawai
	INSERT INTO anak (id, pasangan_id, nama, jenis_kelamin, tanggal_lahir, status_anak, pns_id)
	VALUES (11, 31, 'Anak A', 'M', '2000-01-01', '1', 'pns-1');

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
			name:             "ok: only nip 1c data returned",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": {
					"orang_tua": [
					{
						"id": 21,
						"nama": "Ayah A",
						"hubungan": "Ayah",
						"agama": "Islam",
						"nik": "123",
						"status_hidup": "Masih Hidup"
					},
					{
						"id": 22,
						"nama": "Ibu B",
						"hubungan": "Ibu",
						"agama": "Kristen",
						"nik": "456",
						"status_hidup": "Masih Hidup"
					}
					],

					"pasangan": [
					{
						"agama": null,
						"akte_cerai": null,
						"akte_meninggal": null,
						"akte_nikah": null,
						"id": 31,
						"karsus": "",
						"nama": "Istri A",
						"nik": "",
						"status_nikah": "Menikah",
						"status_pns": "PNS",
						"tanggal_cerai": null,
						"tanggal_menikah": "2000-01-01",
						"tanggal_meninggal": null,
						"tanggal_lahir": null
					}
					],
					"anak": [
					{
						"id": 11,
						"nama": "Anak A",
						"jenis_kelamin": "M",
						"status_anak": "Kandung",
						"anak_ke": 1,
						"nik": "",
						"status_sekolah": "",
						"nama_orang_tua": "Istri A",
						"tanggal_lahir": "2000-01-01T00:00:00Z"
					}
					]

				}
			}`,
		},
		{
			name:             "ok: nip 200 gets empty data",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"orang_tua":[],"pasangan":[],"anak":[]}}`,
		},
		{
			name:             "error: invalid token",
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
			dbRepo := repo.New(db)
			_, err := db.Exec(t.Context(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/keluarga", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, dbRepo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
