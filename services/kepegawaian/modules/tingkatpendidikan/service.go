package tingkatpendidikan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefTingkatPendidikan(ctx context.Context, arg sqlc.ListRefTingkatPendidikanParams) ([]sqlc.ListRefTingkatPendidikanRow, error)
	CountRefTingkatPendidikan(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]tingkatPendidikan, uint, error) {
	rows, err := s.repo.ListRefTingkatPendidikan(ctx, sqlc.ListRefTingkatPendidikanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefTingkatPendidikan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRefTingkatPendidikanRow) tingkatPendidikan {
		return tingkatPendidikan{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), uint(count), nil
}
