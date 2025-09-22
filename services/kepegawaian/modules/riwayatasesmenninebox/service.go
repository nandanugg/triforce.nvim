package riwayatasesmenninebox

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	CountRiwayatAsesmenNineBox(ctx context.Context, pnsNip pgtype.Text) (int64, error)
	ListRiwayatAsesmenNineBox(ctx context.Context, arg sqlc.ListRiwayatAsesmenNineBoxParams) ([]sqlc.ListRiwayatAsesmenNineBoxRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatAsesmenNineBox, uint, error) {
	pnsNIP := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatAsesmenNineBox(ctx, sqlc.ListRiwayatAsesmenNineBoxParams{
		PnsNip: pnsNIP,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatAsesmenNineBox(ctx, pnsNIP)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row sqlc.ListRiwayatAsesmenNineBoxRow) riwayatAsesmenNineBox {
		return riwayatAsesmenNineBox{
			ID:         int(row.ID),
			Tahun:      row.Tahun,
			Kesimpulan: row.Kesimpulan.String,
		}
	}), uint(count), nil
}
