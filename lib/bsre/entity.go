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

type ErrorResponse struct {
	StatusCode StatusCode `json:"status_code"`
	Message    string     `json:"message"`
}

type StatusCode int16

const (
	StatusCodePassphraseInvalid = 2031
)

func (e StatusCode) Message() string {
	switch e {
	case StatusCodePassphraseInvalid:
		return "Passphrase yang Anda masukkan salah. Silakan coba lagi."
	default:
		return "Terdapat kesalahan pada sistem BSRE. Silakan coba lagi."
	}
}
