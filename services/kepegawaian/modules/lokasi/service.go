package lokasi

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefLokasi(ctx context.Context, arg sqlc.ListRefLokasiParams) ([]sqlc.ListRefLokasiRow, error)
	CountRefLokasi(ctx context.Context, nama string) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nama string, limit, offset uint) ([]lokasi, uint, error) {
	rows, err := s.repo.ListRefLokasi(ctx, sqlc.ListRefLokasiParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Nama:   nama,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo lokasi: %w", err)
	}

	count, err := s.repo.CountRefLokasi(ctx, nama)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(r sqlc.ListRefLokasiRow) lokasi {
		return lokasi{
			ID:   r.ID,
			Nama: r.Nama.String,
		}
	}), uint(count), nil
}
