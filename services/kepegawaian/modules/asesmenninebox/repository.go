package asesmenninebox

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]asesmenninebox, error) {
	rows, err := r.db.Query(ctx, `
		select
			anb."ID",
			anb."TAHUN",
			anb."KESIMPULAN"
		from rwt_nine_box anb
		join users u on anb."PNS_NIP" = u.nip
		where u.id = $1
		order by anb."TAHUN" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []asesmenninebox{}
	for rows.Next() {
		var row asesmenninebox
		err := rows.Scan(
			&row.ID,
			&row.Tahun,
			&row.Kesimpulan,
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
		from rwt_nine_box anb
		join users u on anb."PNS_NIP" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
