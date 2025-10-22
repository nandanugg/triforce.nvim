package user

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
			(id, nama,      is_default, is_aktif, deleted_at) values
			(1,  'admin',   false,      true,     null),
			(2,  'pegawai', true,       true,     null),
			(3,  'guest',   false,      true,     '2000-01-01'),
			(4,  'dokter',  false,      false,    null);
		insert into user_role
			(nip,  role_id, deleted_at) values
			('1a', 1,       null),
			('1a', 2,       null),
			('1a', 3,       null),
			('1a', 4,       null),
			('1b', 1,       '2000-01-01'),
			('1b', 2,       '2000-01-01'),
			('1b', 4,       null),
			('1c', 1,       null);
		insert into "user"
			(nip,  nama,       email,               last_login_at, id,                                     source,     deleted_at) values
			('1a', 'John Doe', 'john.doe@test.com', '2010-01-01',  '00000000-0000-0000-0000-000000000001', 'zimbra',   null),
			('1a', 'Jane Doe', 'jane.doe@test.com', '2012-01-01',  '00000000-0000-0000-0000-000000000001', 'keycloak', null),
			('1a', 'Will Doe', 'will.doe@test.com', null,          '00000000-0000-0000-0000-000000000002', 'zimbra',   null),
			('1a', 'Wish Doe', 'wish.doe@test.com', '2011-01-01',  '00000000-0000-0000-0000-000000000003', 'zimbra',   null),
			('1a', null,       null,                '2009-01-01',  '00000000-0000-0000-0000-000000000004', 'zimbra',   null),
			('1b', 'Will',     'will@test.com',     '2000-01-01',  '00000000-0000-0000-0000-000000000005', 'zimbra',   null),
			('1b', 'Joy',      'joy@test.com',      '2001-01-01',  '00000000-0000-0000-0000-000000000006', 'zimbra',   '2000-01-01'),
			('1b', '',         '',                  null,          '00000000-0000-0000-0000-000000000007', 'zimbra',   null),
			('1c', null,       null,                null,          '00000000-0000-0000-0000-000000000008', 'zimbra',   '2000-01-01'),
			('1d', 'Carl',     'carl@test.com',     null,          '00000000-0000-0000-0000-000000000009', 'zimbra',   null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Read)
	RegisterRoutes(e, db, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		requestHeader    http.Header
		requestQuery     url.Values
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: without filter and pagination",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nip": "1a",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "keycloak",
								"nama": "Jane Doe",
								"email": "jane.doe@test.com",
								"last_login_at": "2012-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000003",
								"source": "zimbra",
								"nama": "Wish Doe",
								"email": "wish.doe@test.com",
								"last_login_at": "2011-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "zimbra",
								"nama": "John Doe",
								"email": "john.doe@test.com",
								"last_login_at": "2010-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000004",
								"source": "zimbra",
								"nama": null,
								"email": null,
								"last_login_at": "2009-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000002",
								"source": "zimbra",
								"nama": "Will Doe",
								"email": "will.doe@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 1,
								"nama": "admin",
								"is_default": false,
								"is_aktif": true
							},
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
					},
					{
						"nip": "1b",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000005",
								"source": "zimbra",
								"nama": "Will",
								"email": "will@test.com",
								"last_login_at": "2000-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000007",
								"source": "zimbra",
								"nama": "",
								"email": "",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
					},
					{
						"nip": "1d",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000009",
								"source": "zimbra",
								"nama": "Carl",
								"email": "carl@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
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
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"limit":  []string{"2"},
				"offset": []string{"1"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nip": "1b",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000005",
								"source": "zimbra",
								"nama": "Will",
								"email": "will@test.com",
								"last_login_at": "2000-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000007",
								"source": "zimbra",
								"nama": "",
								"email": "",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
					},
					{
						"nip": "1d",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000009",
								"source": "zimbra",
								"nama": "Carl",
								"email": "carl@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
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
			name:          "ok: with filter nip",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"nip": []string{"1a"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nip": "1a",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "keycloak",
								"nama": "Jane Doe",
								"email": "jane.doe@test.com",
								"last_login_at": "2012-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000003",
								"source": "zimbra",
								"nama": "Wish Doe",
								"email": "wish.doe@test.com",
								"last_login_at": "2011-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "zimbra",
								"nama": "John Doe",
								"email": "john.doe@test.com",
								"last_login_at": "2010-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000004",
								"source": "zimbra",
								"nama": null,
								"email": null,
								"last_login_at": "2009-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000002",
								"source": "zimbra",
								"nama": "Will Doe",
								"email": "will.doe@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 1,
								"nama": "admin",
								"is_default": false,
								"is_aktif": true
							},
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
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
			name:          "ok: with filter default role",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"role_id": []string{"2"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nip": "1a",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "keycloak",
								"nama": "Jane Doe",
								"email": "jane.doe@test.com",
								"last_login_at": "2012-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000003",
								"source": "zimbra",
								"nama": "Wish Doe",
								"email": "wish.doe@test.com",
								"last_login_at": "2011-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "zimbra",
								"nama": "John Doe",
								"email": "john.doe@test.com",
								"last_login_at": "2010-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000004",
								"source": "zimbra",
								"nama": null,
								"email": null,
								"last_login_at": "2009-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000002",
								"source": "zimbra",
								"nama": "Will Doe",
								"email": "will.doe@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 1,
								"nama": "admin",
								"is_default": false,
								"is_aktif": true
							},
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
					},
					{
						"nip": "1b",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000005",
								"source": "zimbra",
								"nama": "Will",
								"email": "will@test.com",
								"last_login_at": "2000-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000007",
								"source": "zimbra",
								"nama": "",
								"email": "",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
					},
					{
						"nip": "1d",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000009",
								"source": "zimbra",
								"nama": "Carl",
								"email": "carl@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
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
			name:          "ok: with filter role",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"role_id": []string{"4"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nip": "1a",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "keycloak",
								"nama": "Jane Doe",
								"email": "jane.doe@test.com",
								"last_login_at": "2012-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000003",
								"source": "zimbra",
								"nama": "Wish Doe",
								"email": "wish.doe@test.com",
								"last_login_at": "2011-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000001",
								"source": "zimbra",
								"nama": "John Doe",
								"email": "john.doe@test.com",
								"last_login_at": "2010-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000004",
								"source": "zimbra",
								"nama": null,
								"email": null,
								"last_login_at": "2009-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000002",
								"source": "zimbra",
								"nama": "Will Doe",
								"email": "will.doe@test.com",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 1,
								"nama": "admin",
								"is_default": false,
								"is_aktif": true
							},
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
					},
					{
						"nip": "1b",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000005",
								"source": "zimbra",
								"nama": "Will",
								"email": "will@test.com",
								"last_login_at": "2000-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000007",
								"source": "zimbra",
								"nama": "",
								"email": "",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
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
			name:          "ok: with filter nip and role",
			requestHeader: http.Header{"Authorization": authHeader},
			requestQuery: url.Values{
				"nip":     []string{"1b"},
				"role_id": []string{"4"},
			},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": [
					{
						"nip": "1b",
						"profiles": [
							{
								"id": "00000000-0000-0000-0000-000000000005",
								"source": "zimbra",
								"nama": "Will",
								"email": "will@test.com",
								"last_login_at": "2000-01-01T00:00:00Z"
							},
							{
								"id": "00000000-0000-0000-0000-000000000007",
								"source": "zimbra",
								"nama": "",
								"email": "",
								"last_login_at": null
							}
						],
						"roles": [
							{
								"id": 4,
								"nama": "dokter",
								"is_default": false,
								"is_aktif": false
							},
							{
								"id": 2,
								"nama": "pegawai",
								"is_default": true,
								"is_aktif": true
							}
						]
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
			name:             "error: invalid token",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/users", nil)
			req.Header = tt.requestHeader
			req.URL.RawQuery = tt.requestQuery.Encode()
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
		insert into role
			(id, nama,      is_default, is_aktif, deleted_at) values
			(1,  'admin',   false,      true,     null),
			(2,  'pegawai', true,       true,     null),
			(3,  'guest',   false,      true,     '2000-01-01'),
			(4,  'dokter',  false,      false,    null);
		insert into user_role
			(nip,  role_id, deleted_at) values
			('1a', 1,       null),
			('1a', 2,       null),
			('1a', 3,       null),
			('1a', 4,       null),
			('1b', 1,       '2000-01-01'),
			('1b', 2,       '2000-01-01'),
			('1b', 4,       null),
			('1c', 1,       null);
		insert into "user"
			(nip,  nama,       email,               last_login_at, id,                                     source,     deleted_at) values
			('1a', 'John Doe', 'john.doe@test.com', '2010-01-01',  '00000000-0000-0000-0000-000000000001', 'zimbra',   null),
			('1a', 'Jane Doe', 'jane.doe@test.com', '2012-01-01',  '00000000-0000-0000-0000-000000000001', 'keycloak', null),
			('1a', 'Will Doe', 'will.doe@test.com', null,          '00000000-0000-0000-0000-000000000002', 'zimbra',   null),
			('1a', 'Wish Doe', 'wish.doe@test.com', '2011-01-01',  '00000000-0000-0000-0000-000000000003', 'zimbra',   null),
			('1a', null,       null,                '2009-01-01',  '00000000-0000-0000-0000-000000000004', 'zimbra',   null),
			('1b', 'Will',     'will@test.com',     '2000-01-01',  '00000000-0000-0000-0000-000000000005', 'zimbra',   null),
			('1b', 'Joy',      'joy@test.com',      '2001-01-01',  '00000000-0000-0000-0000-000000000006', 'zimbra',   '2000-01-01'),
			('1b', '',         '',                  null,          '00000000-0000-0000-0000-000000000007', 'zimbra',   null),
			('1c', null,       null,                null,          '00000000-0000-0000-0000-000000000008', 'zimbra',   '2000-01-01'),
			('1d', 'Carl',     'carl@test.com',     null,          '00000000-0000-0000-0000-000000000009', 'zimbra',   null);
	`
	db := dbtest.New(t, dbmigrations.FS)
	_, err := db.Exec(context.Background(), dbData)
	require.NoError(t, err)

	e, err := api.NewEchoServer(docs.OpenAPIBytes)
	require.NoError(t, err)

	authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Read)
	RegisterRoutes(e, db, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		paramNIP         string
		requestHeader    http.Header
		wantResponseCode int
		wantResponseBody string
	}{
		{
			name:             "ok: user found with roles",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"nip": "1a",
					"profiles": [
						{
							"id": "00000000-0000-0000-0000-000000000001",
							"source": "keycloak",
							"nama": "Jane Doe",
							"email": "jane.doe@test.com",
							"last_login_at": "2012-01-01T00:00:00Z"
						},
						{
							"id": "00000000-0000-0000-0000-000000000003",
							"source": "zimbra",
							"nama": "Wish Doe",
							"email": "wish.doe@test.com",
							"last_login_at": "2011-01-01T00:00:00Z"
						},
						{
							"id": "00000000-0000-0000-0000-000000000001",
							"source": "zimbra",
							"nama": "John Doe",
							"email": "john.doe@test.com",
							"last_login_at": "2010-01-01T00:00:00Z"
						},
						{
							"id": "00000000-0000-0000-0000-000000000004",
							"source": "zimbra",
							"nama": null,
							"email": null,
							"last_login_at": "2009-01-01T00:00:00Z"
						},
						{
							"id": "00000000-0000-0000-0000-000000000002",
							"source": "zimbra",
							"nama": "Will Doe",
							"email": "will.doe@test.com",
							"last_login_at": null
						}
					],
					"roles": [
						{
							"id": 1,
							"nama": "admin",
							"is_default": false,
							"is_aktif": true
						},
						{
							"id": 4,
							"nama": "dokter",
							"is_default": false,
							"is_aktif": false
						},
						{
							"id": 2,
							"nama": "pegawai",
							"is_default": true,
							"is_aktif": true
						}
					]
				}
			}`,
		},
		{
			name:             "ok: user found without deleted roles",
			paramNIP:         "1b",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"nip": "1b",
					"profiles": [
						{
							"id": "00000000-0000-0000-0000-000000000005",
							"source": "zimbra",
							"nama": "Will",
							"email": "will@test.com",
							"last_login_at": "2000-01-01T00:00:00Z"
						},
						{
							"id": "00000000-0000-0000-0000-000000000007",
							"source": "zimbra",
							"nama": "",
							"email": "",
							"last_login_at": null
						}
					],
					"roles": [
						{
							"id": 4,
							"nama": "dokter",
							"is_default": false,
							"is_aktif": false
						},
						{
							"id": 2,
							"nama": "pegawai",
							"is_default": true,
							"is_aktif": true
						}
					]
				}
			}`,
		},
		{
			name:             "ok: user found with default roles only",
			paramNIP:         "1d",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusOK,
			wantResponseBody: `{
				"data": {
					"nip": "1d",
					"profiles": [
						{
							"id": "00000000-0000-0000-0000-000000000009",
							"source": "zimbra",
							"nama": "Carl",
							"email": "carl@test.com",
							"last_login_at": null
						}
					],
					"roles": [
						{
							"id": 2,
							"nama": "pegawai",
							"is_default": true,
							"is_aktif": true
						}
					]
				}
			}`,
		},
		{
			name:             "error: user deleted",
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: user not exists",
			paramNIP:         "1z",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(http.MethodGet, "/v1/users/"+tt.paramNIP, nil)
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}

func Test_handler_update(t *testing.T) {
	t.Parallel()

	seedData := `
		insert into role
			(id, nama,        is_default, is_aktif, deleted_at) values
			(1,  'admin',     false,      true,     null),
			(2,  'pegawai',   true,       true,     null),
			(3,  'delete1',   false,      true,     '2000-01-01'),
			(4,  'delete2',   true,       true,     '2000-01-01'),
			(5,  'inactive1', true,       false,    null),
			(6,  'inactive2', false,      false,    null);
	`

	authHeader := []string{apitest.GenerateAuthHeader("2a")}
	tests := []struct {
		name             string
		dbData           string
		paramNIP         string
		requestHeader    http.Header
		requestBody      string
		wantResponseCode int
		wantResponseBody string
		wantDBUserRoles  dbtest.Rows
	}{
		{
			name: "ok: only create non default role",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1a'),
					('00000000-0000-0000-0000-000000000001', 'zimbre', '1b');
			`,
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 2, 5, 6 ] }`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         "{id}",
					"nip":        "1a",
					"role_id":    int16(1),
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
				{
					"id":         "{id}",
					"nip":        "1a",
					"role_id":    int16(6),
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name: "ok: default roles behavior",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1a');
				insert into user_role
					(nip,  role_id, created_at,   updated_at) values
					('1a',  2,       '2000-01-01', '2000-01-01'),
					('1a',  5,       '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 2 ] }`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1a",
					"role_id":    int16(2),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int32(2),
					"nip":        "1a",
					"role_id":    int16(5),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": "{deleted_at}",
				},
			},
		},
		{
			name: "ok: create & delete user role",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1');
				insert into user_role
					(nip,  role_id, created_at,   updated_at,   deleted_at) values
					('1',  6,       '2000-01-01', '2000-01-01', null),
					('1',  6,       '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1 ] }`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(3),
					"nip":        "1",
					"role_id":    int16(1),
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
				{
					"id":         int32(1),
					"nip":        "1",
					"role_id":    int16(6),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": "{deleted_at}",
				},
				{
					"id":         int32(2),
					"nip":        "1",
					"role_id":    int16(6),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
			},
		},
		{
			name: "ok: without updating any data",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1'),
					('00000000-0000-0000-0000-000000000001', 'zimbre', '1');
				insert into user_role
					(nip,  role_id, created_at,   updated_at,   deleted_at) values
					('1',  1,       '2000-01-01', '2000-01-01', null),
					('1',  2,       '2000-01-01', '2000-01-01', '2000-01-01'),
					('1',  6,       '2000-01-01', '2000-01-01', null);
			`,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 6 ] }`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int32(2),
					"nip":        "1",
					"role_id":    int16(2),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":         int32(3),
					"nip":        "1",
					"role_id":    int16(6),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name: "ok: success create user role that previously being deleted",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip,  deleted_at) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1c', null),
					('00000000-0000-0000-0000-000000000001', 'zimbre', '1c', '2000-01-01');
				insert into user_role
					(nip,  role_id, created_at,   updated_at,   deleted_at) values
					('1c', 1,       '2000-01-01', '2000-01-01', null),
					('1c', 2,       '2000-01-01', '2000-01-01', null),
					('1c', 6,       '2000-01-01', '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 6 ] }`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1c",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int32(2),
					"nip":        "1c",
					"role_id":    int16(2),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": "{deleted_at}",
				},
				{
					"id":         int32(3),
					"nip":        "1c",
					"role_id":    int16(6),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
				},
				{
					"id":         int32(4),
					"nip":        "1c",
					"role_id":    int16(6),
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
			},
		},
		{
			name: "ok: increment user_role.id should be consistent",
			dbData: `
				insert into role
					(id, nama,    is_default, is_aktif) values
					(1,  'role1', false,      true),
					(2,  'role2', false,      true),
					(3,  'role3', true,       true),
					(4,  'role4', true,       false),
					(5,  'role5', false,      false),
					(6,  'role6', false,      false);
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1c');
				insert into user_role
					(nip,  role_id, created_at,   updated_at) values
					('1c', 1,       '2000-01-01', '2000-01-01'),
					('1c', 2,       '2000-01-01', '2000-01-01'),
					('1c', 3,       '2000-01-01', '2000-01-01'),
					('1c', 6,       '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 2, 3, 4, 5, 6 ] }`,
			wantResponseCode: http.StatusNoContent,
			wantResponseBody: `null`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1c",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int32(2),
					"nip":        "1c",
					"role_id":    int16(2),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int32(3),
					"nip":        "1c",
					"role_id":    int16(3),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
				{
					"id":         int32(5),
					"nip":        "1c",
					"role_id":    int16(5),
					"created_at": "{created_at}",
					"updated_at": "{updated_at}",
					"deleted_at": nil,
				},
				{
					"id":         int32(4),
					"nip":        "1c",
					"role_id":    int16(6),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name: "error: user not found",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1');
				insert into user_role
					(nip, role_id, created_at,   updated_at) values
					('1', 1,       '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1c",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 6 ] }`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name: "error: have active, deleted and not exists role",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1');
				insert into user_role
					(nip, role_id, created_at,   updated_at) values
					('1', 1,       '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 2, 4, 5, 6, 8 ] }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data role tidak ditemukan"}`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name: "error: user is deleted",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,     nip, deleted_at) values
					('00000000-0000-0000-0000-000000000001', 'zimbra',   '1', '2000-01-01'),
					('00000000-0000-0000-0000-000000000001', 'keycloak', '1', '2000-01-01');
				insert into user_role
					(nip, role_id, created_at,   updated_at) values
					('1', 1,       '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 1, 2 ] }`,
			wantResponseCode: http.StatusNotFound,
			wantResponseBody: `{"message": "data tidak ditemukan"}`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name: "error: don't allow empty json",
			dbData: seedData + `
				insert into "user"
					(id,                                     source,   nip) values
					('00000000-0000-0000-0000-000000000001', 'zimbra', '1');
				insert into user_role
					(nip, role_id, created_at,   updated_at) values
					('1', 1,       '2000-01-01', '2000-01-01');
			`,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{}`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"role_ids\" harus diisi"}`,
			wantDBUserRoles: dbtest.Rows{
				{
					"id":         int32(1),
					"nip":        "1",
					"role_id":    int16(1),
					"created_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"updated_at": time.Date(2000, 1, 1, 0, 0, 0, 0, time.UTC).Local(),
					"deleted_at": nil,
				},
			},
		},
		{
			name:             "error: deleted role",
			dbData:           seedData,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 3 ] }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data role tidak ditemukan"}`,
			wantDBUserRoles:  dbtest.Rows{},
		},
		{
			name:             "error: role not exists",
			dbData:           seedData,
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": [ 8 ] }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "data role tidak ditemukan"}`,
			wantDBUserRoles:  dbtest.Rows{},
		},
		{
			name:             "error: have null values",
			paramNIP:         "1a",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "role_ids": null }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"role_ids\" tidak boleh null"}`,
			wantDBUserRoles:  dbtest.Rows{},
		},
		{
			name:             "error: have additional params and duplicate role ids",
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			requestBody:      `{ "id": "", "role_ids": [ 1, 1 ] }`,
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "parameter \"id\" tidak didukung | parameter \"role_ids\" item tidak boleh duplikat"}`,
			wantDBUserRoles:  dbtest.Rows{},
		},
		{
			name:             "error: body is empty",
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": authHeader},
			wantResponseCode: http.StatusBadRequest,
			wantResponseBody: `{"message": "request body harus diisi"}`,
			wantDBUserRoles:  dbtest.Rows{},
		},
		{
			name:             "error: invalid token",
			paramNIP:         "1",
			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
			wantResponseCode: http.StatusUnauthorized,
			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
			wantDBUserRoles:  dbtest.Rows{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			db := dbtest.New(t, dbmigrations.FS)
			_, err := db.Exec(context.Background(), tt.dbData)
			require.NoError(t, err)

			req := httptest.NewRequest(http.MethodPatch, "/v1/users/"+tt.paramNIP, strings.NewReader(tt.requestBody))
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

			actualUserRoles, err := dbtest.QueryAll(db, "user_role", "nip, role_id, id")
			require.NoError(t, err)
			if len(tt.wantDBUserRoles) == len(actualUserRoles) {
				for i, row := range actualUserRoles {
					if tt.wantDBUserRoles[i]["id"] == "{id}" {
						tt.wantDBUserRoles[i]["id"] = row["id"]
					}
					if tt.wantDBUserRoles[i]["created_at"] == "{created_at}" {
						assert.WithinDuration(t, time.Now(), row["created_at"].(time.Time), 10*time.Second)
						assert.Equal(t, row["created_at"], row["updated_at"])

						tt.wantDBUserRoles[i]["created_at"] = row["created_at"]
						tt.wantDBUserRoles[i]["updated_at"] = row["updated_at"]
					}
					if tt.wantDBUserRoles[i]["deleted_at"] == "{deleted_at}" {
						assert.WithinDuration(t, time.Now(), row["deleted_at"].(time.Time), 10*time.Second)
						tt.wantDBUserRoles[i]["deleted_at"] = row["deleted_at"]
					}
				}
			}
			assert.Equal(t, tt.wantDBUserRoles, actualUserRoles)
		})
	}
}
