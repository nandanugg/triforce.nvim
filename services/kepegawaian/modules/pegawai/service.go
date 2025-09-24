package pegawai

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	GetProfilePegawaiByPNSID(ctx context.Context, pnsID string) (sqlc.GetProfilePegawaiByPNSIDRow, error)
	ListUnitKerjaHierarchy(ctx context.Context, id string) ([]sqlc.ListUnitKerjaHierarchyRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) getProfileByPNSID(ctx context.Context, pnsID string) (*profile, error) {
	data, err := s.repo.GetProfilePegawaiByPNSID(ctx, pnsID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("repo get profile: %w", err)
	}

	unitOrganisasi := make([]string, 0)
	if data.UnorID.Valid {
		rows, err := s.repo.ListUnitKerjaHierarchy(ctx, data.UnorID.String)
		if err != nil {
			return nil, fmt.Errorf("repo list unit kerja hierarchy: %w", err)
		}

		unitOrganisasi = typeutil.FilterMap(rows, func(row sqlc.ListUnitKerjaHierarchyRow) (string, bool) {
			return row.NamaUnor.String, row.NamaUnor.String != ""
		})
	}

	return &profile{
		NIPLama:        data.NipLama.String,
		NIPBaru:        data.NipBaru.String,
		GelarDepan:     data.GelarDepan.String,
		GelarBelakang:  data.GelarBelakang.String,
		Nama:           data.Nama.String,
		Pangkat:        data.Pangkat.String,
		Golongan:       typeutil.Cast[string](data.Golongan),
		Jabatan:        data.Jabatan.String,
		UnitOrganisasi: unitOrganisasi,
	}, nil
}
