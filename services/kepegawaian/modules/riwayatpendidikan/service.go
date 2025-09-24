package riwayatpendidikan

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	CountRiwayatPendidikan(ctx context.Context, nip pgtype.Text) (int64, error)
	ListRiwayatPendidikan(ctx context.Context, arg sqlc.ListRiwayatPendidikanParams) ([]sqlc.ListRiwayatPendidikanRow, error)
	GetBerkasRiwayatPendidikan(ctx context.Context, arg sqlc.GetBerkasRiwayatPendidikanParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPendidikan, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatPendidikan(ctx, sqlc.ListRiwayatPendidikanParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPendidikan(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatPendidikanRow) riwayatPendidikan {
		return riwayatPendidikan{
			ID:                   row.ID,
			JenjangPendidikan:    row.JenjangPendidikan.String,
			Pendidikan:           row.Pendidikan.String,
			NamaSekolah:          row.NamaSekolah.String,
			TahunLulus:           row.TahunLulus,
			NomorIjazah:          row.NoIjazah.String,
			GelarDepan:           row.GelarDepan.String,
			GelarBelakang:        row.GelarBelakang.String,
			TugasBelajar:         tugasBelajar[row.TugasBelajar.Int16],
			KeteranganPendidikan: row.NegaraSekolah.String,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPendidikan(ctx, sqlc.GetBerkasRiwayatPendidikanParams{
		Nip: pgNip,
		ID:  id,
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}
	if len(res.String) == 0 {
		return "", nil, nil
	}

	return api.GetMimeTypeAndDecodedData(res.String)
}
