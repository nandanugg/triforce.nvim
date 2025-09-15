package pelatihanstruktural

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
		INSERT INTO "ref_jenis_diklat" (id, bkn_id, jenis_diklat, kode, status) VALUES
		(1, 1, 'Struktural', '01', 1),
		(2, 2, 'Fungsional', '02', 1),
		(3, 3, 'Teknis', '03', 1);
		INSERT INTO "pegawai" ( id, pns_id, nip_baru, nama) VALUES
		    ( 1, 'uuid-pns-orang-001', '199001012020121001', 'Agus Purnomo');
		INSERT INTO "riwayat_diklat" (
		    id, jenis_diklat_id, nama_diklat, institusi_penyelenggara, no_sertifikat, tanggal_mulai, tanggal_selesai, tahun_diklat, durasi_jam, nip_baru, pns_orang_id) VALUES
		( 1, 1, 'Pelatihan Kepemimpinan Administrator (PKA)', 'Lembaga Administrasi Negara (LAN)', 'LAN-PKA-2023-00123', '2023-03-15', '2023-06-20', 2023, 900, '199001012020121001', 'uuid-pns-orang-001'),
		( 2, 2, 'Diklat Fungsional Analis Kebijakan', 'Pusat Pembinaan Analis Kebijakan', 'PPAK-2022-00456', '2022-09-01', '2022-09-30', 2022, 240, '199001012020121001', 'uuid-pns-orang-001'),
		( 3, 3, 'Pelatihan Teknis Pengadaan Barang dan Jasa Pemerintah', 'Lembaga Kebijakan Pengadaan Barang/Jasa Pemerintah (LKPP)', 'LKPP-PBJ-2024-00789', '2024-02-10', '2024-02-15', 2024, 50, '199001012020121001', 'uuid-pns-orang-001'),
		( 4, 3, 'Workshop Implementasi Sistem Pemerintahan Berbasis Elektronik (SPBE)', 'Kementerian Komunikasi dan Informatika', 'KOMINFO-SPBE-2023-00331', '2023-11-05', '2023-11-08', 2023, 32, '199001012020121001', 'uuid-pns-orang-001');
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
					"durasi": 900,
					"id": 1,
					"istitusi_penyelenggara": "Lembaga Administrasi Negara (LAN)",
					"jenis_diklat": "Struktural",
					"nama_diklat": "Pelatihan Kepemimpinan Administrator (PKA)",
					"nomor_sertifikat": "LAN-PKA-2023-00123",
					"tahun": 2023,
					"tanggal_mulai": "2023-03-15T00:00:00Z",
					"tanggal_selesai": "2023-06-20T00:00:00Z"
				    },
				    {
					"durasi": 240,
					"id": 2,
					"istitusi_penyelenggara": "Pusat Pembinaan Analis Kebijakan",
					"jenis_diklat": "Fungsional",
					"nama_diklat": "Diklat Fungsional Analis Kebijakan",
					"nomor_sertifikat": "PPAK-2022-00456",
					"tahun": 2022,
					"tanggal_mulai": "2022-09-01T00:00:00Z",
					"tanggal_selesai": "2022-09-30T00:00:00Z"
				    },
				    {
					"durasi": 50,
					"id": 3,
					"istitusi_penyelenggara": "Lembaga Kebijakan Pengadaan Barang/Jasa Pemerintah (LKPP)",
					"jenis_diklat": "Teknis",
					"nama_diklat": "Pelatihan Teknis Pengadaan Barang dan Jasa Pemerintah",
					"nomor_sertifikat": "LKPP-PBJ-2024-00789",
					"tahun": 2024,
					"tanggal_mulai": "2024-02-10T00:00:00Z",
					"tanggal_selesai": "2024-02-15T00:00:00Z"
				    },
				    {
					"durasi": 32,
					"id": 4,
					"istitusi_penyelenggara": "Kementerian Komunikasi dan Informatika",
					"jenis_diklat": "Teknis",
					"nama_diklat": "Workshop Implementasi Sistem Pemerintahan Berbasis Elektronik (SPBE)",
					"nomor_sertifikat": "KOMINFO-SPBE-2023-00331",
					"tahun": 2023,
					"tanggal_mulai": "2023-11-05T00:00:00Z",
					"tanggal_selesai": "2023-11-08T00:00:00Z"
				    }
				]
			    }`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [] }`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/pelatihan-struktural", nil)
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
