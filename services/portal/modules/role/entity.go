package role

type role struct {
	ID                  int16                 `json:"id"`
	Nama                string                `json:"nama"`
	Deskripsi           string                `json:"deskripsi"`
	IsDefault           bool                  `json:"is_default"`
	IsAktif             bool                  `json:"is_aktif"`
	JumlahUser          int32                 `json:"jumlah_user"`
	ResourcePermissions *[]resourcePermission `json:"resource_permissions,omitempty"`
}

type resourcePermission struct {
	ID             int32  `json:"id"`
	Kode           string `json:"kode"`
	NamaResource   string `json:"nama_resource"`
	NamaPermission string `json:"nama_permission"`
}
