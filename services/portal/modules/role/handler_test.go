package role

import (
	"context"
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
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/migrations"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/db/repository"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/portal/docs"
)

func Test_handler_list(t *testing.T) {
	t.Parallel()

	dbData := `
		insert into role
			(id, nama,      deskripsi,     is_default, deleted_at) values
			(1,  'admin',   'deskripsi 1', false,      null),
			(2,  'pegawai', 'deskripsi 2', true,       null),
			(3,  'guest',   'deskripsi 3', false,      '2000-01-01'),
			(4,  'dokter',  'deskripsi 4', false,      null);
		insert into user_role
			(role_id, nip,  deleted_at) values
			(1,       '1a', null),
			(1,       '1b', null),
			(1,       '1c', null),
			(1,       '1d', '2000-01-01'),
			(1,       '1e', null),
			(2,       '1a', null),
			(3,       '1a', null);
		insert into "user"
			(id,                                     source,     nip,  deleted_at) values
			('00000000-0000-0000-0000-000000000001', 'zimbra',   '1a', null),
			('00000000-0000-0000-0000-000000000001', 'keycloak', '1a', null),
			('00000000-0000-0000-0000-000000000002', 'zimbra',   '1b', null),
			('00000000-0000-0000-0000-000000000003', 'zimbra',   '1c', null),
			('00000000-0000-0000-0000-000000000004', 'zimbra',   '1d', null),
			('00000000-0000-0000-0000-000000000005', 'zimbra',   '1e', '2000-01-01'),
			('00000000-0000-0000-0000-000000000006', 'zimbra',   '1f', null),
			('00000000-0000-0000-0000-000000000007', 'zimbra',   '1g', null),
			('00000000-0000-0000-0000-000000000008', 'zimbra',   '1a', null),
			('00000000-0000-0000-0000-000000000009', 'zimbra',   '1a', '2000-01-01');
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
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 1,
						"nama": "admin",
						"deskripsi": "deskripsi 1",
						"jumlah_user": 3,
						"is_default": false
					},
					{
						"id": 4,
						"nama": "dokter",
						"deskripsi": "deskripsi 4",
						"jumlah_user": 0,
						"is_default": false
					},
					{
						"id": 2,
						"nama": "pegawai",
						"deskripsi": "deskripsi 2",
						"jumlah_user": 6,
						"is_default": true
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
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("1c")}},
			requestQuery: url.Values{
				"limit":  []string{"2"},
				"offset": []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"id": 4,
						"nama": "dokter",
						"deskripsi": "deskripsi 4",
						"jumlah_user": 0,
						"is_default": false
					},
					{
						"id": 2,
						"nama": "pegawai",
						"deskripsi": "deskripsi 2",
						"jumlah_user": 6,
						"is_default": true
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

			req := httptest.NewRequest(http.MethodGet, "/v1/roles", nil)
			req.Header = tt.requestHeader
			req.URL.RawQuery = tt.requestQuery.Encode()
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Read)
			RegisterRoutes(e, db, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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
		insert into role
			(id, nama,      deskripsi,     is_default, deleted_at) values
			(1,  'admin',   'deskripsi 1', false,      null),
			(2,  'pegawai', 'deskripsi 2', true,       null),
			(3,  'guest',   'deskripsi 3', false,      '2000-01-01'),
			(4,  'dokter',  'deskripsi 4', false,      null);
		insert into user_role
			(role_id, nip,  deleted_at) values
			(1,       '1a', null),
			(1,       '1b', null),
			(1,       '1c', null),
			(1,       '1d', '2000-01-01'),
			(1,       '1e', null),
			(2,       '1a', null),
			(3,       '1a', null);
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page1', 'Page 1', null),
			(2,  'portal', 'page2', 'Page 2', null),
			(3,  'portal', 'page3', 'Page 3', '2000-01-01'),
			(4,  'portal', 'page4', 'Page 4', null);
		insert into permission
			(id, kode,    nama,    deleted_at) values
			(1,  'read',  'Read',  null),
			(2,  'write', 'Write', null),
			(3,  'del',   'Del',   '2000-01-01');
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           2,             null),
			(2,  2,           1,             null),
			(3,  4,           1,             null),
			(4,  2,           2,             null),
			(5,  1,           1,             null),
			(6,  1,           3,             null),
			(7,  2,           3,             null),
			(8,  3,           1,             null),
			(9,  4,           2,             '2000-01-01');
		insert into role_resource_permission
			(role_id, resource_permission_id, deleted_at) values
			(1,       1,                      null),
			(1,       2,                      null),
			(1,       3,                      null),
			(1,       4,                      null),
			(1,       5,                      null),
			(1,       6,                      null),
			(1,       7,                      null),
			(1,       8,                      null),
			(1,       9,                      null),
			(2,       5,                      null),
			(2,       2,                      '2000-01-01'),
			(2,       3,                      null),
			(2,       8,                      null),
			(3,       1,                      null);
		insert into "user"
			(id,                                     source,     nip,  deleted_at) values
			('00000000-0000-0000-0000-000000000001', 'zimbra',   '1a', null),
			('00000000-0000-0000-0000-000000000001', 'keycloak', '1a', null),
			('00000000-0000-0000-0000-000000000002', 'zimbra',   '1b', null),
			('00000000-0000-0000-0000-000000000003', 'zimbra',   '1c', null),
			('00000000-0000-0000-0000-000000000004', 'zimbra',   '1d', null),
			('00000000-0000-0000-0000-000000000005', 'zimbra',   '1e', '2000-01-01'),
			('00000000-0000-0000-0000-000000000006', 'zimbra',   '1f', null),
			('00000000-0000-0000-0000-000000000007', 'zimbra',   '1g', null),
			('00000000-0000-0000-0000-000000000008', 'zimbra',   '1a', null),
			('00000000-0000-0000-0000-000000000009', 'zimbra',   '1a', '2000-01-01');
	`
	tests := []struct {
		name             string
		dbData           string
		paramID          string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: with role resource permissions",
			dbData:           dbData,
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 1,
					"nama": "admin",
					"deskripsi": "deskripsi 1",
					"jumlah_user": 3,
					"is_default": false,
					"resource_permissions": [
						{
							"id": 5,
							"kode": "portal.page1.read",
							"nama_resource": "Page 1",
							"nama_permission": "Read"
						},
						{
							"id": 1,
							"kode": "portal.page1.write",
							"nama_resource": "Page 1",
							"nama_permission": "Write"
						},
						{
							"id": 2,
							"kode": "portal.page2.read",
							"nama_resource": "Page 2",
							"nama_permission": "Read"
						},
						{
							"id": 4,
							"kode": "portal.page2.write",
							"nama_resource": "Page 2",
							"nama_permission": "Write"
						},
						{
							"id": 3,
							"kode": "portal.page4.read",
							"nama_resource": "Page 4",
							"nama_permission": "Read"
						}
					]
				}
			}`,
		},
		{
			name:             "ok: default role and without deleted role resource permissions",
			dbData:           dbData,
			paramID:          "2",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 2,
					"nama": "pegawai",
					"deskripsi": "deskripsi 2",
					"jumlah_user": 6,
					"is_default": true,
					"resource_permissions": [
						{
							"id": 5,
							"kode": "portal.page1.read",
							"nama_resource": "Page 1",
							"nama_permission": "Read"
						},
						{
							"id": 3,
							"kode": "portal.page4.read",
							"nama_resource": "Page 4",
							"nama_permission": "Read"
						}
					]
				}
			}`,
		},
		{
			name:             "ok: without role resource permissions",
			dbData:           dbData,
			paramID:          "4",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"id": 4,
					"nama": "dokter",
					"deskripsi": "deskripsi 4",
					"jumlah_user": 0,
					"is_default": false,
					"resource_permissions": []
				}
			}`,
		},
		{
			name:             "error: roles deleted",
			dbData:           dbData,
			paramID:          "3",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: invalid param id",
			dbData:           dbData,
			paramID:          "1a",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai"}`,
		},
		{
			name:             "error: invalid token",
			dbData:           dbData,
			paramID:          "1",
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

			req := httptest.NewRequest(http.MethodGet, "/v1/roles/"+tt.paramID, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Read)
			RegisterRoutes(e, db, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_create(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page1', 'Page 1', null),
			(2,  'portal', 'page2', 'Page 2', null),
			(3,  'portal', 'page3', 'Page 3', '2000-01-01');
		insert into permission
			(id, kode,    nama,    deleted_at) values
			(1,  'read',  'Read',  null),
			(2,  'write', 'Write', null),
			(3,  'del',   'Del',   '2000-01-01');
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           1,             null),
			(2,  1,           2,             null),
			(3,  2,           1,             null),
			(4,  2,           2,             null),
			(5,  3,           1,             null),
			(6,  1,           3,             null),
			(7,  1,           1,             '2000-01-01');
	`
	tests := []struct {
		name                          string
		dbData                        string
		requestHeader                 http.Header
		requestBody                   string
		wantResponseCode              int
		wantResponseBody              string
		wantDBRoles                   dbtest.Rows
		wantDBRoleResourcePermissions dbtest.Rows
	}{
		{
			name:          "ok: with resource permissions",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Super Admin",
				"deskripsi": "Deskripsi 1",
				"is_default": false,
				"resource_permission_ids": [ 1, 2, 3, 4 ]
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 1
				}
			}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "Super Admin",
					"deskripsi":  "Deskripsi 1",
					"is_default": false,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(3),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(4),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
			},
		},
		{
			name: "ok: with existings data in database",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at) values
					(1,       1,                      '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01');
			`,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Pegawai",
				"deskripsi": "deskripsi",
				"is_default": true,
				"resource_permission_ids": [ 3 ]
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 2
				}
			}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int16(2),
					"service":    nil,
					"nama":       "Pegawai",
					"deskripsi":  "deskripsi",
					"is_default": true,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(3),
					"role_id":                int16(2),
					"resource_permission_id": int32(3),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
			},
		},
		{
			name:          "ok: required params only and without resource permissions",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Admin Kepegawaian"
			}`,
			wantResponseCode: http.StatusCreated,
			wantResponseBody: `{
				"data": {
					"id": 1
				}
			}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "Admin Kepegawaian",
					"deskripsi":  "",
					"is_default": false,
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: deleted resource",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Admin",
				"resource_permission_ids": [5]
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: deleted permission",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Admin",
				"resource_permission_ids": [6]
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: deleted resource permission",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Admin",
				"resource_permission_ids": [7]
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: resource permission not exists",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Admin",
				"resource_permission_ids": [8]
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: active, deleted and not exists resource permission",
			dbData:        seedData,
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Admin",
				"deskripsi": "desc",
				"is_default": false,
				"resource_permission_ids": [1,2,3,4,5,6,7,8]
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:             "error: have additional and missing required params, duplicate resource permission ids",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:      `{"id": 1, "resource_permission_ids": [ 1,2,3,2 ]}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung` +
				` | parameter \"resource_permission_ids\" item tidak boleh duplikat` +
				` | parameter \"nama\" harus diisi"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: nama is empty string",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:                   `{"nama": ""}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "parameter \"nama\" harus 1 karakter atau lebih"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: nama & deskripsi exceed character length",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 101) + `",
				"deskripsi": "` + strings.Repeat(".", 256) + `"
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "parameter \"deskripsi\" harus 255 karakter atau kurang | parameter \"nama\" harus 100 karakter atau kurang"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: body is empty",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "request body harus diisi"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: invalid token",
			requestHeader:                 http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:              http.StatusUnauthorized,
			wantResponseBody:              `{"message": "token otentikasi tidak valid"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/v1/roles", strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Write)
			RegisterRoutes(e, db, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRoles, err := dbtest.QueryAll(db, "role", "id")
			require.NoError(t, err)
			if len(tt.wantDBRoles) == len(actualRoles) {
				for i, row := range actualRoles {
					if tt.wantDBRoles[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRoles[i]["created_at"] = row["created_at"]
						tt.wantDBRoles[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRoles, actualRoles)

			actualRoleResourcePermissions, err := dbtest.QueryAll(db, "role_resource_permission", "role_id, resource_permission_id")
			require.NoError(t, err)
			if len(tt.wantDBRoleResourcePermissions) == len(actualRoleResourcePermissions) {
				for i, row := range actualRoleResourcePermissions {
					if tt.wantDBRoleResourcePermissions[i]["id"] == "{id}" {
						tt.wantDBRoleResourcePermissions[i]["id"] = row["id"]
					}
					if tt.wantDBRoleResourcePermissions[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRoleResourcePermissions[i]["created_at"] = row["created_at"]
						tt.wantDBRoleResourcePermissions[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRoleResourcePermissions, actualRoleResourcePermissions)
		})
	}
}

func Test_handler_update(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into resource
			(id, service,  kode,    nama,     deleted_at) values
			(1,  'portal', 'page1', 'Page 1', null),
			(2,  'portal', 'page2', 'Page 2', null),
			(3,  'portal', 'page3', 'Page 3', '2000-01-01');
		insert into permission
			(id, kode,    nama,    deleted_at) values
			(1,  'read',  'Read',  null),
			(2,  'write', 'Write', null),
			(3,  'del',   'Del',   '2000-01-01');
		insert into resource_permission
			(id, resource_id, permission_id, deleted_at) values
			(1,  1,           1,             null),
			(2,  1,           2,             null),
			(3,  2,           1,             null),
			(4,  2,           2,             null),
			(5,  3,           1,             null),
			(6,  1,           3,             null),
			(7,  1,           1,             '2000-01-01');
	`
	tests := []struct {
		name                          string
		dbData                        string
		paramID                       string
		requestHeader                 http.Header
		requestBody                   string
		wantResponseCode              int
		wantResponseBody              string
		wantDBRoles                   dbtest.Rows
		wantDBRoleResourcePermissions dbtest.Rows
	}{
		{
			name: "ok: with all params",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "super_admin",
				"deskripsi": "new deskripsi",
				"is_default": true,
				"resource_permission_ids": [ 1, 2 ]
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "super_admin",
					"deskripsi":  "new deskripsi",
					"is_default": true,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
			},
		},
		{
			name: "ok: with empty resource permissions",
			dbData: seedData + `
				insert into role
					(nama,    deskripsi, service, is_default, created_at,   updated_at) values
					('admin', 'desc',    'svc',   true,       '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at,   deleted_at) values
					(1,       1,                      '2000-01-01', '2000-01-01', null),
					(1,       2,                      '2000-01-01', '2000-01-01', null),
					(1,       2,                      '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "new_admin",
				"deskripsi": "new desc",
				"is_default": false,
				"resource_permission_ids": []
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    "svc",
					"nama":       "new_admin",
					"deskripsi":  "new desc",
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             "{deleted_at}",
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             "{deleted_at}",
				},
				{
					"id":                     int32(3),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "ok: partial update role and create & delete resource permissions",
			dbData: seedData + `
				insert into role
					(nama,      deskripsi, is_default, created_at,   updated_at) values
					('admin',   'desc',    false,      '2000-01-01', '2000-01-01'),
					('pegawai', 'desc',    true,       '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at) values
					(1,       1,                      '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01'),
					(2,       1,                      '2000-01-01', '2000-01-01'),
					(2,       2,                      '2000-01-01', '2000-01-01'),
					(2,       4,                      '2000-01-01', '2000-01-01');
			`,
			paramID:       "2",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"deskripsi": "new desc",
				"resource_permission_ids": [ 1, 3 ]
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  "desc",
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int16(2),
					"service":    nil,
					"nama":       "pegawai",
					"deskripsi":  "new desc",
					"is_default": true,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(3),
					"role_id":                int16(2),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(4),
					"role_id":                int16(2),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             "{deleted_at}",
				},
				{
					"id":                     int32(6),
					"role_id":                int16(2),
					"resource_permission_id": int32(3),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     int32(5),
					"role_id":                int16(2),
					"resource_permission_id": int32(4),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             "{deleted_at}",
				},
			},
		},
		{
			name: "ok: without updating any data",
			dbData: seedData + `
				insert into role
					(nama,      deskripsi, is_default, created_at,   updated_at) values
					('admin',   'desc',    false,      '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at) values
					(1,       1,                      '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01'),
					(1,       3,                      '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "admin",
				"deskripsi": "desc",
				"is_default": false,
				"resource_permission_ids": [ 1, 2, 3 ]
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  "desc",
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(3),
					"role_id":                int16(1),
					"resource_permission_id": int32(3),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
			},
		},
		{
			name: "ok: success create resource permission that previously being deleted",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at,   deleted_at) values
					(1,       1,                      '2000-01-01', '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01', '2000-01-01'),
					(1,       3,                      '2000-01-01', '2000-01-01', '2000-01-01'),
					(1,       4,                      '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"resource_permission_ids": [ 1, 4 ]
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":                     int32(3),
					"role_id":                int16(1),
					"resource_permission_id": int32(3),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":                     int32(4),
					"role_id":                int16(1),
					"resource_permission_id": int32(4),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":                     "{id}",
					"role_id":                int16(1),
					"resource_permission_id": int32(4),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
			},
		},
		{
			name: "ok: without resource permissions",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at) values
					(1,       1,                      '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "super_admin",
				"deskripsi": "new desc",
				"is_default": false
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "super_admin",
					"deskripsi":  "new desc",
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
			},
		},
		{
			name: "ok: increment role_resource_permission.id should be consistent",
			dbData: `
				insert into resource
					(id, service,  kode,    nama) values
					(1,  'portal', 'page2', 'Page 2'),
					(2,  'portal', 'page1', 'Page 1'),
					(3,  'portal', 'page3', 'Page 3');
				insert into permission
					(id, kode,    nama) values
					(1,  'read',  'Read'),
					(2,  'write', 'Write'),
					(3,  'del',   'Del');
				insert into resource_permission
					(id, resource_id, permission_id) values
					(1,  1,           1),
					(2,  1,           2),
					(3,  2,           2),
					(4,  2,           1),
					(5,  3,           1),
					(6,  3,           2),
					(7,  3,           3),
					(8,  2,           3),
					(9,  1,           3);
				insert into role
					(nama,    is_default, created_at,   updated_at) values
					('admin', true,       '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at) values
					(1,       1,                      '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01'),
					(1,       3,                      '2000-01-01', '2000-01-01'),
					(1,       4,                      '2000-01-01', '2000-01-01'),
					(1,       5,                      '2000-01-01', '2000-01-01'),
					(1,       9,                      '2000-01-01', '2000-01-01'),
					(1,       7,                      '2000-01-01', '2000-01-01'),
					(1,       8,                      '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"is_default": true,
				"resource_permission_ids": [ 1, 2, 3, 4, 5, 6, 7, 8, 9 ]
			}`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": true,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(3),
					"role_id":                int16(1),
					"resource_permission_id": int32(3),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(4),
					"role_id":                int16(1),
					"resource_permission_id": int32(4),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(5),
					"role_id":                int16(1),
					"resource_permission_id": int32(5),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(9),
					"role_id":                int16(1),
					"resource_permission_id": int32(6),
					"created_at":             "{created_at}",
					"updated_at":             "{updated_at}",
					"deleted_at":             nil,
				},
				{
					"id":                     int32(7),
					"role_id":                int16(1),
					"resource_permission_id": int32(7),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(8),
					"role_id":                int16(1),
					"resource_permission_id": int32(8),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(6),
					"role_id":                int16(1),
					"resource_permission_id": int32(9),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
			},
		},
		{
			name: "error: role not found",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
			`,
			paramID:          "0",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:      `{"nama": "super_admin"}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: deleted resource",
			dbData:                        seedData,
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:                   `{"resource_permission_ids": [5]}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: deleted permission",
			dbData:                        seedData,
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:                   `{"resource_permission_ids": [6]}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: deleted resource permission",
			dbData:                        seedData,
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:                   `{"resource_permission_ids": [7]}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: resource permission not exists",
			dbData:                        seedData,
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:                   `{"resource_permission_ids": [8]}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name: "error: active, deleted and not exists resource permission",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "Pegawai",
				"deskripsi": "desc",
				"is_default": true,
				"resource_permission_ids": [1,2,3,4,5,6,7,8]
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data resource permission tidak ditemukan"}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name: "error: role is deleted",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at,   deleted_at) values
					('admin', '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "super_admin",
				"deskripsi": "desc",
				"is_default": true,
				"resource_permission_ids": [ 1, 2 ]
			}`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name: "error: don't allow empty json",
			dbData: seedData + `
				insert into role
					(nama,    created_at,   updated_at) values
					('admin', '2000-01-01', '2000-01-01');
				insert into role_resource_permission
					(role_id, resource_permission_id, created_at,   updated_at) values
					(1,       1,                      '2000-01-01', '2000-01-01'),
					(1,       2,                      '2000-01-01', '2000-01-01');
			`,
			paramID:          "1",
			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:      `{}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus 1 property atau lebih"}`,
			wantDBRoles: dbtest.Rows{
				{
					"id":         int16(1),
					"service":    nil,
					"nama":       "admin",
					"deskripsi":  nil,
					"is_default": false,
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
			wantDBRoleResourcePermissions: dbtest.Rows{
				{
					"id":                     int32(1),
					"role_id":                int16(1),
					"resource_permission_id": int32(1),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
				{
					"id":                     int32(2),
					"role_id":                int16(1),
					"resource_permission_id": int32(2),
					"created_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at":             time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at":             nil,
				},
			},
		},
		{
			name:          "error: invalid id, have additional params, and have null values",
			paramID:       "1a",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"role_id": 1,
				"nama": null,
				"deskripsi": null,
				"is_default": null,
				"resource_permission_ids": null
			}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" harus dalam format yang sesuai` +
				` | parameter \"deskripsi\" tidak boleh null` +
				` | parameter \"is_default\" tidak boleh null` +
				` | parameter \"nama\" tidak boleh null` +
				` | parameter \"resource_permission_ids\" tidak boleh null` +
				` | parameter \"role_id\" tidak didukung"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: nama is empty string and duplicate resource permission ids",
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody:                   `{"nama": "", "resource_permission_ids": [ 1,1 ]}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "parameter \"nama\" harus 1 karakter atau lebih | parameter \"resource_permission_ids\" item tidak boleh duplikat"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:          "error: nama & deskripsi exceed character length",
			paramID:       "1",
			requestHeader: http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			requestBody: `{
				"nama": "` + strings.Repeat(".", 101) + `",
				"deskripsi": "` + strings.Repeat(".", 256) + `"
			}`,
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "parameter \"deskripsi\" harus 255 karakter atau kurang | parameter \"nama\" harus 100 karakter atau kurang"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: body is empty",
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{apitest.GenerateAuthHeader("2a")}},
			wantResponseCode:              http.StatusBadRequest,
			wantResponseBody:              `{"message": "request body harus diisi"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
		{
			name:                          "error: invalid token",
			paramID:                       "1",
			requestHeader:                 http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode:              http.StatusUnauthorized,
			wantResponseBody:              `{"message": "token otentikasi tidak valid"}`,
			wantDBRoles:                   dbtest.Rows{},
			wantDBRoleResourcePermissions: dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/v1/roles/"+tt.paramID, strings.NewReader(tt.requestBody))
			req.Header = tt.requestHeader
			req.Header.Set("Content-Type", "application/json")
			rec := httptest.NewRecorder()

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Write)
			RegisterRoutes(e, db, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, typeutil.Coalesce(rec.Body.String(), "null"))
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))

			actualRoles, err := dbtest.QueryAll(db, "role", "id")
			require.NoError(t, err)
			if len(tt.wantDBRoles) == len(actualRoles) {
				for i, row := range actualRoles {
					if tt.wantDBRoles[i]["updated_at"] == "{updated_at}" {
						assert.WithinDuration(t, time.Now(), row["updated_at"].(time.Time), 10*time.Second)
						tt.wantDBRoles[i]["updated_at"] = row["updated_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRoles, actualRoles)

			actualRoleResourcePermissions, err := dbtest.QueryAll(db, "role_resource_permission", "role_id, resource_permission_id, id")
			require.NoError(t, err)
			if len(tt.wantDBRoleResourcePermissions) == len(actualRoleResourcePermissions) {
				for i, row := range actualRoleResourcePermissions {
					if tt.wantDBRoleResourcePermissions[i]["id"] == "{id}" {
						tt.wantDBRoleResourcePermissions[i]["id"] = row["id"]
					}
					if tt.wantDBRoleResourcePermissions[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBRoleResourcePermissions[i]["created_at"] = row["created_at"]
						tt.wantDBRoleResourcePermissions[i]["updated_at"] = row["updated_at"]
					}
					if tt.wantDBRoleResourcePermissions[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBRoleResourcePermissions[i]["deleted_at"] = row["deleted_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBRoleResourcePermissions, actualRoleResourcePermissions)
		})
	}
}
