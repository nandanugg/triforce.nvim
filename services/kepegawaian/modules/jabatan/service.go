package jabatan

import (
	"context"
	"fmt"
	"time"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
	utility "gitlab.com/wartek-id/matk/nexus/nexus-be/utils"
)

type repository interface {
	ListRefJabatan(ctx context.Context, arg repo.ListRefJabatanParams) ([]repo.ListRefJabatanRow, error)
	CountRefJabatan(ctx context.Context) (int64, error)
	ListRiwayatJabatan(ctx context.Context, arg repo.ListRiwayatJabatanParams) ([]repo.ListRiwayatJabatanRow, error)
	CountRiwayatJabatan(ctx context.Context) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

type listParams struct {
	Limit  uint
	Offset uint
}

type listRiwayatJabatanParams struct {
	Limit  uint
	Offset uint
	NIP    string
}

func (s *service) listJabatan(ctx context.Context, arg listParams) ([]jabatan, int64, error) {
	data, err := s.repo.ListRefJabatan(ctx, repo.ListRefJabatanParams{
		Limit:  int32(arg.Limit),
		Offset: int32(arg.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRefJabatan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := utility.SlimMap(data, func(row repo.ListRefJabatanRow) jabatan {
		return jabatan{
			ID:          row.ID,
			NamaJabatan: row.NamaJabatan.String,
			KodeJabatan: row.KodeJabatan,
		}
	})

	return result, count, nil
}

func (s *service) listRiwayatJabatan(ctx context.Context, arg listRiwayatJabatanParams) ([]riwayatJabatan, int64, error) {
	data, err := s.repo.ListRiwayatJabatan(ctx, repo.ListRiwayatJabatanParams{
		Limit:  int32(arg.Limit),
		Offset: int32(arg.Offset),
		PnsNip: arg.NIP,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatJabatan(ctx)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := utility.SlimMap(data, func(row repo.ListRiwayatJabatanRow) riwayatJabatan {
		return riwayatJabatan{
			ID:                      row.ID,
			JenisJabatan:            row.JenisJabatan.String,
			NamaJabatan:             row.NamaJabatan.String,
			TmtJabatan:              row.TmtJabatan.Time.Format(time.DateOnly),
			NoSk:                    row.NoSk.String,
			TanggalSk:               row.TanggalSk.Time.Format(time.DateOnly),
			SatuanKerja:             row.SatuanKerja.String,
			UnitOrganisasi:          row.UnitOrganisasi.String,
			StatusPlt:               row.StatusPlt.Bool,
			KelasJabatan:            row.KelasJabatan.String,
			PeriodeJabatanStartDate: row.PeriodeJabatanStartDate.Time.Format(time.DateOnly),
			PeriodeJabatanEndDate:   row.PeriodeJabatanEndDate.Time.Format(time.DateOnly),
		}
	})

	return result, count, nil
}
