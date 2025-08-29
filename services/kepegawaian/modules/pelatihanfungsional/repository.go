package pelatihanfungsional

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]pelatihanFungsional, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rdf."DIKLAT_FUNGSIONAL_ID",
			rdf."JENIS_DIKLAT",
			rdf."NAMA_KURSUS",
			rdf."TANGGAL_KURSUS",
			rdf."TAHUN",
			rdf."INSTITUSI_PENYELENGGARA",
			rdf."NOMOR_SERTIPIKAT"
		from rwt_diklat_fungsional rdf
		join pegawai p on rdf."NIP_BARU" = p."NIP_BARU"
		join users u on p."NIP_BARU" = u.nip
		where u.id = $1
		order by rdf."TANGGAL_KURSUS" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pelatihanFungsional{}
	for rows.Next() {
		var row pelatihanFungsional
		err := rows.Scan(
			&row.ID,
			&row.JenisDiklat,
			&row.NamaDiklat,
			&row.Tanggal,
			&row.Tahun,
			&row.InstitusiPenyelenggara,
			&row.NomorSertifikat,
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
		from rwt_diklat_fungsional rdf
		join pegawai p on rdf."NIP_BARU" = p."NIP_BARU"
		join users u on p."NIP_BARU" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
