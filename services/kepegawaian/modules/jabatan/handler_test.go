package jabatan_test

import (
	"context"
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
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/jabatan"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan(id, nama_jabatan, kode_jabatan, deleted_at) values
			(11, 'Jabatan 11', '11', null),
			(12, 'Jabatan 12', '12', null),
			(13, 'Jabatan 13', '13', null),
			(14, 'Jabatan 14', '14', '2000-01-01'),
			(15, 'Nama Jabatan 15', '15', null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("41")}
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
			wantResponseBody: `{
				"data": [
					{
						"id": "11",
						"nama": "Jabatan 11"
					},
					{
						"id": "12",
						"nama": "Jabatan 12"
					},
					{
						"id": "13",
						"nama": "Jabatan 13"
					},
					{
						"id": "15",
						"nama": "Nama Jabatan 15"
					}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: dengan filter nama dan pgination",
			requestQuery:     url.Values{"nama": []string{"jabatan"}, "limit": []string{"2"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "12",
						"nama": "Jabatan 12"
					},
					{
						"id": "13",
						"nama": "Jabatan 13"
					}
				],
				"meta": {"limit": 2, "offset": 1, "total": 3}
			}`,
		},
		{
			name:             "ok: dengan parameter pagination",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "12",
						"nama": "Jabatan 12"
					}
				],
				"meta": {"limit": 1, "offset": 1, "total": 4}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/jabatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(db)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			jabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminListJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jenis_jabatan(id, nama, deleted_at) values
		(1, 'Jabatan Struktural', null),
		(2, 'Jabatan Fungsional', null),
		(3, 'Jabatan Deleted', '2000-01-01');

		insert into ref_jabatan
		(kode_jabatan, id, nama_jabatan, nama_jabatan_full, jenis_jabatan, kelas, pensiun, kode_bkn, nama_jabatan_bkn, kategori_jabatan, bkn_id, tunjangan_jabatan, created_at, updated_at, deleted_at) values
		('K001', 101, 'Jabatan A', 'Jabatan A Full', 1, 2, 60, 'BKN101', 'Jabatan A BKN', 'Struktural', 'BKNID101', 1000000, '2023-01-01T00:00:00', '2023-01-01T00:00:00', null),
		('K002', 102, 'Jabatan B', 'Jabatan B Full', 2, 3, 61, 'BKN102', 'Jabatan B BKN', 'Fungsional', 'BKNID102', 1000000, '2023-02-01T00:00:00', '2023-02-01T00:00:00', null),
		('K003', 103, 'Jabatan C', 'Jabatan C Full', 1, 4, 62, 'BKN103', 'Jabatan C BKN', 'Struktural', 'BKNID103', 1000000, '2023-03-01T00:00:00', '2023-03-01T00:00:00', null),
		('K004', 104, 'Jabatan D', 'Jabatan D Full', 2, 5, 63, 'BKN104', 'Jabatan D BKN', 'Fungsional', 'BKNID104', 1000000, '2023-04-01T00:00:00', '2023-04-01T00:00:00', null),
		('K005', 105, 'Jabatan E', 'Jabatan E Full', 3, 6, 64, 'BKN105', 'Jabatan E BKN', 'Struktural', 'BKNID105', 1000000, '2023-05-01T00:00:00', '2023-05-01T00:00:00', null),
		('K006', 106, 'Jabatan F', 'Jabatan F Full', 1, 6, 64, 'BKN106', 'Jabatan F BKN', 'Struktural', 'BKNID106', 1000000, '2023-05-01T00:00:00', '2023-05-01T00:00:00', now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get data with default pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"kode": "K001", "id": 101, "nama": "Jabatan A", "nama_full": "Jabatan A Full", "jenis": 1, "nama_jenis" : "Jabatan Struktural", "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_bkn": "Jabatan A BKN", "kategori": "Struktural", "bkn_id": "BKNID101", "tunjangan" : 1000000},
					{"kode": "K002", "id": 102, "nama": "Jabatan B", "nama_full": "Jabatan B Full", "jenis": 2, "nama_jenis" : "Jabatan Fungsional", "kelas": 3, "pensiun": 61, "kode_bkn": "BKN102", "nama_bkn": "Jabatan B BKN", "kategori": "Fungsional", "bkn_id": "BKNID102", "tunjangan" : 1000000},
					{"kode": "K003", "id": 103, "nama": "Jabatan C", "nama_full": "Jabatan C Full", "jenis": 1, "nama_jenis" : "Jabatan Struktural", "kelas": 4, "pensiun": 62, "kode_bkn": "BKN103", "nama_bkn": "Jabatan C BKN", "kategori": "Struktural", "bkn_id": "BKNID103", "tunjangan" : 1000000},
					{"kode": "K004", "id": 104, "nama": "Jabatan D", "nama_full": "Jabatan D Full", "jenis": 2, "nama_jenis" : "Jabatan Fungsional", "kelas": 5, "pensiun": 63, "kode_bkn": "BKN104", "nama_bkn": "Jabatan D BKN", "kategori": "Fungsional", "bkn_id": "BKNID104", "tunjangan" : 1000000},
					{"kode": "K005", "id": 105, "nama": "Jabatan E", "nama_full": "Jabatan E Full", "jenis": 3, "nama_jenis" : "", "kelas": 6, "pensiun": 64, "kode_bkn": "BKN105", "nama_bkn": "Jabatan E BKN", "kategori": "Struktural", "bkn_id": "BKNID105", "tunjangan" : 1000000}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 5
				}
			}`,
		},
		{
			name: "ok: with pagination limit 1 and offset 1",
			requestQuery: url.Values{
				"limit":  []string{"1"},
				"offset": []string{"1"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"kode": "K002", "id": 102, "nama": "Jabatan B", "nama_full": "Jabatan B Full", "jenis": 2, "nama_jenis" : "Jabatan Fungsional", "kelas": 3, "pensiun": 61, "kode_bkn": "BKN102", "nama_bkn": "Jabatan B BKN", "kategori": "Fungsional", "bkn_id": "BKNID102", "tunjangan" : 1000000}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 5
				}
			}`,
		},
		{
			name: "ok: with filter",
			requestQuery: url.Values{
				"limit":   []string{"10"},
				"offset":  []string{"0"},
				"keyword": []string{"C"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"kode": "K003", "id": 103, "nama": "Jabatan C", "nama_full": "Jabatan C Full", "jenis": 1, "nama_jenis" : "Jabatan Struktural", "kelas": 4, "pensiun": 62, "kode_bkn": "BKN103", "nama_bkn": "Jabatan C BKN", "kategori": "Struktural", "bkn_id": "BKNID103", "tunjangan" : 1000000}
				],
				"meta": {
					"limit": 10,
					"offset": 0,
					"total": 1
				}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/jabatan", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			jabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateJabatan(t *testing.T) {
	t.Parallel()

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: create jabatan",
			requestBody: `{"kode_jabatan": "K001", "nama_jabatan": "Jabatan A", "nama_jabatan_full": "Jabatan A Full", "jenis_jabatan": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_jabatan_bkn": "Jabatan A BKN", "kategori_jabatan": "Struktural", "bkn_id": "BKNID101", "tunjangan_jabatan" : 1000000}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id" : 1, "kode": "K001", "nama": "Jabatan A", "nama_full": "Jabatan A Full", "jenis": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_bkn": "Jabatan A BKN", "kategori": "Struktural", "bkn_id": "BKNID101", "tunjangan" : 1000000}}`,
		},
		{
			name:        "ok: create jabatan with minimalis data",
			requestBody: `{"kode_jabatan": "K001", "nama_jabatan": "Jabatan A", "nama_jabatan_full": "Jabatan A Full"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id" : 1, "kode": "K001", "nama": "Jabatan A", "nama_full": "Jabatan A Full", "jenis": null, "kelas": null, "pensiun": null, "kode_bkn": "", "nama_bkn": "", "kategori": "", "bkn_id": "", "tunjangan" : 0}}`,
		},
		{
			name:        "error: auth header tidak valid",
			requestBody: `{"kode_jabatan": "K001", "nama_jabatan": "Jabatan A", "nama_jabatan_full": "Jabatan A Full", "jenis_jabatan": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_jabatan_bkn": "Jabatan A BKN", "kategori_jabatan": "Struktural", "bkn_id": "BKNID101"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/jabatan", strings.NewReader(tt.requestBody))
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			jabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan
		("kode_jabatan", "id", "nama_jabatan", "nama_jabatan_full", "jenis_jabatan", "kelas", "pensiun", "kode_bkn", "nama_jabatan_bkn", "kategori_jabatan", "bkn_id", "tunjangan_jabatan", "created_at", "updated_at", "deleted_at") values
		('K001', 101, 'Jabatan A', 'Jabatan A Full', 1, 2, 60, 'BKN101', 'Jabatan A BKN', 'Struktural', 'BKNID101', 1000000, '2023-01-01T00:00:00', '2023-01-01T00:00:00', null),
		('K002', 102, 'Jabatan B', 'Jabatan B Full', 2, 3, 61, 'BKN102', 'Jabatan B BKN', 'Fungsional', 'BKNID102', 1000000, '2023-02-01T00:00:00', '2023-02-01T00:00:00', null),
		('K003', 103, 'Jabatan C', 'Jabatan C Full', 1, 4, 62, 'BKN103', 'Jabatan C BKN', 'Struktural', 'BKNID103', 1000000, '2023-03-01T00:00:00', '2023-03-01T00:00:00', now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestPath      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get jabatan by id",
			requestPath:      "/v1/admin/jabatan/101",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"kode": "K001", "id": 101, "nama": "Jabatan A", "nama_full": "Jabatan A Full", "jenis": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_bkn": "Jabatan A BKN", "kategori": "Struktural", "bkn_id": "BKNID101", "tunjangan" : 1000000}}`,
		},
		{
			name:             "error: jabatan not found",
			requestPath:      "/v1/admin/jabatan/999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jabatan deleted",
			requestPath:      "/v1/admin/jabatan/103",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestPath:      "/v1/admin/jabatan/101",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, tt.requestPath, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			jabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan
		("kode_jabatan", "id", "nama_jabatan", "nama_jabatan_full", "jenis_jabatan", "kelas", "pensiun", "kode_bkn", "nama_jabatan_bkn", "kategori_jabatan", "bkn_id", "tunjangan_jabatan", "created_at", "updated_at", "deleted_at") values
		('K001', 101, 'Jabatan A', 'Jabatan A Full', 1, 2, 60, 'BKN101', 'Jabatan A BKN', 'Struktural', 'BKNID101', 1000000, '2023-01-01T00:00:00', '2023-01-01T00:00:00', null),
		('K002', 102, 'Jabatan B', 'Jabatan B Full', 2, 3, 61, 'BKN102', 'Jabatan B BKN', 'Fungsional', 'BKNID102', 1000000, '2023-02-01T00:00:00', '2023-02-01T00:00:00', null),
		('K003', 103, 'Jabatan C', 'Jabatan C Full', 1, 4, 62, 'BKN103', 'Jabatan C BKN', 'Struktural', 'BKNID103', 1000000, '2023-03-01T00:00:00', '2023-03-01T00:00:00', now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestPath      string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: update jabatan",
			requestPath: "/v1/admin/jabatan/101",
			requestBody: `{"kode_jabatan": "K001", "nama_jabatan": "Jabatan A Updated", "nama_jabatan_full": "Jabatan A Full Updated", "jenis_jabatan": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_jabatan_bkn": "Jabatan A BKN", "kategori_jabatan": "Struktural", "bkn_id": "BKNID101", "tunjangan_jabatan" : 1000001}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"kode": "K001", "id": 101, "nama": "Jabatan A Updated", "nama_full": "Jabatan A Full Updated", "jenis": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_bkn": "Jabatan A BKN", "kategori": "Struktural", "bkn_id": "BKNID101", "tunjangan" : 1000001}}`,
		},
		{
			name:        "error: jabatan not found",
			requestPath: "/v1/admin/jabatan/999",
			requestBody: `{"kode_jabatan": "K001", "nama_jabatan": "Jabatan Baru", "nama_jabatan_full": "Jabatan Baru Full", "jenis_jabatan": 1, "kelas": 1, "pensiun": 60, "kode_bkn": "BKN999", "nama_jabatan_bkn": "Jabatan Baru BKN", "kategori_jabatan": "Struktural", "bkn_id": "BKNID999"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: jabatan deleted",
			requestPath: "/v1/admin/jabatan/103",
			requestBody: `{"kode_jabatan": "K003", "nama_jabatan": "Jabatan C Baru", "nama_jabatan_full": "Jabatan C Full Baru", "jenis_jabatan": 1, "kelas": 4, "pensiun": 62, "kode_bkn": "BKN103", "nama_jabatan_bkn": "Jabatan C BKN", "kategori_jabatan": "Struktural", "bkn_id": "BKNID103"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: auth header tidak valid",
			requestPath: "/v1/admin/jabatan/101",
			requestBody: `{"kode_jabatan": "K001", "nama_jabatan": "Jabatan A Updated", "nama_jabatan_full": "Jabatan A Full Updated", "jenis_jabatan": 1, "kelas": 2, "pensiun": 60, "kode_bkn": "BKN101", "nama_jabatan_bkn": "Jabatan A BKN", "kategori_jabatan": "Struktural", "bkn_id": "BKNID101"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, tt.requestPath, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			jabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteJabatan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_jabatan
		("kode_jabatan", "id", "nama_jabatan", "nama_jabatan_full", "jenis_jabatan", "kelas", "pensiun", "kode_bkn", "nama_jabatan_bkn", "kategori_jabatan", "bkn_id", "created_at", "updated_at", "deleted_at") values
		('K001', 101, 'Jabatan A', 'Jabatan A Full', 1, 2, 60, 'BKN101', 'Jabatan A BKN', 'Struktural', 'BKNID101', '2023-01-01T00:00:00', '2023-01-01T00:00:00', null),
		('K002', 102, 'Jabatan B', 'Jabatan B Full', 2, 3, 61, 'BKN102', 'Jabatan B BKN', 'Fungsional', 'BKNID102', '2023-02-01T00:00:00', '2023-02-01T00:00:00', null),
		('K003', 103, 'Jabatan C', 'Jabatan C Full', 1, 4, 62, 'BKN103', 'Jabatan C BKN', 'Struktural', 'BKNID103', '2023-03-01T00:00:00', '2023-03-01T00:00:00', now());
	`
	pgxconn := dbtest.New(t, dbmigrations.FS)
	_, err := pgxconn.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestPath      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: delete jabatan",
			requestPath:      "/v1/admin/jabatan/101",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: ``,
		},
		{
			name:             "error: jabatan not found",
			requestPath:      "/v1/admin/jabatan/999",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: jabatan already deleted",
			requestPath:      "/v1/admin/jabatan/103",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestPath:      "/v1/admin/jabatan/101",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, tt.requestPath, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			jabatan.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			} else {
				assert.Empty(t, tt.wantResponseBody)
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
