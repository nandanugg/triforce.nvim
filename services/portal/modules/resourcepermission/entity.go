package resourcepermission

type resource struct {
	Nama                string               `json:"nama"`
	ResourcePermissions []resourcePermission `json:"resource_permissions"`
}

type resourcePermission struct {
	ID             int32  `json:"id"`
	Kode           string `json:"kode"`
	NamaPermission string `json:"nama_permission"`
}
