package riwayatjabatan

import (
	"context"
	"fmt"
	"time"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatJabatan(ctx context.Context, arg repo.ListRiwayatJabatanParams) ([]repo.ListRiwayatJabatanRow, error)
	CountRiwayatJabatan(ctx context.Context, pnsNip string) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatJabatan, int64, error) {
	data, err := s.repo.ListRiwayatJabatan(ctx, repo.ListRiwayatJabatanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		PnsNip: nip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatJabatan(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatJabatanRow) riwayatJabatan {
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
