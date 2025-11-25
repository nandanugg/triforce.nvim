package bsre

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"net/url"
	"time"
)

const (
	signURL = "/api/sign/pdf"
)

// Config holds the API client credentials and host info.
type config struct {
	Host              string
	BasicAuthUsername string
	BasicAuthPassword string
}

// Client defines a simple API sender interface.
type Client interface {
	Sign(query SignParams, files []UploadFile) (string, int, error)
}

// APIClient implements Client using net/http.
type APIClient struct {
	cfg  config
	http *http.Client
}

// New creates a new API client instance.
func New(host, basicAuthUsername, basicAuthPassword string) *APIClient {
	return &APIClient{
		cfg: config{
			Host:              host,
			BasicAuthUsername: basicAuthUsername,
			BasicAuthPassword: basicAuthPassword,
		},
		http: &http.Client{Timeout: 60 * time.Second},
	}
}

func (c *APIClient) Sign(query SignParams, files []UploadFile) (string, int, error) {
	var (
		body        io.Reader
		contentType string
		fullURL     = c.cfg.Host + signURL
	)

	u, err := url.Parse(fullURL)
	if err != nil {
		return "", 0, fmt.Errorf("parse URL: %w", err)
	}
	q := u.Query()
	q.Set("passphrase", query.Passphrase)
	q.Set("nik", query.NIK)
	q.Set("tampilan", query.Tampilan)
	u.RawQuery = q.Encode()
	fullURL = u.String()

	if len(files) > 0 {
		var b bytes.Buffer
		w := multipart.NewWriter(&b)

		for _, file := range files {
			data, err := base64.StdEncoding.DecodeString(file.ContentBase64)
			if err != nil {
				return "", 0, fmt.Errorf("decode base64: %w", err)
			}
			name := file.Name
			if name == "" {
				name = "file.pdf"
			}
			field := file.Field
			if field == "" {
				field = "file"
			}
			partHeader := make(textproto.MIMEHeader)
			partHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, field, name))
			partHeader.Set("Content-Type", "application/pdf")
			fw, err := w.CreatePart(partHeader)
			if err != nil {
				return "", 0, err
			}
			if _, err := io.Copy(fw, bytes.NewReader(data)); err != nil {
				return "", 0, err
			}
		}

		w.Close()
		body = &b
		contentType = w.FormDataContentType()
	}

	req, err := http.NewRequest(http.MethodPost, fullURL, body)
	if err != nil {
		return "", 0, err
	}

	auth := base64.StdEncoding.EncodeToString([]byte(c.cfg.BasicAuthUsername + ":" + c.cfg.BasicAuthPassword))
	req.Header.Set("Authorization", "Basic "+auth)
	req.Header.Set("Cache-Control", "no-cache")
	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}

	resp, err := c.http.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}

	encoded := base64.StdEncoding.EncodeToString(respBody)
	return encoded, resp.StatusCode, nil
}
