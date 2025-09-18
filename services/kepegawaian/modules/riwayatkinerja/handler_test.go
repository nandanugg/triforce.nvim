package riwayatkinerja

// import (
// 	"net/http"
// 	"net/http/httptest"
// 	"net/url"
// 	"testing"

// 	"github.com/stretchr/testify/assert"
// 	"github.com/stretchr/testify/require"

// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api/apitest"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db/dbtest"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/config"
// 	dbmigrations "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/migrations"
// 	"gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/docs"
// )

// func Test_handler_list(t *testing.T) {
// 	t.Parallel()

// 	dbData := `
// 		insert into users
// 			(id, role_id, email, username, password_hash, reset_hash, last_login,  last_ip, created_on,  deleted, reset_by, banned, ban_message, display_name, display_name_changed, timezone, language, active, activate_hash, password_iterations, force_password_reset, nip,  satkers, admin_nomor, imei, token, real_imei, fcm,  banned_asigo) values
// 			(41, 41,      '41a', '41b',    '41c',         '41d',      '2001-01-02','41f',   '2001-01-03',1,       1,        1,      '41k',       '41l',        '2001-01-04',         '41n',    '41o',    1,      '41q',         1,                   1,                    '1c', '41u',   1,           '41w','41x', '41y',     '41z',1);
// 		insert into rwt_kinerja
// 			(id, id_simarin, tahun, periode_mulai, periode_selesai, format_skp, jenis_skp, idp_pegawai, nip,   nama,  panggol, jabatan, penugasan, id_unit_kerja, unit_kerja, idp_penilai, nip_penilai, nama_penilai, panggol_penilai, jabatan_penilai, penugasan_penilai, id_unit_kerja_penilai, unit_kerja_penilai, idp_atasan_penilai, nip_atasan_penilai, nama_atasan_penilai, panggol_atasan_penilai, jabatan_atasan_penilai, penugasan_atasan_penilai, id_unit_kerja_atasan_penilai, unit_kerja_atasan_penilai, idp_penilai_simarin, nip_penilai_simarin, nama_penilai_simarin, panggol_penilai_simarin, jabatan_penilai_simarin, penugasan_penilai_simarin, id_unit_kerja_penilai_simarin, unit_kerja_penilai_simarin, idp_penilai_realisasi, nip_penilai_realisasi, nama_penilai_realisasi, panggol_penilai_realisasi, jabatan_penilai_realisasi, penugasan_penilai_realisasi, id_unit_kerja_penilai_realisasi, unit_kerja_penilai_realisasi, idp_atasan_penilai_realisasi, nip_atasan_penilai_realisasi, nama_atasan_penilai_realisasi, panggol_atasan_penilai_realisasi, jabatan_atasan_penilai_realisasi, penugasan_atasan_penilai_realisasi, id_unit_kerja_atasan_penilai_realisasi, unit_kerja_atasan_penilai_realisasi, idp_penilai_realisasi_simarin, nip_penilai_realisasi_simarin, nama_penilai_realisasi_simarin, panggol_penilai_realisasi_simarin, jabatan_penilai_realisasi_simarin, penugasan_penilai_realisasi_simarin, id_unit_kerja_penilai_realisasi_simarin, unit_kerja_penilai_realisasi_simarin, nama_realisasi, nip_realisasi, panggol_realisasi, jabatan_realisasi, penugasan_realisasi, id_unit_kerja_realisasi, unit_kerja_realisasi, skp_instansi_lama, capaian_kinerja_org, pola_distribusi_img, nilai_akhir_hasil_kerja, rating_hasil_kerja, nilai_akhir_perilaku_kerja, rating_perilaku_kerja, predikat_kinerja, tunjangan_kinerja, catatan_rekomendasi, is_keberatan, keberatan, penjelasan_pejabat_penilai, keputusan_rekomendasi_atasan_pejabat, url_skp_instansi_lama, is_keberatan_date, ref,                                    id_arsip, created_date) values
// 			(11, 111,        112,   '2000-01-01',  '2000-01-02',    '11a',      '11b',     '11c',       '1c',  '11e', '11f',   '11g',   '11h',     113,           '11j',      '11k',       '11l',       '11m',        '11n',           '11o',           '11p',             114,                   '11q',              '11r',              '11s',              '11t',               '11u',                  '11v',                  '11w',                    115,                          '11x',                     '11y',               '11z',               '11aa',               '11ab',                  '11ac',                  '11ad',                    116,                           '11ae',                     '11af',                '11ag',                '11ah',                 '11ai',                    '11aj',                    '11ak',                      117,                             '11al',                       '11am',                       '11an',                       '11ao',                        '11ap',                           '11aq',                           '11ar',                             118,                                    '11as',                              '11at',                        '11au',                        '11av',                         '11aw',                            '11ax',                            '11ay',                              119,                                     '11az',                               '11ba',         '11baa',       '11bb',            '11bc',            '11bd',              1111,                    '11be',               '11bf',            '11bg',              '11bh',              1112,                    '11bi',             1113,                       '11bj',                '11bk',           1114,              '11bl',              '11bm',       '11bn',    '11bo',                     '11bp',                               '11bq',                '11br',            '11000000-0000-0000-0000-000000000001', 1115,     '2000-01-03'),
// 			(12, 121,        132,   '2001-01-01',  '2001-01-02',    '12a',      '12b',     '12c',       '1c',  '12e', '12f',   '12g',   '12h',     123,           '12j',      '12k',       '12l',       '12m',        '12n',           '12o',           '12p',             124,                   '12q',              '12r',              '12s',              '12t',               '12u',                  '12v',                  '12w',                    125,                          '12x',                     '12y',               '12z',               '12aa',               '12ab',                  '12ac',                  '12ad',                    126,                           '12ae',                     '12af',                '12ag',                '12ah',                 '12ai',                    '12aj',                    '12ak',                      127,                             '12al',                       '12am',                       '12an',                       '12ao',                        '12ap',                           '12aq',                           '12ar',                             128,                                    '12as',                              '12at',                        '12au',                        '12av',                         '12aw',                            '12ax',                            '12ay',                              129,                                     '12az',                               '12ba',         '12baa',       '12bb',            '12bc',            '12bd',              1212,                    '12be',               '12bf',            '12bg',              '12bh',              1212,                    '12bi',             1213,                       '12bj',                '12bk',           1214,              '12bl',              '12bm',       '12bn',    '12bo',                     '12bp',                               '12bq',                '12br',            '12000000-0000-0000-0000-000000000001', 1215,     '2001-01-03'),
// 			(13, 131,        122,   '2002-01-01',  '2002-01-02',    '13a',      '13b',     '13c',       '1c',  '13e', '13f',   '13g',   '13h',     133,           '13j',      '13k',       '13l',       '13m',        '13n',           '13o',           '13p',             134,                   '13q',              '13r',              '13s',              '13t',               '13u',                  '13v',                  '13w',                    135,                          '13x',                     '13y',               '13z',               '13aa',               '13ab',                  '13ac',                  '13ad',                    136,                           '13ae',                     '13af',                '13ag',                '13ah',                 '13ai',                    '13aj',                    '13ak',                      137,                             '13al',                       '13am',                       '13an',                       '13ao',                        '13ap',                           '13aq',                           '13ar',                             138,                                    '13as',                              '13at',                        '13au',                        '13av',                         '13aw',                            '13ax',                            '13ay',                              139,                                     '13az',                               '13ba',         '13baa',       '13bb',            '13bc',            '13bd',              1313,                    '13be',               '13bf',            '13bg',              '13bh',              1312,                    '13bi',             1313,                       '13bj',                '13bk',           1314,              '13bl',              '13bm',       '13bn',    '13bo',                     '13bp',                               '13bq',                '13br',            '13000000-0000-0000-0000-000000000001', 1315,     '2002-01-03'),
// 			(14, 141,        142,   '2003-01-01',  '2003-01-02',    '14a',      '14b',     '14c',       '2c',  '14e', '14f',   '14g',   '14h',     143,           '14j',      '14k',       '14l',       '14m',        '14n',           '14o',           '14p',             144,                   '14q',              '14r',              '14s',              '14t',               '14u',                  '14v',                  '14w',                    145,                          '14x',                     '14y',               '14z',               '14aa',               '14ab',                  '14ac',                  '14ad',                    146,                           '14ae',                     '14af',                '14ag',                '14ah',                 '14ai',                    '14aj',                    '14ak',                      147,                             '14al',                       '14am',                       '14an',                       '14ao',                        '14ap',                           '14aq',                           '14ar',                             148,                                    '14as',                              '14at',                        '14au',                        '14av',                         '14aw',                            '14ax',                            '14ay',                              149,                                     '14az',                               '14ba',         '14baa',       '14bb',            '14bc',            '14bd',              1414,                    '14be',               '14bf',            '14bg',              '14bh',              1412,                    '14bi',             1413,                       '14bj',                '14bk',           1414,              '14bl',              '14bm',       '14bn',    '14bo',                     '14bp',                               '14bq',                '14br',            '14000000-0000-0000-0000-000000000001', 1415,     '2003-01-03');
// 	`

