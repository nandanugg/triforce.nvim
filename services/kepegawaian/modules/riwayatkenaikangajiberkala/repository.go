package riwayatkenaikangajiberkala

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]riwayatKenaikanGajiBerkala, error) {
	rows, err := r.db.Query(ctx, `
		select
			rk.id,
			rk.tmt_sk,
			rk.no_sk,
			rk.tgl_sk,
			rk.n_gol_ruang,
			rk.n_gapok
		from rwt_kgb rk
		join users u on rk.pegawai_nip = u.nip
		where u.id = $1
		order by rk.tmt_sk
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []riwayatKenaikanGajiBerkala{}
	for rows.Next() {
		var row riwayatKenaikanGajiBerkala
		err := rows.Scan(
			&row.ID,
			&row.TMTKGB,
			&row.NoSK,
			&row.TglSK,
			&row.GolRuang,
			&row.GajiPokok,
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
		from rwt_kgb rk
		join users u on rk.pegawai_nip = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
