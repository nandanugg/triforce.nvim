package riwayatpenugasan

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

func (r *repository) list(ctx context.Context, limit, offset uint) ([]riwayatPenugasan, error) {
	rows, err := r.db.Query(ctx, `
		select
			rp.id,
			rp.tipe_jabatan,
			rp.deskripsi_jabatan,
			rp.tanggal_mulai,
			rp.tanggal_selesai
		from rwt_penugasan rp
		order by rp.tanggal_mulai desc
		limit $1 offset $2
		`, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []riwayatPenugasan{}
	for rows.Next() {
		var row riwayatPenugasan
		err := rows.Scan(
			&row.ID,
			&row.TipeJabatan,
			&row.DeskripsiJabatan,
			&row.TanggalMulai,
			&row.TanggalSelesai,
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

func (r *repository) count(ctx context.Context) (uint, error) {
	var result uint
	err := r.db.QueryRow(ctx, `
		select count(1)
		from rwt_penugasan rp
		`).Scan(&result)

	return result, err
}
