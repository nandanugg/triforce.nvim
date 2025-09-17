package jenisdiklat

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJenisDiklat(ctx context.Context, arg sqlc.ListRefJenisDiklatParams) ([]sqlc.ListRefJenisDiklatRow, error)
	CountRefJenisDiklat(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisDiklat, uint, error) {
	rows, err := s.repo.ListRefJenisDiklat(ctx, sqlc.ListRefJenisDiklatParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisDiklat(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	data := typeutil.Map(rows, func(row sqlc.ListRefJenisDiklatRow) jenisDiklat {
		return jenisDiklat{
			ID:   row.ID,
			Nama: row.JenisDiklat.String,
		}
	})
	return data, uint(count), nil
}
