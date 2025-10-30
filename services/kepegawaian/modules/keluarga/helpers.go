package keluarga

import "github.com/jackc/pgx/v5/pgtype"

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

func statusPNS(isPns bool) pgtype.Int2 {
	if isPns {
		return pgtype.Int2{Int16: 1, Valid: true}
	}
	return pgtype.Int2{Int16: 0, Valid: true}
}
