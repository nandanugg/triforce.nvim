package bsre

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func NewMockBSreServer(t *testing.T) *APIClient {
	t.Helper()

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case strings.Contains(r.URL.Path, "/sign"):
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"result":"mocked-base64"}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"error":"not found"}`))
		}
	})

	server := httptest.NewServer(handler)

	client := New(server.URL, "", "")

	return client
}
