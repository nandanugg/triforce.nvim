package riwayathukumandisiplin

import (
	"bytes"
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
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

		INSERT INTO riwayat_hukuman_disiplin (
			pns_id, pns_nip, nama, golongan_id, nama_golongan,
			jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal,
			tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman,
			no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, bkn_id, file_base64, keterangan_berkas,
			deleted_at
		)
		VALUES
		('id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK1','2023-01-15','2023-01-20',0,1,'2023-02-20','PP-1','DEL-1','2023-01-16',NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK2','2023-03-10','2023-03-15',0,2,'2023-05-15',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK3','2023-06-01','2023-06-10',1,0,'2024-06-10',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK4','2023-09-01','2023-09-15',2,0,'2025-09-15',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',2,NULL,'SK5','2023-12-01','2023-12-10',3,0,'2026-12-10','PP-5','DEL-5','2023-12-02',NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',3,NULL,'SK6','2023-12-02','2023-12-11',3,0,'2026-12-11',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',3,NULL,'SK7','2023-12-03','2023-12-12',3,0,'2026-12-12',NULL,NULL,NULL,NULL,NULL,NULL,'2023-02-20'),
		('id1','198765432100001','Budi',2,'I/a',1,NULL,'SK8','2023-12-04','2023-12-13',3,0,'2026-12-13','PP-8','DEL-8',NULL,NULL,NULL,NULL,NULL);
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("198765432100001")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: tanpa parameter apapun",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 8,
						"golongan_id": 2,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "",
						"nama_pangkat": "",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK8",
						"tanggal_akhir": "2026-12-13",
						"tanggal_mulai": "2023-12-13",
						"tanggal_sk": "2023-12-04",
						"nomor_pp": "PP-8",
						"nomor_sk_pembatalan": "DEL-8",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 6,
						"golongan_id": 1,
						"jenis_hukuman_id": 3,
						"jenis_hukuman": "",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK6",
						"tanggal_akhir": "2026-12-11",
						"tanggal_mulai": "2023-12-11",
						"tanggal_sk": "2023-12-02",
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 5,
						"golongan_id": 1,
						"jenis_hukuman_id": 2,
						"jenis_hukuman": "Jenis Hukuman 2",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK5",
						"tanggal_akhir": "2026-12-10",
						"tanggal_mulai": "2023-12-10",
						"tanggal_sk": "2023-12-01",
						"nomor_pp": "PP-5",
						"nomor_sk_pembatalan": "DEL-5",
						"tanggal_sk_pembatalan": "2023-12-02"
					},
					{
						"id": 4,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 2,
						"nomor_sk": "SK4",
						"tanggal_akhir": "2025-09-15",
						"tanggal_mulai": "2023-09-15",
						"tanggal_sk": "2023-09-01",
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 3,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 1,
						"nomor_sk": "SK3",
						"tanggal_akhir": "2024-06-10",
						"tanggal_mulai": "2023-06-10",
						"tanggal_sk": "2023-06-01",
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 2,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 2,
						"masa_tahun": 0,
						"nomor_sk": "SK2",
						"tanggal_akhir": "2023-05-15",
						"tanggal_mulai": "2023-03-15",
						"tanggal_sk": "2023-03-10",
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 1,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Snapshotted Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 1,
						"masa_tahun": 0,
						"nomor_sk": "SK1",
						"tanggal_akhir": "2023-02-20",
						"tanggal_mulai": "2023-01-20",
						"tanggal_sk": "2023-01-15",
						"nomor_pp": "PP-1",
						"nomor_sk_pembatalan": "DEL-1",
						"tanggal_sk_pembatalan": "2023-01-16"
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
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"2"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 5,
						"golongan_id": 1,
						"jenis_hukuman_id": 2,
						"jenis_hukuman": "Jenis Hukuman 2",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"masa_bulan": 0,
						"masa_tahun": 3,
						"nomor_sk": "SK5",
						"tanggal_akhir": "2026-12-10",
						"tanggal_mulai": "2023-12-10",
						"tanggal_sk": "2023-12-01",
						"nomor_pp": "PP-5",
						"nomor_sk_pembatalan": "DEL-5",
						"tanggal_sk_pembatalan": "2023-12-02"
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("198765432100002")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/riwayat-hukuman-disiplin", nil)
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

		INSERT INTO riwayat_hukuman_disiplin (
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
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Self)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("198765432100001")}
	tests := []struct {
		name              string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid png",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "image/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 berkas riwayat hukuman disiplin tidak valid",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: berkas riwayat hukuman disiplin sudah dihapus",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat hukuman disiplin berisi null value",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berkas riwayat hukuman disiplin berupa string kosong",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat hukuman disiplin bukan milik user login",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("198765432100002")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: berkas riwayat hukuman disiplin tidak ditemukan",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/riwayat-hukuman-disiplin/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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
		INSERT INTO ref_jenis_hukuman (id, nama, tingkat_hukuman, nama_tingkat_hukuman, deleted_at)
		VALUES
		(1, 'Jenis Hukuman 1', NULL, NULL, NULL),
		(2, 'Jenis Hukuman 2', NULL, NULL, NULL),
		(3, 'Jenis Hukuman 2', NULL, NULL, '2023-02-20');

		INSERT INTO pegawai (pns_id, nip_baru, nama, deleted_at)
		VALUES ('id1', '198765432100001', 'Budi', NULL),
		('id2', '198765432100002', 'Ani', '2023-02-20');

		INSERT INTO ref_golongan (id, nama, nama_pangkat, nama_2, gol, gol_pppk, deleted_at)
		VALUES
			(1, 'I/a', 'Juru Muda', 'Ia', 1, 'I', NULL),
			(2, 'II/a', 'Pengatur Muda', 'IIa', 2, 'II', '2023-02-20');

		INSERT INTO riwayat_hukuman_disiplin (
			pns_id, pns_nip, nama, golongan_id, nama_golongan,
			jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal,
			tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman,
			no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, bkn_id, file_base64, keterangan_berkas,
			deleted_at
		)
		VALUES
		('id1','198765432100001','Budi',1,'I/a',1,'Snapshotted Jenis Hukuman 1','SK1','2023-01-15','2023-01-20',0,1,'2023-02-20','PP-1','DEL-1','2023-01-16',NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK2','2023-03-10','2023-03-15',0,2,'2023-05-15',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK3','2023-06-01','2023-06-10',1,0,'2024-06-10',NULL,NULL,NULL,NULL,NULL,NULL,NULL),
		('id1','198765432100001','Budi',1,'I/a',1,NULL,'SK7','2023-12-03','2023-12-12',3,0,'2026-12-12',NULL,NULL,NULL,NULL,NULL,NULL,'2023-02-20');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(db)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

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
			name:             "ok: nip 198765432100001 data returned",
			nip:              "198765432100001",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 3,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"nomor_sk": "SK3",
						"tanggal_sk": "2023-06-01",
						"tanggal_mulai": "2023-06-10",
						"tanggal_akhir": "2024-06-10",
						"masa_tahun": 1,
						"masa_bulan": 0,
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 2,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"nomor_sk": "SK2",
						"tanggal_sk": "2023-03-10",
						"tanggal_mulai": "2023-03-15",
						"tanggal_akhir": "2023-05-15",
						"masa_tahun": 0,
						"masa_bulan": 2,
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					},
					{
						"id": 1,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Snapshotted Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"nomor_sk": "SK1",
						"tanggal_sk": "2023-01-15",
						"tanggal_mulai": "2023-01-20",
						"tanggal_akhir": "2023-02-20",
						"masa_tahun": 0,
						"masa_bulan": 1,
						"nomor_pp": "PP-1",
						"nomor_sk_pembatalan": "DEL-1",
						"tanggal_sk_pembatalan": "2023-01-16"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 3}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			nip:              "198765432100001",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `
			{
				"data": [
					{
						"id": 2,
						"golongan_id": 1,
						"jenis_hukuman_id": 1,
						"jenis_hukuman": "Jenis Hukuman 1",
						"nama_golongan": "I/a",
						"nama_pangkat": "Juru Muda",
						"nomor_sk": "SK2",
						"tanggal_sk": "2023-03-10",
						"tanggal_mulai": "2023-03-15",
						"tanggal_akhir": "2023-05-15",
						"masa_tahun": 0,
						"masa_bulan": 2,
						"nomor_pp": "",
						"nomor_sk_pembatalan": "",
						"tanggal_sk_pembatalan": null
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: nip 200 gets empty data",
			nip:              "200",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "ok: nip 198765432100002 gets empty data (deleted pegawai)",
			nip:              "198765432100002",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			nip:              "198765432100001",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pegawai/"+tt.nip+"/riwayat-hukuman-disiplin", nil)
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

func Test_handler_getBerkasAdmin(t *testing.T) {
	t.Parallel()

	filePath := "../../../../lib/api/sample/hello.pdf"
	pdfBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

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

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		INSERT INTO riwayat_hukuman_disiplin
			(id, pns_nip, deleted_at, file_base64) VALUES
			(1, '1c', null, 'data:application/pdf;base64,` + pdfBase64 + `'),
			(2, '1c', null, '` + pdfBase64 + `'),
			(3, '1c', null, 'data:images/png;base64,` + pngBase64 + `'),
			(4, '1c', null, 'data:application/pdf;base64,invalid'),
			(5, '1c', '2020-01-02', 'data:application/pdf;base64,` + pdfBase64 + `'),
			(6, '1c', null, null),
			(7, '1c', null, ''),
			(8, '2a', null, 'data:application/pdf;base64,` + pdfBase64 + `');
		`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err = pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgxconn)
	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name              string
		nip               string
		paramID           string
		requestHeader     http.Header
		wantResponseCode  int
		wantContentType   string
		wantResponseBytes []byte
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			nip:               "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			nip:               "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "ok: admin can access other user's berkas",
			nip:               "2a",
			paramID:           "8",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "error: base64 tidak valid",
			nip:               "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: riwayat sudah dihapus",
			nip:               "1c",
			paramID:           "5",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			nip:               "1c",
			paramID:           "6",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			nip:               "1c",
			paramID:           "7",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat with wrong nip",
			nip:               "wrong-nip",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: riwayat tidak ditemukan",
			nip:               "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas riwayat hukuman disiplin tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			nip:               "1c",
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": authHeader},
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBytes: []byte(`{"message": "parameter \"id\" harus dalam format yang sesuai"}`),
		},
		{
			name:              "error: auth header tidak valid",
			nip:               "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBytes: []byte(`{"message": "token otentikasi tidak valid"}`),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/pegawai/%s/riwayat-hukuman-disiplin/%s/berkas", tt.nip, tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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

func Test_handler_adminCreate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai
			(pns_id,  nip_baru, nama, deleted_at) VALUES
			('id_1a', '1a',     'Pegawai 1', NULL),
			('id_1c', '1c',     'Pegawai 2', NULL),
			('id_1d', '1d',     'Pegawai 3', '2000-01-01'),
			('id_1e', '1e',     'Pegawai 4', NULL),
			('id_1f', '1f',     'Pegawai 5', NULL),
			('id_1g', '1g',     'Pegawai 6', NULL);
		INSERT INTO ref_jenis_hukuman
			(id, nama, deleted_at) VALUES
			(1, 'Jenis Hukuman 1', NULL),
			(2, 'Jenis Hukuman 2', '2000-01-01');
		INSERT INTO ref_golongan
			(id, nama, nama_pangkat, nama_2, gol, gol_pppk, deleted_at) VALUES
			(1, 'I/a', 'Juru Muda', 'Ia', 1, 'I', NULL),
			(2, 'II/a', 'Pengatur Muda', 'IIa', 2, 'II', '2000-01-01');
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-001",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20",
				"nomor_pp": "PP-001",
				"nomor_sk_pembatalan": "SK-PEMBATALAN-001",
				"tanggal_sk_pembatalan": "2023-01-16"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK-001",
					"sk_tanggal":            time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(1),
					"masa_bulan":            int16(0),
					"no_pp":                 "PP-001",
					"no_sk_pembatalan":      "SK-PEMBATALAN-001",
					"tanggal_sk_pembatalan": time.Date(2023, 1, 16, 0, 0, 0, 0, time.UTC),
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Pegawai 2",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with minimal required data",
			paramNIP:      "1e",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-002",
				"tanggal_sk": "2023-02-15",
				"tanggal_mulai": "2023-02-20",
				"tanggal_akhir": "2023-05-20"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK-002",
					"sk_tanggal":            time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(3),
					"no_pp":                 "",
					"no_sk_pembatalan":      "",
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Pegawai 4",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1e",
					"pns_nip":               "1e",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with empty string optional fields",
			paramNIP:      "1f",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-003",
				"tanggal_sk": "2023-03-15",
				"tanggal_mulai": "2023-03-20",
				"tanggal_akhir": "2023-06-20",
				"nomor_pp": "",
				"nomor_sk_pembatalan": ""
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": { "id": {id} }
			}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    "{id}",
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK-003",
					"sk_tanggal":            time.Date(2023, 3, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2023, 3, 20, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2023, 6, 20, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(3),
					"no_pp":                 "",
					"no_sk_pembatalan":      "",
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Pegawai 5",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1f",
					"pns_nip":               "1f",
					"created_at":            "{created_at}",
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "error: pegawai is not found",
			paramNIP:      "1b",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-001",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: pegawai is deleted",
			paramNIP:      "1d",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-001",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data pegawai tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: golongan or jenis hukuman is not found",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 0,
				"golongan_id": 0,
				"nomor_sk": "SK-001",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis hukuman tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: golongan or jenis hukuman is deleted",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 2,
				"golongan_id": 2,
				"nomor_sk": "SK-001",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis hukuman tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: tanggal akhir sebelum tanggal mulai",
			paramNIP:      "1a",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-001",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2023-01-10"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "masa hukuman tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: missing required params",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"jenis_hukuman_id\" harus diisi` +
				` | parameter \"golongan_id\" harus diisi` +
				` | parameter \"nomor_sk\" harus diisi` +
				` | parameter \"tanggal_sk\" harus diisi` +
				` | parameter \"tanggal_mulai\" harus diisi` +
				` | parameter \"tanggal_akhir\" harus diisi"}`,
			wantDBRows: dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-hukuman-disiplin", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_hukuman_disiplin", "where pns_nip = $1", tt.paramNIP)
			require.NoError(t, err)
			if len(tt.wantDBRows) == len(actualRows) {
				for i, row := range actualRows {
					if tt.wantDBRows[i]["id"] == "{id}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRows[i]["id"] = row["id"]
						tt.wantDBRows[i]["created_at"] = row["created_at"]
						tt.wantDBRows[i]["updated_at"] = row["updated_at"]

						tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{id}", fmt.Sprintf("%d", row["id"]))
					}
				}
			}
			assert.Equal(t, tt.wantDBRows, actualRows)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai
			(pns_id,  nip_baru, deleted_at) VALUES
			('id_1c', '1c',     NULL),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     NULL);
		INSERT INTO ref_jenis_hukuman
			(id, nama, deleted_at) VALUES
			(1, 'Jenis Hukuman 1', NULL),
			(2, 'Jenis Hukuman 2', '2000-01-01');
		INSERT INTO ref_golongan
			(id, nama, nama_pangkat, nama_2, gol, gol_pppk, deleted_at) VALUES
			(1, 'I/a', 'Juru Muda', 'Ia', 1, 'I', NULL),
			(2, 'II/a', 'Pengatur Muda', 'IIa', 2, 'II', '2000-01-01');
		INSERT INTO riwayat_hukuman_disiplin
			(id, pns_id, pns_nip, nama, golongan_id, nama_golongan, jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal, tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman, no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, created_at, updated_at, deleted_at) VALUES
			(1, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK1', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL),
			(2, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK2', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL),
			(5, 'id_1e', '1e', 'Ani', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK5', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL),
			(6, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK6', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', '2000-01-01'),
			(7, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK7', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                    int64(7),
			"jenis_hukuman_id":      int16(1),
			"golongan_id":           int16(1),
			"sk_nomor":              "SK7",
			"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
			"masa_tahun":            int16(0),
			"masa_bulan":            int16(1),
			"no_pp":                 nil,
			"no_sk_pembatalan":      nil,
			"tanggal_sk_pembatalan": nil,
			"bkn_id":                nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"nama":                  "Budi",
			"nama_golongan":         "I/a",
			"nama_jenis_hukuman":    "Jenis Hukuman 1",
			"pns_id":                "id_1c",
			"pns_nip":               "1c",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:          "ok: with all data",
			paramNIP:      "1c",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-UPDATED",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20",
				"nomor_pp": "PP-UPDATED",
				"nomor_sk_pembatalan": "SK-PEMBATALAN-UPDATED",
				"tanggal_sk_pembatalan": "2023-01-16"
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(1),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK-UPDATED",
					"sk_tanggal":            time.Date(2023, 1, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2023, 1, 20, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2024, 1, 20, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(1),
					"masa_bulan":            int16(0),
					"no_pp":                 "PP-UPDATED",
					"no_sk_pembatalan":      "SK-PEMBATALAN-UPDATED",
					"tanggal_sk_pembatalan": time.Date(2023, 1, 16, 0, 0, 0, 0, time.UTC),
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "ok: with empty string optional fields",
			paramNIP:      "1c",
			paramID:       "2",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-UPDATED-2",
				"tanggal_sk": "2023-02-15",
				"tanggal_mulai": "2023-02-20",
				"tanggal_akhir": "2023-05-20",
				"nomor_pp": "",
				"nomor_sk_pembatalan": ""
			}`,
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(2),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK-UPDATED-2",
					"sk_tanggal":            time.Date(2023, 2, 15, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2023, 2, 20, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2023, 5, 20, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(3),
					"no_pp":                 "",
					"no_sk_pembatalan":      "",
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "error: riwayat hukuman disiplin is not found",
			paramNIP:      "1c",
			paramID:       "0",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-UPDATED",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:          "error: riwayat hukuman disiplin is owned by different pegawai",
			paramNIP:      "1c",
			paramID:       "5",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-UPDATED",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(5),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK5",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Ani",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1e",
					"pns_nip":               "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
		},
		{
			name:          "error: riwayat hukuman disiplin is deleted",
			paramNIP:      "1c",
			paramID:       "6",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-UPDATED",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(6),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK6",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:          "error: golongan or jenis hukuman is not found",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 0,
				"golongan_id": 0,
				"nomor_sk": "SK-UPDATED",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2024-01-20"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data golongan tidak ditemukan | data jenis hukuman tidak ditemukan"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:          "error: tanggal akhir sebelum tanggal mulai",
			paramNIP:      "1c",
			paramID:       "7",
			requestHeader: http.Header{"Authorization": authHeader},
			requestBody: `{
				"jenis_hukuman_id": 1,
				"golongan_id": 1,
				"nomor_sk": "SK-UPDATED",
				"tanggal_sk": "2023-01-15",
				"tanggal_mulai": "2023-01-20",
				"tanggal_akhir": "2023-01-10"
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "masa hukuman tidak valid"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: missing required params",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": 1 }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"jenis_hukuman_id\" harus diisi` +
				` | parameter \"golongan_id\" harus diisi` +
				` | parameter \"nomor_sk\" harus diisi` +
				` | parameter \"tanggal_sk\" harus diisi` +
				` | parameter \"tanggal_mulai\" harus diisi` +
				` | parameter \"tanggal_akhir\" harus diisi"}`,
			wantDBRows: defaultRows,
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBRows:       defaultRows,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "7",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-hukuman-disiplin/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_hukuman_disiplin", "where id = $1", tt.paramID)
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

func Test_handler_adminDelete(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman
			(id, nama, deleted_at) VALUES
			(1, 'Jenis Hukuman 1', NULL);
		INSERT INTO ref_golongan
			(id, nama, nama_pangkat, nama_2, gol, gol_pppk, deleted_at) VALUES
			(1, 'I/a', 'Juru Muda', 'Ia', 1, 'I', NULL);
		INSERT INTO pegawai
			(pns_id,  nip_baru, deleted_at) VALUES
			('id_1c', '1c',     NULL),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     NULL);
		INSERT INTO riwayat_hukuman_disiplin
			(id, pns_id, pns_nip, nama, golongan_id, nama_golongan, jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal, tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman, no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, created_at, updated_at, deleted_at) VALUES
			(1, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK1', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL),
			(2, 'id_1e', '1e', 'Ani', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK2', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL),
			(3, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK3', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', '2000-01-01'),
			(4, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK4', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                    int64(4),
			"jenis_hukuman_id":      int16(1),
			"golongan_id":           int16(1),
			"sk_nomor":              "SK4",
			"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
			"masa_tahun":            int16(0),
			"masa_bulan":            int16(1),
			"no_pp":                 nil,
			"no_sk_pembatalan":      nil,
			"tanggal_sk_pembatalan": nil,
			"bkn_id":                nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"nama":                  "Budi",
			"nama_golongan":         "I/a",
			"nama_jenis_hukuman":    "Jenis Hukuman 1",
			"pns_id":                "id_1c",
			"pns_nip":               "1c",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		wantDBRows       dbtest.Rows
	}{
		{
			name:             "ok: success delete",
			paramNIP:         "1c",
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(1),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK1",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            "{deleted_at}",
				},
			},
		},
		{
			name:             "error: riwayat hukuman disiplin is owned by other pegawai",
			paramNIP:         "1c",
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(2),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK2",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Ani",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1e",
					"pns_nip":               "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
		},
		{
			name:             "error: riwayat hukuman disiplin is not found",
			paramNIP:         "1c",
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows:       dbtest.Rows{},
		},
		{
			name:             "error: riwayat hukuman disiplin is deleted",
			paramNIP:         "1c",
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(3),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK3",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1c",
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:       defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-hukuman-disiplin/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_hukuman_disiplin", "where id = $1", tt.paramID)
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

func Test_handler_adminUploadBerkas(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO ref_jenis_hukuman
			(id, nama, deleted_at) VALUES
			(1, 'Jenis Hukuman 1', NULL);
		INSERT INTO ref_golongan
			(id, nama, nama_pangkat, nama_2, gol, gol_pppk, deleted_at) VALUES
			(1, 'I/a', 'Juru Muda', 'Ia', 1, 'I', NULL);
		INSERT INTO pegawai
			(pns_id,  nip_baru, deleted_at) VALUES
			('id_1c', '1c',     NULL),
			('id_1d', '1d',     '2000-01-01'),
			('id_1e', '1e',     NULL);
		INSERT INTO riwayat_hukuman_disiplin
			(id, pns_id, pns_nip, nama, golongan_id, nama_golongan, jenis_hukuman_id, nama_jenis_hukuman, sk_nomor, sk_tanggal, tanggal_mulai_hukuman, masa_tahun, masa_bulan, tanggal_akhir_hukuman, no_pp, no_sk_pembatalan, tanggal_sk_pembatalan, file_base64, keterangan_berkas, created_at, updated_at, deleted_at) VALUES
			(1, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK1', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, 'data:abc', 'abc', '2000-01-01', '2000-01-01', NULL),
			(2, 'id_1e', '1e', 'Ani', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK2', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL),
			(3, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK3', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, NULL, NULL, '2000-01-01', '2000-01-01', '2000-01-01'),
			(4, 'id_1c', '1c', 'Budi', 1, 'I/a', 1, 'Jenis Hukuman 1', 'SK4', '2000-01-01', '2000-01-10', 0, 1, '2000-02-10', NULL, NULL, NULL, NULL, NULL, '2000-01-01', '2000-01-01', NULL);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	defaultRows := dbtest.Rows{
		{
			"id":                    int64(4),
			"jenis_hukuman_id":      int16(1),
			"golongan_id":           int16(1),
			"sk_nomor":              "SK4",
			"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
			"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
			"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
			"masa_tahun":            int16(0),
			"masa_bulan":            int16(1),
			"no_pp":                 nil,
			"no_sk_pembatalan":      nil,
			"tanggal_sk_pembatalan": nil,
			"bkn_id":                nil,
			"file_base64":           nil,
			"s3_file_id":            nil,
			"keterangan_berkas":     nil,
			"nama":                  "Budi",
			"nama_golongan":         "I/a",
			"nama_jenis_hukuman":    "Jenis Hukuman 1",
			"pns_id":                "id_1c",
			"pns_nip":               "1c",
			"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
			"deleted_at":            nil,
		},
	}

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_Pegawai_Write)
	RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	defaultRequestBody := func(writer *multipart.Writer) error {
		part, err := writer.CreateFormFile("file", "file.txt")
		if err != nil {
			return err
		}
		_, err = io.WriteString(part, "Hello World!!")
		return err
	}

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name              string
		paramNIP          string
		paramID           string
		requestHeader     http.Header
		appendRequestBody func(writer *multipart.Writer) error
		wantResponseCode  int
		wantResponseBody  string
		wantDBRows        dbtest.Rows
	}{
		{
			name:              "ok: success upload",
			paramNIP:          "1c",
			paramID:           "1",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNoContent,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(1),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK1",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           "data:text/plain; charset=utf-8;base64,SGVsbG8gV29ybGQhIQ==",
					"s3_file_id":            nil,
					"keterangan_berkas":     "abc",
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            "{updated_at}",
					"deleted_at":            nil,
				},
			},
		},
		{
			name:              "error: riwayat hukuman disiplin is not found",
			paramNIP:          "1c",
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows:        dbtest.Rows{},
		},
		{
			name:              "error: riwayat hukuman disiplin is owned by different pegawai",
			paramNIP:          "1c",
			paramID:           "2",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(2),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK2",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Ani",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1e",
					"pns_nip":               "1e",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            nil,
				},
			},
		},
		{
			name:              "error: riwayat hukuman disiplin is deleted",
			paramNIP:          "1c",
			paramID:           "3",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: defaultRequestBody,
			wantResponseCode:  http.StatusNotFound,
			wantResponseBody:  `{"message": "data tidak ditemukan"}`,
			wantDBRows: dbtest.Rows{
				{
					"id":                    int64(3),
					"jenis_hukuman_id":      int16(1),
					"golongan_id":           int16(1),
					"sk_nomor":              "SK3",
					"sk_tanggal":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC),
					"tanggal_mulai_hukuman": time.Date(2000, 1, 10, 0, 0, 0, 0, time.UTC),
					"tanggal_akhir_hukuman": time.Date(2000, 2, 10, 0, 0, 0, 0, time.UTC),
					"masa_tahun":            int16(0),
					"masa_bulan":            int16(1),
					"no_pp":                 nil,
					"no_sk_pembatalan":      nil,
					"tanggal_sk_pembatalan": nil,
					"bkn_id":                nil,
					"file_base64":           nil,
					"s3_file_id":            nil,
					"keterangan_berkas":     nil,
					"nama":                  "Budi",
					"nama_golongan":         "I/a",
					"nama_jenis_hukuman":    "Jenis Hukuman 1",
					"pns_id":                "id_1c",
					"pns_nip":               "1c",
					"created_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":            time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name:              "error: missing file",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": authHeader},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusBadRequest,
			wantResponseBody:  `{"message": "parameter \"file\" harus diisi"}`,
			wantDBRows:        defaultRows,
		},
		{
			name:              "error: invalid token",
			paramNIP:          "1c",
			paramID:           "4",
			requestHeader:     http.Header{"Authorization": []string{"Bearer some-token"}},
			appendRequestBody: func(*multipart.Writer) error { return nil },
			wantResponseCode:  http.StatusUnauthorized,
			wantResponseBody:  `{"message": "token otentikasi tidak valid"}`,
			wantDBRows:        defaultRows,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)
			require.NoError(t, tt.appendRequestBody(writer))
			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pegawai/"+tt.paramNIP+"/riwayat-hukuman-disiplin/"+tt.paramID+"/berkas", &buf)
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", writer.FormDataContentType())
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRows, err := dbtest.QueryWithClause(db, "riwayat_hukuman_disiplin", "where id = $1", tt.paramID)
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
