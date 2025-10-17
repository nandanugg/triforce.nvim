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
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(t.Context(), dbData)
	require.NoError(t, err)

	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: only nip 1c data returned",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
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
							"agama": "",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
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
							"tanggal_lahir": "2000-01-01"
						}
					]
				}
			}`,
		},
		{
			name:             "ok: nip 200 gets empty data",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"orang_tua":[],"pasangan":[],"anak":[]}}`,
		},
		{
			name:             "error: invalid token",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/keluarga", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			dbRepo := repo.New(db)
			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
			RegisterRoutes(e, dbRepo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_listAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
	-- Seed ref_agama
	INSERT INTO ref_agama (id, nama, deleted_at) VALUES
	(1, 'Islam', null),
	(2, 'Kristen', null),
	(3, 'Hindu', now());

	-- pegawai with nip_baru = 1c
	INSERT INTO pegawai (pns_id, nip_baru, nama, deleted_at)
	VALUES ('pns-1', '1c', 'Pegawai Test', null),
	('pns-2', '1d', 'Pegawai Test 2', now());

	-- orang_tua linked to pegawai
	INSERT INTO orang_tua (id, hubungan, nama, no_dokumen, agama_id, pns_id, deleted_at)
	VALUES (21, 1, 'Ayah A', '123', 1, 'pns-1', null),
	       (22, 2, 'Ibu B', '456', 3, 'pns-1', null),
	       (23, 3, 'Ayah C', '789', 2, 'pns-1', now());

	-- pasangan linked to pegawai
	INSERT INTO pasangan (id, pns, nama, tanggal_menikah, status, pns_id, deleted_at)
	VALUES (31, 1, 'Istri A', '2000-01-01', 1, 'pns-1', null),
	       (32, 2, 'Istri B', '2000-01-01', 1, 'pns-1', now());

	-- anak linked to pegawai
	INSERT INTO anak (id, pasangan_id, nama, jenis_kelamin, tanggal_lahir, status_anak, pns_id, deleted_at)
	VALUES (11, 31, 'Anak A', 'M', '2000-01-01', '1', 'pns-1', null),
	       (12, 32, 'Anak B', 'F', '2000-01-01', '1', 'pns-1', null),
	       (13, 31, 'Anak C', 'M', '2000-01-01', '1', 'pns-1', now());
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(t.Context(), dbData)
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
			name:             "ok: only nip 1c data returned",
			nip:              "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
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
							"agama": "",
							"nik": "456",
							"status_hidup": "Masih Hidup"
						}
					],
					"pasangan": [
						{
							"agama": "",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
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
							"tanggal_lahir": "2000-01-01"
						},
						{
							"anak_ke": 2,
							"id": 12,
							"jenis_kelamin": "F",
							"nama": "Anak B",
							"nama_orang_tua": "",
							"nik": "",
							"status_anak": "Kandung",
							"status_sekolah": "",
							"tanggal_lahir": "2000-01-01"
						}
					]
				}
			}`,
		},
		{
			name:             "ok: nip 200 gets empty data",
			nip:              "200",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"orang_tua":[],"pasangan":[],"anak":[]}}`,
		},
		{
			name:             "ok: nip 1d gets empty data",
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"orang_tua":[],"pasangan":[],"anak":[]}}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/keluarga", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			dbRepo := repo.New(db)
			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
			RegisterRoutes(e, dbRepo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
