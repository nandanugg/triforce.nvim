package jabatan

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]jabatan, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rj."ID",
			rj."NAMA_JABATAN",
			'TODO: Unit Kerja',
			rj."TMT_JABATAN"
		from rwt_jabatan rj
		join users u on rj."PNS_NIP" = u.nip
		where u.id = $1
		order by rj."TMT_JABATAN" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []jabatan{}
	for rows.Next() {
		var row jabatan
		err := rows.Scan(
			&row.ID,
			&row.Jabatan,
			&row.UnitKerja,
			&row.TMT,
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

func (r *repository) listJenis(ctx context.Context) ([]jenisJabatan, error) {
	rows, err := r.db.QueryContext(ctx, `select "ID", "NAMA" from jenis_jabatan order by 2 asc`)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []jenisJabatan{}
	for rows.Next() {
		var row jenisJabatan
		err := rows.Scan(&row.ID, &row.Nama)
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
		from rwt_jabatan rj
		join users u on rj."PNS_NIP" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
