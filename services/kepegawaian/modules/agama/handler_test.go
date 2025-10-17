package agama_test

import (
	"context"
	"encoding/json"
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
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/modules/agama"
)

func Test_handler_ListRefAgama(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_agama ("id", "nama", "created_at", "updated_at", "deleted_at") values
		(1, 'Islam', '2024-01-01', '2024-01-01', null),
		(2, 'Kristen', '2024-01-01', '2024-01-01', null),
		(3, 'Katolik', '2024-01-01', '2024-01-01', null),
		(4, 'Hindu', '2024-01-01', '2024-01-01', null),
		(5, 'Budha', '2024-01-01', '2024-01-01', null),
		(6, 'Konghucu', '2024-01-01', '2024-01-01', null),
		(7, 'Kepercayaan', '2024-01-01', '2024-01-01', null),
		(8, 'Lainnya', '2024-01-01', '2024-01-01', null),
		(9, 'Test', '2024-01-01', '2024-01-01', now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	defaulTimestamptz := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: admin get all agama with default pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "Islam", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 2, "nama": "Kristen", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 3, "nama": "Katolik", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 4, "nama": "Hindu", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 5, "nama": "Budha", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 6, "nama": "Konghucu", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 7, "nama": "Kepercayaan", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 8, "nama": "Lainnya", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 8}
			}`,
		},
		{
			name: "ok: pagination limit 3 offset 2",
			requestQuery: url.Values{
				"limit":  []string{"3"},
				"offset": []string{"2"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 3, "nama": "Katolik", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 4, "nama": "Hindu", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 5, "nama": "Budha", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"}
				],
				"meta": {"limit": 3, "offset": 2, "total": 8}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/agama", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgx)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			agama.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminListRefAgama(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_agama ("id", "nama", "created_at", "updated_at", "deleted_at") values
		(1, 'Islam', '2024-01-01', '2024-01-01', null),
		(2, 'Kristen', '2024-01-01', '2024-01-01', null),
		(3, 'Katolik', '2024-01-01', '2024-01-01', null),
		(4, 'Hindu', '2024-01-01', '2024-01-01', null),
		(5, 'Budha', '2024-01-01', '2024-01-01', null),
		(6, 'Konghucu', '2024-01-01', '2024-01-01', null),
		(7, 'Kepercayaan', '2024-01-01', '2024-01-01', null),
		(8, 'Lainnya', '2024-01-01', '2024-01-01', null),
		(9, 'Test', '2024-01-01', '2024-01-01', now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	defaulTimestamptz := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: admin get all agama with default pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 1, "nama": "Islam", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 2, "nama": "Kristen", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 3, "nama": "Katolik", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 4, "nama": "Hindu", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 5, "nama": "Budha", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 6, "nama": "Konghucu", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 7, "nama": "Kepercayaan", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 8, "nama": "Lainnya", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"}
				],
				"meta": {"limit": 10, "offset": 0, "total": 8}
			}`,
		},
		{
			name: "ok: pagination limit 3 offset 2",
			requestQuery: url.Values{
				"limit":  []string{"3"},
				"offset": []string{"2"},
			},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{"id": 3, "nama": "Katolik", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 4, "nama": "Hindu", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"},
					{"id": 5, "nama": "Budha", "created_at": "` + defaulTimestamptz + `", "updated_at": "` + defaulTimestamptz + `"}
				],
				"meta": {"limit": 3, "offset": 2, "total": 8}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/agama", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgx)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			agama.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetRefAgama(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_agama ("id", "nama", "created_at", "updated_at", "deleted_at") values
		(1, 'Islam', '2024-01-01', '2024-01-01', null),
		(2, 'Kristen', '2024-01-01', '2024-01-01', now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	defaulTimestamptz := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: get existing agama by id",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": authHeader,
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"nama": "Islam",
					"created_at": "` + defaulTimestamptz + `",
					"updated_at": "` + defaulTimestamptz + `"
				}
			}`,
		},
		{
			name: "error: id not found",
			id:   "99",
			requestHeader: http.Header{
				"Authorization": authHeader,
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: already deleted",
			id:   "2",
			requestHeader: http.Header{
				"Authorization": authHeader,
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: invalid token",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": []string{"Bearer invalid"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/agama/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgx)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			agama.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminCreateRefAgama(t *testing.T) {
	t.Parallel()

	pgx := dbtest.New(t, dbmigrations.FS)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: admin create agama",
			requestBody: `{"nama": "Zoroastrian"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"id": 1, "nama": "Zoroastrian", "created_at": "{created_at}", "updated_at":"{updated_at}"}}`,
		},
		{
			name:        "error: invalid token",
			requestBody: `{"nama": "Zoroastrian"}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer invalid-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/agama", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgx)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			agama.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody != "" {
				gotBody := rec.Body.String()
				wantBody := tt.wantResponseBody

				var got map[string]any
				require.NoError(t, json.Unmarshal([]byte(gotBody), &got))

				data, ok := got["data"].(map[string]any)
				if ok {
					for _, field := range []string{"created_at", "updated_at"} {
						if val, ok := data[field].(string); ok {
							parsed, err := time.Parse(time.RFC3339, val)
							require.NoErrorf(t, err, "invalid timestamp for %s", field)

							diff := time.Since(parsed)
							if diff < 0 {
								diff = -diff
							}
							assert.LessOrEqualf(t, diff.Seconds(), 10.0, "%s difference too large: %v", field, diff)

							wantBody = strings.ReplaceAll(wantBody, "{"+field+"}", val)
						}
					}
				}

				assert.JSONEq(t, wantBody, gotBody)
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateRefAgama(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_agama ("id", "nama", "created_at", "updated_at", "deleted_at") values
		(1, 'Islam', '2024-01-01', now(), null),
		(2, 'Kristen', '2024-01-01', now(), now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	defaulTimestamptz := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
	tests := []struct {
		name             string
		id               string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: admin update agama",
			id:          "1",
			requestBody: `{"nama": "Islam Updated"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data":
				{"id": 1, "nama": "Islam Updated", "created_at": "` + defaulTimestamptz + `", "updated_at":"{updated_at}"}
			}`,
		},
		{
			name:        "error: id not found if deleted",
			id:          "2",
			requestBody: `{"nama": "Unknown"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
		},
		{
			name:        "error: id not found if not exists",
			id:          "100",
			requestBody: `{"nama": "Unknown"}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/agama/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgx)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			agama.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if tt.wantResponseBody != "" && rec.Body.Len() > 0 {
				gotBody := rec.Body.String()
				wantBody := tt.wantResponseBody

				var got map[string]any
				require.NoError(t, json.Unmarshal([]byte(gotBody), &got))

				data, ok := got["data"].(map[string]any)
				if ok {
					for _, field := range []string{"updated_at"} {
						if val, ok := data[field].(string); ok {
							parsed, err := time.Parse(time.RFC3339, val)
							require.NoErrorf(t, err, "invalid timestamp for %s", field)

							diff := time.Since(parsed)
							if diff < 0 {
								diff = -diff
							}
							assert.LessOrEqualf(t, diff.Seconds(), 10.0, "%s difference too large: %v", field, diff)

							wantBody = strings.ReplaceAll(wantBody, "{"+field+"}", val)
						}
					}
				}

				assert.JSONEq(t, wantBody, gotBody)
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteRefAgama(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_agama ("id", "nama", "created_at", "updated_at", "deleted_at") values
		(1, 'Islam', now(), now(), null),
		(2, 'Kristen', now(), now(), now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: admin delete agama",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": authHeader,
			},
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `{}`,
		},
		{
			name: "error: id not found if deleted",
			id:   "2",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message":"data tidak ditemukan"}`,
		},
		{
			name: "error: id not found if not exists",
			id:   "100",
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message":"data tidak ditemukan"}`,
		},
		{
			name: "error: invalid token",
			id:   "1",
			requestHeader: http.Header{
				"Authorization": []string{"Bearer invalid"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/agama/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgx)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			agama.RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if rec.Body.Len() > 0 {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
