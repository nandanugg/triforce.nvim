package jenisjabatan

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
	ListRefJenisJabatan(ctx context.Context, arg repo.ListRefJenisJabatanParams) ([]repo.ListRefJenisJabatanRow, error)
	CountRefJenisJabatan(ctx context.Context) (int64, error)
	CreateRefJenisJabatan(ctx context.Context, nama pgtype.Text) (repo.CreateRefJenisJabatanRow, error)
	DeleteRefJenisJabatan(ctx context.Context, id int32) (int64, error)
	GetRefJenisJabatan(ctx context.Context, id int32) (repo.GetRefJenisJabatanRow, error)
	UpdateRefJenisJabatan(ctx context.Context, arg repo.UpdateRefJenisJabatanParams) (repo.UpdateRefJenisJabatanRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisJabatan, int64, error) {
	data, err := s.repo.ListRefJenisJabatan(ctx, repo.ListRefJenisJabatanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisJabatan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row repo.ListRefJenisJabatanRow) jenisJabatan {
		return jenisJabatan{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), count, nil
}

func (s *service) get(ctx context.Context, id int32) (*jenisJabatan, error) {
	row, err := s.repo.GetRefJenisJabatan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetRefJenisJabatan: %w", err)
	}

	result := &jenisJabatan{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type createParams struct {
	nama string
}

func (s *service) create(ctx context.Context, params createParams) (*jenisJabatan, error) {
	row, err := s.repo.CreateRefJenisJabatan(ctx, pgtype.Text{String: params.nama, Valid: params.nama != ""})
	if err != nil {
		return nil, fmt.Errorf("[create] error createJenisJabatan: %w", err)
	}

	result := &jenisJabatan{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type updateParams struct {
	id   int32
	nama string
}

func (s *service) update(ctx context.Context, params updateParams) (*jenisJabatan, error) {
	row, err := s.repo.UpdateRefJenisJabatan(ctx, repo.UpdateRefJenisJabatanParams{
		ID:   params.id,
		Nama: pgtype.Text{String: params.nama, Valid: params.nama != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error updateJenisJabatan: %w", err)
	}

	result := &jenisJabatan{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteRefJenisJabatan(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error deleteJenisJabatan: %w", err)
	}
	return affected > 0, nil
}