// 	tests := []struct {
// 		name             string
// 		dbData           string
// 		requestQuery     url.Values
// 		requestHeader    http.Header
// 		wantResponseCode int
// 		wantResponseBody string
// 	}{
// 		{
// 			name:             "ok: tanpa parameter apapun",
// 			dbData:           dbData,
// 			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
// 			wantResponseCode: http.StatusOK,
// 			wantResponseBody: `{
// 				"data": [
// 					{
// 						"id":              11,
// 						"hasil_kinerja":   "11bi",
// 						"kuadran_kinerja": "11bk",
// 						"perilaku_kerja":  "11bj",
// 						"tahun":           112
// 					},
// 					{
// 						"id":              13,
// 						"hasil_kinerja":   "13bi",
// 						"kuadran_kinerja": "13bk",
// 						"perilaku_kerja":  "13bj",
// 						"tahun":           122
// 					},
// 					{
// 						"id":              12,
// 						"hasil_kinerja":   "12bi",
// 						"kuadran_kinerja": "12bk",
// 						"perilaku_kerja":  "12bj",
// 						"tahun":           132
// 					}
// 				],
// 				"meta": {"limit": 10, "offset": 0, "total": 3}
// 			}`,
// 		},
// 		{
// 			name:             "ok: dengan parameter pagination",
// 			dbData:           dbData,
// 			requestQuery:     url.Values{"limit": []string{"1"}, "offset": []string{"1"}},
// 			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "41")}},
// 			wantResponseCode: http.StatusOK,
// 			wantResponseBody: `{
// 				"data": [
// 					{
// 						"id":              13,
// 						"hasil_kinerja":   "13bi",
// 						"kuadran_kinerja": "13bk",
// 						"perilaku_kerja":  "13bj",
// 						"tahun":           122
// 					}
// 				],
// 				"meta": {"limit": 1, "offset": 1, "total": 3}
// 			}`,
// 		},
// 		{
// 			name:             "ok: tidak ada data milik user",
// 			dbData:           dbData,
// 			requestHeader:    http.Header{"Authorization": []string{apitest.GenerateAuthHeader(config.Service, "200")}},
// 			wantResponseCode: http.StatusOK,
// 			wantResponseBody: `{"data": [], "meta": {"limit": 10, "offset": 0, "total": 0}}`,
// 		},
// 		{
// 			name:             "error: auth header tidak valid",
// 			dbData:           dbData,
// 			requestHeader:    http.Header{"Authorization": []string{"Bearer some-token"}},
// 			wantResponseCode: http.StatusUnauthorized,
// 			wantResponseBody: `{"message": "token otentikasi tidak valid"}`,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()

// 			db := dbtest.New(t, dbmigrations.FS)
// 			_, err := db.Exec(tt.dbData)
// 			require.NoError(t, err)

// 			req := httptest.NewRequest(http.MethodGet, "/v1/kinerja", nil)
// 			req.URL.RawQuery = tt.requestQuery.Encode()
// 			req.Header = tt.requestHeader
// 			rec := httptest.NewRecorder()

// 			e, err := api.NewEchoServer(docs.OpenAPIBytes)
// 			require.NoError(t, err)
// 			RegisterRoutes(e, db, api.NewAuthMiddleware(config.Service, apitest.Keyfunc))
// 			e.ServeHTTP(rec, req)

// 			assert.Equal(t, tt.wantResponseCode, rec.Code)
// 			assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
// 			assert.NoError(t, apitest.ValidateResponseSchema(rec, req, e))
// 		})
// 	}
// }
