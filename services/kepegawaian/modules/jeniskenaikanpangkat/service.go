package jeniskenaikanpangkat

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListJenisKenaikanPangkat(ctx context.Context, arg repo.ListJenisKenaikanPangkatParams) ([]repo.ListJenisKenaikanPangkatRow, error)
	CountJenisKenaikanPangkat(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisKenaikanPangkat, int64, error) {
	rows, err := s.repo.ListJenisKenaikanPangkat(ctx, repo.ListJenisKenaikanPangkatParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error listJenisKenaikanPangkat: %w", err)
	}

	total, err := s.repo.CountJenisKenaikanPangkat(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error countJenisKenaikanPangkat: %w", err)
	}

	return typeutil.Map(rows, func(row repo.ListJenisKenaikanPangkatRow) jenisKenaikanPangkat {
		return jenisKenaikanPangkat{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), total, nil
}
