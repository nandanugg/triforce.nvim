package samplelogharian

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context, userID string, limit, offset uint) ([]logHarian, error) {
	// Contoh menggunakan raw SQL:
	rows, err := r.db.QueryContext(ctx, `
		select tanggal, aktivitas
		from sample_log_harian
		where user_id = $1
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []logHarian{}
	for rows.Next() {
		var lh logHarian
		if err := rows.Scan(&lh.Tanggal, &lh.Aktivitas); err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		result = append(result, lh)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}

func (r *repository) listCount(ctx context.Context, userID string) (uint, error) {
	// Contoh menggunakan query builder:
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(1)").
		From("sample_log_harian").
		Where("user_id = ?", userID)

	q, args, _ := qb.ToSql()
	var result uint

	err := r.db.QueryRowContext(ctx, q, args...).Scan(&result)
	return result, err
}
