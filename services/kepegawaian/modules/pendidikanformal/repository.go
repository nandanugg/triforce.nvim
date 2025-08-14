package pendidikanformal

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

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]pendidikanFormal, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			rp."ID",
			tp."NAMA",
			rp."NAMA_SEKOLAH",
			'TODO: Jurusan',
			rp."KETERANGAN_BERKAS",
			rp."TAHUN_LULUS",
			rp."NOMOR_IJASAH"
		from kepegawaian.rwt_pendidikan rp
		join kepegawaian.tkpendidikan tp on rp."TINGKAT_PENDIDIKAN_ID" = tp."ID"
		join kepegawaian.pegawai p on rp."PNS_ID" = p."PNS_ID"
		join kepegawaian.users u on p."NIP_BARU" = u.nip
		where u.id = $1
		order by rp."TAHUN_LULUS" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pendidikanFormal{}
	for rows.Next() {
		var row pendidikanFormal
		err := rows.Scan(
			&row.ID,
			&row.JenjangPendidikan,
			&row.NamaSekolah,
			&row.Jurusan,
			&row.KeteranganPendidikan,
			&row.TahunLulus,
			&row.NomorIjazah,
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
		from kepegawaian.rwt_pendidikan rp
		join kepegawaian.tkpendidikan tp on rp."TINGKAT_PENDIDIKAN_ID" = tp."ID"
		join kepegawaian.pegawai p on rp."PNS_ID" = p."PNS_ID"
		join kepegawaian.users u on p."NIP_BARU" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
