package kinerja

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]kinerja, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rk.id,
			rk.tahun,
			rk.rating_hasil_kerja,
			rk.rating_perilaku_kerja,
			rk.predikat_kinerja
		from kepegawaian.rwt_kinerja rk
		join kepegawaian.users u on rk.nip = u.nip
		where u.id = $1
		order by rk.tahun asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []kinerja{}
	for rows.Next() {
		var row kinerja
		err := rows.Scan(
			&row.ID,
			&row.Tahun,
			&row.HasilKinerja,
			&row.PerilakuKerja,
			&row.KuadranKinerja,
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
		from kepegawaian.rwt_kinerja rk
		join kepegawaian.users u on rk.nip = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
