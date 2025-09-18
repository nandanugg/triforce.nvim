package riwayatpekerjaan

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context, userID int64, limit, offset uint) ([]riwayatPekerjaan, error) {
	rows, err := r.db.Query(ctx, `
		select
			rp."ID",
			rp."PNS_NIP",
			rp."JENIS_PERUSAHAAN",
			rp."NAMA_PERUSAHAAN",
			rp."SEBAGAI",
			rp."DARI_TANGGAL",
			rp."SAMPAI_TANGGAL",
			rp."PNS_ID",
			rp."KETERANGAN_BERKAS"
		from rwt_pekerjaan rp
		join users u on rp."PNS_NIP" = u.nip
		where u.id = $1
		order by rp."DARI_TANGGAL" asc
		limit $2 offset $3
		`, userID, limit, offset,
	)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []riwayatPekerjaan{}
	for rows.Next() {
		var row riwayatPekerjaan
		err := rows.Scan(
			&row.ID,
			&row.PNSNIP,
			&row.JenisPerusahaan,
			&row.NamaPerusahaan,
			&row.Sebagai,
			&row.DariTanggal,
			&row.SampaiTanggal,
			&row.PNSID,
			&row.KeteranganBerkas,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		for _, toTrim := range []*string{
			&row.PNSNIP,
			&row.JenisPerusahaan,
			&row.NamaPerusahaan,
			&row.Sebagai,
			&row.PNSID,
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
	err := r.db.QueryRow(ctx, `
		select count(1)
		from rwt_pekerjaan rp
		join users u on rp."PNS_NIP" = u.nip
		where u.id = $1
		`, userID,
	).Scan(&result)

	return result, err
}
