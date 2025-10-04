package riwayatpelatihanstruktural

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.ListRiwayatPelatihanStrukturalParams) ([]sqlc.ListRiwayatPelatihanStrukturalRow, error)
	CountRiwayatPelatihanStruktural(ctx context.Context, pnsNip pgtype.Text) (int64, error)
	GetBerkasRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.GetBerkasRiwayatPelatihanStrukturalParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPelatihanStruktural, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatPelatihanStruktural(ctx, sqlc.ListRiwayatPelatihanStrukturalParams{
		PnsNip: pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPelatihanStruktural(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatPelatihanStrukturalRow) riwayatPelatihanStruktural {
		return riwayatPelatihanStruktural{
			ID:         row.ID,
			NamaDiklat: row.NamaDiklat.String,
			Tanggal:    db.Date(row.Tanggal.Time),
			Nomor:      row.Nomor.String,
			Lama:       row.Lama,
			Tahun:      row.Tahun,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id string) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatPelatihanStruktural(ctx, sqlc.GetBerkasRiwayatPelatihanStrukturalParams{
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
