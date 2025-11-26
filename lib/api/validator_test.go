package api

import (
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_validateRequest(t *testing.T) {
	t.Parallel()

	openapiSchema := []byte(`
openapi: 3.1.1
paths:
  /some-endpoint:
    get:
      parameters:
        - name: uuid_satu
          in: query
          schema:
            type: string
            format: uuid
        - name: string_dua
          in: query
          schema:
            type: string
        - name: limit
          in: query
          schema:
            type: integer
            minimum: 1
            maximum: 50
        - name: offset
          in: query
          schema:
            type: integer
            minimum: 0
      responses:
        "200":
          description: Sukses.
        "400":
          $ref: "#/components/responses/badRequest"
  /other-endpoint:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              additionalProperties: false
              properties:
                string_empat:
                  type: string
                  minLength: 1
                  maxLength: 10
                deskripsi:
                  type: [string, null]
                enum_lima:
                  type: string
                  enum:
                    - Enum value satu
                    - enum_value_dua
                array_tiga:
                  type: array
                  items:
                    type: string
                  uniqueItems: true
                  minItems: 1
                  maxItems: 2
              required:
                - string_empat
                - enum_lima
      responses:
        "200":
          description: Sukses.
        "400":
          $ref: "#/components/responses/badRequest"
  /optional-endpoint:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              additionalProperties: false
              minProperties: 1
              maxProperties: 2
              properties:
                var1:
                  type: string
                var2:
                  type: object
                  additionalProperties: false
                  properties:
                    var0:
                      type: string
                var3:
                  type: object
                  minProperties: 1
                  properties:
                    var0:
                      type: string
      responses:
        "200":
          description: Sukses.
        "400":
          $ref: "#/components/responses/badRequest"
components:
  responses:
    badRequest:
      description: Request dari client tidak valid.
      content:
        application/json:
          schema:
            additionalProperties: false
            properties:
              message:
                type: string
`)

	tests := []struct {
		name               string
		requestMethod      string
		requestPath        string
		requestQuery       url.Values
		requestHeader      http.Header
		requestBody        io.Reader
		wantResponseStatus int
		wantResponseBody   string
	}{
		{
			name:          "ok: get request with optional params",
			requestMethod: http.MethodGet,
			requestPath:   "/some-endpoint",
			requestQuery: url.Values{
				"uuid_satu": []string{uuid.NewString()},
				"offset":    []string{"0"},
				"limit":     []string{"50"},
			},
			wantResponseStatus: 200,
		},
		{
			name:               "ok: get request without optional params",
			requestMethod:      http.MethodGet,
			requestPath:        "/some-endpoint",
			wantResponseStatus: 200,
		},
		{
			name:               "error: path not found",
			requestMethod:      http.MethodGet,
			requestPath:        "/schema-not-found",
			wantResponseStatus: 404,
			wantResponseBody:   `{"message": "route tidak ditemukan pada openapi"}`,
		},
		{
			name:          "error: invalid parameter values",
			requestMethod: http.MethodGet,
			requestPath:   "/some-endpoint",
			requestQuery: url.Values{
				"uuid_satu":  []string{"test"},
				"offset":     []string{"-1"},
				"limit":      []string{"51"},
				"string_dua": []string{""},
			},
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "parameter \"uuid_satu\" harus dalam format uuid` +
				` | parameter \"string_dua\" tidak boleh kosong` +
				` | parameter \"limit\" harus tidak lebih dari 50` +
				` | parameter \"offset\" harus tidak kurang dari 0"}`,
		},
		{
			name:          "error: query params not number",
			requestMethod: http.MethodGet,
			requestPath:   "/some-endpoint",
			requestQuery: url.Values{
				"offset": []string{"offset"},
				"limit":  []string{"limit"},
			},
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "parameter \"limit\" harus dalam format yang sesuai` +
				` | parameter \"offset\" harus dalam format yang sesuai"}`,
		},
		{
			name:          "ok: post request",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": "Testing",
				"deskripsi": "Deskripsi",
				"enum_lima": "enum_value_dua"
			}`),
			wantResponseStatus: 200,
		},
		{
			name:          "error: undefined content type",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestBody: strings.NewReader(`{
				"string_empat": "Testing",
				"deskripsi": "Deskripsi",
				"enum_lima": "enum_value_dua"
			}`),
			wantResponseStatus: 400,
			wantResponseBody:   `{"message": "header Content-Type harus dalam format yang sesuai"}`,
		},
		{
			name:          "error: malformatted request body",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": "Testing"
				"deskripsi": "Deskripsi"
				"enum_lima": "enum_value_dua"
			}`),
			wantResponseStatus: 400,
			wantResponseBody:   `{"message": "request body harus dalam format yang sesuai"}`,
		},
		{
			name:          "error: missing required body params",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody:        strings.NewReader(`{}`),
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "parameter \"string_empat\" harus diisi` +
				` | parameter \"enum_lima\" harus diisi"}`,
		},
		{
			name:          "error: missing required body",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			wantResponseStatus: 400,
			wantResponseBody:   `{"message": "request body harus diisi"}`,
		},
		{
			name:          "error: invalid data type in request body",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": 123,
				"deskripsi": 123,
				"enum_lima": "Enum value satu"
			}`),
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "parameter \"deskripsi\" harus dalam tipe string, ` +
				` | parameter \"string_empat\" harus dalam tipe string"}`,
		},
		{
			name:          "error: invalid enum value",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": "foo",
				"deskripsi": "bar",
				"enum_lima": "enumune"
			}`),
			wantResponseStatus: 400,
			wantResponseBody:   `{"message": "parameter \"enum_lima\" harus salah satu dari \"Enum value satu\", \"enum_value_dua\""}`,
		},
		{
			name:          "error: string & array size too small",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": "",
				"enum_lima": "Enum value satu",
				"array_tiga": []
			}`),
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "parameter \"array_tiga\" harus 1 item atau lebih` +
				` | parameter \"string_empat\" harus 1 karakter atau lebih"}`,
		},
		{
			name:          "error: string & array size too big & duplicate",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": "This is a long title",
				"enum_lima": "enum_value_dua",
				"array_tiga": ["Testing", "Baru", "Baru"]
			}`),
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "parameter \"array_tiga\" harus 2 item atau kurang` +
				` | parameter \"array_tiga\" item tidak boleh duplikat` +
				` | parameter \"string_empat\" harus 10 karakter atau kurang"}`,
		},
		{
			name:          "error: request body have additional params",
			requestMethod: http.MethodPost,
			requestPath:   "/other-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"string_empat": "Judul",
				"enum_lima": "enum_value_dua",
				"jenis": 12345
			}`),
			wantResponseStatus: 400,
			wantResponseBody:   `{"message": "parameter \"jenis\" tidak didukung"}`,
		},
		{
			name:          "error: non nullable params, exceed maxProperties, and below minProperties",
			requestMethod: http.MethodPost,
			requestPath:   "/optional-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody: strings.NewReader(`{
				"var1": null,
				"var2": { "var1": "abc" },
				"var3": {}
			}`),
			wantResponseStatus: 400,
			wantResponseBody: `{"message": "request body harus 2 property atau kurang` +
				` | parameter \"var1\" tidak boleh null` +
				` | parameter \"var2.var1\" tidak didukung` +
				` | parameter \"var3\" harus 1 property atau lebih"}`,
		},
		{
			name:          "error: empty property",
			requestMethod: http.MethodPost,
			requestPath:   "/optional-endpoint",
			requestHeader: http.Header{
				"Content-Type": []string{"application/json"},
			},
			requestBody:        strings.NewReader(`{}`),
			wantResponseStatus: 400,
			wantResponseBody:   `{"message": "request body harus 1 property atau lebih"}`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			req := httptest.NewRequest(tt.requestMethod, tt.requestPath, tt.requestBody)
			req.URL.RawQuery = tt.requestQuery.Encode()
			req.Header = tt.requestHeader
			rec := httptest.NewRecorder()

			e, err := NewEchoServer(openapiSchema)
			require.NoError(t, err)
			e.Add(tt.requestMethod, tt.requestPath, func(c echo.Context) error {
				return c.Bind(&map[string]any{})
			})
			e.ServeHTTP(rec, req)

			assert.Equal(t, tt.wantResponseStatus, rec.Code)
			if tt.wantResponseStatus != http.StatusOK {
				assert.JSONEq(t, tt.wantResponseBody, rec.Body.String())
			} else {
				assert.Empty(t, rec.Body.String())
			}
		})
	}
}
