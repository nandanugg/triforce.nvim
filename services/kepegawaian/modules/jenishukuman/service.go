package jenishukuman

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
	ListRefJenisHukuman(ctx context.Context, arg sqlc.ListRefJenisHukumanParams) ([]sqlc.ListRefJenisHukumanRow, error)
	CountRefJenisHukuman(ctx context.Context) (int64, error)
	CreateRefJenisHukuman(ctx context.Context, arg sqlc.CreateRefJenisHukumanParams) (sqlc.CreateRefJenisHukumanRow, error)
	DeleteRefJenisHukuman(ctx context.Context, id int32) (int64, error)
	GetRefJenisHukuman(ctx context.Context, id int32) (sqlc.GetRefJenisHukumanRow, error)
	UpdateRefJenisHukuman(ctx context.Context, arg sqlc.UpdateRefJenisHukumanParams) (sqlc.UpdateRefJenisHukumanRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]jenisHukuman, uint, error) {
	rows, err := s.repo.ListRefJenisHukuman(ctx, sqlc.ListRefJenisHukumanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJenisHukuman(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	data := typeutil.Map(rows, func(row sqlc.ListRefJenisHukumanRow) jenisHukuman {
		return jenisHukuman{
			ID:      row.ID,
			Nama:    row.Nama.String,
			Tingkat: row.Tingkat.String,
		}
	})
	return data, uint(count), nil
}

func (s *service) get(ctx context.Context, id int32) (*jenisHukuman, error) {
	row, err := s.repo.GetRefJenisHukuman(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetRefGolongan: %w", err)
	}

	result := &jenisHukuman{
		ID:      row.ID,
		Nama:    row.Nama.String,
		Tingkat: row.Tingkat.String,
	}

	return result, nil
}

type createParams struct {
	nama    string
	tingkat string
}

func (s *service) create(ctx context.Context, params createParams) (*jenisHukuman, error) {
	row, err := s.repo.CreateRefJenisHukuman(ctx, sqlc.CreateRefJenisHukumanParams{
		Nama:    pgtype.Text{String: params.nama, Valid: params.nama != ""},
		Tingkat: pgtype.Text{String: params.tingkat, Valid: params.tingkat != ""},
	})
	if err != nil {
		return nil, fmt.Errorf("[create] error createJenisHukuman: %w", err)
	}

	result := &jenisHukuman{
		ID:      row.ID,
		Nama:    row.Nama.String,
		Tingkat: row.Tingkat.String,
	}

	return result, nil
}

type updateParams struct {
	id      int32
	nama    string
	tingkat string
}

func (s *service) update(ctx context.Context, params updateParams) (*jenisHukuman, error) {
	row, err := s.repo.UpdateRefJenisHukuman(ctx, sqlc.UpdateRefJenisHukumanParams{
		ID:      params.id,
		Nama:    pgtype.Text{String: params.nama, Valid: params.nama != ""},
		Tingkat: pgtype.Text{String: params.tingkat, Valid: params.tingkat != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error updateJenisHukuman: %w", err)
	}

	result := &jenisHukuman{
		ID:      row.ID,
		Nama:    row.Nama.String,
		Tingkat: row.Tingkat.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteRefJenisHukuman(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error deleteJenisHukuman: %w", err)
	}
	return affected > 0, nil
}
