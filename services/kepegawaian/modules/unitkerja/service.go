package unitkerja

import (
	"context"
	"fmt"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetUnitKerjaByNamaOrInduk(ctx context.Context, arg repo.GetUnitKerjaByNamaOrIndukParams) ([]repo.GetUnitKerjaByNamaOrIndukRow, error)
	CountUnitKerja(ctx context.Context, arg repo.CountUnitKerjaParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(repo repository) *service {
	return &service{repo: repo}
}

type listUnitKerjaParams struct {
	Nama      string `db:"nama_unor"`
	UnorInduk string `db:"unor_induk"`
	Limit     uint
	Offset    uint
}

func (s *service) listUnitKerja(ctx context.Context, arg listUnitKerjaParams) ([]unitKerja, int64, error) {
	rows, err := s.repo.GetUnitKerjaByNamaOrInduk(ctx, repo.GetUnitKerjaByNamaOrIndukParams{
		UnorInduk: arg.UnorInduk,
		Limit:     int32(arg.Limit),
		Offset:    int32(arg.Offset),
		Nama:      arg.Nama,
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
		Nama:      arg.Nama,
		UnorInduk: arg.UnorInduk,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("[listUnitKerja] error countUnitKerja: %w", err)
	}
	return result, total, nil
}
