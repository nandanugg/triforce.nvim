package keluarga

import (
	"github.com/jackc/pgx/v5/pgtype"
)

// HubunganToPeran converts hubungan (1=ayah, 2=ibu) to human-readable role
func HubunganToPeran(hubungan pgtype.Int2) string {
	if !hubungan.Valid {
		return "Orang Tua"
	}
	switch hubungan.Int16 {
	case 1:
		return "Ayah"
	case 2:
		return "Ibu"
	default:
		return "Orang Tua"
	}
}

// StatusHidupFromTanggalMeninggal returns "Masih Hidup" or "Sudah Meninggal"
func StatusHidupFromTanggalMeninggal(tglMeninggal pgtype.Date) string {
	if tglMeninggal.Valid {
		return "Sudah Meninggal"
	}
	return "Masih Hidup"
}

// PNSToLabel converts pns (0/1) to "PNS" or "Bukan PNS"
func PNSToLabel(pns pgtype.Int2) string {
	if !pns.Valid || pns.Int16 == 0 {
		return "Bukan PNS"
	}
	return "PNS"
}

// StatusAnakToLabel converts "1"/"2" to "Kandung"/"Angkat"
func StatusAnakToLabel(statusAnak pgtype.Text) string {
	switch statusAnak.String {
	case "1":
		return "Kandung"
	case "2":
		return "Angkat"
	default:
		return "Tidak Diketahui"
	}
}

func nullStringPtr(ns pgtype.Text) *string {
	if ns.Valid {
		return &ns.String
	}
	return nil
}

func StatusPernikahanToString(status pgtype.Int2) string {
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
