package jenisjabatan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJenisJabatan(ctx context.Context, arg repo.ListRefJenisJabatanParams) ([]repo.ListRefJenisJabatanRow, error)
	CountRefJenisJabatan(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listJenisJabatan(ctx context.Context, limit, offset uint) ([]jenisJabatan, int64, error) {
	data, err := s.repo.ListRefJenisJabatan(ctx, repo.ListRefJenisJabatanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisJabatan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row repo.ListRefJenisJabatanRow) jenisJabatan {
		return jenisJabatan{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), count, nil
}
