package pendidikan

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
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
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan (id, nama, deleted_at) values
			(1, 'Jenis Pendidikan 1', null),
			(2, 'Jenis Pendidikan 2', null),
			(3, 'Jenis Pendidikan 3', now());

		insert into ref_pendidikan (id, nama, tingkat_pendidikan_id, deleted_at) values
			('id-1', 'Pendidikan 1', 1, null),
			('id-2', 'Pendidikan 2', 2, null),
			('id-3', 'Pendidikan 3', 3, now()),
			('id-4', 'Pendidikan 4', 3, null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: sucess without parameter",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "id-1",
						"nama": "Pendidikan 1",
						"tingkat_pendidikan": "Jenis Pendidikan 1",
						"tingkat_pendidikan_id": 1
					},
					{
						"id": "id-2",
						"nama": "Pendidikan 2",
						"tingkat_pendidikan": "Jenis Pendidikan 2",
						"tingkat_pendidikan_id": 2
					},
					{
						"id": "id-4",
						"nama": "Pendidikan 4",
						"tingkat_pendidikan": null,
						"tingkat_pendidikan_id": 3
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
			name:             "ok: with parameter pagination",
			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "id-2",
						"nama": "Pendidikan 2",
						"tingkat_pendidikan": "Jenis Pendidikan 2",
						"tingkat_pendidikan_id": 2
					}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 3
				}
			}`,
		},
		{
			name:             "ok: with parameter nama",
			requestQuery:     url.Values{"nama": []string{"1"}},
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": "id-1",
						"nama": "Pendidikan 1",
						"tingkat_pendidikan": "Jenis Pendidikan 1",
						"tingkat_pendidikan_id": 1
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
			name:             "error: auth header tidak valid",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: missing auth header",
			requestHeader:    http.Header{},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pendidikan", nil)
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

func Test_handler_get(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan (id, nama, deleted_at) values
			(1, 'Jenis Pendidikan 1', null),
			(2, 'Jenis Pendidikan 2', null),
			(3, 'Jenis Pendidikan 3', now());

		insert into ref_pendidikan (id, nama, tingkat_pendidikan_id, deleted_at) values
			('id-1', 'Pendidikan 1', 1, null),
			('id-2', 'Pendidikan 2', 2, null),
			('id-3', 'Pendidikan 3', 3, now()),
			('id-4', 'Pendidikan 4', 3, null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: sucess",
			id:               "id-1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": "id-1",
					"nama": "Pendidikan 1",
					"tingkat_pendidikan": "Jenis Pendidikan 1",
					"tingkat_pendidikan_id": 1
				}
			}`,
		},
		{
			name:             "ok: get another pendidikan",
			id:               "id-4",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": "id-4",
					"nama": "Pendidikan 4",
					"tingkat_pendidikan": null,
					"tingkat_pendidikan_id": 3
				}
			}`,
		},
		{
			name:             "error: missing auth header",
			id:               "id-1",
			requestHeader:    http.Header{},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: data tidak ditemukan",
			id:               "id-5",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pendidikan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_create(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan (id, nama, deleted_at) values
			(1, 'Jenis Pendidikan 1', null),
			(2, 'Jenis Pendidikan 2', null),
			(3, 'Jenis Pendidikan 3', now());
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	removeUnassertedKey := func(jsonStr string, keys []string) string {
		var m map[string]any
		if err := json.Unmarshal([]byte(jsonStr), &m); err != nil {
			return jsonStr
		}
		for _, key := range keys {
			parts := strings.Split(key, ".")
			curr := m
			for i, part := range parts {
				if i == len(parts)-1 {
					delete(curr, part)
				} else {
					// Traverse deeper if possible
					if next, ok := curr[part].(map[string]any); ok {
						curr = next
					} else {
						break
					}
				}
			}
		}
		b, err := json.Marshal(m)
		if err != nil {
			return jsonStr
		}
		return string(b)
	}

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
		removedKeys      []string
	}{
		{
			name:        "ok: create pendidikan",
			requestBody: `{"nama": "Pendidikan 1", "tingkat_pendidikan_id": 1}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {"nama": "Pendidikan 1", "tingkat_pendidikan_id": 1, "tingkat_pendidikan": "Jenis Pendidikan 1"}}`,
			removedKeys:      []string{"data.id"}, // remove id key from response body because it's auto generated
		},
		{
			name:        "error: missing required field nama",
			requestBody: `{"tingkat_pendidikan_id": 1}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name:             "error: auth header tidak valid",
			requestBody:      `{"nama": "Pendidikan 1", "tingkat_pendidikan_id": 1}`,
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pendidikan", nil)
			req.Header = tt.requestHeader
			req.Body = io.NopCloser(bytes.NewBufferString(tt.requestBody))
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			responseBody := rec.Body.String()
			if len(tt.removedKeys) > 0 {
				responseBody = removeUnassertedKey(responseBody, tt.removedKeys)
			}
			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, responseBody)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_update(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan (id, nama, deleted_at) values
			(1, 'Jenis Pendidikan 1', null),
			(2, 'Jenis Pendidikan 2', null),
			(3, 'Jenis Pendidikan 3', now());

		insert into ref_pendidikan (id, nama, tingkat_pendidikan_id, deleted_at) values
			('id-1', 'Pendidikan 1', 1, null),
			('id-2', 'Pendidikan 2', 2, null),
			('id-3', 'Pendidikan 3', 3, now()),
			('id-4', 'Pendidikan 4', 3, null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:        "ok: update pendidikan",
			id:          "id-1",
			requestBody: `{"nama": "Pendidikan 1 Updated", "tingkat_pendidikan_id": 2}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {"id": "id-1", "nama": "Pendidikan 1 Updated", "tingkat_pendidikan_id": 2, "tingkat_pendidikan": "Jenis Pendidikan 2"}}`,
		},
		{
			name:        "error: missing required field nama",
			id:          "id-1",
			requestBody: `{"tingkat_pendidikan_id": 2}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"nama\" harus diisi"}`,
		},
		{
			name:        "error: auth header tidak valid",
			id:          "id-1",
			requestBody: `{"nama": "Pendidikan 1 Updated", "tingkat_pendidikan_id": 2}`,
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:        "error: data tidak ditemukan",
			id:          "id-5",
			requestBody: `{"nama": "Pendidikan 1 Updated", "tingkat_pendidikan_id": 2}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:        "error: data deleted",
			id:          "id-3",
			requestBody: `{"nama": "Pendidikan 1 Updated", "tingkat_pendidikan_id": 2}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pendidikan/"+tt.id, nil)
			req.Header = tt.requestHeader
			req.Body = io.NopCloser(bytes.NewBufferString(tt.requestBody))
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_delete(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into ref_tingkat_pendidikan (id, nama, deleted_at) values
			(1, 'Jenis Pendidikan 1', null),
			(2, 'Jenis Pendidikan 2', null),
			(3, 'Jenis Pendidikan 3', now());

		insert into ref_pendidikan (id, nama, tingkat_pendidikan_id, deleted_at) values
			('id-1', 'Pendidikan 1', 1, null),
			('id-2', 'Pendidikan 2', 2, null),
			('id-3', 'Pendidikan 3', 3, now()),
			('id-4', 'Pendidikan 4', 3, null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := dbrepo.New(db)
	authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}
	tests := []struct {
		name             string
		id               string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: delete pendidikan",
			id:   "id-1",
			requestHeader: http.Header{
				"Authorization": authHeader,
			},
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: "",
		},
		{
			name: "error: data tidak ditemukan",
			id:   "id-5",
			requestHeader: http.Header{
				"Authorization": authHeader,
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name: "error: auth header tidak valid",
			id:   "id-1",
			requestHeader: http.Header{
				"Authorization": []string{"Bearer some-token"},
			},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pendidikan/"+tt.id, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

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
