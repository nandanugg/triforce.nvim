package jenissatker

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
	CountRefJenisSatker(ctx context.Context, nama pgtype.Text) (int64, error)
	CreateRefJenisSatker(ctx context.Context, nama pgtype.Text) (repo.CreateRefJenisSatkerRow, error)
	DeleteRefJenisSatker(ctx context.Context, id int32) (int64, error)
	GetRefJenisSatker(ctx context.Context, id int32) (repo.GetRefJenisSatkerRow, error)
	ListRefJenisSatker(ctx context.Context, arg repo.ListRefJenisSatkerParams) ([]repo.ListRefJenisSatkerRow, error)
	UpdateRefJenisSatker(ctx context.Context, arg repo.UpdateRefJenisSatkerParams) (repo.UpdateRefJenisSatkerRow, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

type listParams struct {
	limit  uint
	offset uint
	nama   string
}

func (s *service) list(ctx context.Context, params listParams) ([]jenisSatker, int64, error) {
	pgNama := pgtype.Text{String: params.nama, Valid: params.nama != ""}
	rows, err := s.repo.ListRefJenisSatker(ctx, repo.ListRefJenisSatkerParams{
		Limit:  int32(params.limit),
		Offset: int32(params.offset),
		Nama:   pgNama,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error listRefJenisSatker: %w", err)
	}

	total, err := s.repo.CountRefJenisSatker(ctx, pgNama)
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error countRefJenisSatker: %w", err)
	}

	return typeutil.Map(rows, func(row repo.ListRefJenisSatkerRow) jenisSatker {
		return jenisSatker{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), total, nil
}

func (s *service) get(ctx context.Context, id int32) (*jenisSatker, error) {
	row, err := s.repo.GetRefJenisSatker(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetRefGolongan: %w", err)
	}

	result := &jenisSatker{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type createParams struct {
	nama string
}

func (s *service) create(ctx context.Context, params createParams) (*jenisSatker, error) {
	row, err := s.repo.CreateRefJenisSatker(ctx, pgtype.Text{String: params.nama, Valid: params.nama != ""})
	if err != nil {
		return nil, fmt.Errorf("[create] error createJenisSatker: %w", err)
	}

	result := &jenisSatker{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

type updateParams struct {
	id   int32
	nama string
}

func (s *service) update(ctx context.Context, params updateParams) (*jenisSatker, error) {
	row, err := s.repo.UpdateRefJenisSatker(ctx, repo.UpdateRefJenisSatkerParams{
		ID:   params.id,
		Nama: pgtype.Text{String: params.nama, Valid: params.nama != ""},
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error updateJenisSatker: %w", err)
	}

	result := &jenisSatker{
		ID:   row.ID,
		Nama: row.Nama.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteRefJenisSatker(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error deleteJenisSatker: %w", err)
	}
	return affected > 0, nil
}
