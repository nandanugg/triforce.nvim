package api

// These constants are intentionally written using a mix of PascalCase and snake_case
// to improve readability within route definitions.
//
//revive:disable:var-naming
const (
	Kode_Pegawai_Self  = "kepegawaian.pegawai.self"
	Kode_Pegawai_Read  = "kepegawaian.pegawai.read"
	Kode_Pegawai_Write = "kepegawaian.pegawai.write"

	Kode_SuratKeputusan_Self = "kepegawaian.surat-keputusan.self"
	Kode_SuratKeputusan_Read = "kepegawaian.surat-keputusan.read"

	Kode_SuratKeputusanApproval_Read   = "kepegawaian.surat-keputusan-approval.read"
	Kode_SuratKeputusanApproval_Review = "kepegawaian.surat-keputusan-approval.review"
	Kode_SuratKeputusanApproval_Sign   = "kepegawaian.surat-keputusan-approval.sign"

	Kode_DataMaster_Public = "kepegawaian.data-master.public"
	Kode_DataMaster_Read   = "kepegawaian.data-master.read"
	Kode_DataMaster_Write  = "kepegawaian.data-master.write"

	Kode_ManajemenAkses_Self  = "portal.manajemen-akses.self"
	Kode_ManajemenAkses_Read  = "portal.manajemen-akses.read"
	Kode_ManajemenAkses_Write = "portal.manajemen-akses.write"

	Kode_Informasi_Read = "portal.informasi.read"
) //revive:enable:var-naming
