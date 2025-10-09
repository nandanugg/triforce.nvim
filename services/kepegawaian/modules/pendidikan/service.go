package pendidikan

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
	ListRefPendidikanWithTingkatPendidikan(ctx context.Context, arg sqlc.ListRefPendidikanWithTingkatPendidikanParams) ([]sqlc.ListRefPendidikanWithTingkatPendidikanRow, error)
	CountRefPendidikan(ctx context.Context, nama pgtype.Text) (int64, error)
	GetRefPendidikan(ctx context.Context, id string) (sqlc.GetRefPendidikanRow, error)
	CreateRefPendidikan(ctx context.Context, arg sqlc.CreateRefPendidikanParams) (string, error)
	DeleteRefPendidikan(ctx context.Context, id string) (int64, error)
	UpdateRefPendidikan(ctx context.Context, arg sqlc.UpdateRefPendidikanParams) (string, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

func (s *service) list(ctx context.Context, limit, offset uint, nama string) ([]pendidikan, uint, error) {
	data, err := s.repo.ListRefPendidikanWithTingkatPendidikan(ctx, sqlc.ListRefPendidikanWithTingkatPendidikanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		Nama:   pgtype.Text{String: nama, Valid: nama != ""},
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefPendidikan(ctx, pgtype.Text{String: nama, Valid: nama != ""})
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row sqlc.ListRefPendidikanWithTingkatPendidikanRow) pendidikan {
		return pendidikan{
			ID:                  row.ID,
			Nama:                row.Nama,
			TingkatPendidikan:   row.TingkatPendidikan,
			TingkatPendidikanID: row.TingkatPendidikanID,
		}
	})

	return result, uint(count), nil
}

func (s *service) get(ctx context.Context, id string) (*pendidikan, error) {
	row, err := s.repo.GetRefPendidikan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get: %w", err)
	}

	return &pendidikan{
		ID:                  row.ID,
		Nama:                row.Nama,
		TingkatPendidikan:   row.TingkatPendidikan,
		TingkatPendidikanID: row.TingkatPendidikanID,
	}, nil
}

type createParams struct {
	nama                string
	tingkatPendidikanID int32
}

func (s *service) create(ctx context.Context, params createParams) (*pendidikan, error) {
	id, err := generateID()
	if err != nil {
		return nil, fmt.Errorf("repo create: %w", err)
	}
	createdID, err := s.repo.CreateRefPendidikan(ctx, sqlc.CreateRefPendidikanParams{
		ID:                  id,
		Nama:                params.nama,
		TingkatPendidikanID: pgtype.Int2{Int16: int16(params.tingkatPendidikanID), Valid: params.tingkatPendidikanID != 0},
	})
	if err != nil {
		return nil, fmt.Errorf("repo create: %w", err)
	}

	row, err := s.repo.GetRefPendidikan(ctx, createdID)
	if err != nil {
		return nil, fmt.Errorf("repo get: %w", err)
	}

	return &pendidikan{
		ID:                  row.ID,
		Nama:                row.Nama,
		TingkatPendidikanID: row.TingkatPendidikanID,
		TingkatPendidikan:   row.TingkatPendidikan,
	}, nil
}

type updateParams struct {
	nama                string
	tingkatPendidikanID int32
}

func (s *service) update(ctx context.Context, id string, params updateParams) (*pendidikan, error) {
	_, err := s.repo.UpdateRefPendidikan(ctx, sqlc.UpdateRefPendidikanParams{
		ID:                  id,
		Nama:                params.nama,
		TingkatPendidikanID: pgtype.Int2{Int16: int16(params.tingkatPendidikanID), Valid: params.tingkatPendidikanID != 0},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo update: %w", err)
	}

	row, err := s.repo.GetRefPendidikan(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("repo get: %w", err)
	}

	return &pendidikan{
		ID:                  row.ID,
		Nama:                row.Nama,
		TingkatPendidikanID: row.TingkatPendidikanID,
		TingkatPendidikan:   row.TingkatPendidikan,
	}, nil
}

func (s *service) delete(ctx context.Context, id string) (bool, error) {
	affected, err := s.repo.DeleteRefPendidikan(ctx, id)
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}
	return affected > 0, nil
}
