package api

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

func GetMimetypeAndDecodedData(fileBase64 string) (string, []byte, error) {
	parts := strings.SplitN(fileBase64, ",", 2)
	rawBase64 := parts[len(parts)-1]

	decoded, err := base64.StdEncoding.DecodeString(rawBase64)
	if err != nil {
		return "", nil, fmt.Errorf("decode file base64: %w", err)
	}

	var mimeType string
	if strings.HasPrefix(fileBase64, "data:") {
		header := strings.Split(parts[0], ";")[0]
		mimeType = strings.TrimPrefix(header, "data:")
	} else {
		mimeType = http.DetectContentType(decoded)
	}

	return mimeType, decoded, nil
}
