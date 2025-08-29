package pemberitahuan

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

func (r *repository) list(ctx context.Context, limit, offset uint, cari string) ([]pemberitahuan, error) {
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			"p.id",
			"p.judul_berita",
			"p.deskripsi_berita",
			"p.status",
			"p.updated_by",
			"p.updated_at",
		).
		From("pemberitahuan p").
		Limit(uint64(limit)).
		Offset(uint64(offset))

	if cari != "" {
		cari = cari + "%"
		qb = qb.Where("(p.judul_berita like ? or p.deskripsi_berita like ?)", cari, cari)
	}

	q, args, _ := qb.ToSql()
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pemberitahuan{}
	for rows.Next() {
		var row pemberitahuan
		err := rows.Scan(
			&row.ID,
			&row.JudulBerita,
			&row.DeskripsiBerita,
			&row.Status,
			&row.DiperbaruiOleh,
			&row.TerakhirDiperbarui,
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

func (r *repository) count(ctx context.Context, cari string) (uint, error) {
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(1)").
		From("pemberitahuan p")

	if cari != "" {
		cari = cari + "%"
		qb = qb.Where("(p.judul_berita like ? or p.deskripsi_berita like ?)", cari, cari)
	}

	q, args, _ := qb.ToSql()
	var result uint
	err := r.db.QueryRowContext(ctx, q, args...).Scan(&result)

	return result, err
}
