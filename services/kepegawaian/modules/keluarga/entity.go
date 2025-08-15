package keluarga

type keluarga struct {
	ID           int64   `json:"id"`
	Peran        *string `json:"peran,omitempty"`
	Nama         *string `json:"nama,omitempty"`
	TanggalLahir *string `json:"tanggal_lahir,omitempty"`
}
