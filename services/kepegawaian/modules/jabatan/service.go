package jabatan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJabatan(ctx context.Context, arg repo.ListRefJabatanParams) ([]repo.ListRefJabatanRow, error)
	CountRefJabatan(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listJabatan(ctx context.Context, limit, offset uint) ([]jabatan, int64, error) {
	data, err := s.repo.ListRefJabatan(ctx, repo.ListRefJabatanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJabatan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRefJabatanRow) jabatan {
		return jabatan{
			NamaJabatan: row.NamaJabatan.String,
			KodeJabatan: row.KodeJabatan,
		}
	})

	return result, count, nil
}
