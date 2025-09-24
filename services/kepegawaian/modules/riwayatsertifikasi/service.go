package riwayatsertifikasi

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatSertifikasi(ctx context.Context, arg sqlc.ListRiwayatSertifikasiParams) ([]sqlc.ListRiwayatSertifikasiRow, error)
	CountRiwayatSertifikasi(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatSertifikasi(ctx context.Context, arg sqlc.GetBerkasRiwayatSertifikasiParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatSertifikasi, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatSertifikasi(ctx, sqlc.ListRiwayatSertifikasiParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatSertifikasi(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatSertifikasiRow) riwayatSertifikasi {
		return riwayatSertifikasi{
			ID:              row.ID,
			NamaSertifikasi: row.NamaSertifikasi.String,
			Tahun:           row.Tahun,
			Deskripsi:       row.Deskripsi.String,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatSertifikasi(ctx, sqlc.GetBerkasRiwayatSertifikasiParams{
		Nip: pgNip,
		ID:  id,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}
	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}
