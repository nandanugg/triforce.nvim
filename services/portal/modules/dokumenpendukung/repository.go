package dokumenpendukung

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

func (r *repository) list(ctx context.Context) ([]dokumenPendukung, error) {
	rows, err := r.db.QueryContext(ctx, `
		select
			dp.id,
			dp.nama_tombol,
			dp.nama_halaman,
			dp.updated_by,
			dp.updated_at,
			case when dp.file is null then 'Belum Upload' else 'Sudah Upload' end as status
		from portal.dokumen_pendukung dp
	`)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []dokumenPendukung{}
	for rows.Next() {
		var row dokumenPendukung
		err := rows.Scan(
			&row.ID,
			&row.NamaTombol,
			&row.NamaHalaman,
			&row.DiperbaruiOleh,
			&row.TerakhirDiperbarui,
			&row.Status,
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
