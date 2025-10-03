package resourcepermission

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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/config"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func Test_handler_listResourcePermissions(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into permission
			(id, kode,     nama,     deleted_at) values
			(1,  'write',  'Write',  null),
			(2,  'read',   'Read',   null),
			(3,  'delete', 'Delete', '2000-01-01');
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page2', 'Page 2', null),
			(2,  'portal', 'page1', 'Page 1', null),
			(3,  'portal', 'page3', 'Page 3', null),
			(4,  'portal', 'page4', 'Page 4', '2000-01-01');
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           1,             null),
			(2,  2,           1,             null),
			(3,  3,           1,             null),
			(4,  4,           1,             null),
			(5,  1,           2,             null),
			(6,  2,           2,             null),
			(7,  3,           2,             null),
			(8,  1,           3,             null),
			(9,  1,           2,             '2000-01-01');
		insert into role
			(id, nama,      deleted_at) values
			(1,  'admin',   null),
			(2,  'pegawai', null),
			(3,  'guest',   null),
			(4,  'deleted', '2000-01-01');
		insert into role_resource_permission
			(role_id, resource_permission_id, deleted_at) values
			(1,       1,                      null),
			(1,       2,                      null),
			(4,       3,                      null),
			(1,       4,                      null),
			(2,       1,                      null),
			(2,       2,                      null),
			(2,       5,                      null),
			(2,       6,                      null),
			(2,       9,                      null),
			(3,       7,                      '2000-01-01'),
			(3,       8,                      null);
		insert into user_role
			(nip,  role_id, deleted_at) values
			('1c', 1,       null),
			('1c', 2,       null),
			('1c', 4,       null),
			('1d', 3,       null),
			('1d', 2,       '2000-01-01'),
			('1e', 1,       null);
	`
	tests := []struct {
		name             string
		dbData           string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: nip with multiple roles",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					"portal.page1.read",
					"portal.page1.write",
					"portal.page2.read",
					"portal.page2.write"
				]
			}`,
		},
		{
			name:             "ok: nip with empty roles",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1d")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{ "data": [] }`,
		},
		{
			name:             "ok: nip with one roles",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1e")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					"portal.page1.write",
					"portal.page2.write"
				]
			}`,
		},
		{
			name:             "error: invalid token",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/resource-permissions/me", nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_listResourcesAdmin(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into permission
			(id, kode,     nama,     deleted_at) values
			(1,  'write',  'Write',  null),
			(2,  'read',   'Read',   null),
			(3,  'delete', 'Delete', '2000-01-01');
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page2', 'Page 2', null),
			(2,  'portal', 'page1', 'Page 1', null),
			(3,  'portal', 'page3', 'Page 3', null),
			(4,  'portal', 'page4', 'Page 4', '2000-01-01');
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           1,             null),
			(2,  2,           1,             null),
			(3,  3,           2,             null),
			(4,  4,           1,             null),
			(5,  1,           2,             null),
			(6,  2,           2,             null),
			(7,  3,           1,             null),
			(8,  1,           3,             null),
			(9,  1,           2,             '2000-01-01');
	`
	tests := []struct {
		name             string
		dbData           string
		requestHeader    http.Header
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nama": "Page 1",
						"resource_permissions": [
							{
								"id": 6,
								"kode": "portal.page1.read",
								"nama_permission": "Read"
							},
							{
								"id": 2,
								"kode": "portal.page1.write",
								"nama_permission": "Write"
							}
						]
					},
					{
						"nama": "Page 2",
						"resource_permissions": [
							{
								"id": 5,
								"kode": "portal.page2.read",
								"nama_permission": "Read"
							},
							{
								"id": 1,
								"kode": "portal.page2.write",
								"nama_permission": "Write"
							}
						]
					},
					{
						"nama": "Page 3",
						"resource_permissions": [
							{
								"id": 3,
								"kode": "portal.page3.read",
								"nama_permission": "Read"
							},
							{
								"id": 7,
								"kode": "portal.page3.write",
								"nama_permission": "Write"
							}
						]
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
			name:          "ok: with limit offset",
			dbData:        dbData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c", api.RoleAdmin)}},
			requestQuery: url.Values{
				"limit":  []string{"2"},
				"offset": []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nama": "Page 2",
						"resource_permissions": [
							{
								"id": 5,
								"kode": "portal.page2.read",
								"nama_permission": "Read"
							},
							{
								"id": 1,
								"kode": "portal.page2.write",
								"nama_permission": "Write"
							}
						]
					},
					{
						"nama": "Page 3",
						"resource_permissions": [
							{
								"id": 3,
								"kode": "portal.page3.read",
								"nama_permission": "Read"
							},
							{
								"id": 7,
								"kode": "portal.page3.write",
								"nama_permission": "Write"
							}
						]
					}
				],
				"meta": {
					"limit": 2,
					"offset": 1,
					"total": 3
				}
			}`,
		},
		{
			name:             "error: user bukan admin",
			dbData:           dbData,
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "1c")}},
			wantResponseCode: http.StatusForbidden,
			wantResponseBody: `{"message": "akses ditolak"}`,
		},
		{
			name:             "error: invalid token",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/resources", nil)
			req.Header = tt.requestHeader
			req.URL.RawQuery = tt.requestQuery.Encode()
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
