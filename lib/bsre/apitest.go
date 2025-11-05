package bsre

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NewMockBSREServer(t *testing.T) *APIClient {
	t.Helper()

	handler := handler(http.StatusOK, `{"result":"mocked-base64"}`)

	server := httptest.NewServer(handler)

	client := New(server.URL, "", "")

	return client
}

func NewMockBSREServerWithCustomResponse(t *testing.T, statusCode int, responseBody string) *APIClient {
	t.Helper()

	handler := handler(statusCode, responseBody)

	server := httptest.NewServer(handler)

	client := New(server.URL, "", "")
	return client
}

func handler(statusCode int, responseBody string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/sign"):
			w.WriteHeader(statusCode)
			_, _ = w.Write([]byte(responseBody))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found"}`))
		}
	})
}
