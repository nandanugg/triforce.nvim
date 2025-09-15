package kepangkatan

import (
	"context"
	"fmt"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatKepangkatan(ctx context.Context, arg dbrepo.ListRiwayatKepangkatanParams) ([]dbrepo.ListRiwayatKepangkatanRow, error)
	CountRiwayatKepangkatan(ctx context.Context, pnsNip string) (int64, error)
}

type service struct {
	repo repository
}

type listRiwayatParams struct {
	Limit  uint
	Offset uint
	PnsNip string
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, params listRiwayatParams) ([]kepangkatan, uint, error) {
	data, err := s.repo.ListRiwayatKepangkatan(ctx, dbrepo.ListRiwayatKepangkatanParams{
		PnsNip: params.PnsNip,
		Limit:  int32(params.Limit),
		Offset: int32(params.Offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo ListRiwayatKepangkatan: %w", err)
	}

	count, err := s.repo.CountRiwayatKepangkatan(ctx, params.PnsNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo CountRiwayatKepangkatan: %w", err)
	}

	result := typeutil.Map(data, func(row dbrepo.ListRiwayatKepangkatanRow) kepangkatan {
		return kepangkatan{
			ID:                        row.ID,
			IDJenisKP:                 row.JenisKpID.Int32,
			NamaJenisKP:               row.NamaJenisKp.String,
			IDGolongan:                row.GolonganID.Int32,
			NamaGolongan:              row.NamaGolongan.String,
			NamaGolonganPangkat:       row.NamaGolonganPangkat.String,
			TMTGolongan:               db.Date(row.TmtGolongan.Time),
			SKNomor:                   row.SkNomor.String,
			SKTanggal:                 db.Date(row.SkTanggal.Time),
			MKGolonganTahun:           row.MkGolonganTahun.Int16,
			MKGolonganBulan:           row.MkGolonganBulan.Int16,
			NoBKN:                     row.NoBkn.String,
			TanggalBKN:                db.Date(row.TanggalBkn.Time),
			JumlahAngkaKreditTambahan: row.JumlahAngkaKreditTambahan.Int16,
			JumlahAngkaKreditUtama:    row.JumlahAngkaKreditUtama.Int16,
		}
	})

	return result, uint(count), nil
}
