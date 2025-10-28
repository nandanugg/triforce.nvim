package pemberitahuan

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
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func getDate(day int) string {
	t := time.Now().AddDate(0, 0, day)
	return time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339)
}

func Test_handler_ListPemberitahuan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pemberitahuan (id, judul_berita, deskripsi_berita, pinned, diterbitkan_pada, ditarik_pada, updated_by, updated_at, deleted_at) values
		  (1, 'Notice over', 'Desc 1', false, current_date - interval '3 days', current_date - interval '2 days', 'admin', current_date - interval '5 days', null),
		  (2, 'Notice active', 'Desc 1', false, current_date - interval '3 days', current_date + interval '3 days', 'admin', current_date - interval '4 days', null),
		  (3, 'Notice waiting', 'Desc 1', false, current_date + interval '3 days', current_date + interval '4 days', 'admin', current_date, null),
		  (4, 'Notice pinned', 'Desc 1', true, current_date - interval '3 days', current_date + interval '3 days', 'admin', current_date - interval '3 days', null),
		  (5, 'Notice 3', 'Desc 3', false, current_date - interval '3 days', current_date + interval '3 days', 'admin', now(), now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgx)
	authSvc := apitest.NewAuthService(api.Kode_Pemberitahuan_Public)
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
			name:             "ok: list all pemberitahuan with default pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 4,
					"judul_berita": "Notice pinned",
					"deskripsi_berita": "Desc 1",
					"pinned": true,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-3) + `"
				},
				{
					"id": 3,
					"judul_berita": "Notice waiting",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "WAITING",
					"diterbitkan_pada": "` + getDate(3) + `",
					"ditarik_pada": "` + getDate(4) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(0) + `"
				},
				{
					"id": 1,
					"judul_berita": "Notice over",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "OVER",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(-2) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-5) + `"
				},
				{
					"id": 2,
					"judul_berita": "Notice active",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-4) + `"
				}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: list pemberitahuan with limit=2",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestQuery:     url.Values{"limit": []string{"2"}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 4,
					"judul_berita": "Notice pinned",
					"deskripsi_berita": "Desc 1",
					"pinned": true,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-3) + `"
				},
				{
					"id": 3,
					"judul_berita": "Notice waiting",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "WAITING",
					"diterbitkan_pada": "` + getDate(3) + `",
					"ditarik_pada": "` + getDate(4) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(0) + `"
				}
				],
				"meta": {"limit": 2, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: list pemberitahuan with limit=2 and offset=2",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"2"}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 1,
					"judul_berita": "Notice over",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "OVER",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(-2) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-5) + `"
				},
				{
					"id": 2,
					"judul_berita": "Notice active",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-4) + `"
				}
				],
				"meta": {"limit": 2, "offset": 2, "total": 4}
			}`,
		},
		{
			name:             "error: invalid token",
			requestHeader:    http.Header{"Authorization": []string{"Bearer invalid"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/v1/pemberitahuan", nil)
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

func Test_handler_adminListPemberitahuan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pemberitahuan (id, judul_berita, deskripsi_berita, pinned, diterbitkan_pada, ditarik_pada, updated_by, updated_at, deleted_at) values
		  (1, 'Notice over', 'Desc 1', false, current_date - interval '3 days', current_date - interval '2 days', 'admin', current_date - interval '5 days', null),
		  (2, 'Notice active', 'Desc 1', false, current_date - interval '3 days', current_date + interval '3 days', 'admin', current_date - interval '4 days', null),
		  (3, 'Notice waiting', 'Desc 1', false, current_date + interval '3 days', current_date + interval '4 days', 'admin', current_date, null),
		  (4, 'Notice pinned', 'Desc 1', true, current_date - interval '3 days', current_date + interval '3 days', 'admin', current_date - interval '3 days', null),
		  (5, 'Notice 3', 'Desc 3', false, current_date - interval '3 days', current_date + interval '3 days', 'admin', now(), now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	repo := sqlc.New(pgx)
	authSvc := apitest.NewAuthService(api.Kode_Pemberitahuan_Read)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("123456789")}

	tests := []struct {
		name             string
		requestHeader    http.Header
		wantResponseCode int
		requestQuery     url.Values
		wantResponseBody string
	}{
		{
			name:             "ok: admin get all pemberitahuan",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 4,
					"judul_berita": "Notice pinned",
					"deskripsi_berita": "Desc 1",
					"pinned": true,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-3) + `"
				},
				{
					"id": 3,
					"judul_berita": "Notice waiting",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "WAITING",
					"diterbitkan_pada": "` + getDate(3) + `",
					"ditarik_pada": "` + getDate(4) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(0) + `"
				},
				{
					"id": 1,
					"judul_berita": "Notice over",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "OVER",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(-2) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-5) + `"
				},
				{
					"id": 2,
					"judul_berita": "Notice active",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-4) + `"
				}
				],
				"meta": {"limit": 10, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: list pemberitahuan with limit=2",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestQuery:     url.Values{"limit": []string{"2"}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 4,
					"judul_berita": "Notice pinned",
					"deskripsi_berita": "Desc 1",
					"pinned": true,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-3) + `"
				},
				{
					"id": 3,
					"judul_berita": "Notice waiting",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "WAITING",
					"diterbitkan_pada": "` + getDate(3) + `",
					"ditarik_pada": "` + getDate(4) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(0) + `"
				}
				],
				"meta": {"limit": 2, "offset": 0, "total": 4}
			}`,
		},
		{
			name:             "ok: list pemberitahuan with limit=2 and offset=2",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestQuery:     url.Values{"limit": []string{"2"}, "offset": []string{"2"}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
				{
					"id": 1,
					"judul_berita": "Notice over",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "OVER",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(-2) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-5) + `"
				},
				{
					"id": 2,
					"judul_berita": "Notice active",
					"deskripsi_berita": "Desc 1",
					"pinned": false,
					"status": "ACTIVE",
					"diterbitkan_pada": "` + getDate(-3) + `",
					"ditarik_pada": "` + getDate(3) + `",
					"diperbarui_oleh": "admin",
					"terakhir_diperbarui": "` + getDate(-4) + `"
				}
				],
				"meta": {"limit": 2, "offset": 2, "total": 4}
			}`,
		},
		{
			name:             "error: invalid token",
			requestHeader:    http.Header{"Authorization": []string{"Bearer invalid"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			req := httptest.NewRequest(http.MethodGet, "/v1/admin/pemberitahuan", nil)
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

func Test_handler_adminCreatePemberitahuan(t *testing.T) {
	t.Parallel()

	pgx := dbtest.New(t, dbmigrations.FS)
	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)
	repo := sqlc.New(pgx)
	authSvc := apitest.NewAuthService(api.Kode_Pemberitahuan_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
	authHeader := []string{apitest.GenerateAuthHeader("123456789")}

	tests := []struct {
		name             string
		requestBody      string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: create pemberitahuan",
			requestBody: `{
				"judul_berita":"New Notice",
				"deskripsi_berita":"Some desc",
				"pinned":false,
				"diterbitkan_pada":"2024-01-01T00:00:00Z",
				"ditarik_pada":"2024-01-02T00:00:00Z"
			}`,
			requestHeader: http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			},
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"judul_berita": "New Notice",
					"deskripsi_berita": "Some desc",
					"pinned": false,
					"status": "OVER",
					"diterbitkan_pada":"` + time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `",
					"ditarik_pada":"` + time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `",
					"diperbarui_oleh": "123456789",
					"terakhir_diperbarui": "{updated_at}"
				}
			}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPost, "/v1/admin/pemberitahuan", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if rec.Code == http.StatusCreated {
				var resp createUpdateResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.WithinDuration(t, time.Now(), resp.Data.TerakhirDiperbarui, 10*time.Second)
				tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{updated_at}", resp.Data.TerakhirDiperbarui.Format(time.RFC3339Nano))
			}
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_adminUpdatePemberitahuan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pemberitahuan (
		  id, judul_berita, deskripsi_berita, pinned,
		  diterbitkan_pada, ditarik_pada, updated_by, updated_at, deleted_at
		)
		values
		  (1, 'Old', 'Desc', false, now(), now(), 'admin', now(), null),
		  (2, 'Deleted', 'Desc', false, now(), now(), 'admin', now(), now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)
	e, _ := api.NewEchoServer(docs.OpenAPIBytes)
	repo := sqlc.New(pgx)
	authSvc := apitest.NewAuthService(api.Kode_Pemberitahuan_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
	authHeader := []string{apitest.GenerateAuthHeader("123456789")}

	tests := []struct {
		name             string
		id               string
		requestBody      string
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name: "ok: update pemberitahuan",
			id:   "1",
			requestBody: `{
				"judul_berita":"New Notice",
				"deskripsi_berita":"Some desc",
				"pinned":false,
				"diterbitkan_pada":"2024-01-01T00:00:00Z",
				"ditarik_pada":"2024-01-02T00:00:00Z"
			}`,
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"judul_berita": "New Notice",
					"deskripsi_berita": "Some desc",
					"pinned": false,
					"status": "OVER",
					"diterbitkan_pada":"` + time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `",
					"ditarik_pada":"` + time.Date(2024, 1, 2, 0, 0, 0, 0, time.UTC).Local().Format(time.RFC3339) + `",
					"diperbarui_oleh": "123456789",
					"terakhir_diperbarui": "{updated_at}"
				}
			}`,
		},
		{
			name: "error: deleted id",
			id:   "2",
			requestBody: `{
				"judul_berita":"New Notice",
				"deskripsi_berita":"Some desc",
				"pinned":false,
				"diterbitkan_pada":"2024-01-01T00:00:00Z",
				"ditarik_pada":"2024-01-02T00:00:00Z"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message":"data tidak ditemukan"}`,
		},
		{
			name: "error: not exists",
			id:   "99",
			requestBody: `{
				"judul_berita":"New Notice",
				"deskripsi_berita":"Some desc",
				"pinned":false,
				"diterbitkan_pada":"2024-01-01T00:00:00Z",
				"ditarik_pada":"2024-01-02T00:00:00Z"
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message":"data tidak ditemukan"}`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodPut, "/v1/admin/pemberitahuan/"+tt.id, strings.NewReader(tt.requestBody))
			req.Header = http.Header{
				"Authorization": authHeader,
				"Content-Type":  []string{"application/json"},
			}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			if rec.Code == http.StatusOK {
				var resp createUpdateResponse
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				require.NoError(t, err)

				assert.WithinDuration(t, time.Now(), resp.Data.TerakhirDiperbarui, 10*time.Second)
				tt.wantResponseBody = strings.ReplaceAll(tt.wantResponseBody, "{updated_at}", resp.Data.TerakhirDiperbarui.Format(time.RFC3339Nano))
			}
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
		})
	}
}

