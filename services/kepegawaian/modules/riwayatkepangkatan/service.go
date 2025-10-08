package riwayatkepangkatan

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	dbrepo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatKepangkatan(ctx context.Context, arg dbrepo.ListRiwayatKepangkatanParams) ([]dbrepo.ListRiwayatKepangkatanRow, error)
	CountRiwayatKepangkatan(ctx context.Context, pnsNip string) (int64, error)
	GetBerkasRiwayatKepangkatan(ctx context.Context, arg dbrepo.GetBerkasRiwayatKepangkatanParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatKepangkatan, uint, error) {
	data, err := s.repo.ListRiwayatKepangkatan(ctx, dbrepo.ListRiwayatKepangkatanParams{
		PnsNip: nip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo ListRiwayatKepangkatan: %w", err)
	}

	count, err := s.repo.CountRiwayatKepangkatan(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo CountRiwayatKepangkatan: %w", err)
	}

	result := typeutil.Map(data, func(row dbrepo.ListRiwayatKepangkatanRow) riwayatKepangkatan {
		return riwayatKepangkatan{
			ID:                        row.ID,
			IDJenisKP:                 row.JenisKpID,
			NamaJenisKP:               row.NamaJenisKp.String,
			IDGolongan:                row.GolonganID,
			NamaGolongan:              row.NamaGolongan.String,
			NamaGolonganPangkat:       row.NamaGolonganPangkat.String,
			TMTGolongan:               db.Date(row.TmtGolongan.Time),
			SKNomor:                   row.SkNomor.String,
			SKTanggal:                 db.Date(row.SkTanggal.Time),
			MKGolonganTahun:           row.MkGolonganTahun,
			MKGolonganBulan:           row.MkGolonganBulan,
			NoBKN:                     row.NoBkn.String,
			TanggalBKN:                db.Date(row.TanggalBkn.Time),
			JumlahAngkaKreditTambahan: row.JumlahAngkaKreditTambahan,
			JumlahAngkaKreditUtama:    row.JumlahAngkaKreditUtama,
		}
	})

	return result, uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id string) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatKepangkatan(ctx, dbrepo.GetBerkasRiwayatKepangkatanParams{
		PnsNip: pgtype.Text{String: nip, Valid: true},
		ID:     id,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}
	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}
