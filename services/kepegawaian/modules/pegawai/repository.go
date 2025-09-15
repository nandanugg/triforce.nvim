package pegawai

import (
	"context"
	"fmt"
	"strings"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"
)

type repository struct {
	db *pgxpool.Pool
}

func newRepository(db *pgxpool.Pool) *repository {
	return &repository{db: db}
}

func (r *repository) list(ctx context.Context, limit, offset uint64, opts listOptions) ([]pegawai, error) {
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select(
			`p."ID"`,
			`p."NIP_BARU"`,
			`p."NAMA"`,
			`p."GELAR_DEPAN"`,
			`p."GELAR_BELAKANG"`,
			`g."NAMA"`,
			`p."JABATAN_NAMA"`,
			`uk."NAMA_UNOR"`,
			`p."STATUS_CPNS_PNS"`,
		).
		From("pegawai p").
		Join(`golongan g on g."ID" = p."GOL_ID"`).
		Join(`unitkerja uk on uk."ID" = p."UNOR_ID"`).
		Where(`p."TMT_PENSIUN" is null`).
		Where(`p."NO_SK_PEMBERHENTIAN" is null`).
		Where(`p."TGL_MENINGGAL" is null`).
		OrderBy(`p."ID"`).
		Limit(limit).
		Offset(offset)

	if opts.cari != "" {
		cari := opts.cari + "%"
		qb = qb.Where(`(p."NAMA" like ? or p."NIP_BARU" like ? or p."JABATAN_NAMA" like ?)`, cari, cari, cari)
	}
	if opts.unitID != "" {
		qb = qb.Where(`p."UNOR_ID" = ?`, opts.unitID)
	}
	if opts.golonganID != 0 {
		qb = qb.Where(`p."GOL_ID" = ?`, opts.golonganID)
	}
	if opts.jabatanID != "" {
		qb = qb.Where(`p."JABATAN_ID" = ?`, opts.jabatanID)
	}
	if opts.status != "" {
		qb = qb.Where(`p."STATUS_CPNS_PNS" = ?`, opts.status)
	}

	q, args, _ := qb.ToSql()
	rows, err := r.db.Query(ctx, q, args...)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []pegawai{}
	for rows.Next() {
		var row pegawai
		err := rows.Scan(
			&row.ID,
			&row.NIP,
			&row.NamaPegawai,
			&row.GelarDepan,
			&row.GelarBelakang,
			&row.Golongan,
			&row.Jabatan,
			&row.UnitKerja,
			&row.StatusPegawai,
		)
		if err != nil {
			return nil, fmt.Errorf("row scan: %w", err)
		}

		row.Jabatan = strings.TrimSpace(row.Jabatan)

		result = append(result, row)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows scan: %w", err)
	}

	return result, nil
}

func (r *repository) count(ctx context.Context, opts listOptions) (uint, error) {
	qb := squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).
		Select("count(1)").
		From("pegawai p").
		Where(`p."TMT_PENSIUN" is null`).
		Where(`p."NO_SK_PEMBERHENTIAN" is null`).
		Where(`p."TGL_MENINGGAL" is null`)

	if opts.cari != "" {
		cari := opts.cari + "%"
		qb = qb.Where(`(p."NAMA" like ? or p."NIP_BARU" like ? or p."JABATAN_NAMA" like ?)`, cari, cari, cari)
	}
	if opts.unitID != "" {
		qb = qb.Where(`p."UNOR_ID" = ?`, opts.unitID)
	}
	if opts.golonganID != 0 {
		qb = qb.Where(`p."GOL_ID" = ?`, opts.golonganID)
	}
	if opts.jabatanID != "" {
		qb = qb.Where(`p."JABATAN_ID" = ?`, opts.jabatanID)
	}
	if opts.status != "" {
		qb = qb.Where(`p."STATUS_CPNS_PNS" = ?`, opts.status)
	}

	var result uint
	q, args, _ := qb.ToSql()
	err := r.db.QueryRow(ctx, q, args...).Scan(&result)

	return result, err
}

func (r *repository) listStatusPegawai(ctx context.Context) ([]statusPegawai, error) {
	rows, err := r.db.Query(ctx, `select "ID", "NAMA" from jenis_pegawai order by 2 asc`)
	if err != nil {
		return nil, fmt.Errorf("sql select: %w", err)
	}
	defer rows.Close()

	result := []statusPegawai{}
	for rows.Next() {
		var row statusPegawai
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
