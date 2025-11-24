package jabatan

import "errors"

var (
	errJabatanReferenced = errors.New("jabatan masih digunakan oleh pegawai")
	errJabatanExists     = errors.New("kode jabatan sudah ada")
)
