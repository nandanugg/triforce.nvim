package hukumandisiplin

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]hukumanDisiplin, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rh."ID",
			rh."NAMA_JENIS_HUKUMAN",
			rh."SK_NOMOR",
			rh."SK_TANGGAL",
			rh."TANGGAL_MULAI_HUKUMAN",
			rh."MASA_TAHUN",
			rh."MASA_BULAN"
		from kepegawaian.rwt_hukdis rh
		join kepegawaian.users u on rh."PNS_NIP" = u.nip
		where u.id = $1
		order by rh."TANGGAL_MULAI_HUKUMAN" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []hukumanDisiplin{}
	for rows.Next() {
		var row hukumanDisiplin
		err := rows.Scan(
			&row.ID,
			&row.JenisHukuman,
			&row.NomorSK,
			&row.TanggalSK,
			&row.TanggalMulai,
			&row.MasaTahun,
			&row.MasaBulan,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		for _, toTrim := range []*string{
			&row.JenisHukuman,
			&row.NomorSK,
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
		from kepegawaian.rwt_hukdis rh
		join kepegawaian.users u on rh."PNS_NIP" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
