package riwayatkenaikangajiberkala

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.ListRiwayatKenaikanGajiBerkalaParams) ([]repo.ListRiwayatKenaikanGajiBerkalaRow, error)
	CountRiwayatKenaikanGajiBerkala(ctx context.Context, nipBaru pgtype.Text) (int64, error)
	GetBerkasRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.GetBerkasRiwayatKenaikanGajiBerkalaParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatKenaikanGajiBerkala, int64, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatKenaikanGajiBerkala(ctx, repo.ListRiwayatKenaikanGajiBerkalaParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		NipBaru: pgNip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatKenaikanGajiBerkala(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatKenaikanGajiBerkalaRow) riwayatKenaikanGajiBerkala {
		return riwayatKenaikanGajiBerkala{
			ID:                     row.ID,
			IDGolongan:             row.GolonganID.Int32,
			NamaGolongan:           row.GolonganNama.String,
			NamaGolonganPangkat:    row.GolonganNamaPangkat.String,
			NomorSK:                row.NoSk.String,
			TanggalSK:              db.Date(row.TanggalSk.Time),
			TMTGolongan:            db.Date(row.TmtGolongan.Time),
			MasaKerjaGolonganTahun: row.MasaKerjaGolonganTahun,
			MasaKerjaGolonganBulan: row.MasaKerjaGolonganBulan,
			TMTKenaikanGajiBerkala: db.Date(row.TmtKenaikanGajiBerkala.Time),
			GajiPokok:              row.GajiPokok,
			Jabatan:                row.Jabatan.String,
			TMTJabatan:             db.Date(row.TmtJabatan.Time),
			Pendidikan:             row.Pendidikan.String,
			TanggalLulus:           db.Date(row.TanggalLulus.Time),
			KantorPembayaran:       row.KantorPembayaran.String,
			UnitKerjaIndukID:       row.UnitKerjaIndukID.String,
			UnitKerjaInduk:         row.UnitKerjaInduk.String,
			Pejabat:                row.Pejabat.String,
		}
	})

	return result, count, nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatKenaikanGajiBerkala(ctx, repo.GetBerkasRiwayatKenaikanGajiBerkalaParams{
		NipBaru: pgtype.Text{String: nip, Valid: true},
		ID:      id,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}
	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}
