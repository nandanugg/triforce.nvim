package tingkatpendidikan

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
	ListRefTingkatPendidikan(ctx context.Context, arg sqlc.ListRefTingkatPendidikanParams) ([]sqlc.ListRefTingkatPendidikanRow, error)
	CountRefTingkatPendidikan(ctx context.Context) (int64, error)
	GetRefTingkatPendidikan(ctx context.Context, id int32) (sqlc.GetRefTingkatPendidikanRow, error)
	CreateRefTingkatPendidikan(ctx context.Context, arg sqlc.CreateRefTingkatPendidikanParams) (sqlc.CreateRefTingkatPendidikanRow, error)
	UpdateRefTingkatPendidikan(ctx context.Context, arg sqlc.UpdateRefTingkatPendidikanParams) (sqlc.UpdateRefTingkatPendidikanRow, error)
	DeleteRefTingkatPendidikan(ctx context.Context, id int32) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) listPublic(ctx context.Context, limit, offset uint) ([]tingkatPendidikanPublic, uint, error) {
	rows, err := s.repo.ListRefTingkatPendidikan(ctx, sqlc.ListRefTingkatPendidikanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefTingkatPendidikan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRefTingkatPendidikanRow) tingkatPendidikanPublic {
		return tingkatPendidikanPublic{
			ID:   row.ID,
			Nama: row.Nama.String,
		}
	}), uint(count), nil
}

func (s *service) list(ctx context.Context, limit, offset uint) ([]tingkatPendidikan, uint, error) {
	rows, err := s.repo.ListRefTingkatPendidikan(ctx, sqlc.ListRefTingkatPendidikanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefTingkatPendidikan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRefTingkatPendidikanRow) tingkatPendidikan {
		return tingkatPendidikan{
			ID:               row.ID,
			Nama:             row.Nama.String,
			Abbreviation:     row.Abbreviation,
			GolonganID:       row.GolonganID,
			GolonganAwalID:   row.GolonganAwalID,
			Tingkat:          row.Tingkat,
			NamaGolongan:     &row.NamaGolongan,
			NamaGolonganAwal: &row.NamaGolonganAwal,
		}
	}), uint(count), nil
}

func (s *service) get(ctx context.Context, id int32) (*tingkatPendidikan, error) {
	row, err := s.repo.GetRefTingkatPendidikan(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("error GetRefTingkatPendidikan: %w", err)
	}

	return &tingkatPendidikan{
		ID:             row.ID,
		Nama:           row.Nama.String,
		Abbreviation:   row.Abbreviation,
		GolonganID:     row.GolonganID,
		GolonganAwalID: row.GolonganAwalID,
		Tingkat:        row.Tingkat,
	}, nil
}

type createParams struct {
	nama           string
	abbreviation   *string
	golonganID     *int32
	golonganAwalID *int32
	tingkat        *int16
}

func (s *service) create(ctx context.Context, params createParams) (*tingkatPendidikan, error) {
	row, err := s.repo.CreateRefTingkatPendidikan(ctx, sqlc.CreateRefTingkatPendidikanParams{
		Nama:           pgtype.Text{String: params.nama, Valid: params.nama != ""},
		Abbreviation:   typeutil.PointerToPgtype(params.abbreviation).(pgtype.Text),
		GolonganID:     typeutil.PointerToPgtype(params.golonganID).(pgtype.Int4),
		GolonganAwalID: typeutil.PointerToPgtype(params.golonganAwalID).(pgtype.Int4),
		Tingkat:        typeutil.PointerToPgtype(params.tingkat).(pgtype.Int2),
	})
	if err != nil {
		return nil, fmt.Errorf("[create] error CreateRefTingkatPendidikan: %w", err)
	}

	return &tingkatPendidikan{
		ID:             row.ID,
		Nama:           row.Nama.String,
		Abbreviation:   row.Abbreviation,
		GolonganID:     row.GolonganID,
		GolonganAwalID: row.GolonganAwalID,
		Tingkat:        row.Tingkat,
	}, nil
}

type updateParams struct {
	id             int32
	nama           string
	abbreviation   *string
	golonganID     *int32
	golonganAwalID *int32
	tingkat        *int16
}

func (s *service) update(ctx context.Context, params updateParams) (*tingkatPendidikan, error) {
	row, err := s.repo.UpdateRefTingkatPendidikan(ctx, sqlc.UpdateRefTingkatPendidikanParams{
		ID:             params.id,
		Nama:           pgtype.Text{String: params.nama, Valid: params.nama != ""},
		Abbreviation:   typeutil.PointerToPgtype(params.abbreviation).(pgtype.Text),
		GolonganID:     typeutil.PointerToPgtype(params.golonganID).(pgtype.Int4),
		GolonganAwalID: typeutil.PointerToPgtype(params.golonganAwalID).(pgtype.Int4),
		Tingkat:        typeutil.PointerToPgtype(params.tingkat).(pgtype.Int2),
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("[update] error UpdateRefTingkatPendidikan: %w", err)
	}

	return &tingkatPendidikan{
		ID:             row.ID,
		Nama:           row.Nama.String,
		Abbreviation:   row.Abbreviation,
		GolonganID:     row.GolonganID,
		GolonganAwalID: row.GolonganAwalID,
		Tingkat:        row.Tingkat,
	}, nil
}

func (s *service) delete(ctx context.Context, id int32) (bool, error) {
	rowsAffected, err := s.repo.DeleteRefTingkatPendidikan(ctx, id)
	if err != nil {
		return false, fmt.Errorf("[delete] error DeleteRefTingkatPendidikan: %w", err)
	}
	return rowsAffected > 0, nil
}
