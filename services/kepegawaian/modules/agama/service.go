package agama

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefAgama(ctx context.Context, arg sqlc.ListRefAgamaParams) ([]sqlc.ListRefAgamaRow, error)
	CountRefAgama(ctx context.Context) (int64, error)
	GetRefAgama(ctx context.Context, id int32) (sqlc.GetRefAgamaRow, error)
	CreateRefAgama(ctx context.Context, pg pgtype.Text) (sqlc.CreateRefAgamaRow, error)
	UpdateRefAgama(ctx context.Context, arg sqlc.UpdateRefAgamaParams) (sqlc.UpdateRefAgamaRow, error)
	DeleteRefAgama(ctx context.Context, id int32) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]agama, uint, error) {
	rows, err := s.repo.ListRefAgama(ctx, sqlc.ListRefAgamaParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefAgama(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(r sqlc.ListRefAgamaRow) agama {
		return agama{
			ID:        r.ID,
			Nama:      r.Nama.String,
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		}
	}), uint(count), nil
}

func (s *service) get(ctx context.Context, id int32) (*agama, error) {
	r, err := s.repo.GetRefAgama(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get: %w", err)
	}
	return &agama{
		ID:        r.ID,
		Nama:      r.Nama.String,
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func (s *service) create(ctx context.Context, nama string) (*agama, error) {
	r, err := s.repo.CreateRefAgama(ctx, pgtype.Text{Valid: true, String: nama})
	if err != nil {
		return nil, fmt.Errorf("repo create: %w", err)
	}
	return &agama{ID: r.ID, Nama: r.Nama.String, CreatedAt: r.CreatedAt, UpdatedAt: r.UpdatedAt}, nil
}

func (s *service) update(ctx context.Context, id int32, nama string) (*agama, error) {
	r, err := s.repo.UpdateRefAgama(ctx, sqlc.UpdateRefAgamaParams{
		ID: id, Nama: pgtype.Text{Valid: true, String: nama},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo update: %w", err)
	}
	return &agama{ID: r.ID, Nama: r.Nama.String, UpdatedAt: r.UpdatedAt, CreatedAt: r.CreatedAt}, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteRefAgama(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error DeleteRefAgama: %w", err)
	}
	return affected > 0, nil
}
