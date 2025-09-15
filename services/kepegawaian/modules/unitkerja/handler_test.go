package unitkerja

import (
	"context"
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

func Test_handler_ListUnitKerja(t *testing.T) {
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
			repo := repo.New(pgxconn)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
