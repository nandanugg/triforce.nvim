package unitkerja

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListAkarUnitKerja(ctx context.Context, arg repo.ListAkarUnitKerjaParams) ([]repo.ListAkarUnitKerjaRow, error)
	ListUnitKerjaByDiatasanID(ctx context.Context, arg repo.ListUnitKerjaByDiatasanIDParams) ([]repo.ListUnitKerjaByDiatasanIDRow, error)
	ListUnitKerjaByNamaOrInduk(ctx context.Context, arg repo.ListUnitKerjaByNamaOrIndukParams) ([]repo.ListUnitKerjaByNamaOrIndukRow, error)
	CountAkarUnitKerja(ctx context.Context) (int64, error)
	CountUnitKerjaByDiatasanID(ctx context.Context, diatasanID pgtype.Text) (int64, error)
	CountUnitKerja(ctx context.Context, arg repo.CountUnitKerjaParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

type listUnitKerjaParams struct {
	nama      string
	unorInduk string
	limit     uint
	offset    uint
}

func (s *service) listUnitKerja(ctx context.Context, arg listUnitKerjaParams) ([]unitKerja, int64, error) {
	rows, err := s.repo.ListUnitKerjaByNamaOrInduk(ctx, repo.ListUnitKerjaByNamaOrIndukParams{
		UnorInduk: arg.unorInduk,
		Limit:     int32(arg.limit),
		Offset:    int32(arg.offset),
		Nama:      arg.nama,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listUnitKerja] error getUnitKerjaByNamaOrInduk: %w", err)
	}

	result := typeutil.Map(rows, func(row repo.ListUnitKerjaByNamaOrIndukRow) unitKerja {
		return unitKerja{
			ID:   row.ID,
			Nama: row.NamaUnor.String,
		}
	})

	total, err := s.repo.CountUnitKerja(ctx, repo.CountUnitKerjaParams{
		Nama:      arg.nama,
		UnorInduk: arg.unorInduk,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listUnitKerja] error countUnitKerja: %w", err)
	}
	return result, total, nil
}

type listAkarUnitKerjaParams struct {
	limit  uint
	offset uint
}

func (s *service) listAkarUnitKerja(ctx context.Context, arg listAkarUnitKerjaParams) ([]unitKerja, int64, error) {
	rows, err := s.repo.ListAkarUnitKerja(ctx, repo.ListAkarUnitKerjaParams{
		Limit:  int32(arg.limit),
		Offset: int32(arg.offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listAkarUnitKerja] error listAkarUnitKerja: %w", err)
	}

	result := typeutil.Map(rows, func(row repo.ListAkarUnitKerjaRow) unitKerja {
		return unitKerja{
			ID:   row.ID,
			Nama: row.NamaUnor.String,
		}
	})

	total, err := s.repo.CountAkarUnitKerja(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("[listAkarUnitKerja] error countAkarUnitKerja: %w", err)
	}
	return result, total, nil
}

type listAnakUnitKerjaParams struct {
	limit  uint
	offset uint
	id     string
}

func (s *service) listAnakUnitKerja(ctx context.Context, arg listAnakUnitKerjaParams) ([]anakUnitKerja, int64, error) {
	pgID := pgtype.Text{Valid: arg.id != "", String: arg.id}
	rows, err := s.repo.ListUnitKerjaByDiatasanID(ctx, repo.ListUnitKerjaByDiatasanIDParams{
		Limit:      int32(arg.limit),
		Offset:     int32(arg.offset),
		DiatasanID: pgID,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listAnakUnitKerja] error listUnitKerjaByDiatasanID: %w", err)
	}

	result := typeutil.Map(rows, func(row repo.ListUnitKerjaByDiatasanIDRow) anakUnitKerja {
		return anakUnitKerja{
			ID:      row.ID,
			Nama:    row.NamaUnor.String,
			HasAnak: row.HasAnak,
		}
	})

	total, err := s.repo.CountUnitKerjaByDiatasanID(ctx, pgID)
	if err != nil {
		return nil, 0, fmt.Errorf("[listAkarUnitKerja] error countUnitKerjaByDiatasanID: %w", err)
	}
	return result, total, nil
}
