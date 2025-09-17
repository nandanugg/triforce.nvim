package tingkatpendidikan

import (
	"context"
	"fmt"

	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefTingkatPendidikan(ctx context.Context) ([]sqlc.ListRefTingkatPendidikanRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context) ([]tingkatPendidikan, error) {
	rows, err := s.repo.ListRefTingkatPendidikan(ctx)
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	data := make([]tingkatPendidikan, 0, len(rows))
	for _, row := range rows {
		data = append(data, tingkatPendidikan{
			ID:   row.ID,
			Nama: row.Nama.String,
		})
	}
	return data, nil
}
