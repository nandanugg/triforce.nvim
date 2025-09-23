package riwayatkinerja

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatKinerja(ctx context.Context, arg sqlc.ListRiwayatKinerjaParams) ([]sqlc.ListRiwayatKinerjaRow, error)
	CountRiwayatKinerja(ctx context.Context, nip pgtype.Text) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatKinerja, uint, error) {
	pgNIP := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatKinerja(ctx, sqlc.ListRiwayatKinerjaParams{
		Nip:    pgNIP,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatKinerja(ctx, pgNIP)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row sqlc.ListRiwayatKinerjaRow) riwayatKinerja {
		return riwayatKinerja{
			ID:             row.ID,
			Tahun:          row.Tahun,
			HasilKinerja:   row.RatingHasilKerja.String,
			PerilakuKerja:  row.RatingPerilakuKerja.String,
			KuadranKinerja: row.PredikatKinerja.String,
		}
	}), uint(count), nil
}
