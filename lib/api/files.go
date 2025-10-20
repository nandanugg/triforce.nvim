package api

import (
	"encoding/base64"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"
)

const maxFileSize = 50 << 20

func GetFileBase64(c echo.Context) (b64 string, filename string, err error) {
	file, err := c.FormFile("file")
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "file is required")
	}

	if file.Size > maxFileSize {
		return "", "", echo.NewHTTPError(http.StatusRequestEntityTooLarge, fmt.Sprintf("file size exceeds limit of %d MB", maxFileSize>>20))
	}

	src, err := file.Open()
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "failed to open uploaded file")
	}
	defer src.Close()

	dst, err := io.ReadAll(src)
	if err != nil {
		return "", "", echo.NewHTTPError(http.StatusBadRequest, "failed to read uploaded file")
	}

	mimeType := http.DetectContentType(dst)
	if mimeType == "application/octet-stream" {
		if ext := filepath.Ext(file.Filename); ext != "" {
			if typ := mime.TypeByExtension(ext); typ != "" {
				mimeType = typ
			}
		}
	}

	base64Data := base64.StdEncoding.EncodeToString(dst)

	return fmt.Sprintf("data:%s;base64,%s", mimeType, base64Data), file.Filename, nil
}
