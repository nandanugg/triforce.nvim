package hukumandisiplin

import (
	"context"
	"encoding/base64"
	"fmt"
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
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman (id, nama, tingkat_hukuman, nama_tingkat_hukuman, deleted_at)
		VALUES
		(1, 'Jenis Hukuman 1', NULL, NULL, NULL),
		(2, 'Jenis Hukuman 2', NULL, NULL, NULL),
		(3, 'Jenis Hukuman 2', NULL, NULL, '2023-02-20');

		INSERT INTO pegawai (pns_id, nip_baru, nama)
		VALUES ('id1', '198765432100001', 'Budi');

		INSERT INTO ref_golongan (id, nama, nama_pangkat, nama_2, gol, gol_pppk, deleted_at)
		VALUES
			(1, 'I/a', 'Juru Muda', 'Ia', 1, 'I', NULL),
			(2, 'II/a', 'Pengatur Muda', 'IIa', 2, 'II', '2023-02-20');

		INSERT INTO riwayat_hukdis (
			pns_id, pns_nip, nama, golongan_id, nama_golongan,
			jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal,
			tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman,
			no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, bkn_id, file_base64, keterangan_berkas,
			deleted_at
		)
		VALUES
		('id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK1','2023-01-15','2023-01-20',0,1,'2023-02-20',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK2','2023-03-10','2023-03-15',0,2,'2023-05-15',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK3','2023-06-01','2023-06-10',1,0,'2024-06-10',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK4','2023-09-01','2023-09-15',2,0,'2025-09-15',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',2,NULL,'SK5','2023-12-01','2023-12-10',3,0,'2026-12-10',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',3,NULL,'SK6','2023-12-02','2023-12-11',3,0,'2026-12-11',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',3,NULL,'SK7','2023-12-03','2023-12-12',3,0,'2026-12-12',NULL,NULL,NULL,NULL,NULL,NULL,'2023-02-20'),
		('id1','198765432100001','Budi',2,'I/a',1,NULL,'SK8','2023-12-04','2023-12-13',3,0,'2026-12-13',NULL,NULL,NULL,NULL,NULL,NULL,NULL);
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 8,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "",
						"nama_pangkat": "",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK8",
						"tanggal_akhir": "2026-12-13",
						"tanggal_mulai": "2023-12-13",
						"tanggal_sk": "2023-12-04"
					},
					{
						"id": 6,
						"jenis_hukuman": "",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK6",
						"tanggal_akhir": "2026-12-11",
						"tanggal_mulai": "2023-12-11",
						"tanggal_sk": "2023-12-02"
					},
					{
						"id": 5,
						"jenis_hukuman": "Jenis Hukuman 2",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK5",
						"tanggal_akhir": "2026-12-10",
						"tanggal_mulai": "2023-12-10",
						"tanggal_sk": "2023-12-01"
					},
					{
						"id": 4,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 2,
						"nomor_sk": "SK4",
						"tanggal_akhir": "2025-09-15",
						"tanggal_mulai": "2023-09-15",
						"tanggal_sk": "2023-09-01"
					},
					{
						"id": 3,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 1,
						"nomor_sk": "SK3",
						"tanggal_akhir": "2024-06-10",
						"tanggal_mulai": "2023-06-10",
						"tanggal_sk": "2023-06-01"
					},
					{
						"id": 2,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 2,
						"masa_tahun": 0,
						"nomor_sk": "SK2",
						"tanggal_akhir": "2023-05-15",
						"tanggal_mulai": "2023-03-15",
						"tanggal_sk": "2023-03-10"
					},
					{
						"id": 1,
						"jenis_hukuman": "Snapshotted Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 1,
						"masa_tahun": 0,
						"nomor_sk": "SK1",
						"tanggal_akhir": "2023-02-20",
						"tanggal_mulai": "2023-01-20",
						"tanggal_sk": "2023-01-15"
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
			name:             "ok: dengan parameter pagination",
			dbData:           dbData,
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"2"}},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 5,
						"jenis_hukuman": "Jenis Hukuman 2",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK5",
						"tanggal_akhir": "2026-12-10",
						"tanggal_mulai": "2023-12-10",
						"tanggal_sk": "2023-12-01"
					}
				],
				"meta": {
					"limit": 1,
					"offset": 2,
					"total": 7
				}
			}`,
		},
		{
			name:             "ok: tidak ada data milik user",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100002")}},
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

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-hukuman-disiplin", nil)
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

	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		INSERT INTO ref_jenis_hukuman (id, nama, tingkat_hukuman, nama_tingkat_hukuman)
		VALUES
		(1, 'Jenis Hukuman 1', NULL, NULL);

		INSERT INTO pegawai (pns_id, nip_baru, nama)
		VALUES ('id1', '198765432100001', 'Budi');

		INSERT INTO riwayat_hukdis (
			id, pns_id, pns_nip, nama, golongan_id, nama_golongan,
			jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal,
			tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman,
			no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, bkn_id, file_base64, keterangan_berkas,
			deleted_at
		)
		VALUES
		(1, 'id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK1','2023-01-15','2023-01-20',0,1,'2023-02-20',NULL,NULL,NULL,NULL,'data:image/png;base64,` + pngBase64 + `',NULL,NULL),
		(2, 'id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK2','2023-01-15','2023-01-20',0,1,'2023-02-20',NULL,NULL,NULL,NULL,'data:image/png;base64,invalid',NULL,NULL),
		(3, 'id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK2','2023-01-15','2023-01-20',0,1,'2023-02-20',NULL,NULL,NULL,NULL,'data:image/png;base64,invalid',NULL,'2023-02-20'),
		(4, 'id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK2','2023-01-15','2023-01-20',0,1,'2023-02-20',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		(5, 'id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK2','2023-01-15','2023-01-20',0,1,'2023-02-20',NULL,NULL,NULL,NULL,'',NULL,NULL);
		`

	tests := []struct {
		name              string
		dbData            string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 berkas riwayat hukuman disiplin tidak valid",
			dbData:            dbData,
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas riwayat hukuman disiplin sudah dihapus",
			dbData:            dbData,
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat hukuman disiplin berisi null value",
			dbData:            dbData,
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat hukuman disiplin berupa string kosong",
			dbData:            dbData,
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat hukuman disiplin bukan milik user login",
			dbData:            dbData,
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100002")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat hukuman disiplin tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "198765432100001")}},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			dbData:            dbData,
			paramID:           "1",
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-hukuman-disiplin/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := sqlc.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
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
