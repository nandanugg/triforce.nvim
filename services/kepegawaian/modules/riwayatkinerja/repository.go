package riwayatkinerja

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]riwayatKinerja, error) {
	rows, err := r.db.Query(ctx, `
		select
			rk.id,
			rk.tahun,
			rk.rating_hasil_kerja,
			rk.rating_perilaku_kerja,
			rk.predikat_kinerja
		from rwt_kinerja rk
		join users u on rk.nip = u.nip
		where u.id = $1
		order by rk.tahun asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []riwayatKinerja{}
	for rows.Next() {
		var row riwayatKinerja
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
	err := r.db.QueryRow(ctx, `
		select count(1)
		from rwt_kinerja rk
		join users u on rk.nip = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
