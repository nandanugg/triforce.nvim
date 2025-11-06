package pdfcpu

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"

	pdfcpuapi "github.com/pdfcpu/pdfcpu/pkg/api"
)

func AddSignatureToPDF(pdfBase64, sigBase64 string, x, y, scale float64, pages string) (string, error) {
	pdfBytes, err := base64.StdEncoding.DecodeString(pdfBase64)
	if err != nil {
		return "", err
	}

	imgBytes, err := base64.StdEncoding.DecodeString(sigBase64)
	if err != nil {
		return "", err
	}

	imgFile, err := os.CreateTemp("", "sig-*.png")
	if err != nil {
		return "", err
	}
	defer func() {
		imgFile.Close()
		os.Remove(imgFile.Name())
	}()
	if _, err := imgFile.Write(imgBytes); err != nil {
		return "", err
	}

	pdfReader := bytes.NewReader(pdfBytes)

	// jika letak ttd di isi 2, maka taruh ttd di halaman terakhir
	if pages == "2" {
		totalPages, err := pdfcpuapi.PageCount(bytes.NewReader(pdfBytes), nil)
		if err != nil {
			return "", err
		}
		pages = strconv.Itoa(totalPages)
	}

	wmDesc := fmt.Sprintf("pos:bl, off:%.1f %.1f, rot:0, scale:%.2f", x, y, scale)
	wm, err := pdfcpuapi.ImageWatermark(imgFile.Name(), wmDesc, true, true, 0)
	if err != nil {
		return "", err
	}

	var outBuf bytes.Buffer

	conf := pdfcpuapi.LoadConfiguration()
	err = pdfcpuapi.AddWatermarks(pdfReader, &outBuf, []string{pages}, wm, conf)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(outBuf.Bytes()), nil
}
