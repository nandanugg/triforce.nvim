package kepangkatan

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]kepangkatan, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rg."ID",
			rg."PANGKAT",
			rg."GOLONGAN",
			rg."TMT_GOLONGAN",
			rg."MK_GOLONGAN_TAHUN",
			rg."MK_GOLONGAN_BULAN"
		from rwt_golongan rg
		join pegawai p on rg."PNS_NIP" = p."NIP_BARU"
		join users u on p."NIP_BARU" = u.nip
		where u.id = $1
		order by rg."TMT_GOLONGAN" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []kepangkatan{}
	for rows.Next() {
		var row kepangkatan
		err := rows.Scan(
			&row.ID,
			&row.Pangkat,
			&row.Golongan,
			&row.TMT,
			&row.MKGolonganTahun,
			&row.MKGolonganBulan,
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
		from rwt_golongan rg
		join pegawai p on rg."PNS_NIP" = p."NIP_BARU"
		join users u on p."NIP_BARU" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
