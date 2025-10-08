package jenispenghargaan

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
	ListRefJenisPenghargaan(ctx context.Context, arg sqlc.ListRefJenisPenghargaanParams) ([]sqlc.ListRefJenisPenghargaanRow, error)
	CountRefJenisPenghargaan(ctx context.Context) (int64, error)
	CreateRefJenisPenghargaan(ctx context.Context, nama pgtype.Text) (sqlc.CreateRefJenisPenghargaanRow, error)
	DeleteRefJenisPenghargaan(ctx context.Context, id int32) (int64, error)
	GetRefJenisPenghargaan(ctx context.Context, id int32) (sqlc.GetRefJenisPenghargaanRow, error)
	UpdateRefJenisPenghargaan(ctx context.Context, arg sqlc.UpdateRefJenisPenghargaanParams) (sqlc.UpdateRefJenisPenghargaanRow, error)
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

func (s *service) get(ctx context.Context, id int32) (*jenisPenghargaan, error) {
	row, err := s.repo.GetRefJenisPenghargaan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error getRefJenisPenghargaan: %w", err)
	}

	result := &jenisPenghargaan{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type createParams struct {
	nama string
}

func (s *service) create(ctx context.Context, params createParams) (*jenisPenghargaan, error) {
	row, err := s.repo.CreateRefJenisPenghargaan(ctx, pgtype.Text{String: params.nama, Valid: params.nama != ""})
	if err != nil {
		return nil, fmt.Errorf("[create] error createRefJenisPenghargaan: %w", err)
	}

	result := &jenisPenghargaan{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type updateParams struct {
	id   int32
	nama string
}

func (s *service) update(ctx context.Context, params updateParams) (*jenisPenghargaan, error) {
	row, err := s.repo.UpdateRefJenisPenghargaan(ctx, sqlc.UpdateRefJenisPenghargaanParams{
		ID:   params.id,
		Nama: pgtype.Text{String: params.nama, Valid: params.nama != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error updateRefJenisPenghargaan: %w", err)
	}

	result := &jenisPenghargaan{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteRefJenisPenghargaan(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error deleteRefJenisPenghargaan: %w", err)
	}
	return affected > 0, nil
}
