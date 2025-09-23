package jabatan

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefJabatan(ctx context.Context, arg sqlc.ListRefJabatanParams) ([]sqlc.ListRefJabatanRow, error)
	CountRefJabatan(ctx context.Context, nama pgtype.Text) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listJabatan(ctx context.Context, nama string, limit, offset uint) ([]jabatan, int64, error) {
	pgNama := pgtype.Text{Valid: nama != "", String: nama}
	data, err := s.repo.ListRefJabatan(ctx, sqlc.ListRefJabatanParams{
		Nama:   pgNama,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJabatan(ctx, pgNama)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row sqlc.ListRefJabatanRow) jabatan {
		return jabatan{
			ID:   row.KodeJabatan,
			Nama: row.NamaJabatan.String,
		}
	})

	return result, count, nil
}
