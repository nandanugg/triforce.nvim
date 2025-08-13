package pelatihanteknis

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]pelatihanTeknis, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			'TODO: Tipe Diklat',
			rd.jenis_diklat,
			rd.nama_diklat,
			rd.durasi_jam,
			rd.tanggal_mulai,
			rd.nomor_sertifikat
		from kepegawaian.rwt_diklat rd
		join kepegawaian.pegawai p on rd.nip_baru = p."NIP_BARU"
		join kepegawaian.users u on p."NIP_BARU" = u.nip
		where u.id = $1
		order by rd.tanggal_mulai asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pelatihanTeknis{}
	for rows.Next() {
		var row pelatihanTeknis
		err := rows.Scan(
			&row.TipeDiklat,
			&row.JenisDiklat,
			&row.NamaDiklat,
			&row.JumlahJam,
			&row.TanggalDiklat,
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
		from kepegawaian.rwt_diklat rd
		join kepegawaian.pegawai p on rd.nip_baru = p."NIP_BARU"
		join kepegawaian.users u on p."NIP_BARU" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
