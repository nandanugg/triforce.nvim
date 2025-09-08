package api

import (
	"net/http"
	"time"
)

// NewHTTPClient http.Client initialization with timeout and idle connection increase
func NewHTTPClient() *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 100
	t.MaxIdleConnsPerHost = 100

	return &http.Client{
		Timeout:   10 * time.Second,
		Transport: t,
	}
}
