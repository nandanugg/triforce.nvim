package jeniskenaikanpangkat

import (
	"context"
	"fmt"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListJenisKP(ctx context.Context, arg repo.ListJenisKPParams) ([]repo.ListJenisKPRow, error)
	CountJenisKP(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) listJenisKP(ctx context.Context, limit, offset uint) ([]jenisKp, int64, error) {
	rows, err := s.repo.ListJenisKP(ctx, repo.ListJenisKPParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listJenisKP] error getJenisKP: %w", err)
	}

	result := []jenisKp{}
	for _, row := range rows {
		result = append(result, jenisKp{
			ID:   row.ID,
			Nama: row.Nama.String,
		})
	}

	total, err := s.repo.CountJenisKP(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[listJenisKP] error countJenisKP: %w", err)
	}

	return result, total, nil
}
