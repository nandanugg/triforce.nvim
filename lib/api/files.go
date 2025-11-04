package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

const maxFileSize = 50 << 20

func GetFileBase64(c echo.Context) (b64 string, filename string, err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "file harus diunggah")
	}

	if file.Size > maxFileSize {
		return "", "", echo.NewHTTPError(http.StatusRequestEntityTooLarge, fmt.Sprintf("ukuran file melebihi batas %d MB", maxFileSize>>20))
	}

	src, err := file.Open()
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "gagal membuka file yang diunggah")
	}
	defer src.Close()

	dst, err := io.ReadAll(src)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "gagal membaca file yang diunggah")
	}

	mimeType := http.DetectContentType(dst)
	if mimeType == "application/octet-stream" || strings.HasPrefix(mimeType, "text/plain") {
		if ext := filepath.Ext(file.Filename); ext != "" {
			if typ := mime.TypeByExtension(ext); typ != "" {
				mimeType = typ
			}
		}
	}

	base64Data := base64.StdEncoding.EncodeToString(dst)
	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), file.Filename, nil
}
