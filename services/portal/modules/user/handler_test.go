package user

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

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Read)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
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

			e, err := api.NewEchoServer(docs.OpenAPIBytes)
			require.NoError(t, err)

			authSvc := apitest.NewAuthService(api.Kode_ManajemenAkses_Read)
			RegisterRoutes(e, sqlc.New(db), api.NewAuthMiddleware(authSvc, apitest.Keyfunc))
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseCode, rec.Code)
			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
		})
	}
}
