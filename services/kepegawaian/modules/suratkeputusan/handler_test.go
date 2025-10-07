package suratkeputusan_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	dbrepository "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	suratkeputusan "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/suratkeputusan"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, nama_jabatan,    pemimpin_pns_id, deleted_at) values
			('unor-1', null, 'Paling Atas', 'Atasan 1', null, null),
			('unor-2', 'unor-1', 'Tengah', 'Atasan 2', null, null),
			('unor-3', 'unor-2', 'Bawah', 'Atasan 3', null, null),
			('unor-4', 'unor-1', 'Tengah deleted', 'Atasan 4', null, now()),
			('unor-5', 'unor-4', 'Bawah 2', 'Atasan 5', null, null);
		
		insert into pegawai 
			("nip_baru","pns_id","unor_id") values
			('123456789','123456789','unor-3'),
			('123456788','123456788','unor-4'),
			('123456787','123456787','unor-5');

		INSERT INTO file_digital_signature
				("file_id", "nip_sk", "kategori", "no_sk", "tanggal_sk", "status_sk", "created_at", "deleted_at") VALUES
				('sk-001', '123456789', 'Kenaikan Pangkat', 'SK-001/2024', '2024-01-15', 1, '2024-01-15', NULL),
				('sk-002', '123456789', 'Mutasi', 'SK-002/2024', '2024-02-20', 0, '2024-02-20', NULL),
				('sk-003', '123456789', 'Kenaikan Gaji', 'SK-003/2024', '2024-03-10', 2, '2024-03-10', NULL),
				('sk-004', '123456789', 'Kenaikan Gaji', 'SK-004/2024', '2024-03-10', 2, '2024-03-10', NOW()),
				('sk-005', '123456788', 'Mutasi', 'SK-005/2024', '2024-03-10', 2, '2024-03-10', NULL),
				('sk-006', '123456787', 'Mutasi', 'SK-006/2024', '2024-03-10', 2, '2024-03-10', NULL);
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
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-003", "kategori_sk": "Kenaikan Gaji", "no_sk": "SK-003/2024", "tanggal_sk": "2024-03-10", "status_sk": "Sudah Dikoreksi & Dikembalikan", "unit_kerja": "Bawah - Tengah - Paling Atas"},
					{"id_sk": "sk-002", "kategori_sk": "Mutasi", "no_sk": "SK-002/2024", "tanggal_sk": "2024-02-20", "status_sk": "Belum Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"},
					{"id_sk": "sk-001", "kategori_sk": "Kenaikan Pangkat", "no_sk": "SK-001/2024", "tanggal_sk": "2024-01-15", "status_sk": "Sedang Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok with filter kategori_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"kategori_sk": []string{"Kenaikan Gaji"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-003", "kategori_sk": "Kenaikan Gaji", "no_sk": "SK-003/2024", "tanggal_sk": "2024-03-10", "status_sk": "Sudah Dikoreksi & Dikembalikan", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok with filter no_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"no_sk": []string{"SK-002/2024"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-002", "kategori_sk": "Mutasi", "no_sk": "SK-002/2024", "tanggal_sk": "2024-02-20", "status_sk": "Belum Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok with filter status_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"status_sk": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-001", "kategori_sk": "Kenaikan Pangkat", "no_sk": "SK-001/2024", "tanggal_sk": "2024-01-15", "status_sk": "Sedang Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok with filter multiple status_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"status_sk": []string{"1", "2"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-003", "kategori_sk": "Kenaikan Gaji", "no_sk": "SK-003/2024", "tanggal_sk": "2024-03-10", "status_sk": "Sudah Dikoreksi & Dikembalikan", "unit_kerja": "Bawah - Tengah - Paling Atas"},
					{"id_sk": "sk-001", "kategori_sk": "Kenaikan Pangkat", "no_sk": "SK-001/2024", "tanggal_sk": "2024-01-15", "status_sk": "Sedang Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 2}
			}`,
		},
		{
			name:             "ok with limit 2",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"2"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-003", "kategori_sk": "Kenaikan Gaji", "no_sk": "SK-003/2024", "tanggal_sk": "2024-03-10", "status_sk": "Sudah Dikoreksi & Dikembalikan", "unit_kerja": "Bawah - Tengah - Paling Atas"},
					{"id_sk": "sk-002", "kategori_sk": "Mutasi", "no_sk": "SK-002/2024", "tanggal_sk": "2024-02-20", "status_sk": "Belum Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 2, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok with limit 2 and offset 1",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-002", "kategori_sk": "Mutasi", "no_sk": "SK-002/2024", "tanggal_sk": "2024-02-20", "status_sk": "Belum Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"},
					{"id_sk": "sk-001", "kategori_sk": "Kenaikan Pangkat", "no_sk": "SK-001/2024", "tanggal_sk": "2024-01-15", "status_sk": "Sedang Dikoreksi", "unit_kerja": "Bawah - Tengah - Paling Atas"}
				],
				"meta": {"limit": 2, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok with unor utama deleted",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456788")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-005", "kategori_sk": "Mutasi", "no_sk": "SK-005/2024", "tanggal_sk": "2024-03-10", "status_sk": "Sudah Dikoreksi & Dikembalikan", "unit_kerja": ""}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok with unor parent deleted",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456787")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id_sk": "sk-006", "kategori_sk": "Mutasi", "no_sk": "SK-006/2024", "tanggal_sk": "2024-03-10", "status_sk": "Sudah Dikoreksi & Dikembalikan", "unit_kerja": "Bawah 2"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 1}
			}`,
		},
		{
			name:             "ok with empty data",
			dbData:           ``,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
		},
		{
			name:             "ok with different user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [],
				"meta": {"limit": 10, "offset": 0, "total": 0}
			}`,
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
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/surat-keputusan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := dbrepository.New(db)
			suratkeputusan.RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_get(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, nama_jabatan,    pemimpin_pns_id, deleted_at) values
			('unor-1', null, 'Paling Atas', 'Atasan 1', null, null),
			('unor-2', 'unor-1', 'Tengah', 'Atasan 2', null, null),
			('unor-3', 'unor-2', 'Bawah', 'Atasan 3', null, null),
			('unor-4', 'unor-1', 'Tengah deleted', 'Atasan 4', null, now()),
			('unor-5', 'unor-4', 'Bawah 2', 'Atasan 5', null, null);
		
		insert into pegawai 
			("nip_baru","pns_id","unor_id","nama") values
			('123456789','123456789','unor-3','pemilik_sk'),
			('123456788','123456788','unor-4','pemilik_sk_2'),
			('123456787','123456787','unor-5','pemilik_sk_3'),
			('12345678','12345678','unor-3','Jane Smith'),
			('12345677','12345677','unor-3','Korektor');

		INSERT INTO file_digital_signature
				("file_id", "nip_sk", "kategori", "no_sk", "tanggal_sk", "status_sk","nip_pemroses", "created_at", "deleted_at") VALUES
				('sk-001', '123456789', 'Kenaikan Pangkat', 'SK-001/2024', '2024-01-15', 1, '12345678', '2024-01-15', NULL),
				('sk-002', '123456789', 'Mutasi', 'SK-002/2024', '2024-02-20', 0, '12345678', '2024-02-20', NULL),
				('sk-003', '123456789', 'Kenaikan Gaji', 'SK-003/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-004', '123456789', 'Kenaikan Gaji', 'SK-004/2024', '2024-03-10', 2,'12345678', '2024-03-10', NOW()),
				('sk-005', '123456788', 'Mutasi', 'SK-005/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-006', '123456787', 'Mutasi', 'SK-006/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL);
		INSERT INTO file_digital_signature_riwayat
			("file_id","nip_pemroses","tindakan","created_at","deleted_at") VALUES
			('sk-002','12345677','Mengkoreksi','2025-09-24T01:02:03',NULL),
			('sk-002','12345678','Mengkoreksi',NOW(),NOW());
	`

	tests := []struct {
		name             string
		dbData           string
		requestPath      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-001",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-001",
					"kategori_sk": "Kenaikan Pangkat",
					"no_sk": "SK-001/2024",
					"tanggal_sk": "2024-01-15",
					"status_sk": "Sedang Dikoreksi",
					"unit_kerja": "Bawah - Tengah - Paling Atas",
					"nama_pemilik": "pemilik_sk",
					"nama_penandatangan": "Jane Smith",
					"logs": []
				}
			}`,
		},
		{
			name:             "ok with different sk and has logs",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-002",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-002",
					"kategori_sk": "Mutasi",
					"no_sk": "SK-002/2024",
					"tanggal_sk": "2024-02-20",
					"status_sk": "Belum Dikoreksi",
					"unit_kerja": "Bawah - Tengah - Paling Atas",
					"nama_pemilik": "pemilik_sk",
					"nama_penandatangan": "Jane Smith",
					"logs": [
						{
							"log" : "Mengkoreksi",
							"actor" : "Korektor",
							"timestamp" : "` + time.Date(2025, time.September, 24, 1, 2, 3, 0, time.UTC).Local().Format(time.RFC3339) + `"
						}
					]
				}
			}`,
		},
		{
			name:             "ok with unor utama deleted",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-005",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456788")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-005",
					"kategori_sk": "Mutasi",
					"no_sk": "SK-005/2024",
					"tanggal_sk": "2024-03-10",
					"status_sk": "Sudah Dikoreksi & Dikembalikan",
					"unit_kerja": "",
					"nama_pemilik": "pemilik_sk_2",
					"nama_penandatangan": "Jane Smith",
					"logs": []
				}
			}`,
		},
		{
			name:             "ok with parent unor deleted",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-006",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456787")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-006",
					"kategori_sk": "Mutasi",
					"no_sk": "SK-006/2024",
					"tanggal_sk": "2024-03-10",
					"status_sk": "Sudah Dikoreksi & Dikembalikan",
					"unit_kerja": "Bawah 2",
					"nama_pemilik": "pemilik_sk_3",
					"nama_penandatangan": "Jane Smith",
					"logs": []
				}
			}`,
		},
		{
			name:             "error: sk not found",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: sk deleted",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-004",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: sk belongs to different user",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-001",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			requestPath:      "/v1/surat-keputusan/sk-001",
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
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := dbrepository.New(db)
			suratkeputusan.RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_getBerkas(t *testing.T) {
	t.Parallel()

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

	signedBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xde, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
		0x00, 0x04, 0x00, 0x01, 0xe2, 0x26, 0x05, 0x9b,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
		0xae, 0x42, 0x60, 0x82,
	}

	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)
	signedPngBase64 := base64.StdEncoding.EncodeToString(signedBytes)

	dbData := `

		INSERT INTO file_digital_signature
				("file_id","file_base64","file_base64_sign", "nip_sk", "kategori", "no_sk", "tanggal_sk", "status_sk","nip_pemroses", "created_at", "deleted_at") VALUES
				('sk-001','data:image/png;base64,` + pngBase64 + `',NULL, '123456789', 'Kenaikan Pangkat', 'SK-001/2024', '2024-01-15', 1, '12345678', '2024-01-15', NULL),
				('sk-002','data:image/png;base64,invalid',NULL, '123456789', 'Mutasi', 'SK-002/2024', '2024-02-20', 0, '12345678', '2024-02-20', NULL),
				('sk-003',NULL,NULL, '123456789', 'Kenaikan Gaji', 'SK-003/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-004','',NULL, '123456789', 'Kenaikan Gaji', 'SK-004/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-005','data:image/png;base64,` + pngBase64 + `',NULL, '123456789', 'Kenaikan Gaji', 'SK-005/2024', '2024-03-10', 2,'12345678', '2024-03-10', NOW()),
				('sk-006',NULL,'data:image/png;base64,` + signedPngBase64 + `', '123456789', 'Kenaikan Gaji', 'SK-006/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-007','data:image/png;base64,` + pngBase64 + `','data:image/png;base64,` + signedPngBase64 + `', '123456789', 'Kenaikan Gaji', 'SK-007/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL);
			
	`

	tests := []struct {
		name              string
		dbData            string
		requestPath       string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-001/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "ok: valid signed png",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-006/berkas?signed=true",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: signedBytes,
		},
		{
			name:              "ok: valid unsigned png",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-007/berkas?signed=false",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "ok: valid signed png with params 1",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-007/berkas?signed=1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: signedBytes,
		},
		{
			name:              "error: empty value params signed",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-007/berkas?signed=",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message":"parameter \"signed\" tidak boleh kosong"}`),
		},
		{
			name:              "error: invalid params signed",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-007/berkas?signed=''",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message":"parameter \"signed\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: base64 berkas SK tidak valid",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-002/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas SK sudah dihapus",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-005/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas SK berisi null value",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-003/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas SK berupa string kosong",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-004/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: berkas SK bukan milik user login",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-001/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456788")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: berkas SK tidak ditemukan",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-009/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			requestPath:       "/v1/surat-keputusan/sk-001/berkas",
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

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepository.New(pgxconn)
			suratkeputusan.RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, nama_jabatan, pemimpin_pns_id, eselon_1, eselon_2, eselon_3, eselon_4, deleted_at) values
			('unor-1', null, 'Paling Atas', 'Atasan 1', null, null, null, null, null, null),
			('unor-2', 'unor-1', 'Tengah', 'Atasan 2', null, null, 'unor-1', null, null, null),
			('unor-3', 'unor-2', 'Bawah', 'Atasan 3', null, 'unor-1', null, null, null, null),
			('unor-4', 'unor-1', 'Tengah deleted', 'Atasan 4', null, null, null, null, null, now()),
			('unor-5', 'unor-4', 'Bawah 2', 'Atasan 5', null, null, null, null, 'unor-1', null),
			('unor-6', 'unor-1', 'Bawah 3', 'Atasan 6', null, null, null, 'unor-1', null, null),
			('unor-7', null, 'Paling Atas 2', 'Atasan 7', null, null, null, null, null, null);

		INSERT INTO ref_golongan (id, nama, nama_pangkat, nama_2, gol, gol_pppk)
		VALUES
			(1, 'I/a', 'Juru Muda', 'I/a', 1, 'I/a'),
			(2, 'I/b', 'Juru Muda Tingkat I', 'I/b', 2, 'I/b'),
			(3, 'II/a', 'Pengatur Muda', 'II/a', 3, 'II/a'),
			(4, 'II/b', 'Pengatur Muda Tingkat I', 'II/b', 4, 'II/b'),
			(5, 'III/a', 'Penata Muda', 'III/a', 5, 'III/a'),
			(6, 'III/b', 'Penata Muda Tingkat I', 'III/b', 6, 'III/b'),
			(7, 'IV/a', 'Pembina', 'IV/a', 7, 'IV/a');

		INSERT INTO ref_jabatan (kode_jabatan, id, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id)
		VALUES
			('JAB-01', 1, 'Kepala Bagian', 'Kepala Bagian Administrasi', 1, 3, 60, 'BKN-01', 'Kepala Bagian BKN', 'Struktural', 'BKN-01'),
			('JAB-02', 2, 'Kepala Subbag', 'Kepala Subbag Umum', 1, 2, 60, 'BKN-02', 'Kepala Subbag BKN', 'Struktural', 'BKN-02'),
			('JAB-03', 3, 'Pengawas', 'Pengawas Pelayanan', 2, 2, 58, 'BKN-03', 'Pengawas BKN', 'Fungsional', 'BKN-03'),
			('JAB-04', 4, 'Staf Ahli', 'Staf Ahli Bidang Kepegawaian', 2, 3, 58, 'BKN-04', 'Staf Ahli BKN', 'Fungsional', 'BKN-04'),
			('JAB-05', 5, 'Kepala Seksi', 'Kepala Seksi Administrasi', 1, 2, 60, 'BKN-05', 'Kepala Seksi BKN', 'Struktural', 'BKN-05'),
			('JAB-06', 6, 'Analis', 'Analis Kepegawaian', 2, 1, 58, 'BKN-06', 'Analis BKN', 'Fungsional', 'BKN-06'),
			('JAB-07', 7, 'Operator', 'Operator Sistem', 2, 1, 58, 'BKN-07', 'Operator BKN', 'Fungsional', 'BKN-07');

		INSERT INTO pegawai ("nip_baru", "pns_id", "unor_id", "nama", "gol_id", "jabatan_instansi_id", deleted_at) 
		VALUES
			('123456789','123456789','unor-3','Andi', 1, 'JAB-01', NULL),
			('123456788','123456788','unor-4','Budi', 2, 'JAB-02', NULL),
			('123456787','123456787','unor-5','Candi', 3, 'JAB-03', NULL),
			('123456786','123456786','unor-6','Dedi', 4, 'JAB-04', NULL),
			('123456785','123456785','unor-7','Edi', 5, 'JAB-05', NULL),
			('123456784','123456784','unor-1','Fredi', 6, 'JAB-06', NULL),
			('123456783','123456783','unor-2','Gledi', 7, 'JAB-07', NULL),
			('123456782','123456782','unor-2','Hadi', 7, 'JAB-07', NOW());

		INSERT INTO file_digital_signature
			("file_id", "nip_sk", "kategori", "no_sk", "tanggal_sk", "status_sk", "created_at", "deleted_at") VALUES
			('sk-001', '123456789', 'Kenaikan Pangkat', 'SK-001/2024', '2024-01-15', 1, '2024-01-15', NULL),
			('sk-002', '123456789', 'Mutasi', 'SK-002/2024', '2024-02-20', 0, '2024-02-20', NULL),
			('sk-003', '123456789', 'Kenaikan Gaji', 'SK-003/2024', '2024-03-10', 2, '2024-03-10', NULL),
			('sk-004', '123456789', 'Kenaikan Gaji', 'SK-004/2024', '2024-03-10', 2, '2024-03-10', NOW()),
			('sk-005', '123456788', 'Mutasi', 'SK-005/2024', '2024-03-11', 3, '2024-03-11', NULL),
			('sk-006', '123456787', 'Mutasi', 'SK-006/2024', '2024-03-12', 4, '2024-03-12', NULL),
			('sk-007', '123456786', 'Mutasi', 'SK-007/2024', '2024-03-13', 5, '2024-03-13', NULL),
			('sk-008', '123456785', 'Mutasi', 'SK-008/2024', '2024-03-14', 5, '2024-03-14', NULL),
			('sk-009', '123456784', 'Mutasi', 'SK-009/2024', '2024-03-15', 5, '2024-03-15', NULL),
			('sk-010', '123456783', 'Mutasi', 'SK-010/2024', '2024-03-16', 5, '2024-03-16', NULL),
			('sk-011', '123456782', 'Kenaikan Pangkat', 'SK-011/2024', '2024-01-17', 1, '2024-01-15', NULL);
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
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-010",
						"nama_pemilik": "Gledi",
						"nip_pemilik": "123456783",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-010/2024",
						"tanggal_sk": "2024-03-16",
						"unit_kerja": "Tengah - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-009",
						"nama_pemilik": "Fredi",
						"nip_pemilik": "123456784",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-009/2024",
						"tanggal_sk": "2024-03-15",
						"unit_kerja": "Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-008",
						"nama_pemilik": "Edi",
						"nip_pemilik": "123456785",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-008/2024",
						"tanggal_sk": "2024-03-14",
						"unit_kerja": "Paling Atas 2",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-006",
						"nama_pemilik": "Candi",
						"nip_pemilik": "123456787",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-006/2024",
						"tanggal_sk": "2024-03-12",
						"unit_kerja": "Bawah 2",
						"status_sk": "Sudah TTD"
					},
					{
						"id_sk": "sk-005",
						"nama_pemilik": "Budi",
						"nip_pemilik": "123456788",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-005/2024",
						"tanggal_sk": "2024-03-11",
						"unit_kerja": "",
						"status_sk": "Sudah dikoreksi dan menunggu TTD"
					},
					{
						"id_sk": "sk-003",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Kenaikan Gaji",
						"no_sk": "SK-003/2024",
						"tanggal_sk": "2024-03-10",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Sudah Dikoreksi & Dikembalikan"
					},
					{
						"id_sk": "sk-002",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-002/2024",
						"tanggal_sk": "2024-02-20",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Belum Dikoreksi"
					},
					{
						"id_sk": "sk-001",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Kenaikan Pangkat",
						"no_sk": "SK-001/2024",
						"tanggal_sk": "2024-01-15",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Sedang Dikoreksi"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 9
				}
			}`,
		},
		{
			name:             "ok with filter unit_kerja_id",
			dbData:           dbData,
			requestQuery:     url.Values{"unit_kerja_id": []string{"unor-1"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-010",
						"nama_pemilik": "Gledi",
						"nip_pemilik": "123456783",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-010/2024",
						"tanggal_sk": "2024-03-16",
						"unit_kerja": "Tengah - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-009",
						"nama_pemilik": "Fredi",
						"nip_pemilik": "123456784",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-009/2024",
						"tanggal_sk": "2024-03-15",
						"unit_kerja": "Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-006",
						"nama_pemilik": "Candi",
						"nip_pemilik": "123456787",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-006/2024",
						"tanggal_sk": "2024-03-12",
						"unit_kerja": "Bawah 2",
						"status_sk": "Sudah TTD"
					},
					{
						"id_sk": "sk-003",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Kenaikan Gaji",
						"no_sk": "SK-003/2024",
						"tanggal_sk": "2024-03-10",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Sudah Dikoreksi & Dikembalikan"
					},
					{
						"id_sk": "sk-002",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-002/2024",
						"tanggal_sk": "2024-02-20",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Belum Dikoreksi"
					},
					{
						"id_sk": "sk-001",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Kenaikan Pangkat",
						"no_sk": "SK-001/2024",
						"tanggal_sk": "2024-01-15",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Sedang Dikoreksi"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 7
				}
			}`,
		},
		{
			name:             "ok with filter nama_pemilik",
			dbData:           dbData,
			requestQuery:     url.Values{"nama_pemilik": []string{"Edi"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-010",
						"nama_pemilik": "Gledi",
						"nip_pemilik": "123456783",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-010/2024",
						"tanggal_sk": "2024-03-16",
						"unit_kerja": "Tengah - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-009",
						"nama_pemilik": "Fredi",
						"nip_pemilik": "123456784",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-009/2024",
						"tanggal_sk": "2024-03-15",
						"unit_kerja": "Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-008",
						"nama_pemilik": "Edi",
						"nip_pemilik": "123456785",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-008/2024",
						"tanggal_sk": "2024-03-14",
						"unit_kerja": "Paling Atas 2",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 4
				}
			}`,
		},
		{
			name:             "ok with filter nip",
			dbData:           dbData,
			requestQuery:     url.Values{"nip": []string{"123456783"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-010",
						"nama_pemilik": "Gledi",
						"nip_pemilik": "123456783",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-010/2024",
						"tanggal_sk": "2024-03-16",
						"unit_kerja": "Tengah - Paling Atas",
						"status_sk": "Sudah dikirim"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:             "ok with filter golongan_id",
			dbData:           dbData,
			requestQuery:     url.Values{"golongan_id": []string{"5"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-008",
						"nama_pemilik": "Edi",
						"nip_pemilik": "123456785",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-008/2024",
						"tanggal_sk": "2024-03-14",
						"unit_kerja": "Paling Atas 2",
						"status_sk": "Sudah dikirim"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:             "ok with filter jabatan_id",
			dbData:           dbData,
			requestQuery:     url.Values{"jabatan_id": []string{"JAB-04"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:             "ok with filter kategori_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"kategori_sk": []string{"Kenaikan Gaji"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-003",
						"nama_pemilik": "Andi",
						"nip_pemilik": "123456789",
						"kategori_sk": "Kenaikan Gaji",
						"no_sk": "SK-003/2024",
						"tanggal_sk": "2024-03-10",
						"unit_kerja": "Bawah - Tengah - Paling Atas",
						"status_sk": "Sudah Dikoreksi & Dikembalikan"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:             "ok with filter tanggal_sk_mulai & tanggal_sk_akhir",
			dbData:           dbData,
			requestQuery:     url.Values{"tanggal_sk_mulai": []string{"2024-03-12"}, "tanggal_sk_akhir": []string{"2024-03-14"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-008",
						"nama_pemilik": "Edi",
						"nip_pemilik": "123456785",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-008/2024",
						"tanggal_sk": "2024-03-14",
						"unit_kerja": "Paling Atas 2",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-006",
						"nama_pemilik": "Candi",
						"nip_pemilik": "123456787",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-006/2024",
						"tanggal_sk": "2024-03-12",
						"unit_kerja": "Bawah 2",
						"status_sk": "Sudah TTD"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 3
				}
			}`,
		},
		{
			name:             "ok with filter status_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"status_sk": []string{"4"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-006",
						"nama_pemilik": "Candi",
						"nip_pemilik": "123456787",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-006/2024",
						"tanggal_sk": "2024-03-12",
						"unit_kerja": "Bawah 2",
						"status_sk": "Sudah TTD"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:             "ok with filter multiple status_sk",
			dbData:           dbData,
			requestQuery:     url.Values{"status_sk": []string{"3", "4"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-006",
						"nama_pemilik": "Candi",
						"nip_pemilik": "123456787",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-006/2024",
						"tanggal_sk": "2024-03-12",
						"unit_kerja": "Bawah 2",
						"status_sk": "Sudah TTD"
					},
					{
						"id_sk": "sk-005",
						"nama_pemilik": "Budi",
						"nip_pemilik": "123456788",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-005/2024",
						"tanggal_sk": "2024-03-11",
						"unit_kerja": "",
						"status_sk": "Sudah dikoreksi dan menunggu TTD"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 2
				}
			}`,
		},
		{
			name:   "ok with all filter",
			dbData: dbData,
			requestQuery: url.Values{
				"unit_kerja_id":    []string{"unor-1"},
				"nama_pemilik":     []string{"Dedi"},
				"nip":              []string{"123456786"},
				"golongan_id":      []string{"4"},
				"jabatan_id":       []string{"JAB-04"},
				"kategori_sk":      []string{"Mutasi"},
				"tanggal_sk_mulai": []string{"2024-03-13"},
				"tanggal_sk_akhir": []string{"2024-03-15"},
				"status_sk":        []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:             "ok with limit and offset",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"3"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id_sk": "sk-007",
						"nama_pemilik": "Dedi",
						"nip_pemilik": "123456786",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-007/2024",
						"tanggal_sk": "2024-03-13",
						"unit_kerja": "Bawah 3 - Paling Atas",
						"status_sk": "Sudah dikirim"
					},
					{
						"id_sk": "sk-006",
						"nama_pemilik": "Candi",
						"nip_pemilik": "123456787",
						"kategori_sk": "Mutasi",
						"no_sk": "SK-006/2024",
						"tanggal_sk": "2024-03-12",
						"unit_kerja": "Bawah 2",
						"status_sk": "Sudah TTD"
					}
				],
				"meta": {
					"limit": 2,
					"offset": 3,
					"total": 9
				}
			}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
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
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/surat-keputusan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := dbrepository.New(db)
			suratkeputusan.RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_getAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
			(id,  diatasan_id, nama_unor, nama_jabatan, pemimpin_pns_id, eselon_1, eselon_2, eselon_3, eselon_4, deleted_at) values
			('unor-1', null, 'Paling Atas', 'Atasan 1', null, null, null, null, null, null),
			('unor-2', 'unor-1', 'Tengah', 'Atasan 2', null, null, 'unor-1', null, null, null),
			('unor-3', 'unor-2', 'Bawah', 'Atasan 3', null, 'unor-1', null, null, null, null),
			('unor-4', 'unor-1', 'Tengah deleted', 'Atasan 4', null, null, null, null, null, now()),
			('unor-5', 'unor-4', 'Bawah 2', 'Atasan 5', null, null, null, null, 'unor-1', null),
			('unor-6', 'unor-1', 'Bawah 3', 'Atasan 6', null, null, null, 'unor-1', null, null),
			('unor-7', null, 'Paling Atas 2', 'Atasan 7', null, null, null, null, null, null);

		INSERT INTO ref_golongan (id, nama, nama_pangkat, nama_2, gol, gol_pppk)
		VALUES
			(1, 'I/a', 'Juru Muda', 'I/a', 1, 'I/a'),
			(2, 'I/b', 'Juru Muda Tingkat I', 'I/b', 2, 'I/b'),
			(3, 'II/a', 'Pengatur Muda', 'II/a', 3, 'II/a'),
			(4, 'II/b', 'Pengatur Muda Tingkat I', 'II/b', 4, 'II/b'),
			(5, 'III/a', 'Penata Muda', 'III/a', 5, 'III/a'),
			(6, 'III/b', 'Penata Muda Tingkat I', 'III/b', 6, 'III/b'),
			(7, 'IV/a', 'Pembina', 'IV/a', 7, 'IV/a');

		INSERT INTO ref_jabatan (kode_jabatan, id, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id)
		VALUES
			('JAB-01', 1, 'Kepala Bagian', 'Kepala Bagian Administrasi', 1, 3, 60, 'BKN-01', 'Kepala Bagian BKN', 'Struktural', 'BKN-01'),
			('JAB-02', 2, 'Kepala Subbag', 'Kepala Subbag Umum', 1, 2, 60, 'BKN-02', 'Kepala Subbag BKN', 'Struktural', 'BKN-02'),
			('JAB-03', 3, 'Pengawas', 'Pengawas Pelayanan', 2, 2, 58, 'BKN-03', 'Pengawas BKN', 'Fungsional', 'BKN-03'),
			('JAB-04', 4, 'Staf Ahli', 'Staf Ahli Bidang Kepegawaian', 2, 3, 58, 'BKN-04', 'Staf Ahli BKN', 'Fungsional', 'BKN-04'),
			('JAB-05', 5, 'Kepala Seksi', 'Kepala Seksi Administrasi', 1, 2, 60, 'BKN-05', 'Kepala Seksi BKN', 'Struktural', 'BKN-05'),
			('JAB-06', 6, 'Analis', 'Analis Kepegawaian', 2, 1, 58, 'BKN-06', 'Analis BKN', 'Fungsional', 'BKN-06'),
			('JAB-07', 7, 'Operator', 'Operator Sistem', 2, 1, 58, 'BKN-07', 'Operator BKN', 'Fungsional', 'BKN-07');

		INSERT INTO pegawai ("nip_baru", "pns_id", "unor_id", "nama", "gol_id", "jabatan_instansi_id", deleted_at) 
		VALUES
			('123456789','123456789','unor-3','Andi', 1, 'JAB-01', NULL),
			('123456788','123456788','unor-4','Budi', 2, 'JAB-02', NULL),
			('123456787','123456787','unor-5','Candi', 3, 'JAB-03', NULL),
			('123456786','123456786','unor-6','Dedi', 4, 'JAB-04', NULL),
			('123456785','123456785','unor-7','Edi', 5, 'JAB-05', NULL),
			('123456784','123456784','unor-1','Fredi', 6, 'JAB-06', NULL),
			('123456783','123456783','unor-2','Gledi', 7, 'JAB-07', NULL),
			('123456782','123456782','unor-2','Hadi', 7, 'JAB-07', NULL),
			('123456781','123456781','unor-2','Ikhdi', 7, 'JAB-07', NOW());

		INSERT INTO file_digital_signature
			("file_id", "nip_sk", "kategori", "no_sk", "tanggal_sk", "status_sk", "nip_pemroses", "created_at", "deleted_at") VALUES
			('sk-001', '123456789', 'Kenaikan Pangkat', 'SK-001/2024', '2024-01-15', 1, NULL, '2024-01-15', NULL),
			('sk-002', '123456789', 'Mutasi', 'SK-002/2024', '2024-02-20', 0, NULL, '2024-02-20', NULL),
			('sk-003', '123456789', 'Kenaikan Gaji', 'SK-003/2024', '2024-03-10', 2, NULL, '2024-03-10', NULL),
			('sk-004', '123456789', 'Kenaikan Gaji', 'SK-004/2024', '2024-03-10', 2, NULL, '2024-03-10', NOW()),
			('sk-005', '123456788', 'Mutasi', 'SK-005/2024', '2024-03-11', 3, NULL, '2024-03-11', NULL),
			('sk-006', '123456787', 'Mutasi', 'SK-006/2024', '2024-03-12', 4, '123456782', '2024-03-12', NULL),
			('sk-007', '123456786', 'Mutasi', 'SK-007/2024', '2024-03-13', 5, '123456782', '2024-03-13', NULL),
			('sk-008', '123456785', 'Mutasi', 'SK-008/2024', '2024-03-14', 5, '123456782', '2024-03-14', NULL),
			('sk-009', '123456784', 'Mutasi', 'SK-009/2024', '2024-03-15', 5, '123456782', '2024-03-15', NULL),
			('sk-010', '123456783', 'Mutasi', 'SK-010/2024', '2024-03-16', 5, '123456782', '2024-03-16', NULL),
			('sk-011', '123456781', 'Kenaikan Pangkat', 'SK-011/2024', '2024-01-17', 1, NULL, '2024-01-15', NULL);

		INSERT INTO file_digital_signature_riwayat
			("file_id","nip_pemroses","tindakan","created_at","deleted_at") VALUES
			('sk-010','123456782','Mengkoreksi','2025-09-24T01:02:03',NULL),
			('sk-010','123456782','Mengkoreksi',NOW(),NOW()),
			('sk-010','123456781','Mengkoreksi','2025-09-25T01:02:03',NULL);
	`

	tests := []struct {
		name             string
		dbData           string
		requestPath      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-010",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-010",
					"nama_pemilik": "Gledi",
					"nip_pemilik": "123456783",
					"kategori_sk": "Mutasi",
					"no_sk": "SK-010/2024",
					"tanggal_sk": "2024-03-16",
					"unit_kerja": "Tengah - Paling Atas",
					"status_sk": "Sudah dikirim",
					"nama_penandatangan": "Hadi",
					"jabatan_penandatangan": "Operator",
					"logs": [
						{
							"log" : "Mengkoreksi",
							"actor" : "Hadi",
							"timestamp" : "` + time.Date(2025, time.September, 24, 1, 2, 3, 0, time.UTC).Local().Format(time.RFC3339) + `"
						},
						{
							"log" : "Mengkoreksi",
							"actor" : "",
							"timestamp" : "` + time.Date(2025, time.September, 25, 1, 2, 3, 0, time.UTC).Local().Format(time.RFC3339) + `"
						}
					]
				}
			}`,
		},
		{
			name:             "ok with no signer",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-003",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-003",
					"nama_pemilik": "Andi",
					"nip_pemilik": "123456789",
					"kategori_sk": "Kenaikan Gaji",
					"no_sk": "SK-003/2024",
					"tanggal_sk": "2024-03-10",
					"unit_kerja": "Bawah - Tengah - Paling Atas",
					"status_sk": "Sudah Dikoreksi & Dikembalikan",
					"logs": []
				}
			}`,
		},
		{
			name:             "ok with unor utama deleted",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-005",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-005",
					"nama_pemilik": "Budi",
					"nip_pemilik": "123456788",
					"kategori_sk": "Mutasi",
					"no_sk": "SK-005/2024",
					"tanggal_sk": "2024-03-11",
					"unit_kerja": "",
					"status_sk": "Sudah dikoreksi dan menunggu TTD",
					"logs": []
				}
			}`,
		},
		{
			name:             "ok with parent unor deleted",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-006",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id_sk": "sk-006",
					"nama_pemilik": "Candi",
					"nip_pemilik": "123456787",
					"kategori_sk": "Mutasi",
					"no_sk": "SK-006/2024",
					"tanggal_sk": "2024-03-12",
					"unit_kerja": "Bawah 2",
					"status_sk": "Sudah TTD",
					"nama_penandatangan": "Hadi",
					"jabatan_penandatangan": "Operator",
					"logs": []
				}
			}`,
		},
		{
			name:             "error: sk not found",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-999",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: sk deleted",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-004",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: pegawai deleted",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-011",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-004",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			requestPath:      "/v1/admin/surat-keputusan/sk-004",
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
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := dbrepository.New(db)
			suratkeputusan.RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_getBerkasAdmin(t *testing.T) {
	t.Parallel()

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

	signedBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x02, 0x00, 0x00, 0x00, 0x90, 0x77, 0x53,
		0xde, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x08, 0xd7, 0x63, 0xf8, 0xcf, 0xc0, 0x00,
		0x00, 0x04, 0x00, 0x01, 0xe2, 0x26, 0x05, 0x9b,
		0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e, 0x44,
		0xae, 0x42, 0x60, 0x82,
	}

	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)
	signedPngBase64 := base64.StdEncoding.EncodeToString(signedBytes)

	dbData := `

		INSERT INTO file_digital_signature
				("file_id","file_base64","file_base64_sign", "nip_sk", "kategori", "no_sk", "tanggal_sk", "status_sk","nip_pemroses", "created_at", "deleted_at") VALUES
				('sk-001','data:image/png;base64,` + pngBase64 + `',NULL, '123456789', 'Kenaikan Pangkat', 'SK-001/2024', '2024-01-15', 1, '12345678', '2024-01-15', NULL),
				('sk-002','data:image/png;base64,invalid',NULL, '123456789', 'Mutasi', 'SK-002/2024', '2024-02-20', 0, '12345678', '2024-02-20', NULL),
				('sk-003',NULL,NULL, '123456789', 'Kenaikan Gaji', 'SK-003/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-004','',NULL, '123456789', 'Kenaikan Gaji', 'SK-004/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-005','data:image/png;base64,` + pngBase64 + `',NULL, '123456789', 'Kenaikan Gaji', 'SK-005/2024', '2024-03-10', 2,'12345678', '2024-03-10', NOW()),
				('sk-006',NULL,'data:image/png;base64,` + signedPngBase64 + `', '123456789', 'Kenaikan Gaji', 'SK-006/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL),
				('sk-007','data:image/png;base64,` + pngBase64 + `','data:image/png;base64,` + signedPngBase64 + `', '123456789', 'Kenaikan Gaji', 'SK-007/2024', '2024-03-10', 2,'12345678', '2024-03-10', NULL);
			
	`

	tests := []struct {
		name              string
		dbData            string
		requestPath       string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-001/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "ok: valid signed png",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-006/berkas?signed=true",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: signedBytes,
		},
		{
			name:              "ok: valid unsigned png",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-007/berkas?signed=false",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "ok: valid signed png with params 1",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-007/berkas?signed=1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: signedBytes,
		},
		{
			name:              "error: empty value params signed",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-007/berkas?signed=",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message":"parameter \"signed\" tidak boleh kosong"}`),
		},
		{
			name:              "error: invalid params signed",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-007/berkas?signed=''",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message":"parameter \"signed\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: base64 berkas SK tidak valid",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-002/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas SK sudah dihapus",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-005/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas SK berisi null value",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-003/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas SK berupa string kosong",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-004/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: berkas SK tidak ditemukan",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-009/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas SK tidak ditemukan"}`),
		},
		{
			name:              "error: user bukan admin",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-001/berkas",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode:  http.StatusForbidden,
			wantResponseBytes: []byte(`{"message": "akses ditolak"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			requestPath:       "/v1/admin/surat-keputusan/sk-001/berkas",
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

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := dbrepository.New(pgxconn)
			suratkeputusan.RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
