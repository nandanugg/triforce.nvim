package jenispegawai

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJenisPegawai(ctx context.Context, arg sqlc.ListRefJenisPegawaiParams) ([]sqlc.ListRefJenisPegawaiRow, error)
	CountRefJenisPegawai(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisPegawai, uint, error) {
	rows, err := s.repo.ListRefJenisPegawai(ctx, sqlc.ListRefJenisPegawaiParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisPegawai(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	data := typeutil.Map(rows, func(row sqlc.ListRefJenisPegawaiRow) jenisPegawai {
		return jenisPegawai{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	})
	return data, uint(count), nil
}
