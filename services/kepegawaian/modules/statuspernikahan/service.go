package statuspernikahan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJenisKawin(ctx context.Context, arg sqlc.ListRefJenisKawinParams) ([]sqlc.ListRefJenisKawinRow, error)
	CountRefJenisKawin(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]statusPernikahan, uint, error) {
	rows, err := s.repo.ListRefJenisKawin(ctx, sqlc.ListRefJenisKawinParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisKawin(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	data := typeutil.Map(rows, func(row sqlc.ListRefJenisKawinRow) statusPernikahan {
		return statusPernikahan{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	})
	return data, uint(count), nil
}
