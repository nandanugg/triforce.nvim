package pelatihanstruktural

import (
	"context"
	"database/sql"
	"fmt"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]pelatihanStruktural, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rds."ID",
			rds."NAMA_DIKLAT",
			rds."NOMOR",
			rds."TANGGAL",
			rds."TAHUN"
		from kepegawaian.rwt_diklat_struktural rds
		join kepegawaian.pegawai p on rds."PNS_NIP" = p."NIP_BARU"
		join kepegawaian.users u on p."NIP_BARU" = u.nip
		where u.id = $1
		order by rds."TANGGAL" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pelatihanStruktural{}
	for rows.Next() {
		var row pelatihanStruktural
		err := rows.Scan(
			&row.ID,
			&row.NamaDiklat,
			&row.Nomor,
			&row.Tanggal,
			&row.Tahun,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}

func (r *repository) count(ctx context.Context, userID int64) (uint, error) {
	var result uint
	err := r.db.QueryRowContext(ctx, `
		select count(1)
		from kepegawaian.rwt_diklat_struktural rds
		join kepegawaian.pegawai p on rds."PNS_NIP" = p."NIP_BARU"
		join kepegawaian.users u on p."NIP_BARU" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
