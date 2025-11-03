package keluarga

import (
	"context"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	INSERT INTO ref_agama (id, nama, deleted_at) VALUES
		(1, 'Islam', null),
		(2, 'Kristen', null),
		(3, 'Katolik', '2000-01-01');

	INSERT INTO ref_jenis_kawin (id, nama, deleted_at) VALUES
		(1, 'Menikah', null),
		(2, 'Cerai', null),
		(3, 'Duda', '2000-01-01');

	INSERT INTO pegawai (pns_id, nip_baru, nama, deleted_at) VALUES
		('pns-1', '1c', 'Pegawai Test', null),
		('pns-2', '200', 'Pegawai Test 2', '2000-01-01'),
		('pns-3', '1d', 'Pegawai 1d', null);

	INSERT INTO orang_tua (id, hubungan, nama, no_dokumen, agama_id, pns_id, tanggal_meninggal, akte_meninggal, deleted_at) VALUES
		(21, 1, 'Ayah A', '123', 1, 'pns-1', null, null, null),
		(22, 2, 'Ibu B', '456', 2, 'pns-1', null, '', null),
		(23, 2, 'Ibu C', '789', 1, 'pns-1', null, null, '2000-01-01'),
		(24, 3, 'Ayah D', null, 3, 'pns-1', '2000-01-01', '1ab', null),
		(25, 1, 'Ayah E', '000', 1, 'pns-2', null, null, null),
		(26, 1, 'Ayah F', null, null, 'pns-3', null, null, null);

	INSERT INTO pasangan (id, pns, nama, nik, tanggal_lahir, tanggal_menikah, tanggal_cerai, tanggal_meninggal, status, agama_id, karsus, pns_id, deleted_at) VALUES
		(31, 1, 'Istri A', '1a', '1980-01-01', '2000-01-01', null, null, 1, 2, 'abc', 'pns-1', null),
		(32, 0, 'Istri B', '1b', null, '1999-01-01', '2000-01-01', null, 2, 1, null, 'pns-1', null),
		(33, 0, 'Istri C', '1c', null, '2000-01-01', null, null, 1, 3, null, 'pns-1', '2000-01-01'),
		(34, 1, 'Istri D', '1d', null, null, null, '2000-01-01', 3, 3, 'def', 'pns-1', null),
		(35, 1, 'Istri E', '1e', '1980-01-01', '2000-01-01', null, null, 1, 2, 'def', 'pns-2', null),
		(36, null, 'Istri F', null, null, null, null, null, null, null, null, 'pns-3', null);

	INSERT INTO anak (id, pasangan_id, nama, nik, jenis_kelamin, tanggal_lahir, status_anak, status_sekolah, jenis_kawin_id, agama_id, anak_ke, pns_id, deleted_at) VALUES
		(11, 31, 'Anak A', '1a', 'M', '2000-01-01', '1', 1, 1, 1, 2, 'pns-1', null),
		(12, 31, 'Anak B', '1b', 'M', '2001-01-01', '1', 1, 1, 1, 0, 'pns-1', '2000-01-01'),
		(13, 32, 'Anak C', '1c', 'F', '2001-01-01', '2', 2, 2, 2, null, 'pns-1', null),
		(14, 33, 'Anak D', '1d', 'F', '2002-01-01', '3', 3, 3, 3, 1, 'pns-1', null),
		(15, 35, 'Anak E', '1e', 'M', '1999-01-01', '1', 1, 1, 1, 0, 'pns-2', null),
		(16, null, 'Anak F', null, null, null, null, null, null, null, 0, 'pns-3', null),
		(17, 34, 'Anak G', null, null, '1990-01-01', '', 1, 1, 1, null, 'pns-1', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(t.Context(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	dbRepo := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, dbRepo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			wantResponseBody: `{
				"data": {
					"orang_tua": [
						{
							"id": 21,
							"nama": "Ayah A",
							"hubungan": "Ayah",
							"agama_id": 1,
							"agama": "Islam",
							"nik": "123",
							"status_hidup": "Masih Hidup",
							"tanggal_meninggal": null,
							"akte_meninggal": ""
						},
						{
							"id": 22,
							"nama": "Ibu B",
							"hubungan": "Ibu",
							"agama_id": 2,
							"agama": "Kristen",
							"nik": "456",
							"status_hidup": "Masih Hidup",
							"tanggal_meninggal": null,
							"akte_meninggal": ""
						},
						{
							"id": 24,
							"nama": "Ayah D",
							"hubungan": "",
							"agama_id": 3,
							"agama": "",
							"nik": "",
							"status_hidup": "Sudah Meninggal",
							"tanggal_meninggal": "2000-01-01",
							"akte_meninggal": "1ab"
						}
					],
					"pasangan": [
						{
							"id": 31,
							"agama_id": 2,
							"agama": "Kristen",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "abc",
							"nama": "Istri A",
							"nik": "1a",
							"status_pernikahan_id": 1,
							"status_nikah": "Menikah",
							"status_pns": "PNS",
							"tanggal_cerai": null,
							"tanggal_menikah": "2000-01-01",
							"tanggal_meninggal": null,
							"tanggal_lahir": "1980-01-01"
						},
						{
							"id": 32,
							"agama_id": 1,
							"agama": "Islam",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "",
							"nama": "Istri B",
							"nik": "1b",
							"status_pernikahan_id": 2,
							"status_nikah": "Cerai",
							"status_pns": "Bukan PNS",
							"tanggal_cerai": "2000-01-01",
							"tanggal_menikah": "1999-01-01",
							"tanggal_meninggal": null,
							"tanggal_lahir": null
						},
						{
							"id": 34,
							"agama_id": 3,
							"agama": "",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "def",
							"nama": "Istri D",
							"nik": "1d",
							"status_pernikahan_id": 3,
							"status_nikah": "",
							"status_pns": "PNS",
							"tanggal_cerai": null,
							"tanggal_menikah": null,
							"tanggal_meninggal": "2000-01-01",
							"tanggal_lahir": null
						}
					],
					"anak": [
						{
							"id": 14,
							"nama": "Anak D",
							"jenis_kelamin": "F",
							"status_anak": "",
							"anak_ke": 1,
							"nik": "1d",
							"status_sekolah": "",
							"pasangan_orang_tua_id": 33,
							"nama_orang_tua": "",
							"tanggal_lahir": "2002-01-01",
							"agama_id": 3,
							"agama": "",
							"status_pernikahan_id": 3,
							"status_pernikahan": ""
						},
						{
							"id": 11,
							"nama": "Anak A",
							"jenis_kelamin": "M",
							"status_anak": "Kandung",
							"anak_ke": 2,
							"nik": "1a",
							"status_sekolah": "Masih Sekolah",
							"pasangan_orang_tua_id": 31,
							"nama_orang_tua": "Istri A",
							"tanggal_lahir": "2000-01-01",
							"agama_id": 1,
							"agama": "Islam",
							"status_pernikahan_id": 1,
							"status_pernikahan": "Menikah"
						},
						{
							"id": 17,
							"nama": "Anak G",
							"jenis_kelamin": "",
							"status_anak": "",
							"anak_ke": null,
							"nik": "",
							"status_sekolah": "Masih Sekolah",
							"pasangan_orang_tua_id": 34,
							"nama_orang_tua": "Istri D",
							"tanggal_lahir": "1990-01-01",
							"agama_id": 1,
							"agama": "Islam",
							"status_pernikahan_id": 1,
							"status_pernikahan": "Menikah"
						},
						{
							"id": 13,
							"nama": "Anak C",
							"jenis_kelamin": "F",
							"status_anak": "Angkat",
							"anak_ke": null,
							"nik": "1c",
							"status_sekolah": "Sudah Bekerja",
							"pasangan_orang_tua_id": 32,
							"nama_orang_tua": "Istri B",
							"tanggal_lahir": "2001-01-01",
							"agama_id": 2,
							"agama": "Kristen",
							"status_pernikahan_id": 2,
							"status_pernikahan": "Cerai"
						}
					]
				}
			}`,
		},
		{
			name:             "ok: with null references",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1d")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"orang_tua": [
						{
							"id": 26,
							"nama": "Ayah F",
							"hubungan": "Ayah",
							"agama_id": null,
							"agama": "",
							"nik": "",
							"status_hidup": "Masih Hidup",
							"tanggal_meninggal": null,
							"akte_meninggal": ""
						}
					],
					"pasangan": [
						{
							"id": 36,
							"agama_id": null,
							"agama": "",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "",
							"nama": "Istri F",
							"nik": "",
							"status_pernikahan_id": null,
							"status_nikah": "",
							"status_pns": "Bukan PNS",
							"tanggal_cerai": null,
							"tanggal_menikah": null,
							"tanggal_meninggal": null,
							"tanggal_lahir": null
						}
					],
					"anak": [
						{
							"id": 16,
							"nama": "Anak F",
							"jenis_kelamin": "",
							"status_anak": "",
							"anak_ke": 0,
							"nik": "",
							"status_sekolah": "",
							"pasangan_orang_tua_id": null,
							"nama_orang_tua": "",
							"tanggal_lahir": null,
							"agama_id": null,
							"agama": "",
							"status_pernikahan_id": null,
							"status_pernikahan": ""
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
	INSERT INTO ref_agama (id, nama, deleted_at) VALUES
		(1, 'Islam', null),
		(2, 'Kristen', null),
		(3, 'Katolik', '2000-01-01');

	INSERT INTO ref_jenis_kawin (id, nama, deleted_at) VALUES
		(1, 'Menikah', null),
		(2, 'Cerai', null),
		(3, 'Duda', '2000-01-01');

	INSERT INTO pegawai (pns_id, nip_baru, nama, deleted_at) VALUES
		('pns-1', '1c', 'Pegawai Test', null),
		('pns-2', '200', 'Pegawai Test 2', '2000-01-01'),
		('pns-3', '1d', 'Pegawai 1d', null);

	INSERT INTO orang_tua (id, hubungan, nama, no_dokumen, agama_id, pns_id, tanggal_meninggal, akte_meninggal, deleted_at) VALUES
		(21, 1, 'Ayah A', '123', 1, 'pns-1', null, null, null),
		(22, 2, 'Ibu B', '456', 2, 'pns-1', null, '', null),
		(23, 2, 'Ibu C', '789', 1, 'pns-1', null, null, '2000-01-01'),
		(24, 3, 'Ayah D', null, 3, 'pns-1', '2000-01-01', '1ab', null),
		(25, 1, 'Ayah E', '000', 1, 'pns-2', null, null, null),
		(26, 1, 'Ayah F', null, null, 'pns-3', null, null, null);

	INSERT INTO pasangan (id, pns, nama, nik, tanggal_lahir, tanggal_menikah, tanggal_cerai, tanggal_meninggal, status, agama_id, karsus, pns_id, deleted_at) VALUES
		(31, 1, 'Istri A', '1a', '1980-01-01', '2000-01-01', null, null, 1, 2, 'abc', 'pns-1', null),
		(32, 0, 'Istri B', '1b', null, '1999-01-01', '2000-01-01', null, 2, 1, null, 'pns-1', null),
		(33, 0, 'Istri C', '1c', null, '2000-01-01', null, null, 1, 3, null, 'pns-1', '2000-01-01'),
		(34, 1, 'Istri D', '1d', null, null, null, '2000-01-01', 3, 3, 'def', 'pns-1', null),
		(35, 1, 'Istri E', '1e', '1980-01-01', '2000-01-01', null, null, 1, 2, 'def', 'pns-2', null),
		(36, null, 'Istri F', null, null, null, null, null, null, null, null, 'pns-3', null);

	INSERT INTO anak (id, pasangan_id, nama, nik, jenis_kelamin, tanggal_lahir, status_anak, status_sekolah, jenis_kawin_id, agama_id, anak_ke, pns_id, deleted_at) VALUES
		(11, 31, 'Anak A', '1a', 'M', '2000-01-01', '1', 1, 1, 1, 2, 'pns-1', null),
		(12, 31, 'Anak B', '1b', 'M', '2001-01-01', '1', 1, 1, 1, 0, 'pns-1', '2000-01-01'),
		(13, 32, 'Anak C', '1c', 'F', '2001-01-01', '2', 2, 2, 2, null, 'pns-1', null),
		(14, 33, 'Anak D', '1d', 'F', '2002-01-01', '3', 3, 3, 3, 1, 'pns-1', null),
		(15, 35, 'Anak E', '1e', 'M', '1999-01-01', '1', 1, 1, 1, 0, 'pns-2', null),
		(16, null, 'Anak F', null, null, null, null, null, null, null, 0, 'pns-3', null),
		(17, 34, 'Anak G', null, null, '1990-01-01', '', 1, 1, 1, null, 'pns-1', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(t.Context(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	dbRepo := repo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, dbRepo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			wantResponseBody: `{
				"data": {
					"orang_tua": [
						{
							"id": 21,
							"nama": "Ayah A",
							"hubungan": "Ayah",
							"agama_id": 1,
							"agama": "Islam",
							"nik": "123",
							"status_hidup": "Masih Hidup",
							"tanggal_meninggal": null,
							"akte_meninggal": ""
						},
						{
							"id": 22,
							"nama": "Ibu B",
							"hubungan": "Ibu",
							"agama_id": 2,
							"agama": "Kristen",
							"nik": "456",
							"status_hidup": "Masih Hidup",
							"tanggal_meninggal": null,
							"akte_meninggal": ""
						},
						{
							"id": 24,
							"nama": "Ayah D",
							"hubungan": "",
							"agama_id": 3,
							"agama": "",
							"nik": "",
							"status_hidup": "Sudah Meninggal",
							"tanggal_meninggal": "2000-01-01",
							"akte_meninggal": "1ab"
						}
					],
					"pasangan": [
						{
							"id": 31,
							"agama_id": 2,
							"agama": "Kristen",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "abc",
							"nama": "Istri A",
							"nik": "1a",
							"status_pernikahan_id": 1,
							"status_nikah": "Menikah",
							"status_pns": "PNS",
							"tanggal_cerai": null,
							"tanggal_menikah": "2000-01-01",
							"tanggal_meninggal": null,
							"tanggal_lahir": "1980-01-01"
						},
						{
							"id": 32,
							"agama_id": 1,
							"agama": "Islam",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "",
							"nama": "Istri B",
							"nik": "1b",
							"status_pernikahan_id": 2,
							"status_nikah": "Cerai",
							"status_pns": "Bukan PNS",
							"tanggal_cerai": "2000-01-01",
							"tanggal_menikah": "1999-01-01",
							"tanggal_meninggal": null,
							"tanggal_lahir": null
						},
						{
							"id": 34,
							"agama_id": 3,
							"agama": "",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "def",
							"nama": "Istri D",
							"nik": "1d",
							"status_pernikahan_id": 3,
							"status_nikah": "",
							"status_pns": "PNS",
							"tanggal_cerai": null,
							"tanggal_menikah": null,
							"tanggal_meninggal": "2000-01-01",
							"tanggal_lahir": null
						}
					],
					"anak": [
						{
							"id": 14,
							"nama": "Anak D",
							"jenis_kelamin": "F",
							"status_anak": "",
							"anak_ke": 1,
							"nik": "1d",
							"status_sekolah": "",
							"pasangan_orang_tua_id": 33,
							"nama_orang_tua": "",
							"tanggal_lahir": "2002-01-01",
							"agama_id": 3,
							"agama": "",
							"status_pernikahan_id": 3,
							"status_pernikahan": ""
						},
						{
							"id": 11,
							"nama": "Anak A",
							"jenis_kelamin": "M",
							"status_anak": "Kandung",
							"anak_ke": 2,
							"nik": "1a",
							"status_sekolah": "Masih Sekolah",
							"pasangan_orang_tua_id": 31,
							"nama_orang_tua": "Istri A",
							"tanggal_lahir": "2000-01-01",
							"agama_id": 1,
							"agama": "Islam",
							"status_pernikahan_id": 1,
							"status_pernikahan": "Menikah"
						},
						{
							"id": 17,
							"nama": "Anak G",
							"jenis_kelamin": "",
							"status_anak": "",
							"anak_ke": null,
							"nik": "",
							"status_sekolah": "Masih Sekolah",
							"pasangan_orang_tua_id": 34,
							"nama_orang_tua": "Istri D",
							"tanggal_lahir": "1990-01-01",
							"agama_id": 1,
							"agama": "Islam",
							"status_pernikahan_id": 1,
							"status_pernikahan": "Menikah"
						},
						{
							"id": 13,
							"nama": "Anak C",
							"jenis_kelamin": "F",
							"status_anak": "Angkat",
							"anak_ke": null,
							"nik": "1c",
							"status_sekolah": "Sudah Bekerja",
							"pasangan_orang_tua_id": 32,
							"nama_orang_tua": "Istri B",
							"tanggal_lahir": "2001-01-01",
							"agama_id": 2,
							"agama": "Kristen",
							"status_pernikahan_id": 2,
							"status_pernikahan": "Cerai"
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
			name:             "ok: nip 1d with null references",
			nip:              "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"orang_tua": [
						{
							"id": 26,
							"nama": "Ayah F",
							"hubungan": "Ayah",
							"agama_id": null,
							"agama": "",
							"nik": "",
							"status_hidup": "Masih Hidup",
							"tanggal_meninggal": null,
							"akte_meninggal": ""
						}
					],
					"pasangan": [
						{
							"id": 36,
							"agama_id": null,
							"agama": "",
							"akte_cerai": "",
							"akte_meninggal": "",
							"akte_nikah": "",
							"karsus": "",
							"nama": "Istri F",
							"nik": "",
							"status_pernikahan_id": null,
							"status_nikah": "",
							"status_pns": "Bukan PNS",
							"tanggal_cerai": null,
							"tanggal_menikah": null,
							"tanggal_meninggal": null,
							"tanggal_lahir": null
						}
					],
					"anak": [
						{
							"id": 16,
							"nama": "Anak F",
							"jenis_kelamin": "",
							"status_anak": "",
							"anak_ke": 0,
							"nik": "",
							"status_sekolah": "",
							"pasangan_orang_tua_id": null,
							"nama_orang_tua": "",
							"tanggal_lahir": null,
							"agama_id": null,
							"agama": "",
							"status_pernikahan_id": null,
							"status_pernikahan": ""
						}
					]
				}
			}`,
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

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateOrangTua(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01');
		insert into ref_agama
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: with all data",
			dbData: seedData + `
				insert into orang_tua
					(nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					('John Doe', 1,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": 1,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 2 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "John Doe",
					"hubungan":          int16(1),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
				{
					"id":                int32(2),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          int16(1),
					"tanggal_meninggal": time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					"akte_meninggal":    "akte_01",
					"jenis_dokumen":     "KTP",
					"no_dokumen":        "123",
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"hubungan": "Ibu",
				"nik": "",
				"agama_id": null,
				"tanggal_meninggal": null,
				"akte_meninggal": ""
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"hubungan": "Ayah"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "John Doe",
					"hubungan":          int16(1),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			dbData:        seedData,
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": 1,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			dbData:        seedData,
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": null,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: agama is not found",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"agama_id": 0,
				"hubungan": "Ayah",
				"tanggal_meninggal": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: agama is deleted",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "123456",
				"agama_id": 2,
				"hubungan": "Ayah",
				"tanggal_meninggal": null,
				"akte_meninggal": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: invalid format date, exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 256) + `",
				"nik": "` + strings.Repeat(".", 21) + `",
				"agama_id": "Islam",
				"hubungan": "Paman",
				"tanggal_meninggal": "",
				"akte_meninggal": "` + strings.Repeat(".", 101) + `"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"agama_id\" harus dalam tipe integer` +
				` | parameter \"akte_meninggal\" harus 100 karakter atau kurang` +
				` | parameter \"hubungan\" harus salah satu dari \"Ayah\", \"Ibu\"` +
				` | parameter \"nama\" harus 255 karakter atau kurang` +
				` | parameter \"nik\" harus 20 karakter atau kurang` +
				` | parameter \"tanggal_meninggal\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null on required params",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": null,
				"hubungan": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"hubungan\" harus salah satu dari \"Ayah\", \"Ibu\"` +
				` | parameter \"nama\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama\" harus diisi` +
				` | parameter \"hubungan\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/orang-tua", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "orang_tua", "id")
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminUpdateOrangTua(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into ref_agama
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: with all data",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, gelar_depan, gelar_belakang, tempat_lahir, tanggal_lahir, email,           pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 1,        'Dr.',       'S.Kom',        'Jakarta',    '1990-01-01',  'test@test.com', 'id_1c', '1c', '2000-01-01', '2000-01-01'),
					(2,  'Jane Wee', 2,        null,        null,           null,         null,          null,            'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": 1,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          int16(1),
					"tanggal_meninggal": time.Date(2020, 1, 2, 0, 0, 0, 0, time.UTC),
					"akte_meninggal":    "akte_01",
					"jenis_dokumen":     "KTP",
					"no_dokumen":        "123",
					"gelar_depan":       "Dr.",
					"gelar_belakang":    "S.Kom",
					"tempat_lahir":      "Jakarta",
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"email":             "test@test.com",
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
				{
					"id":                int32(2),
					"nama":              "Jane Wee",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "ok: with null values",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, agama_id, tanggal_meninggal, akte_meninggal, jenis_dokumen, no_dokumen, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        1,        '2000-01-01',      'akte-01',      'KTP',         '123',      'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"hubungan": "Ibu",
				"nik": "",
				"agama_id": null,
				"tanggal_meninggal": null,
				"akte_meninggal": ""
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "ok: required data only",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, agama_id, tanggal_meninggal, akte_meninggal, jenis_dokumen, no_dokumen, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        1,        '2000-01-01',      'akte-01',      'KTP',         '123',      'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"hubungan": "Ayah"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "John Doe",
					"hubungan":          int16(1),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: orang tua is not found",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": 1,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: orang tua is owned by different pegawai",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1e', '1e', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": 1,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          nil,
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1e",
					"nip":               "1e",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: orang tua is deleted",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": 1,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "error: pegawai is deleted",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        'id_1d', '1d', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1d",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"agama_id": null,
				"hubungan": "Ibu",
				"tanggal_meninggal": "2020-01-02",
				"akte_meninggal": "akte_01"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1d",
					"nip":               "1d",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: agama is not found",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, agama_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"agama_id": 0,
				"hubungan": "Ayah",
				"tanggal_meninggal": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          int16(2),
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: agama is deleted",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, agama_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        1,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "123456",
				"agama_id": 2,
				"hubungan": "Ayah",
				"tanggal_meninggal": null,
				"akte_meninggal": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          int16(1),
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: invalid format date, exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 256) + `",
				"nik": "` + strings.Repeat(".", 21) + `",
				"agama_id": "Islam",
				"hubungan": "Paman",
				"tanggal_meninggal": "",
				"akte_meninggal": "` + strings.Repeat(".", 101) + `"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai` +
				` | parameter \"agama_id\" harus dalam tipe integer` +
				` | parameter \"akte_meninggal\" harus 100 karakter atau kurang` +
				` | parameter \"hubungan\" harus salah satu dari \"Ayah\", \"Ibu\"` +
				` | parameter \"nama\" harus 255 karakter atau kurang` +
				` | parameter \"nik\" harus 20 karakter atau kurang` +
				` | parameter \"tanggal_meninggal\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null on required params",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": null,
				"hubungan": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"hubungan\" harus salah satu dari \"Ayah\", \"Ibu\"` +
				` | parameter \"nama\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama\" harus diisi` +
				` | parameter \"hubungan\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/orang-tua/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "orang_tua", "id")
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

func Test_handler_adminDeleteOrangTua(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: success delete",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       jenis_dokumen, no_dokumen, hubungan, gelar_depan, gelar_belakang, tempat_lahir, tanggal_lahir, email,           pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 'KTP',         '123',      1,        'Dr.',       'S.Kom',        'Jakarta',    '1990-01-01',  'test@test.com', 'id_1c', '1c', '2000-01-01', '2000-01-01'),
					(2,  'Jane Doe', null,          null,       2,        null,        null,           null,         null,          null,            'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "John Doe",
					"hubungan":          int16(1),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     "KTP",
					"no_dokumen":        "123",
					"gelar_depan":       "Dr.",
					"gelar_belakang":    "S.Kom",
					"tempat_lahir":      "Jakarta",
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"email":             "test@test.com",
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        "{deleted_at}",
				},
				{
					"id":                int32(2),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: orang tua is not found",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: orang tua is deleted",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "error: pegawai is deleted",
			dbData: seedData + `
				insert into orang_tua
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        'id_1d', '1d', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1d",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int32(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"agama_id":          nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"jenis_dokumen":     nil,
					"no_dokumen":        nil,
					"gelar_depan":       nil,
					"gelar_belakang":    nil,
					"tempat_lahir":      nil,
					"tanggal_lahir":     nil,
					"email":             nil,
					"pns_id":            "id_1d",
					"nip":               "1d",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:             "error: unexpected data type",
			paramNIP:         "1c",
			paramID:          "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/orang-tua/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "orang_tua", "id")
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

func Test_handler_adminCreatePasangan(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01');
		insert into ref_agama
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into ref_jenis_kawin
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: with all data",
			dbData: seedData + `
				insert into pasangan
					(nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					('John Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"tanggal_lahir": "1990-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Istri",
				"nik": "123456",
				"is_pns": true,
				"no_karsus": "121212",
				"agama_id": 1,
				"tanggal_menikah": "2000-01-01",
				"akte_nikah": "akte-01",
				"tanggal_meninggal": "2001-02-02",
				"akte_meninggal": "akte-02",
				"tanggal_cerai": "2002-03-03",
				"akte_cerai": "akte-03"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 2 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "John Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
				{
					"id":                int64(2),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(1),
					"nik":               "123456",
					"agama_id":          int16(1),
					"tanggal_menikah":   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"akte_nikah":        "akte-01",
					"tanggal_cerai":     time.Date(2002, 3, 3, 0, 0, 0, 0, time.UTC),
					"akte_cerai":        "akte-03",
					"tanggal_meninggal": time.Date(2001, 2, 2, 0, 0, 0, 0, time.UTC),
					"akte_meninggal":    "akte-02",
					"karsus":            "121212",
					"status":            int16(1),
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"tanggal_lahir": "1990-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Istri",
				"nik": "",
				"is_pns": false,
				"no_karsus": "",
				"agama_id": null,
				"tanggal_menikah": null,
				"akte_nikah": "",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(0),
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            int16(1),
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"tanggal_lahir": "1990-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Suami"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "John Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(0),
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            int16(1),
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        "{created_at}",
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			dbData:        seedData,
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"tanggal_lahir": "2000-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Suami"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			dbData:        seedData,
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "",
				"is_pns": true,
				"tanggal_lahir": "2000-01-01",
				"no_karsus": "",
				"agama_id": null,
				"status_pernikahan_id": 1,
				"hubungan": "Istri",
				"tanggal_menikah": null,
				"akte_nikah": "",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: agama or status pernikahan is not found",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"tanggal_lahir": "2000-01-01",
				"agama_id": 0,
				"status_pernikahan_id": 0,
				"hubungan": "Suami"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan | data status pernikahan tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: agama or status pernikahan is deleted",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "123456",
				"is_pns": false,
				"tanggal_lahir": "2000-01-01",
				"no_karsus": "",
				"agama_id": 2,
				"status_pernikahan_id": 2,
				"hubungan": "Suami",
				"tanggal_menikah": "2000-01-01",
				"akte_nikah": "akte_01",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan | data status pernikahan tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: invalid format date, exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 101) + `",
				"nik": "` + strings.Repeat(".", 21) + `",
				"is_pns": "",
				"tanggal_lahir": "",
				"no_karsus": "` + strings.Repeat(".", 101) + `",
				"agama_id": "Islam",
				"status_pernikahan_id": "Menikah",
				"hubungan": "Pelakor",
				"tanggal_menikah": "",
				"akte_nikah": "` + strings.Repeat(".", 101) + `",
				"tanggal_meninggal": "",
				"akte_meninggal": "` + strings.Repeat(".", 101) + `",
				"tanggal_cerai": "",
				"akte_cerai": "` + strings.Repeat(".", 101) + `"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"agama_id\" harus dalam tipe integer` +
				` | parameter \"akte_cerai\" harus 100 karakter atau kurang` +
				` | parameter \"akte_meninggal\" harus 100 karakter atau kurang` +
				` | parameter \"akte_nikah\" harus 100 karakter atau kurang` +
				` | parameter \"hubungan\" harus salah satu dari \"Istri\", \"Suami\"` +
				` | parameter \"is_pns\" harus dalam tipe boolean` +
				` | parameter \"nama\" harus 100 karakter atau kurang` +
				` | parameter \"nik\" harus 20 karakter atau kurang` +
				` | parameter \"no_karsus\" harus 100 karakter atau kurang` +
				` | parameter \"status_pernikahan_id\" harus dalam tipe integer` +
				` | parameter \"tanggal_cerai\" harus dalam format date` +
				` | parameter \"tanggal_lahir\" harus dalam format date` +
				` | parameter \"tanggal_menikah\" harus dalam format date` +
				` | parameter \"tanggal_meninggal\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null on required params",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": null,
				"hubungan": null,
				"tanggal_lahir": null,
				"status_pernikahan_id": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"hubungan\" harus salah satu dari \"Istri\", \"Suami\"` +
				` | parameter \"nama\" tidak boleh null` +
				` | parameter \"status_pernikahan_id\" tidak boleh null` +
				` | parameter \"tanggal_lahir\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama\" harus diisi` +
				` | parameter \"tanggal_lahir\" harus diisi` +
				` | parameter \"status_pernikahan_id\" harus diisi` +
				` | parameter \"hubungan\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/pasangan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "pasangan", "id")
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminUpdatePasangan(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into ref_agama
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into ref_jenis_kawin
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: with all data",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01'),
					(2,  'John Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"is_pns": true,
				"tanggal_lahir": "1990-01-01",
				"no_karsus": "k.01",
				"agama_id": 1,
				"status_pernikahan_id": 1,
				"hubungan": "Istri",
				"tanggal_menikah": "2000-01-01",
				"akte_nikah": "akte-01",
				"tanggal_meninggal": "2001-01-01",
				"akte_meninggal": "akte-02",
				"tanggal_cerai": "2002-01-01",
				"akte_cerai": "akte-03"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(1),
					"nik":               "123",
					"agama_id":          int16(1),
					"tanggal_menikah":   time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"akte_nikah":        "akte-01",
					"tanggal_cerai":     time.Date(2002, 1, 1, 0, 0, 0, 0, time.UTC),
					"akte_cerai":        "akte-03",
					"tanggal_meninggal": time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC),
					"akte_meninggal":    "akte-02",
					"karsus":            "k.01",
					"status":            int16(1),
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
				{
					"id":                int64(2),
					"nama":              "John Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "ok: with null values",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, agama_id, nik,   karsus, tanggal_menikah, akte_nikah, tanggal_cerai, akte_cerai, tanggal_meninggal, akte_meninggal, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        1,        '123', '12',   '2000-01-01',    'akte-01',  '2000-01-01',  'akte-01',  '2000-01-01',      'akte-01',      'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "",
				"is_pns": false,
				"tanggal_lahir": "2000-01-01",
				"no_karsus": "",
				"agama_id": null,
				"status_pernikahan_id": 1,
				"hubungan": "Suami",
				"tanggal_menikah": null,
				"akte_nikah": "",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "John Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(0),
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            int16(1),
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "ok: required data only",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, agama_id, nik,   karsus, tanggal_menikah, akte_nikah, tanggal_cerai, akte_cerai, tanggal_meninggal, akte_meninggal, pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 2,        1,        '123', '12',   '2000-01-01',    'akte-01',  '2000-01-01',  'akte-01',  '2000-01-01',      'akte-01',      'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"tanggal_lahir": "2000-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Istri"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(0),
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            int16(1),
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        "{updated_at}",
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: pasangan is not found",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"tanggal_lahir": "2000-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Suami"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: pasangan is owned by different pegawai",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1e', '1e', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"tanggal_lahir": "2000-01-01",
				"status_pernikahan_id": 1,
				"hubungan": "Suami"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          nil,
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1e",
					"nip":               "1e",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: pasangan is deleted",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "",
				"is_pns": true,
				"tanggal_lahir": "2000-01-01",
				"no_karsus": "",
				"agama_id": null,
				"status_pernikahan_id": 1,
				"hubungan": "Istri",
				"tanggal_menikah": null,
				"akte_nikah": "",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "error: pegawai is deleted",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        'id_1d', '1d', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1d",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "",
				"is_pns": true,
				"tanggal_lahir": "2000-01-01",
				"no_karsus": "",
				"agama_id": null,
				"status_pernikahan_id": 1,
				"hubungan": "Istri",
				"tanggal_menikah": null,
				"akte_nikah": "",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1d",
					"nip":               "1d",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: agama or status pernikahan is not found",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, agama_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"tanggal_lahir": "2000-01-01",
				"agama_id": 0,
				"status_pernikahan_id": 0,
				"hubungan": "Suami"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan | data status pernikahan tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          int16(2),
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: agama or status pernikahan is deleted",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, agama_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        1,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "123456",
				"is_pns": false,
				"tanggal_lahir": "2000-01-01",
				"no_karsus": "",
				"agama_id": 2,
				"status_pernikahan_id": 2,
				"hubungan": "Suami",
				"tanggal_menikah": "2000-01-01",
				"akte_nikah": "akte_01",
				"tanggal_meninggal": null,
				"akte_meninggal": "",
				"tanggal_cerai": null,
				"akte_cerai": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan | data status pernikahan tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          int16(1),
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:          "error: invalid format date, exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 101) + `",
				"nik": "` + strings.Repeat(".", 21) + `",
				"is_pns": "",
				"tanggal_lahir": "",
				"no_karsus": "` + strings.Repeat(".", 101) + `",
				"agama_id": "Islam",
				"status_pernikahan_id": "Menikah",
				"hubungan": "Pelakor",
				"tanggal_menikah": "",
				"akte_nikah": "` + strings.Repeat(".", 101) + `",
				"tanggal_meninggal": "",
				"akte_meninggal": "` + strings.Repeat(".", 101) + `",
				"tanggal_cerai": "",
				"akte_cerai": "` + strings.Repeat(".", 101) + `"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"agama_id\" harus dalam tipe integer` +
				` | parameter \"akte_cerai\" harus 100 karakter atau kurang` +
				` | parameter \"akte_meninggal\" harus 100 karakter atau kurang` +
				` | parameter \"akte_nikah\" harus 100 karakter atau kurang` +
				` | parameter \"hubungan\" harus salah satu dari \"Istri\", \"Suami\"` +
				` | parameter \"is_pns\" harus dalam tipe boolean` +
				` | parameter \"nama\" harus 100 karakter atau kurang` +
				` | parameter \"nik\" harus 20 karakter atau kurang` +
				` | parameter \"no_karsus\" harus 100 karakter atau kurang` +
				` | parameter \"status_pernikahan_id\" harus dalam tipe integer` +
				` | parameter \"tanggal_cerai\" harus dalam format date` +
				` | parameter \"tanggal_lahir\" harus dalam format date` +
				` | parameter \"tanggal_menikah\" harus dalam format date` +
				` | parameter \"tanggal_meninggal\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null on required params",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": null,
				"hubungan": null,
				"tanggal_lahir": null,
				"status_pernikahan_id": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"hubungan\" harus salah satu dari \"Istri\", \"Suami\"` +
				` | parameter \"nama\" tidak boleh null` +
				` | parameter \"status_pernikahan_id\" tidak boleh null` +
				` | parameter \"tanggal_lahir\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama\" harus diisi` +
				` | parameter \"tanggal_lahir\" harus diisi` +
				` | parameter \"status_pernikahan_id\" harus diisi` +
				` | parameter \"hubungan\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/pasangan/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "pasangan", "id")
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

func Test_handler_adminDeletePasangan(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: success delete",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns,  nik,   tanggal_lahir, tanggal_menikah, akte_nikah, pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 2,        1,    '123', '1990-01-01',  '2020-01-01',    'akte_01',  'id_1c', '1c', '2000-01-01', '2000-01-01'),
					(2,  'Jane Doe', 1,        null, null,  null,          null,            null,       'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "John Doe",
					"hubungan":          int16(2),
					"tanggal_lahir":     time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"pns":               int16(1),
					"nik":               "123",
					"agama_id":          nil,
					"tanggal_menikah":   time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC),
					"akte_nikah":        "akte_01",
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        "{deleted_at}",
				},
				{
					"id":                int64(2),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: pasangan is not found",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name: "error: pasangan is deleted",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at,   deleted_at) values
					(1,  'Jane Doe', 1,        'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1c",
					"nip":               "1c",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "error: pegawai is deleted",
			dbData: seedData + `
				insert into pasangan
					(id, nama,       hubungan, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        'id_1d', '1d', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1d",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                int64(1),
					"nama":              "Jane Doe",
					"hubungan":          int16(1),
					"tanggal_lahir":     nil,
					"pns":               nil,
					"nik":               nil,
					"agama_id":          nil,
					"tanggal_menikah":   nil,
					"akte_nikah":        nil,
					"tanggal_cerai":     nil,
					"akte_cerai":        nil,
					"tanggal_meninggal": nil,
					"akte_meninggal":    nil,
					"karsus":            nil,
					"status":            nil,
					"pns_id":            "id_1d",
					"nip":               "1d",
					"created_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":        time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":        nil,
				},
			},
		},
		{
			name:             "error: unexpected data type",
			paramNIP:         "1c",
			paramID:          "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/pasangan/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "pasangan", "id")
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

func Test_handler_adminCreateAnak(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01');
		insert into ref_agama
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into ref_jenis_kawin
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into pasangan
			(id, pns_id,  deleted_at) values
			(1,  'id_1c', null),
			(2,  'id_1c', '2000-01-01'),
			(3,  'id_1d', null);
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: with all data",
			dbData: seedData + `
				insert into anak
					(nama,       pns_id,  nip,  created_at,   updated_at) values
					('John Doe', 'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"agama_id": 1,
				"status_anak": "Angkat",
				"status_sekolah": "Masih Sekolah",
				"anak_ke": 1
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 2 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "John Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
				{
					"id":             int64(2),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(1),
					"nik":            "123",
					"jenis_kelamin":  "F",
					"tempat_lahir":   nil,
					"tanggal_lahir":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "2",
					"status_sekolah": int16(1),
					"agama_id":       int16(1),
					"jenis_kawin_id": int16(1),
					"anak_ke":        int16(1),
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     "{created_at}",
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name:          "ok: with different enum data and anak_ke = 0",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Will Doe",
				"nik": "123",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"agama_id": 1,
				"status_anak": "Kandung",
				"status_sekolah": "Sudah Bekerja",
				"anak_ke": 0
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Will Doe",
					"pasangan_id":    int64(1),
					"nik":            "123",
					"jenis_kelamin":  "M",
					"tempat_lahir":   nil,
					"tanggal_lahir":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "1",
					"status_sekolah": int16(2),
					"agama_id":       int16(1),
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     "{created_at}",
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name:          "ok: with null values",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"agama_id": null,
				"status_anak": "Angkat",
				"status_sekolah": "",
				"anak_ke": null
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(1),
					"nik":            nil,
					"jenis_kelamin":  "F",
					"tempat_lahir":   nil,
					"tanggal_lahir":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "2",
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     "{created_at}",
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name:          "ok: required data only",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": 1 }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "John Doe",
					"pasangan_id":    int64(1),
					"nik":            nil,
					"jenis_kelamin":  "M",
					"tempat_lahir":   nil,
					"tanggal_lahir":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "1",
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     "{created_at}",
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			dbData:        seedData,
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 2,
				"status_anak": "Kandung"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			dbData:        seedData,
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 3,
				"agama_id": 2,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung",
				"status_sekolah": "Masih Sekolah",
				"anak_ke": 0
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pasangan orang tua is owned by different pegawai",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 3,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data pasangan orang tua tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: agama or status pernikahan or pasangan orang tua is not found",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 0,
				"agama_id": 0,
				"status_pernikahan_id": 0,
				"status_anak": "Angkat"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan | data status pernikahan tidak ditemukan | data pasangan orang tua tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: agama or status pernikahan or pasangan orang tua is deleted",
			dbData:        seedData,
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 2,
				"agama_id": 2,
				"status_pernikahan_id": 2,
				"status_anak": "Kandung",
				"status_sekolah": "Masih Sekolah",
				"anak_ke": 0
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan | data status pernikahan tidak ditemukan | data pasangan orang tua tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: invalid format date, exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 101) + `",
				"nik": "` + strings.Repeat(".", 21) + `",
				"jenis_kelamin": "L",
				"tanggal_lahir": "",
				"pasangan_orang_tua_id": "Jane Doe",
				"agama_id": "Islam",
				"status_pernikahan_id": "Menikah",
				"status_anak": "Anak",
				"status_sekolah": "Pelajar",
				"anak_ke": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"agama_id\" harus dalam tipe integer` +
				` | parameter \"anak_ke\" harus dalam tipe integer` +
				` | parameter \"jenis_kelamin\" harus salah satu dari \"M\", \"F\"` +
				` | parameter \"nama\" harus 100 karakter atau kurang` +
				` | parameter \"nik\" harus 20 karakter atau kurang` +
				` | parameter \"pasangan_orang_tua_id\" harus dalam tipe integer` +
				` | parameter \"status_anak\" harus salah satu dari \"Kandung\", \"Angkat\"` +
				` | parameter \"status_pernikahan_id\" harus dalam tipe integer` +
				` | parameter \"status_sekolah\" harus salah satu dari \"Masih Sekolah\", \"Sudah Bekerja\", \"\"` +
				` | parameter \"tanggal_lahir\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null on required params",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": null,
				"jenis_kelamin": null,
				"tanggal_lahir": null,
				"pasangan_orang_tua_id": null,
				"status_pernikahan_id": null,
				"status_anak": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"jenis_kelamin\" harus salah satu dari \"M\", \"F\"` +
				` | parameter \"nama\" tidak boleh null` +
				` | parameter \"pasangan_orang_tua_id\" tidak boleh null` +
				` | parameter \"status_anak\" harus salah satu dari \"Kandung\", \"Angkat\"` +
				` | parameter \"status_pernikahan_id\" tidak boleh null` +
				` | parameter \"tanggal_lahir\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama\" harus diisi` +
				` | parameter \"jenis_kelamin\" harus diisi` +
				` | parameter \"tanggal_lahir\" harus diisi` +
				` | parameter \"pasangan_orang_tua_id\" harus diisi` +
				` | parameter \"status_pernikahan_id\" harus diisi` +
				` | parameter \"status_anak\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/anak", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "anak", "id")
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
		})
	}
}

func Test_handler_adminUpdateAnak(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     null);
		insert into ref_agama
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into ref_jenis_kawin
			(id, deleted_at) values
			(1,  null),
			(2,  '2000-01-01');
		insert into pasangan
			(id, pns_id,  deleted_at) values
			(1,  'id_1c', null),
			(2,  'id_1c', '2000-01-01'),
			(3,  'id_1d', null);
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: with all data",
			dbData: seedData + `
				insert into anak
					(id, nama,       tempat_lahir, pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 'Jakarta',    'id_1c', '1c', '2000-01-01', '2000-01-01'),
					(2,  'John Doe', 'Jakarta',    'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"nik": "123",
				"jenis_kelamin": "F",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"agama_id": 1,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat",
				"status_sekolah": "Masih Sekolah",
				"anak_ke": 1
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(1),
					"nik":            "123",
					"jenis_kelamin":  "F",
					"tempat_lahir":   "Jakarta",
					"tanggal_lahir":  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "2",
					"status_sekolah": int16(1),
					"agama_id":       int16(1),
					"jenis_kawin_id": int16(1),
					"anak_ke":        int16(1),
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
				{
					"id":             int64(2),
					"nama":           "John Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   "Jakarta",
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "ok: with different enum data and anak_ke = 0",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Will Doe",
				"nik": "123",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"agama_id": 1,
				"status_anak": "Kandung",
				"status_sekolah": "Sudah Bekerja",
				"anak_ke": 0
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Will Doe",
					"pasangan_id":    int64(1),
					"nik":            "123",
					"jenis_kelamin":  "M",
					"tempat_lahir":   nil,
					"tanggal_lahir":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "1",
					"status_sekolah": int16(2),
					"agama_id":       int16(1),
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "ok: with null values",
			dbData: seedData + `
				insert into anak
					(id, nama,       tempat_lahir, status_sekolah, anak_ke, agama_id, nik,   pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 'Medan',      1,              1,       1,        '123', 'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "",
				"jenis_kelamin": "M",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"agama_id": null,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat",
				"status_sekolah": "",
				"anak_ke": null
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "John Doe",
					"pasangan_id":    int64(1),
					"nik":            nil,
					"jenis_kelamin":  "M",
					"tempat_lahir":   "Medan",
					"tanggal_lahir":  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "2",
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "ok: required data only",
			dbData: seedData + `
				insert into anak
					(id, nama,       tempat_lahir, status_sekolah, anak_ke, agama_id, nik,   pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', 'Medan',      1,              1,       1,        '123', 'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(1),
					"nik":            nil,
					"jenis_kelamin":  "F",
					"tempat_lahir":   "Medan",
					"tanggal_lahir":  time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    "1",
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     "{updated_at}",
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: anak is not found",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: anak is owned by different pegawai",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1e', '1e', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1e",
					"nip":            "1e",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: anak is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at,   deleted_at) values
					(1,  'Jane Doe', 'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "",
				"jenis_kelamin": "M",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"agama_id": null,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung",
				"status_sekolah": "",
				"anak_ke": null
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "error: pegawai is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1d', '1d', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1d",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "",
				"jenis_kelamin": "M",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 3,
				"agama_id": null,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung",
				"status_sekolah": "",
				"anak_ke": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data pasangan orang tua tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1d",
					"nip":            "1d",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: pasangan orang tua is not found",
			dbData: seedData + `
				insert into anak
					(id, nama,       pasangan_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,           'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 0,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data pasangan orang tua tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(2),
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: pasangan orang tua is owned by different pegawai",
			dbData: seedData + `
				insert into anak
					(id, nama,       pasangan_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,           'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"jenis_kelamin": "F",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 3,
				"status_pernikahan_id": 1,
				"status_anak": "Angkat"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data pasangan orang tua tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(1),
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: pasangan orang tua is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       pasangan_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,           'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "M",
				"tanggal_lahir": "2000-01-01",
				"pasangan_orang_tua_id": 2,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data pasangan orang tua tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    int64(1),
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: status pernikahan is not found",
			dbData: seedData + `
				insert into anak
					(id, nama,       jenis_kawin_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,              'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "Jane Doe",
				"jenis_kelamin": "F",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"status_pernikahan_id": 0,
				"status_anak": "Kandung"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data status pernikahan tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": int16(2),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: status pernikahan is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       jenis_kawin_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,              'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "123",
				"jenis_kelamin": "M",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"agama_id": 1,
				"status_pernikahan_id": 2,
				"status_anak": "Kandung",
				"status_sekolah": "Masih Sekolah",
				"anak_ke": 1
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data status pernikahan tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": int16(1),
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: agama is not found",
			dbData: seedData + `
				insert into anak
					(id, nama,       agama_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 2,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"jenis_kelamin": "M",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"agama_id": 0,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       int16(2),
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: agama is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       agama_id, pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 1,        'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "John Doe",
				"nik": "123",
				"jenis_kelamin": "M",
				"tanggal_lahir": "1990-01-01",
				"pasangan_orang_tua_id": 1,
				"agama_id": 2,
				"status_pernikahan_id": 1,
				"status_anak": "Kandung",
				"status_sekolah": "Masih Sekolah",
				"anak_ke": 1
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data agama tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       int16(1),
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name:          "error: invalid format date, exceed length limit, unexpected enum or data type",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 101) + `",
				"nik": "` + strings.Repeat(".", 21) + `",
				"jenis_kelamin": "L",
				"tanggal_lahir": "",
				"pasangan_orang_tua_id": "Jane Doe",
				"agama_id": "Islam",
				"status_pernikahan_id": "Menikah",
				"status_anak": "Anak",
				"status_sekolah": "Pelajar",
				"anak_ke": ""
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"agama_id\" harus dalam tipe integer` +
				` | parameter \"anak_ke\" harus dalam tipe integer` +
				` | parameter \"jenis_kelamin\" harus salah satu dari \"M\", \"F\"` +
				` | parameter \"nama\" harus 100 karakter atau kurang` +
				` | parameter \"nik\" harus 20 karakter atau kurang` +
				` | parameter \"pasangan_orang_tua_id\" harus dalam tipe integer` +
				` | parameter \"status_anak\" harus salah satu dari \"Kandung\", \"Angkat\"` +
				` | parameter \"status_pernikahan_id\" harus dalam tipe integer` +
				` | parameter \"status_sekolah\" harus salah satu dari \"Masih Sekolah\", \"Sudah Bekerja\", \"\"` +
				` | parameter \"tanggal_lahir\" harus dalam format date"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:          "error: null on required params",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"nama": null,
				"jenis_kelamin": null,
				"tanggal_lahir": null,
				"pasangan_orang_tua_id": null,
				"status_pernikahan_id": null,
				"status_anak": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"jenis_kelamin\" harus salah satu dari \"M\", \"F\"` +
				` | parameter \"nama\" tidak boleh null` +
				` | parameter \"pasangan_orang_tua_id\" tidak boleh null` +
				` | parameter \"status_anak\" harus salah satu dari \"Kandung\", \"Angkat\"` +
				` | parameter \"status_pernikahan_id\" tidak boleh null` +
				` | parameter \"tanggal_lahir\" tidak boleh null"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: missing required params & have additional params",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"nama\" harus diisi` +
				` | parameter \"jenis_kelamin\" harus diisi` +
				` | parameter \"tanggal_lahir\" harus diisi` +
				` | parameter \"pasangan_orang_tua_id\" harus diisi` +
				` | parameter \"status_pernikahan_id\" harus diisi` +
				` | parameter \"status_anak\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/anak/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "anak", "id")
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

func Test_handler_adminDeleteAnak(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into pegawai
			(pns_id,  nip_baru, deleted_at) values
			('id_1c', '1c',     null),
			('id_1d', '1d',     '2000-01-01');
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name: "ok: success delete",
			dbData: seedData + `
				insert into anak
					(id, nama,       nik,   tanggal_lahir, pns_id,  nip,  created_at,   updated_at) values
					(1,  'John Doe', '123', '1990-01-01',  'id_1c', '1c', '2000-01-01', '2000-01-01'),
					(2,  'Jane Doe', null,  null,          'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "John Doe",
					"pasangan_id":    nil,
					"nik":            "123",
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     "{deleted_at}",
				},
				{
					"id":             int64(2),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: anak is not found",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1c', '1c', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name: "error: anak is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at,   deleted_at) values
					(1,  'Jane Doe', 'id_1c', '1c', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1c",
					"nip":            "1c",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "error: pegawai is deleted",
			dbData: seedData + `
				insert into anak
					(id, nama,       pns_id,  nip,  created_at,   updated_at) values
					(1,  'Jane Doe', 'id_1d', '1d', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1d",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":             int64(1),
					"nama":           "Jane Doe",
					"pasangan_id":    nil,
					"nik":            nil,
					"jenis_kelamin":  nil,
					"tempat_lahir":   nil,
					"tanggal_lahir":  nil,
					"status_anak":    nil,
					"status_sekolah": nil,
					"agama_id":       nil,
					"jenis_kawin_id": nil,
					"anak_ke":        nil,
					"pns_id":         "id_1d",
					"nip":            "1d",
					"created_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":     time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":     nil,
				},
			},
		},
		{
			name:             "error: unexpected data type",
			paramNIP:         "1c",
			paramID:          "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/anak/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
			RegisterRoutes(e, repo.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, typeutil.Coalesce(tt.wantResponseBody, "null"), typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryAll(db, "anak", "id")
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
