package golongan

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRefGolongan(ctx context.Context, arg repo.ListRefGolonganParams) ([]repo.ListRefGolonganRow, error)
	CountRefGolongan(ctx context.Context) (int64, error)
	GetRefGolongan(ctx context.Context, id int32) (repo.GetRefGolonganRow, error)
	CreateRefGolongan(ctx context.Context, arg repo.CreateRefGolonganParams) (repo.CreateRefGolonganRow, error)
	UpdateRefGolongan(ctx context.Context, arg repo.UpdateRefGolonganParams) (repo.UpdateRefGolonganRow, error)
	DeleteRefGolongan(ctx context.Context, id int32) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listPublic(ctx context.Context, limit, offset uint) ([]golonganPublic, int64, error) {
	rows, err := s.repo.ListRefGolongan(ctx, repo.ListRefGolonganParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error ListRefGolongan: %w", err)
	}

	total, err := s.repo.CountRefGolongan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error CountRefGolongan: %w", err)
	}

	return typeutil.Map(rows, func(row repo.ListRefGolonganRow) golonganPublic {
		result := golonganPublic{
			ID:          row.ID,
			Nama:        row.Nama.String,
			NamaPangkat: row.NamaPangkat.String,
		}
		return result
	}), total, nil
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]golongan, int64, error) {
	rows, err := s.repo.ListRefGolongan(ctx, repo.ListRefGolonganParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error ListRefGolongan: %w", err)
	}

	total, err := s.repo.CountRefGolongan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[list] error CountRefGolongan: %w", err)
	}

	return typeutil.Map(rows, func(row repo.ListRefGolonganRow) golongan {
		result := golongan{
			ID:          row.ID,
			Nama:        row.Nama.String,
			NamaPangkat: row.NamaPangkat.String,
			Nama2:       row.Nama2.String,
			Gol:         row.Gol.Int16,
			GolPppk:     row.GolPppk.String,
		}
		return result
	}), total, nil
}

func (s *service) get(ctx context.Context, id int32) (*golongan, error) {
	row, err := s.repo.GetRefGolongan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[get] error GetRefGolongan: %w", err)
	}

	result := &golongan{
		ID:          row.ID,
		Nama:        row.Nama.String,
		NamaPangkat: row.NamaPangkat.String,
		Nama2:       row.Nama2.String,
		Gol:         row.Gol.Int16,
		GolPppk:     row.GolPppk.String,
	}

	return result, nil
}

type createParams struct {
	nama        string
	namaPangkat string
	nama2       string
	gol         int16
	golPppk     string
}

func (s *service) create(ctx context.Context, params createParams) (*golongan, error) {
	row, err := s.repo.CreateRefGolongan(ctx, repo.CreateRefGolonganParams{
		Nama:        params.nama,
		NamaPangkat: params.namaPangkat,
		Nama2:       params.nama2,
		Gol:         params.gol,
		GolPppk:     params.golPppk,
	})
	if err != nil {
		return nil, fmt.Errorf("[create] error CreateRefGolongan: %w", err)
	}

	result := &golongan{
		ID:          row.ID,
		Nama:        row.Nama.String,
		NamaPangkat: row.NamaPangkat.String,
		Nama2:       row.Nama2.String,
		Gol:         row.Gol.Int16,
		GolPppk:     row.GolPppk.String,
	}

	return result, nil
}

type updateParams struct {
	nama        string
	namaPangkat string
	nama2       string
	gol         int16
	golPppk     string
}

func (s *service) update(ctx context.Context, id int32, params updateParams) (*golongan, error) {
	row, err := s.repo.UpdateRefGolongan(ctx, repo.UpdateRefGolonganParams{
		ID:          id,
		Nama:        params.nama,
		NamaPangkat: params.namaPangkat,
		Nama2:       params.nama2,
		Gol:         params.gol,
		GolPppk:     params.golPppk,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error UpdateRefGolongan: %w", err)
	}

	result := &golongan{
		ID:          row.ID,
		Nama:        row.Nama.String,
		NamaPangkat: row.NamaPangkat.String,
		Nama2:       row.Nama2.String,
		Gol:         row.Gol.Int16,
		GolPppk:     row.GolPppk.String,
	}

	return result, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	affected, err := s.repo.DeleteRefGolongan(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error DeleteRefGolongan: %w", err)
	}
	return affected > 0, nil
}
