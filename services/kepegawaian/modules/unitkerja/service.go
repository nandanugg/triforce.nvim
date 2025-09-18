package unitkerja

import (
	"context"
	"fmt"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListUnitKerjaByNamaOrInduk(ctx context.Context, arg repo.ListUnitKerjaByNamaOrIndukParams) ([]repo.ListUnitKerjaByNamaOrIndukRow, error)
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

	result := []unitKerja{}
	for _, row := range rows {
		result = append(result, unitKerja{
			ID:   row.ID,
			Nama: row.NamaUnor.String,
		})
	}

	total, err := s.repo.CountUnitKerja(ctx, repo.CountUnitKerjaParams{
		Nama:      arg.nama,
		UnorInduk: arg.unorInduk,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listUnitKerja] error countUnitKerja: %w", err)
	}
	return result, total, nil
}
