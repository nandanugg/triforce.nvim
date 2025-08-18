package sertifikasi

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]sertifikasi, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rs.id,
			rs.nama_sertifikasi,
			rs.tahun,
			rs.createddate
		from kepegawaian.rwt_sertifikasi rs
		join kepegawaian.users u on rs.nip = u.nip
		where u.id = $1
		order by rs.tahun asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []sertifikasi{}
	for rows.Next() {
		var row sertifikasi
		err := rows.Scan(
			&row.ID,
			&row.NamaSertifikasi,
			&row.Tahun,
			&row.TanggalUpload,
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
		from kepegawaian.rwt_sertifikasi rs
		join kepegawaian.users u on rs.nip = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
