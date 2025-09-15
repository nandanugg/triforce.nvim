package pemberitahuan

import (
	"context"
	"fmt"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
)

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
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
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pemberitahuan{}
	for rows.Next() {
		var row pemberitahuan
		var terakhirDiperbarui pgtype.Date
		err := rows.Scan(
			&row.ID,
			&row.JudulBerita,
			&row.DeskripsiBerita,
			&row.Status,
			&row.DiperbaruiOleh,
			&terakhirDiperbarui,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		row.TerakhirDiperbarui = db.Date(terakhirDiperbarui.Time)
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
	err := r.db.QueryRow(ctx, q, args...).Scan(&result)

	return result, err
}
