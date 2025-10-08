package jeniskenaikanpangkat

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	CountJenisKenaikanPangkat(ctx context.Context) (int64, error)
	CreateJenisKenaikanPangkat(ctx context.Context, nama pgtype.Text) (repo.CreateJenisKenaikanPangkatRow, error)
	DeleteJenisKenaikanPangkat(ctx context.Context, id int32) (int64, error)
	GetJenisKenaikanPangkat(ctx context.Context, id int32) (repo.GetJenisKenaikanPangkatRow, error)
	ListJenisKenaikanPangkat(ctx context.Context, arg repo.ListJenisKenaikanPangkatParams) ([]repo.ListJenisKenaikanPangkatRow, error)
	UpdateJenisKenaikanPangkat(ctx context.Context, arg repo.UpdateJenisKenaikanPangkatParams) (repo.UpdateJenisKenaikanPangkatRow, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisKenaikanPangkat, int64, error) {
	rows, err := s.repo.ListJenisKenaikanPangkat(ctx, repo.ListJenisKenaikanPangkatParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error listJenisKenaikanPangkat: %w", err)
	}

	total, err := s.repo.CountJenisKenaikanPangkat(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error countJenisKenaikanPangkat: %w", err)
	}

	return typeutil.Map(rows, func(row repo.ListJenisKenaikanPangkatRow) jenisKenaikanPangkat {
		return jenisKenaikanPangkat{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), total, nil
}

func (s *service) get(ctx context.Context, id int32) (*jenisKenaikanPangkat, error) {
	row, err := s.repo.GetJenisKenaikanPangkat(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetRefGolongan: %w", err)
	}

	result := &jenisKenaikanPangkat{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type createParams struct {
	nama string
}

func (s *service) create(ctx context.Context, params createParams) (*jenisKenaikanPangkat, error) {
	row, err := s.repo.CreateJenisKenaikanPangkat(ctx, pgtype.Text{String: params.nama, Valid: params.nama != ""})
	if err != nil {
		return nil, fmt.Errorf("[create] error createJenisKenaikanPangkat: %w", err)
	}

	result := &jenisKenaikanPangkat{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type updateParams struct {
	id   int32
	nama string
}

func (s *service) update(ctx context.Context, params updateParams) (*jenisKenaikanPangkat, error) {
	row, err := s.repo.UpdateJenisKenaikanPangkat(ctx, repo.UpdateJenisKenaikanPangkatParams{
		ID:   params.id,
		Nama: pgtype.Text{String: params.nama, Valid: params.nama != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error updateJenisKenaikanPangkat: %w", err)
	}

	result := &jenisKenaikanPangkat{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteJenisKenaikanPangkat(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error deleteJenisKenaikanPangkat: %w", err)
	}
	return affected > 0, nil
}
