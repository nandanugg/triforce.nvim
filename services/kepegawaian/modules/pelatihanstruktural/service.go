package pelatihanstruktural

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgtype"

	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListPelatihanStruktural(ctx context.Context, nipBaru pgtype.Text) ([]repo.ListPelatihanStrukturalRow, error)
}
type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, NIP string) ([]pelatihanStruktural, error) {
	result := []pelatihanStruktural{}
	data, err := s.repo.ListPelatihanStruktural(ctx, pgtype.Text{String: NIP, Valid: true})
	if err != nil {
		return nil, fmt.Errorf("repo list: %w", err)
	}

	for _, d := range data {
		result = append(result, pelatihanStruktural{
			ID:                    d.ID,
			JenisDiklat:           d.JenisDiklat.String,
			NamaDiklat:            d.NamaDiklat.String,
			IstitusiPenyelenggara: d.InstitusiPenyelenggara.String,
			NomorSertifikat:       d.NoSertifikat.String,
			TanggalMulai:          d.TanggalMulai.Time,
			TanggalSelesai:        d.TanggalSelesai.Time,
			Tahun:                 int(d.TahunDiklat.Int16),
			DurasiJam:             int(d.DurasiJam.Int16),
		})
	}

	return result, nil
}
