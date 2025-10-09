package unitkerja

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func Test_handler_List(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
		("id", "nama_unor", "unor_induk") values
		('001', 'Sekretariat Jenderal', ''),
		('002', 'Direktorat Jenderal Pendidikan Dasar dan Menengah', '001'),
		('003', 'Direktorat Jenderal Pendidikan Tinggi', '001'),
		('004', 'Direktorat Jenderal Guru dan Tenaga Kependidikan', '001'),
		('005', 'Direktorat Jenderal PAUD dan Dikmas', '001'),
		('006', 'Inspektorat Jenderal', '001'),
		('007', 'Badan Pengembangan dan Pembinaan Bahasa', '001'),
		('008', 'Badan Penelitian dan Pengembangan', '001'),
		('009', 'Sekretariat Direktorat Jenderal Pendidikan Dasar dan Menengah', '002'),
		('010', 'Direktorat Sekolah Dasar', '002'),
		('011', 'Direktorat Sekolah Menengah Pertama', '002'),
		('012', 'Direktorat Sekolah Menengah Atas', '002'),
		('013', 'Direktorat Sekolah Menengah Kejuruan', '002'),
		('014', 'Sekretariat Direktorat Jenderal Pendidikan Tinggi', '003'),
		('015', 'Direktorat Pembelajaran dan Kemahasiswaan', '003'),
		('016', 'Direktorat Kelembagaan dan Kerjasama', '003'),
		('017', 'Direktorat Riset dan Pengabdian Masyarakat', '003'),
		('018', 'Sekretariat Direktorat Jenderal Guru dan Tenaga Kependidikan', '004'),
		('019', 'Direktorat Pendidikan Profesi dan Pembinaan Guru', '004'),
		('020', 'Direktorat Pembinaan Tenaga Kependidikan', '004');
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
			name:             "ok: get data with default pagination",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "001", "nama": "Sekretariat Jenderal"},
					{"id": "002", "nama": "Direktorat Jenderal Pendidikan Dasar dan Menengah"},
					{"id": "003", "nama": "Direktorat Jenderal Pendidikan Tinggi"},
					{"id": "004", "nama": "Direktorat Jenderal Guru dan Tenaga Kependidikan"},
					{"id": "005", "nama": "Direktorat Jenderal PAUD dan Dikmas"},
					{"id": "006", "nama": "Inspektorat Jenderal"},
					{"id": "007", "nama": "Badan Pengembangan dan Pembinaan Bahasa"},
					{"id": "008", "nama": "Badan Penelitian dan Pengembangan"},
					{"id": "009", "nama": "Sekretariat Direktorat Jenderal Pendidikan Dasar dan Menengah"},
					{"id": "010", "nama": "Direktorat Sekolah Dasar"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 20
				}
			}`,
		},
		{
			name:   "ok: with pagination limit 5",
			dbData: dbData,
			requestQuery: url.Values{
				"limit": []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "001", "nama": "Sekretariat Jenderal"},
					{"id": "002", "nama": "Direktorat Jenderal Pendidikan Dasar dan Menengah"},
					{"id": "003", "nama": "Direktorat Jenderal Pendidikan Tinggi"},
					{"id": "004", "nama": "Direktorat Jenderal Guru dan Tenaga Kependidikan"},
					{"id": "005", "nama": "Direktorat Jenderal PAUD dan Dikmas"}
				],
				"meta": {
					"limit": 5,
					"offset": 0,
					"total": 20
				}
			}`,
		},
		{
			name:   "ok: with pagination limit 3 offset 5",
			dbData: dbData,
			requestQuery: url.Values{
				"limit":  []string{"3"},
				"offset": []string{"5"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "006", "nama": "Inspektorat Jenderal"},
					{"id": "007", "nama": "Badan Pengembangan dan Pembinaan Bahasa"},
					{"id": "008", "nama": "Badan Penelitian dan Pengembangan"}
				],
				"meta": {
					"limit": 3,
					"offset": 5,
					"total": 20
				}
			}`,
		},
		{
			name:   "ok: with search parameter",
			dbData: dbData,
			requestQuery: url.Values{
				"nama": []string{"Direktorat"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "002", "nama": "Direktorat Jenderal Pendidikan Dasar dan Menengah"},
					{"id": "003", "nama": "Direktorat Jenderal Pendidikan Tinggi"},
					{"id": "004", "nama": "Direktorat Jenderal Guru dan Tenaga Kependidikan"},
					{"id": "005", "nama": "Direktorat Jenderal PAUD dan Dikmas"},
					{"id": "010", "nama": "Direktorat Sekolah Dasar"},
					{"id": "011", "nama": "Direktorat Sekolah Menengah Pertama"},
					{"id": "012", "nama": "Direktorat Sekolah Menengah Atas"},
					{"id": "013", "nama": "Direktorat Sekolah Menengah Kejuruan"},
					{"id": "015", "nama": "Direktorat Pembelajaran dan Kemahasiswaan"},
					{"id": "016", "nama": "Direktorat Kelembagaan dan Kerjasama"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 13
				}
			}`,
		},
		{
			name:   "ok: with unor_induk filter",
			dbData: dbData,
			requestQuery: url.Values{
				"unor_induk": []string{"001"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "002", "nama": "Direktorat Jenderal Pendidikan Dasar dan Menengah"},
					{"id": "003", "nama": "Direktorat Jenderal Pendidikan Tinggi"},
					{"id": "004", "nama": "Direktorat Jenderal Guru dan Tenaga Kependidikan"},
					{"id": "005", "nama": "Direktorat Jenderal PAUD dan Dikmas"},
					{"id": "006", "nama": "Inspektorat Jenderal"},
					{"id": "007", "nama": "Badan Pengembangan dan Pembinaan Bahasa"},
					{"id": "008", "nama": "Badan Penelitian dan Pengembangan"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 7
				}
			}`,
		},
		{
			name:   "ok: with search and unor_induk filter",
			dbData: dbData,
			requestQuery: url.Values{
				"nama":       []string{"Direktorat"},
				"unor_induk": []string{"002"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "010", "nama": "Direktorat Sekolah Dasar"},
					{"id": "011", "nama": "Direktorat Sekolah Menengah Pertama"},
					{"id": "012", "nama": "Direktorat Sekolah Menengah Atas"},
					{"id": "013", "nama": "Direktorat Sekolah Menengah Kejuruan"}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 4
				}
			}`,
		},
		{
			name:             "ok: empty data",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
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
		{
			name:             "error: missing auth header",
			dbData:           dbData,
			requestHeader:    http.Header{},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/unit-kerja", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			r := repo.New(pgxconn)
			RegisterRoutes(e, r, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_ListAkar(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
		("id", "nama_unor", "order", "diatasan_id", deleted_at) values
		('001', 'Tingkat 1', 8, NULL, NULL),
		('002', 'Tingkat 2', 3, '001', NULL),
		('003', 'Tingkat 2 Kedua', 2, '001', NULL),
		('004', 'Tingkat 2 Deleted', 4, '001', NOW()),
		('005', 'Tingkat 1 Deleted', 5, NULL, NOW()),
		('006', 'Tingkat 3', 6, '002', NULL),
		('007', 'Tingkat 3 Deleted', 7, '002', NOW()),
		('008', 'Tingkat 1 Kedua', 1, NULL, NULL);
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
			name:             "ok: get data with default pagination",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "008", "nama": "Tingkat 1 Kedua", "has_anak": false},
					{"id": "001", "nama": "Tingkat 1", "has_anak": true}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 2
				}
			}`,
		},
		{
			name:   "ok: with pagination limit and offset",
			dbData: dbData,
			requestQuery: url.Values{
				"limit":  []string{"1"},
				"offset": []string{"1"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "001", "nama": "Tingkat 1", "has_anak": true}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 2
				}
			}`,
		},
		{
			name:             "ok: empty data",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
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
		{
			name:             "error: missing auth header",
			dbData:           dbData,
			requestHeader:    http.Header{},
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

			req := httptest.NewRequest(http.MethodGet, "/v1/unit-kerja/akar", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := repo.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_ListAnak(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into unit_kerja
		("id", "nama_unor", "order", "diatasan_id", deleted_at) values
		('001', 'Tingkat 1', 8, NULL, NULL),
		('002', 'Tingkat 2', 3, '001', NULL),
		('003', 'Tingkat 2 Kedua', 2, '001', NULL),
		('004', 'Tingkat 2 Deleted', 4, '001', NOW()),
		('005', 'Tingkat 1 Deleted', 5, NULL, NOW()),
		('006', 'Tingkat 3', 6, '002', NULL),
		('007', 'Tingkat 3 Deleted', 7, '002', NOW()),
		('008', 'Tingkat 1 Kedua', 1, NULL, NULL);
	`

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get data with default pagination",
			dbData:           dbData,
			id:               "001",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "003", "nama": "Tingkat 2 Kedua", "has_anak": false},
					{"id": "002", "nama": "Tingkat 2", "has_anak": true}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 2
				}
			}`,
		},
		{
			name:             "ok: get data with another id",
			dbData:           dbData,
			id:               "002",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "006", "nama": "Tingkat 3", "has_anak": false}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
		},
		{
			name:   "ok: with pagination limit and offset",
			dbData: dbData,
			id:     "001",
			requestQuery: url.Values{
				"limit":  []string{"1"},
				"offset": []string{"1"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": "002", "nama": "Tingkat 2", "has_anak": true}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 2
				}
			}`,
		},
		{
			name:             "ok: leaf data",
			dbData:           dbData,
			id:               "003",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "ok: deleted data",
			dbData:           dbData,
			id:               "004",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "ok: empty data",
			id:               "001",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			id:               "001",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: missing auth header",
			dbData:           dbData,
			id:               "001",
			requestHeader:    http.Header{},
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/unit-kerja/%s/anak", tt.id), nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := repo.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGet(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai (
				pns_id, nama, nip_baru, jabatan_nama, instansi_kerja_nama
		) VALUES
		(
			'PNSROOT', 'Andi Prasetyo', '198501012020031001', 'Kepala Kantor', 'Kantor Pusat'
		),
		(
			'PNS001', 'Budi Santoso', '198701012020031002', 'Kepala Bagian Keuangan', 'Direktorat Keuangan'
		),
		(
			'PNS002', 'Siti Aminah', '198901012020031003', 'Kepala Bagian SDM', 'Direktorat SDM'
		);


		INSERT INTO ref_instansi (
				id, nama, created_at, updated_at, deleted_at
		) VALUES
		(
			'INSTROOT', 'Kementerian Contoh', now(), now(), NULL
		),
		(
			'INST001', 'Direktorat Keuangan', now(), now(), NULL
		),
		(
			'INST002', 'Direktorat SDM', now(), now(), NULL
		);

		INSERT INTO unit_kerja (
				id, "no", kode_internal, nama_unor, eselon_id, cepat_kode, nama_jabatan, nama_pejabat,
				diatasan_id, instansi_id, pemimpin_pns_id, jenis_unor_id, unor_induk, jumlah_ideal_staff,
				"order", is_satker, eselon_1, eselon_2, eselon_3, eselon_4, expired_date, keterangan,
				jenis_satker, abbreviation, unor_induk_penyetaraan, jabatan_id, waktu, peraturan, remark,
				aktif, eselon_nama, deleted_at
		) VALUES 
		(
			'00000000-0000-0000-0000-000000000001', 0, 'ROOT001', 'Kantor Pusat', 'E0', 'CKROOT', 'Kepala Kantor', 'Andi Prasetyo',
			NULL, 'INSTROOT', 'PNSROOT', 'JUROOT', NULL, 20,
			0, true, NULL, NULL, NULL, NULL, '2035-12-31', 'Unit induk pusat',
			'Satker', 'KP', '', 'JABROOT', '2025', 'Peraturan ROOT 001', 'Unit pusat', true, 'Eselon ROOT', NULL
		),
		(
			'11111111-1111-1111-1111-111111111111', 1, 'KI001', 'Bagian Keuangan', 'E1', 'CK01', 'Kepala Bagian', 'Budi Santoso',
			'00000000-0000-0000-0000-000000000001', 'INST001', 'PNS001', 'JU001', '00000000-0000-0000-0000-000000000001', 5,
			1, true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '2030-12-31', 'Bagian pengelola keuangan',
			'Satker', 'BK', 'PENY001', 'JAB001', '2025', 'Peraturan 123', 'Remark contoh', true, 'Eselon Nama I',
			NULL
		),
		(
			'22222222-2222-2222-2222-222222222222', 2, 'KI002', 'Bagian SDM', 'E2', 'CK02', 'Kepala Bagian', 'Siti Aminah',
			'11111111-1111-1111-1111-111111111111', 'INST002', 'PNS002', 'JU002', '11111111-1111-1111-1111-111111111111', 8,
			2, false, '00000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', '2031-06-30', 'Bagian pengelola SDM',
			'Satker', 'SDM', 'PENY002', 'JAB002', '2026', 'Peraturan 456', 'Remark contoh 2', false, 'Eselon Nama II',
			'now()'
		);
	`

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get unit kerja",
			dbData:           dbData,
			id:               "00000000-0000-0000-0000-000000000001",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": "00000000-0000-0000-0000-000000000001",
					"no": 0,
					"kode_internal": "ROOT001",
					"nama": "Kantor Pusat",
					"eselon_id": "E0",
					"cepat_kode": "CKROOT",
					"nama_jabatan": "Kepala Kantor",
					"nama_pejabat": "Andi Prasetyo",
					"diatasan_id": "",
					"instansi_id": "INSTROOT",
					"pemimpin_pns_id": "PNSROOT",
					"jenis_unor_id": "JUROOT",
					"unor_induk": "",
					"jumlah_ideal_staff": 20,
					"order": 0,
					"is_satker": true,
					"eselon_1": "",
					"eselon_2": "",
					"eselon_3": "",
					"eselon_4": "",
					"expired_date": "2035-12-31",
					"keterangan": "Unit induk pusat",
					"jenis_satker": "Satker",
					"abbreviation": "KP",
					"unor_induk_penyetaraan": "",
					"jabatan_id": "JABROOT",
					"waktu": "2025",
					"peraturan": "Peraturan ROOT 001",
					"remark": "Unit pusat",
					"aktif": true,
					"eselon_nama": "Eselon ROOT",
					"nama_diatasan": "",
					"nama_unor_induk": ""
				}
			}`,
		},
		{
			name:             "ok: get unit kerja with another id",
			dbData:           dbData,
			id:               "11111111-1111-1111-1111-111111111111",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": "11111111-1111-1111-1111-111111111111",
					"no": 1,
					"kode_internal": "KI001",
					"nama": "Bagian Keuangan",
					"eselon_id": "E1",
					"cepat_kode": "CK01",
					"nama_jabatan": "Kepala Bagian",
					"nama_pejabat": "Budi Santoso",
					"diatasan_id": "00000000-0000-0000-0000-000000000001",
					"instansi_id": "INST001",
					"pemimpin_pns_id": "PNS001",
					"jenis_unor_id": "JU001",
					"unor_induk": "00000000-0000-0000-0000-000000000001",
					"jumlah_ideal_staff": 5,
					"order": 1,
					"is_satker": true,
					"eselon_1": "00000000-0000-0000-0000-000000000001",
					"eselon_2": "00000000-0000-0000-0000-000000000001",
					"eselon_3": "00000000-0000-0000-0000-000000000001",
					"eselon_4": "00000000-0000-0000-0000-000000000001",
					"expired_date": "2030-12-31",
					"keterangan": "Bagian pengelola keuangan",
					"jenis_satker": "Satker",
					"abbreviation": "BK",
					"unor_induk_penyetaraan": "PENY001",
					"jabatan_id": "JAB001",
					"waktu": "2025",
					"peraturan": "Peraturan 123",
					"remark": "Remark contoh",
					"aktif": true,
					"eselon_nama": "Eselon Nama I",
					"nama_unor_induk": "Kantor Pusat",
					"nama_diatasan": "Kantor Pusat"
				}
			}`,
		},
		{
			name:             "error: unit kerja not found",
			dbData:           dbData,
			id:               "11111111-1111-1111-1111-22222222222",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: unit kerja deleted",
			dbData:           dbData,
			id:               "22222222-2222-2222-2222-222222222222",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "111", api.RoleAdmin)}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: user is not an admin",
			dbData:           dbData,
			id:               "11111111-1111-1111-1111-111111111111",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			id:               "11111111-1111-1111-1111-111111111111",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/unit-kerja/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			repo := repo.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai (
				pns_id, nama, nip_baru, jabatan_nama, instansi_kerja_nama
		) VALUES
		(
			'PNSROOT', 'Andi Prasetyo', '198501012020031001', 'Kepala Kantor', 'Kantor Pusat'
		),
		(
			'PNS001', 'Budi Santoso', '198701012020031002', 'Kepala Bagian Keuangan', 'Direktorat Keuangan'
		),
		(
			'PNS002', 'Siti Aminah', '198901012020031003', 'Kepala Bagian SDM', 'Direktorat SDM'
		);


		INSERT INTO ref_instansi (
				id, nama, created_at, updated_at, deleted_at
		) VALUES
		(
			'INSTROOT', 'Kementerian Contoh', now(), now(), NULL
		),
		(
			'INST001', 'Direktorat Keuangan', now(), now(), NULL
		),
		(
			'INST002', 'Direktorat SDM', now(), now(), NULL
		);

		INSERT INTO unit_kerja (
				id, "no", kode_internal, nama_unor, eselon_id, cepat_kode, nama_jabatan, nama_pejabat,
				diatasan_id, instansi_id, pemimpin_pns_id, jenis_unor_id, unor_induk, jumlah_ideal_staff,
				"order", is_satker, eselon_1, eselon_2, eselon_3, eselon_4, expired_date, keterangan,
				jenis_satker, abbreviation, unor_induk_penyetaraan, jabatan_id, waktu, peraturan, remark,
				aktif, eselon_nama, deleted_at
		) VALUES 
		(
			'00000000-0000-0000-0000-000000000001', 0, 'ROOT001', 'Kantor Pusat', 'E0', 'CKROOT', 'Kepala Kantor', 'Andi Prasetyo',
			NULL, 'INSTROOT', 'PNSROOT', 'JUROOT', NULL, 20,
			0, true, NULL, NULL, NULL, NULL, '2035-12-31', 'Unit induk pusat',
			'Satker', 'KP', '', 'JABROOT', '2025', 'Peraturan ROOT 001', 'Unit pusat', true, 'Eselon ROOT', NULL
		);
	`

	tests := []struct {
		name             string
		dbData           string
		requestQuery     url.Values
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:   "ok: create unit kerja only with required value",
			dbData: dbData,
			requestBody: `{
				"id": "11111111-1111-1111-1111-111111111111",
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": "11111111-1111-1111-1111-111111111111",
					"no": 0,
					"kode_internal": "",
					"nama": "Bagian Keuangan",
					"eselon_id": "",
					"cepat_kode": "",
					"nama_jabatan": "",
					"nama_pejabat": "",
					"diatasan_id": "",
					"instansi_id": "",
					"pemimpin_pns_id": "",
					"jenis_unor_id": "",
					"unor_induk": "",
					"jumlah_ideal_staff": 0,
					"order": 0,
					"is_satker": true,
					"eselon_1": "",
					"eselon_2": "",
					"eselon_3": "",
					"eselon_4": "",
					"expired_date": null,
					"keterangan": "",
					"jenis_satker": "",
					"abbreviation": "",
					"unor_induk_penyetaraan": "",
					"jabatan_id": "",
					"waktu": "",
					"peraturan": "",
					"remark": "",
					"aktif": false,
					"eselon_nama": ""
				}
			}`,
		},
		{
			name:   "ok: create unit kerja",
			dbData: dbData,
			requestBody: `{
				"diatasan_id": "00000000-0000-0000-0000-000000000001",
				"id": "22222222-2222-2222-2222-222222222222",
				"nama": "Bagian Keuangan",
				"kode_internal": "KI001",
				"nama_jabatan": "Kepala Bagian",
				"pemimpin_pns_id": "PNS001",
				"is_satker": true,
				"unor_induk": "00000000-0000-0000-0000-000000000001",
				"expired_date": "2030-12-31",
				"keterangan": "Mengelola anggaran dan administrasi keuangan",
				"abbreviation": "BK",
				"waktu": "2025",
				"jenis_satker": "Satker",
				"peraturan": "Peraturan 123"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": "22222222-2222-2222-2222-222222222222",
					"no": 0,
					"kode_internal": "KI001",
					"nama": "Bagian Keuangan",
					"eselon_id": "",
					"cepat_kode": "",
					"nama_jabatan": "Kepala Bagian",
					"nama_pejabat": "Budi Santoso",
					"diatasan_id": "00000000-0000-0000-0000-000000000001",
					"instansi_id": "",
					"pemimpin_pns_id": "PNS001",
					"jenis_unor_id": "",
					"unor_induk": "00000000-0000-0000-0000-000000000001",
					"jumlah_ideal_staff": 0,
					"order": 0,
					"is_satker": true,
					"eselon_1": "",
					"eselon_2": "",
					"eselon_3": "",
					"eselon_4": "",
					"expired_date": "2030-12-31",
					"keterangan": "Mengelola anggaran dan administrasi keuangan",
					"jenis_satker": "Satker",
					"abbreviation": "BK",
					"unor_induk_penyetaraan": "",
					"jabatan_id": "",
					"waktu": "2025",
					"peraturan": "Peraturan 123",
					"remark": "",
					"aktif": false,
					"eselon_nama": ""
				}
			}`,
		},
		{
			name:   "ok: create unit kerja only with required value and expired date is null",
			dbData: dbData,
			requestBody: `{
				"id": "11111111-1111-1111-1111-111111111111",
				"nama": "Bagian Keuangan",
				"is_satker": true,
				"expired_date": null
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": "11111111-1111-1111-1111-111111111111",
					"no": 0,
					"kode_internal": "",
					"nama": "Bagian Keuangan",
					"eselon_id": "",
					"cepat_kode": "",
					"nama_jabatan": "",
					"nama_pejabat": "",
					"diatasan_id": "",
					"instansi_id": "",
					"pemimpin_pns_id": "",
					"jenis_unor_id": "",
					"unor_induk": "",
					"jumlah_ideal_staff": 0,
					"order": 0,
					"is_satker": true,
					"eselon_1": "",
					"eselon_2": "",
					"eselon_3": "",
					"eselon_4": "",
					"expired_date": null,
					"keterangan": "",
					"jenis_satker": "",
					"abbreviation": "",
					"unor_induk_penyetaraan": "",
					"jabatan_id": "",
					"waktu": "",
					"peraturan": "",
					"remark": "",
					"aktif": false,
					"eselon_nama": ""
				}
			}`,
		},
		{
			name:   "error: create unit kerja with existed id",
			dbData: dbData,
			requestBody: `{
				"id": "00000000-0000-0000-0000-000000000001",
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusConflict,
			wantResponseBody: `{"message": "Data dengan ID ini sudah terdaftar"}`,
		},
		{
			name:   "error: auth header tidak valid",
			dbData: dbData,
			requestBody: `{
				"id": "11111111-1111-1111-1111-111111111111",
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:   "error: user is not an admin",
			dbData: dbData,
			requestBody: `{
				"id": "11111111-1111-1111-1111-111111111111",
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)

			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/unit-kerja", strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			r := repo.New(pgxconn)
			RegisterRoutes(e, r, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdate(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai (
				pns_id, nama, nip_baru, jabatan_nama, instansi_kerja_nama
		) VALUES
		(
			'PNSROOT', 'Andi Prasetyo', '198501012020031001', 'Kepala Kantor', 'Kantor Pusat'
		),
		(
			'PNS001', 'Budi Santoso', '198701012020031002', 'Kepala Bagian Keuangan', 'Direktorat Keuangan'
		),
		(
			'PNS002', 'Siti Aminah', '198901012020031003', 'Kepala Bagian SDM', 'Direktorat SDM'
		);


		INSERT INTO ref_instansi (
				id, nama, created_at, updated_at, deleted_at
		) VALUES
		(
			'INSTROOT', 'Kementerian Contoh', now(), now(), NULL
		),
		(
			'INST001', 'Direktorat Keuangan', now(), now(), NULL
		),
		(
			'INST002', 'Direktorat SDM', now(), now(), NULL
		);

		INSERT INTO unit_kerja (
			id, "no", kode_internal, nama_unor, eselon_id, cepat_kode, nama_jabatan, nama_pejabat,
			diatasan_id, instansi_id, pemimpin_pns_id, jenis_unor_id, unor_induk, jumlah_ideal_staff,
			"order", is_satker, eselon_1, eselon_2, eselon_3, eselon_4, expired_date, keterangan,
			jenis_satker, abbreviation, unor_induk_penyetaraan, jabatan_id, waktu, peraturan, remark,
			aktif, eselon_nama, deleted_at
		) VALUES 
		(
			'00000000-0000-0000-0000-000000000001', 0, 'ROOT001', 'Kantor Pusat', 'E0', 'CKROOT', 'Kepala Kantor', 'Andi Prasetyo',
			NULL, 'INSTROOT', 'PNSROOT', 'JUROOT', NULL, 20,
			0, true, NULL, NULL, NULL, NULL, '2035-12-31', 'Unit induk pusat',
			'Satker', 'KP', '', 'JABROOT', '2025', 'Peraturan ROOT 001', 'Unit pusat', true, 'Eselon ROOT', NULL
		),
		(
			'11111111-1111-1111-1111-111111111111', 1, 'KI001', 'Bagian Keuangan', 'E1', 'CK01', 'Kepala Bagian', 'Budi Santoso',
			'00000000-0000-0000-0000-000000000001', 'INST001', 'PNS001', 'JU001', '00000000-0000-0000-0000-000000000001', 5,
			1, true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '2030-12-31', 'Bagian pengelola keuangan',
			'Satker', 'BK', 'PENY001', 'JAB001', '2025', 'Peraturan 123', 'Remark contoh', true, 'Eselon Nama I',
			NULL
		),
		(
			'22222222-2222-2222-2222-222222222222', 2, 'KI002', 'Bagian SDM', 'E2', 'CK02', 'Kepala Bagian', 'Siti Aminah',
			'11111111-1111-1111-1111-111111111111', 'INST002', 'PNS002', 'JU002', '11111111-1111-1111-1111-111111111111', 8,
			2, false, '00000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', '2031-06-30', 'Bagian pengelola SDM',
			'Satker', 'SDM', 'PENY002', 'JAB002', '2026', 'Peraturan 456', 'Remark contoh 2', false, 'Eselon Nama II',
			now()
		),
		(
			'33333333-3333-3333-3333-333333333333', 3, 'KI003', 'Bagian IT', 'E3', 'CK03', 'Kepala Bagian', 'Rina Suryani',
			'22222222-2222-2222-2222-222222222222', 'INST002', 'PNS001', 'JU003', '22222222-2222-2222-2222-222222222222', 10,
			3, true, '00000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', '2032-03-31', 'Bagian pengelola IT',
			'Satker', 'IT', 'PENY003', 'JAB003', '2027', 'Peraturan 789', 'Remark IT', true, 'Eselon Nama III',
			NULL
		);
	`

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestQuery     url.Values
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:   "ok: update unit kerja only with required value",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestBody: `{
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": "11111111-1111-1111-1111-111111111111",
					"no": 1,
					"kode_internal": "",
					"nama": "Bagian Keuangan",
					"eselon_id": "E1",
					"cepat_kode": "CK01",
					"nama_jabatan": "",
					"nama_pejabat": "",
					"diatasan_id": "",
					"instansi_id": "INST001",
					"pemimpin_pns_id": "",
					"jenis_unor_id": "JU001",
					"unor_induk": "",
					"jumlah_ideal_staff": 5,
					"order": 1,
					"is_satker": true,
					"eselon_1": "00000000-0000-0000-0000-000000000001",
					"eselon_2": "00000000-0000-0000-0000-000000000001",
					"eselon_3": "00000000-0000-0000-0000-000000000001",
					"eselon_4": "00000000-0000-0000-0000-000000000001",
					"expired_date": null,
					"keterangan": "",
					"jenis_satker": "",
					"abbreviation": "",
					"unor_induk_penyetaraan": "PENY001",
					"jabatan_id": "JAB001",
					"waktu": "",
					"peraturan": "",
					"remark": "Remark contoh",
					"aktif": true,
					"eselon_nama": "Eselon Nama I"
				}
			}`,
		},
		{
			name:   "ok: update unit kerja",
			dbData: dbData,
			id:     "33333333-3333-3333-3333-333333333333",
			requestBody: `{
				"diatasan_id": "11111111-1111-1111-1111-111111111111",
				"nama": "Bagian Legal",
				"kode_internal": "KL001",
				"nama_jabatan": "Kepala Legal",
				"pemimpin_pns_id": "PNS002",
				"is_satker": false,
				"unor_induk": "22222222-2222-2222-2222-222222222222",
				"expired_date": "2035-05-31",
				"keterangan": "Bagian pengelola legal dan regulasi",
				"abbreviation": "BL",
				"waktu": "2030",
				"jenis_satker": "Non-Satker",
				"peraturan": "Peraturan Legal 001"
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
					"data": {
						"id": "33333333-3333-3333-3333-333333333333",
						"no": 3,
						"kode_internal": "KL001",
						"nama": "Bagian Legal",
						"eselon_id": "E3",
						"cepat_kode": "CK03",
						"nama_jabatan": "Kepala Legal",
						"nama_pejabat": "Siti Aminah",
						"diatasan_id": "11111111-1111-1111-1111-111111111111",
						"instansi_id": "INST002",
						"pemimpin_pns_id": "PNS002",
						"jenis_unor_id": "JU003",
						"unor_induk": "22222222-2222-2222-2222-222222222222",
						"jumlah_ideal_staff": 10,
						"order": 3,
						"is_satker": false,
						"eselon_1": "00000000-0000-0000-0000-000000000001",
						"eselon_2": "22222222-2222-2222-2222-222222222222",
						"eselon_3": "22222222-2222-2222-2222-222222222222",
						"eselon_4": "22222222-2222-2222-2222-222222222222",
						"expired_date": "2035-05-31",
						"keterangan": "Bagian pengelola legal dan regulasi",
						"jenis_satker": "Non-Satker",
						"abbreviation": "BL",
						"unor_induk_penyetaraan": "PENY003",
						"jabatan_id": "JAB003",
						"waktu": "2030",
						"peraturan": "Peraturan Legal 001",
						"remark": "Remark IT",
						"aktif": true,
						"eselon_nama": "Eselon Nama III"
					}
				}`,
		},
		{
			name:   "ok: update unit kerja only with required value and expired date is null",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestBody: `{
				"nama": "Bagian Keuangan",
				"is_satker": true,
				"expired_date": null
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": "11111111-1111-1111-1111-111111111111",
					"no": 1,
					"kode_internal": "",
					"nama": "Bagian Keuangan",
					"eselon_id": "E1",
					"cepat_kode": "CK01",
					"nama_jabatan": "",
					"nama_pejabat": "",
					"diatasan_id": "",
					"instansi_id": "INST001",
					"pemimpin_pns_id": "",
					"jenis_unor_id": "JU001",
					"unor_induk": "",
					"jumlah_ideal_staff": 5,
					"order": 1,
					"is_satker": true,
					"eselon_1": "00000000-0000-0000-0000-000000000001",
					"eselon_2": "00000000-0000-0000-0000-000000000001",
					"eselon_3": "00000000-0000-0000-0000-000000000001",
					"eselon_4": "00000000-0000-0000-0000-000000000001",
					"expired_date": null,
					"keterangan": "",
					"jenis_satker": "",
					"abbreviation": "",
					"unor_induk_penyetaraan": "PENY001",
					"jabatan_id": "JAB001",
					"waktu": "",
					"peraturan": "",
					"remark": "Remark contoh",
					"aktif": true,
					"eselon_nama": "Eselon Nama I"
				}
			}`,
		},
		{
			name:   "error: update not found unit kerja",
			dbData: dbData,
			id:     "1234",
			requestBody: `{
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: update deleted unit kerja",
			dbData: dbData,
			id:     "22222222-2222-2222-2222-222222222222",
			requestBody: `{
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: auth header tidak valid",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestBody: `{
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:   "error: user is not an admin",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestBody: `{
				"nama": "Bagian Keuangan",
				"is_satker": true
			}`,
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)

			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/unit-kerja/"+tt.id, strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			r := repo.New(pgxconn)
			RegisterRoutes(e, r, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDelete(t *testing.T) {
	t.Parallel()

	dbData := `
		INSERT INTO pegawai (
				pns_id, nama, nip_baru, jabatan_nama, instansi_kerja_nama
		) VALUES
		(
			'PNSROOT', 'Andi Prasetyo', '198501012020031001', 'Kepala Kantor', 'Kantor Pusat'
		),
		(
			'PNS001', 'Budi Santoso', '198701012020031002', 'Kepala Bagian Keuangan', 'Direktorat Keuangan'
		),
		(
			'PNS002', 'Siti Aminah', '198901012020031003', 'Kepala Bagian SDM', 'Direktorat SDM'
		);


		INSERT INTO ref_instansi (
				id, nama, created_at, updated_at, deleted_at
		) VALUES
		(
			'INSTROOT', 'Kementerian Contoh', now(), now(), NULL
		),
		(
			'INST001', 'Direktorat Keuangan', now(), now(), NULL
		),
		(
			'INST002', 'Direktorat SDM', now(), now(), NULL
		);

		INSERT INTO unit_kerja (
			id, "no", kode_internal, nama_unor, eselon_id, cepat_kode, nama_jabatan, nama_pejabat,
			diatasan_id, instansi_id, pemimpin_pns_id, jenis_unor_id, unor_induk, jumlah_ideal_staff,
			"order", is_satker, eselon_1, eselon_2, eselon_3, eselon_4, expired_date, keterangan,
			jenis_satker, abbreviation, unor_induk_penyetaraan, jabatan_id, waktu, peraturan, remark,
			aktif, eselon_nama, deleted_at
		) VALUES 
		(
			'00000000-0000-0000-0000-000000000001', 0, 'ROOT001', 'Kantor Pusat', 'E0', 'CKROOT', 'Kepala Kantor', 'Andi Prasetyo',
			NULL, 'INSTROOT', 'PNSROOT', 'JUROOT', NULL, 20,
			0, true, NULL, NULL, NULL, NULL, '2035-12-31', 'Unit induk pusat',
			'Satker', 'KP', '', 'JABROOT', '2025', 'Peraturan ROOT 001', 'Unit pusat', true, 'Eselon ROOT', NULL
		),
		(
			'11111111-1111-1111-1111-111111111111', 1, 'KI001', 'Bagian Keuangan', 'E1', 'CK01', 'Kepala Bagian', 'Budi Santoso',
			'00000000-0000-0000-0000-000000000001', 'INST001', 'PNS001', 'JU001', '00000000-0000-0000-0000-000000000001', 5,
			1, true, '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '00000000-0000-0000-0000-000000000001', '2030-12-31', 'Bagian pengelola keuangan',
			'Satker', 'BK', 'PENY001', 'JAB001', '2025', 'Peraturan 123', 'Remark contoh', true, 'Eselon Nama I',
			NULL
		),
		(
			'22222222-2222-2222-2222-222222222222', 2, 'KI002', 'Bagian SDM', 'E2', 'CK02', 'Kepala Bagian', 'Siti Aminah',
			'11111111-1111-1111-1111-111111111111', 'INST002', 'PNS002', 'JU002', '11111111-1111-1111-1111-111111111111', 8,
			2, false, '00000000-0000-0000-0000-000000000001', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', '11111111-1111-1111-1111-111111111111', '2031-06-30', 'Bagian pengelola SDM',
			'Satker', 'SDM', 'PENY002', 'JAB002', '2026', 'Peraturan 456', 'Remark contoh 2', false, 'Eselon Nama II',
			now()
		),
		(
			'33333333-3333-3333-3333-333333333333', 3, 'KI003', 'Bagian IT', 'E3', 'CK03', 'Kepala Bagian', 'Rina Suryani',
			'22222222-2222-2222-2222-222222222222', 'INST002', 'PNS001', 'JU003', '22222222-2222-2222-2222-222222222222', 10,
			3, true, '00000000-0000-0000-0000-000000000001', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', '22222222-2222-2222-2222-222222222222', '2032-03-31', 'Bagian pengelola IT',
			'Satker', 'IT', 'PENY003', 'JAB003', '2027', 'Peraturan 789', 'Remark IT', true, 'Eselon Nama III',
			NULL
		);
	`

	tests := []struct {
		name             string
		dbData           string
		id               string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:   "ok: delete unit kerja",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
		},
		{
			name:   "error: delete not found unit kerja",
			dbData: dbData,
			id:     "1234",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: delete deleted unit kerja",
			dbData: dbData,
			id:     "22222222-2222-2222-2222-222222222222",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "123456789", api.RoleAdmin)},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:   "error: auth header tidak valid",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:   "error: user is not an admin",
			dbData: dbData,
			id:     "11111111-1111-1111-1111-111111111111",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "987654321")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			pgxconn := dbtest.New(t, dbmigrations.FS)

			_, err := pgxconn.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/unit-kerja/"+tt.id, nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			r := repo.New(pgxconn)
			RegisterRoutes(e, r, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			} else {
				assert.Empty(t, rec.Body.String())
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
