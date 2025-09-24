package golongan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefGolongan(ctx context.Context, arg repo.ListRefGolonganParams) ([]repo.ListRefGolonganRow, error)
	CountRefGolongan(ctx context.Context) (int64, error)
}
type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) listRefGolongan(ctx context.Context, limit, offset uint) ([]refGolongan, int64, error) {
	rows, err := s.repo.ListRefGolongan(ctx, repo.ListRefGolonganParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listRefGolongan] error GetRefGolongan: %w", err)
	}

	total, err := s.repo.CountRefGolongan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[listRefGolongan] error CountRefGolongan: %w", err)
	}

	return typeutil.Map(rows, func(row repo.ListRefGolonganRow) refGolongan {
		return refGolongan{
			ID:          row.ID,
			Nama:        row.Nama.String,
			NamaPangkat: row.NamaPangkat.String,
		}
	}), total, nil
}
