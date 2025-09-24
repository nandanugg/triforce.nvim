package jenispenghargaan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJenisPenghargaan(ctx context.Context, arg sqlc.ListRefJenisPenghargaanParams) ([]sqlc.ListRefJenisPenghargaanRow, error)
	CountRefJenisPenghargaan(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisPenghargaan, uint, error) {
	rows, err := s.repo.ListRefJenisPenghargaan(ctx, sqlc.ListRefJenisPenghargaanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisPenghargaan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRefJenisPenghargaanRow) jenisPenghargaan {
		return jenisPenghargaan{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), uint(count), nil
}
