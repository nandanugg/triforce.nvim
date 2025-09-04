package keluarga

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	dbRepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository struct {
	db   *sql.DB
	sqlc *dbRepo.Queries
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

var (
	peranAyah  = "AYAH"
	peranIbu   = "IBU"
	peranIstri = "ISTRI"
	peranSuami = "SUAMI"
	peranAnak  = "ANAK"
)

func (r *repository) listOrangTua(ctx context.Context, userID int64) ([]keluarga, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			ot."ID",
			ot."KODE",
			ot."NAMA",
			ot."TANGGAL_LAHIR"
		from orang_tua ot
		join users u on ot."NIP" = u.nip
		where u.id = $1
		`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []keluarga{}
	for rows.Next() {
		var row keluarga
		var kode int
		err := rows.Scan(
			&row.ID,
			&kode,
			&row.Nama,
			&row.TanggalLahir,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		switch kode {
		case 1:
			row.Peran = &peranAyah
		case 2:
			row.Peran = &peranIbu
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}

func (r *repository) listPasangan(ctx context.Context, userID int64) ([]keluarga, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			i."ID",
			i."HUBUNGAN",
			i."NAMA"
		from istri i
		join users u on i."NIP" = u.nip
		where u.id = $1
		`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []keluarga{}
	for rows.Next() {
		var row keluarga
		var hubungan int
		err := rows.Scan(
			&row.ID,
			&hubungan,
			&row.Nama,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		switch hubungan {
		case 1:
			row.Peran = &peranIstri
		case 2:
			row.Peran = &peranSuami
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}

func (r *repository) listAnak(ctx context.Context, userID int64) ([]keluarga, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			a."ID",
			a."NAMA",
			a."TANGGAL_LAHIR"
		from anak a
		join users u on a."NIP" = u.nip
		where u.id = $1
		`, userID,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []keluarga{}
	for rows.Next() {
		var row keluarga
		err := rows.Scan(
			&row.ID,
			&row.Nama,
			&row.TanggalLahir,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		row.Peran = &peranAnak
		if row.TanggalLahir != nil {
			row.TanggalLahir = &strings.Split(*row.TanggalLahir, "T")[0]
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}

func (r *repository) GetAnakByEmployeeID(ctx context.Context, pnsID sql.NullString) ([]dbRepo.GetChildrenByEmployeeIDRow, error) {
	return r.sqlc.GetChildrenByEmployeeID(ctx, pnsID)
}
