package penghargaan

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type repository struct {
	db *sql.DB
}

func newRepository(db *sql.DB) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]penghargaan, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rp."ID",
			rp."NAMA_JENIS_PENGHARGAAN",
			rp."NAMA",
			rp."KETERANGAN",
			rp."SK_TANGGAL"
		from rwt_penghargaan rp
		join users u on rp."PNS_NIP" = u.nip
		where u.id = $1
		order by rp."SK_TANGGAL" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []penghargaan{}
	for rows.Next() {
		var row penghargaan
		err := rows.Scan(
			&row.ID,
			&row.JenisPenghargaan,
			&row.NamaPenghargaan,
			&row.Deskripsi,
			&row.Tanggal,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		for _, toTrim := range []*string{
			&row.NamaPenghargaan,
			&row.JenisPenghargaan,
		} {
			*toTrim = strings.TrimSpace(*toTrim)
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
		from rwt_penghargaan rp
		join users u on rp."PNS_NIP" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
