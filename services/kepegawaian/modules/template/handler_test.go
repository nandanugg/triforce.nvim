package template

import (
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"net/textproto"
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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
)

func Test_handler_ListTemplate(t *testing.T) {
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

	defaultTimestamptz := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	dbData := `
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan 2', 'data:image/png;base64,invalid', '2001-01-01', '2001-01-01', NULL),
			(13, 'Penghargaan 3', 'data:image/png;base64,invalid', '2001-01-01', '2001-01-01', now());
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 11,
						"nama": "Penghargaan 1",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
					},
					{
						"id": 12,
						"nama": "Penghargaan 2",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
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
			name:   "ok: with pagination limit 1 and offset 1",
			dbData: dbData,
			requestQuery: url.Values{
				"limit":  []string{"1"},
				"offset": []string{"1"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 12,
						"nama": "Penghargaan 2",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
					}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 2
				}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/template", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_GetBerkasTemplate(t *testing.T) {
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
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan png', 'data:images/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan pdf', 'data:application/pdf;base64,` + pdfBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(95, 'Penghargaan kosong', '', '2001-01-01', '2001-01-01', NULL),
			(96, 'Penghargaan null', NULL, '2001-01-01', '2001-01-01', NULL),
			(97, 'Penghargaan deleted', 'data:images/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', now()),
			(98, 'Penghargaan invalid', 'data:images/png;base64,invalid', '2001-01-01', '2001-01-01', NULL)
			;
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
			name:              "ok: valid pdf with data: prefix",
			dbData:            dbData,
			paramID:           "12",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 tidak valid",
			dbData:            dbData,
			paramID:           "98",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: sudah dihapus",
			dbData:            dbData,
			paramID:           "97",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			dbData:            dbData,
			paramID:           "96",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			dbData:            dbData,
			paramID:           "95",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: template tidak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/template/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			queries := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Public)
			RegisterRoutes(e, queries, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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

func Test_handler_adminListTemplate(t *testing.T) {
	t.Parallel()
	defaultTimestamptz := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
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
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan 2', 'data:image/png;base64,invalid', '2001-01-01', '2001-01-01', NULL),
			(13, 'Penghargaan 3', 'data:image/png;base64,invalid', '2001-01-01', '2001-01-01', now());
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 11,
						"nama": "Penghargaan 1",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
					},
					{
						"id": 12,
						"nama": "Penghargaan 2",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
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
			name:   "ok: with pagination limit 1 and offset 1",
			dbData: dbData,
			requestQuery: url.Values{
				"limit":  []string{"1"},
				"offset": []string{"1"},
			},
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 12,
						"nama": "Penghargaan 2",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
					}
				],
				"meta": {
					"limit": 1,
					"offset": 1,
					"total": 2
				}
			}`,
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/template", nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetTemplate(t *testing.T) {
	t.Parallel()
	defaultTimestamptz := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
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
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan 2', 'data:image/png;base64,invalid', '2001-01-01', '2001-01-01', NULL),
			(13, 'Penghargaan 3', 'data:image/png;base64,invalid', '2001-01-01', '2001-01-01', now());
		`

	tests := []struct {
		name             string
		dbData           string
		paramID          string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: get data",
			dbData:           dbData,
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
						"id": 11,
						"nama": "Penghargaan 1",
						"created_at": "` + defaultTimestamptz + `",
						"updated_at": "` + defaultTimestamptz + `"
					}
			}`,
		},
		{
			name:             "error: not found",
			dbData:           dbData,
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("41")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: auth header tidak valid",
			dbData:           dbData,
			paramID:          "11",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
		{
			name:             "error: missing auth header",
			dbData:           dbData,
			paramID:          "11",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/admin/template/"+tt.paramID, nil)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminGetBerkasTemplate(t *testing.T) {
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
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan png', 'data:images/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan pdf', 'data:application/pdf;base64,` + pdfBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(95, 'Penghargaan kosong', '', '2001-01-01', '2001-01-01', NULL),
			(96, 'Penghargaan null', NULL, '2001-01-01', '2001-01-01', NULL),
			(97, 'Penghargaan deleted', 'data:images/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', now()),
			(98, 'Penghargaan invalid', 'data:images/png;base64,invalid', '2001-01-01', '2001-01-01', NULL)
			;
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
			name:              "ok: valid pdf with data: prefix",
			dbData:            dbData,
			paramID:           "12",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			dbData:            dbData,
			paramID:           "11",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusOK,
			wantContentType:   "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:              "error: base64 tidak valid",
			dbData:            dbData,
			paramID:           "98",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusInternalServerError,
			wantResponseBytes: []byte(`{"message": "Internal Server Error"}`),
		},
		{
			name:              "error: sudah dihapus",
			dbData:            dbData,
			paramID:           "97",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berisi null value",
			dbData:            dbData,
			paramID:           "96",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: base64 berupa string kosong",
			dbData:            dbData,
			paramID:           "95",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: template tiiak ditemukan",
			dbData:            dbData,
			paramID:           "0",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
			wantResponseCode:  http.StatusNotFound,
			wantResponseBytes: []byte(`{"message": "berkas template tidak ditemukan"}`),
		},
		{
			name:              "error: invalid id",
			dbData:            dbData,
			paramID:           "abc",
			requestHeader:     http.Header{"Authorization": []string{apitest.GenerateAuthHeader("123456789")}},
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

			req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/v1/admin/template/%s/berkas", tt.paramID), nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			queries := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Read)
			RegisterRoutes(e, queries, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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

func Test_handler_adminCreateTemplate(t *testing.T) {
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

	tests := []struct {
		name             string
		formFields       map[string]string
		files            map[string][]byte
		authHeader       string
		fileContentType  string
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: create template with file",
			formFields: map[string]string{
				"nama": "master 1",
			},
			files: map[string][]byte{
				"file": pngBytes,
			},
			fileContentType:  "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{"data": {
				"id": 1,
				"nama": "master 1",
				"created_at": "{created_at}",
				"updated_at": "{updated_at}"
			}}`,
		},
		{
			name: "error: file with invalid type",
			formFields: map[string]string{
				"nama": "master 1",
			},
			files: map[string][]byte{
				"file": pngBytes,
			},
			fileContentType:  "image/x-xpixmap",
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"file\" harus dalam format yang sesuai"}`,
		},
		{
			name: "error: missing file upload",
			formFields: map[string]string{
				"nama": "no file",
			},
			fileContentType:  "application/pdf",
			files:            nil,
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"file\" harus diisi"}`,
		},
		{
			name: "error: invalid auth header",
			formFields: map[string]string{
				"nama": "bad",
			},
			fileContentType: "application/pdf",
			files: map[string][]byte{
				"file": pngBytes,
			},
			authHeader:       "Bearer some-token",
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pgxconn := dbtest.New(t, dbmigrations.FS)
			_, err := pgxconn.Exec(context.Background(), "")
			require.NoError(t, err)

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			for k, v := range tt.formFields {
				require.NoError(t, writer.WriteField(k, v))
			}

			for fieldName, content := range tt.files {
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, "example.bin"))
				h.Set("Content-Type", tt.fileContentType)

				part, err := writer.CreatePart(h)
				require.NoError(t, err)
				_, err = part.Write(content)
				require.NoError(t, err)
			}

			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/template", &buf)
			req.Header.Set("Authorization", tt.authHeader)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			// --- Assertions ---
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

						assert.WithinDuration(t, time.Now(), parsed, 10*time.Second)

						wantBody = strings.ReplaceAll(wantBody, "{"+field+"}", val)
					}
				}
			}

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, wantBody, gotBody)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminUpdateTemplate(t *testing.T) {
	t.Parallel()

	defaultTimestamptz := time.Date(2001, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
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
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan deleted', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', now());
		`

	tests := []struct {
		name             string
		dbData           string
		paramID          string
		formFields       map[string]string
		fileContentType  string
		files            map[string][]byte
		authHeader       string
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:    "ok: template with file",
			dbData:  dbData,
			paramID: "11",
			formFields: map[string]string{
				"nama": "master 1",
			},
			fileContentType: "application/vnd.openxmlformats-officedocument.wordprocessingml.document",
			files: map[string][]byte{
				"file": pngBytes,
			},
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{"data": {
				"id": 11,
				"nama": "master 1",
				"created_at": "` + defaultTimestamptz + `",
				"updated_at": "{updated_at}"
			}}`,
		},
		{
			name:    "error: file with invalid type",
			dbData:  dbData,
			paramID: "11",
			formFields: map[string]string{
				"nama": "master 1",
			},
			fileContentType: "image/x-xpixmap",
			files: map[string][]byte{
				"file": pngBytes,
			},
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"file\" harus dalam format yang sesuai"}`,
		},
		{
			name:    "error: not found if deleted",
			dbData:  dbData,
			paramID: "12",
			formFields: map[string]string{
				"nama": "master 1",
			},
			fileContentType: "application/pdf",
			files: map[string][]byte{
				"file": pngBytes,
			},
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:    "error: not found",
			dbData:  dbData,
			paramID: "1",
			formFields: map[string]string{
				"nama": "master 1",
			},
			fileContentType: "application/pdf",
			files: map[string][]byte{
				"file": pngBytes,
			},
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:    "error: missing file upload",
			dbData:  dbData,
			paramID: "11",
			formFields: map[string]string{
				"nama": "no file",
			},
			fileContentType:  "application/pdf",
			files:            nil,
			authHeader:       apitest.GenerateAuthHeader("123456789"),
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"file\" harus diisi"}`,
		},
		{
			name:    "error: invalid auth header",
			dbData:  dbData,
			paramID: "11",
			formFields: map[string]string{
				"nama": "bad",
			},
			fileContentType: "application/pdf",
			files: map[string][]byte{
				"file": pngBytes,
			},
			authHeader:       "Bearer some-token",
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

			var buf bytes.Buffer
			writer := multipart.NewWriter(&buf)

			for k, v := range tt.formFields {
				require.NoError(t, writer.WriteField(k, v))
			}

			for fieldName, content := range tt.files {
				h := make(textproto.MIMEHeader)
				h.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, fieldName, "example.bin"))
				h.Set("Content-Type", tt.fileContentType)

				part, err := writer.CreatePart(h)
				require.NoError(t, err)
				_, err = part.Write(content)
				require.NoError(t, err)
			}

			require.NoError(t, writer.Close())

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/template/"+tt.paramID, &buf)
			req.Header.Set("Authorization", tt.authHeader)
			req.Header.Set("Content-Type", writer.FormDataContentType())

			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			// --- Assertions ---
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

						assert.WithinDuration(t, time.Now(), parsed, 10*time.Second)

						wantBody = strings.ReplaceAll(wantBody, "{"+field+"}", val)
					}
				}
			}

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, wantBody, gotBody)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_adminDeleteTemplate(t *testing.T) {
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
		insert into ref_template
			(id, nama, file_base64, created_at, updated_at, deleted_at)
			values
			(11, 'Penghargaan 1', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', NULL),
			(12, 'Penghargaan dihapus', 'data:image/png;base64,` + pngBase64 + `', '2001-01-01', '2001-01-01', now());
		`

	tests := []struct {
		name             string
		paramID          string
		requestQuery     url.Values
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:    "ok: delete",
			paramID: "11",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `{}`,
		},
		{
			name:    "error: data sudah dihapus",
			paramID: "12",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:    "error: id tidak valid",
			paramID: "100",
			requestHeader: http.Header{
				"Authorization": []string{apitest.GenerateAuthHeader("123456789")},
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:    "error: auth header tidak valid",
			paramID: "1",
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

			_, err := pgxconn.Exec(context.Background(), dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/template/"+tt.paramID, nil)

			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			repo := repository.New(pgxconn)
			authSvc := apitest.NewAuthService(api.Kode_DataMaster_Write)
			RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			if rec.Body.String() != "" {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			} else {
				assert.Empty(t, rec.Body.String())
			}
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
