package keluarga

import "github.com/jackc/pgx/v5/pgtype"

func hubunganToPeran(hubungan pgtype.Int2) string {
	switch hubungan.Int16 {
	case 1:
		return "Ayah"
	case 2:
		return "Ibu"
	default:
		return "Orang Tua"
	}
}

func statusHidupFromTanggalMeninggal(tglMeninggal pgtype.Date) string {
	if tglMeninggal.Valid {
		return "Sudah Meninggal"
	}
	return "Masih Hidup"
}

func pnsToLabel(pns pgtype.Int2) string {
	if pns.Int16 == 0 {
		return "Bukan PNS"
	}
	return "PNS"
}

func statusAnakToLabel(statusAnak pgtype.Text) string {
	switch statusAnak.String {
	case "1":
		return "Kandung"
	case "2":
		return "Angkat"
	default:
		return "Tidak Diketahui"
	}
}

func statusPernikahanToString(status pgtype.Int2) string {
	switch status.Int16 {
	case 1:
		return "Menikah"
	case 2:
		return "Cerai Hidup"
	case 3:
		return "Cerai Mati"
	default:
		return "Tidak Diketahui"
	}
}
