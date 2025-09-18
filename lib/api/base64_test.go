package api

import (
	"encoding/base64"
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_base64_GetMimeTypeAndDecodedData(t *testing.T) {
	t.Parallel()

	filePath := "sample/hello.pdf"
	pdfBytes, err := os.ReadFile(filePath)
	require.NoError(t, err)

	pngBytes := []byte{
		0x89, 0x50, 0x4e, 0x47, 0x0d, 0x0a, 0x1a, 0x0a,
		0x00, 0x00, 0x00, 0x0d, 0x49, 0x48, 0x44, 0x52,
		0x00, 0x00, 0x00, 0x01, 0x00, 0x00, 0x00, 0x01,
		0x08, 0x06, 0x00, 0x00, 0x00, 0x1f, 0x15, 0xc4,
		0x89, 0x00, 0x00, 0x00, 0x0a, 0x49, 0x44, 0x41,
		0x54, 0x78, 0x9c, 0x63, 0xf8, 0xff, 0xff, 0x3f,
		0x00, 0x05, 0xfe, 0x02, 0xfe, 0xa7, 0x46, 0x90,
		0x3d, 0x00, 0x00, 0x00, 0x00, 0x49, 0x45, 0x4e,
		0x44, 0xae, 0x42, 0x60, 0x82,
	}

	pdfBase64 := base64.StdEncoding.EncodeToString(pdfBytes)
	pngBase64 := base64.StdEncoding.EncodeToString(pngBytes)

	tests := []struct {
		name              string
		fileBase64        string
		wantMimetype      string
		wantResponseBytes []byte
		wantError         error
	}{
		{
			name:              "ok: valid pdf with data: prefix",
			fileBase64:        "data:application/pdf;base64," + pdfBase64,
			wantMimetype:      "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid pdf without data: prefix",
			fileBase64:        pdfBase64,
			wantMimetype:      "application/pdf",
			wantResponseBytes: pdfBytes,
		},
		{
			name:              "ok: valid png with incorrect content-type",
			fileBase64:        "data:images/png;base64," + pngBase64,
			wantMimetype:      "images/png",
			wantResponseBytes: pngBytes,
		},
		{
			name:       "error: invalid pdf",
			fileBase64: "data:application/pdf;base64,invalid",
			wantError:  errors.New("decode file base64: illegal base64 data at input byte 4"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			mimetype, decodedData, err := GetMimeTypeAndDecodedData(tt.fileBase64)

			if tt.wantError != nil {
				assert.EqualError(t, err, tt.wantError.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.wantMimetype, mimetype)
				assert.Equal(t, tt.wantResponseBytes, decodedData)
			}
		})
	}
}
