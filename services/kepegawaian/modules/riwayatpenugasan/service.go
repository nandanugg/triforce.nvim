package riwayatpenugasan

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatPenugasan(ctx context.Context, arg sqlc.ListRiwayatPenugasanParams) ([]sqlc.ListRiwayatPenugasanRow, error)
	CountRiwayatPenugasan(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatPenugasan(ctx context.Context, arg sqlc.GetBerkasRiwayatPenugasanParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPenugasan, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatPenugasan(ctx, sqlc.ListRiwayatPenugasanParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPenugasan(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row sqlc.ListRiwayatPenugasanRow) riwayatPenugasan {
		isMenjabat := row.IsMenjabat.Bool || row.TanggalSelesai.Time.IsZero() || row.TanggalSelesai.Time.After(time.Now())
		return riwayatPenugasan{
			ID:               row.ID,
			TipeJabatan:      row.TipeJabatan.String,
			NameJabatan:      row.NamaJabatan.String,
			DeskripsiJabatan: row.DeskripsiJabatan.String,
			TanggalMulai:     db.Date(row.TanggalMulai.Time),
			TanggalSelesai:   db.Date(row.TanggalSelesai.Time),
			IsMenjabat:       isMenjabat,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPenugasan(ctx, sqlc.GetBerkasRiwayatPenugasanParams{
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
