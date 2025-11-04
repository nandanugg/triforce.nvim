package bsre

const (
	InvisibleMode = "invisible"
)

// SignParams represents required query parameters for BSrE sign API.
type SignParams struct {
	NIK        string `form:"nik"`
	Passphrase string `form:"passphrase"`
	Tampilan   string `form:"tampilan"`
}

// UploadFile represents one multipart file with optional custom filename.
type UploadFile struct {
	Field         string
	ContentBase64 string
	Name          string
}
