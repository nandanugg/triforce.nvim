package golongan

import (
	"context"
	"fmt"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetRefGolongan(ctx context.Context, arg repo.GetRefGolonganParams) ([]repo.GetRefGolonganRow, error)
	CountRefGolongan(ctx context.Context) (int64, error)
}
type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

type listRefGolonganParams struct {
	Limit  uint
	Offset uint
}

func (s *service) listRefGolongan(ctx context.Context, arg listRefGolonganParams) ([]refGolongan, int64, error) {
	rows, err := s.repo.GetRefGolongan(ctx, repo.GetRefGolonganParams{
		Limit:  int32(arg.Limit),
		Offset: int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listRefGolongan] error GetRefGolongan: %w", err)
	}

	result := []refGolongan{}
	for _, row := range rows {
		result = append(result, refGolongan{
			ID:          row.ID,
			Nama:        row.Nama.String,
			NamaPangkat: row.NamaPangkat.String,
		})
	}

	total, err := s.repo.CountRefGolongan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[listRefGolongan] error CountRefGolongan: %w", err)
	}

	return result, total, nil
}