func Test_handler_adminDeletePemberitahuan(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into pemberitahuan (
		  id, judul_berita, deskripsi_berita, pinned,
		  diterbitkan_pada, ditarik_pada, updated_by, updated_at, deleted_at
		)
		values
		  (1, 'Old', 'Desc', false, now(), now(), 'admin', now(), null),
		  (2, 'Deleted', 'Desc', false, now(), now(), 'admin', now(), now());
	`
	pgx := dbtest.New(t, dbmigrations.FS)
	_, err := pgx.Exec(context.Background(), dbData)
	require.NoError(t, err)
	e, _ := api.NewEchoServer(docs.OpenAPIBytes)
	repo := sqlc.New(pgx)
	authSvc := apitest.NewAuthService(api.Kode_Pemberitahuan_Write)
	RegisterRoutes(e, repo, api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
	authHeader := []string{apitest.GenerateAuthHeader("123456789")}

	tests := []struct {
		name             string
		id               string
		wantResponseCode int
	}{
		{"ok: delete pemberitahuan", "1", http.StatusNoContent},
		{"error: already deleted", "2", http.StatusNotFound},
		{"error: not exists", "99", http.StatusNotFound},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodDelete, "/v1/admin/pemberitahuan/"+tt.id, nil)
			req.Header = http.Header{"Authorization": authHeader}
			rec := httptest.NewRecorder()
			e.ServeHTTP(rec, req)
			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
