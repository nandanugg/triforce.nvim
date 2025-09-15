package dokumenpendukung

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
)

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context) ([]dokumenPendukung, error) {
	rows, err := r.db.Query(ctx, `
		select
			dp.id,
			dp.nama_tombol,
			dp.nama_halaman,
			dp.updated_by,
			dp.updated_at,
			case when dp.file is null then 'Belum Upload' else 'Sudah Upload' end as status
		from dokumen_pendukung dp
	`)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []dokumenPendukung{}
	for rows.Next() {
		var row dokumenPendukung
		var terakhirDiperbarui pgtype.Date
		err := rows.Scan(
			&row.ID,
			&row.NamaTombol,
			&row.NamaHalaman,
			&row.DiperbaruiOleh,
			&terakhirDiperbarui,
			&row.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		if terakhirDiperbarui.Valid {
			row.TerakhirDiperbarui = typeutil.ToPtr(db.Date(terakhirDiperbarui.Time))
		}
		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}
