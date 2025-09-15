package pendidikanformal

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListPendidikanFormal(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListPendidikanFormalRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string) ([]pendidikanFormal, error) {
	result := []pendidikanFormal{}
	data, err := s.repo.ListPendidikanFormal(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	for _, d := range data {
		result = append(result, s.mapList(d))
	}

	return result, nil
}

func (s *service) mapList(d repo.ListPendidikanFormalRow) pendidikanFormal {
	return pendidikanFormal{
		ID:                   int(d.ID),
		JenjangPendidikan:    d.JenjangPendidikan.String,
		Pendidikan:           d.Pendidikan.String,
		NamaSekolah:          d.NamaSekolah.String,
		TahunLulus:           d.TahunLulus.String,
		NomorIjazah:          d.NoIjazah.String,
		GelarDepan:           d.GelarDepan.String,
		GelarBelakang:        d.GelarBelakang.String,
		TugasBelajar:         s.mapTugasBelajar(d.TugasBelajar),
		KeteranganPendidikan: d.NegaraSekolah.String,
	}
}

func (s *service) mapTugasBelajar(tugasBelajar pgtype.Int2) string {
	switch tugasBelajar.Int16 {
	case 1:
		return "Tugas Belajar"
	case 2:
		return "Izin Belajar"
	}
	return ""
}
