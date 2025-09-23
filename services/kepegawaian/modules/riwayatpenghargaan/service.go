package riwayatpenghargaan

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
	ListRiwayatPenghargaan(ctx context.Context, arg repo.ListRiwayatPenghargaanParams) ([]repo.ListRiwayatPenghargaanRow, error)
	CountRiwayatPenghargaan(ctx context.Context, nip string) (int64, error)
	GetBerkasRiwayatPenghargaan(ctx context.Context, arg repo.GetBerkasRiwayatPenghargaanParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPenghargaan, uint, error) {
	data, err := s.repo.ListRiwayatPenghargaan(ctx, repo.ListRiwayatPenghargaanParams{
		Nip:    nip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPenghargaan(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatPenghargaanRow) riwayatPenghargaan {
		return riwayatPenghargaan{
			ID:               int(row.ID),
			JenisPenghargaan: row.JenisPenghargaan.String,
			NamaPenghargaan:  row.NamaPenghargaan.String,
			Deskripsi:        row.DeskripsiPenghargaan.String,
			Tanggal:          db.Date(row.TanggalPenghargaan.Time),
		}
	})

	return result, uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPenghargaan(ctx, repo.GetBerkasRiwayatPenghargaanParams{
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
