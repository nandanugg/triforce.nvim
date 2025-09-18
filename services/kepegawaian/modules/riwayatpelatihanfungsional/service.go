package riwayatpelatihanfungsional

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
	repo "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatPelatihanFungsional(ctx context.Context, arg repo.ListRiwayatPelatihanFungsionalParams) ([]repo.ListRiwayatPelatihanFungsionalRow, error)
	CountRiwayatPelatihanFungsional(ctx context.Context, nipBaru pgtype.Text) (int64, error)
	GetBerkasRiwayatPelatihanFungsional(ctx context.Context, arg repo.GetBerkasRiwayatPelatihanFungsionalParams) (pgtype.Text, error)
}
type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPelatihanFungsional, uint, error) {
	data, err := s.repo.ListRiwayatPelatihanFungsional(ctx, repo.ListRiwayatPelatihanFungsionalParams{
		NipBaru: pgtype.Text{String: nip, Valid: true},
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatPelatihanFungsionalRow) riwayatPelatihanFungsional {
		return riwayatPelatihanFungsional{
			ID:                     row.ID,
			JenisDiklat:            row.JenisDiklat.String,
			NamaDiklat:             row.NamaKursus.String,
			InstitusiPenyelenggara: row.InstitusiPenyelenggara.String,
			NomorSertifikat:        row.NoSertifikat.String,
			TanggalMulai:           db.Date(row.TanggalKursus.Time),
			TanggalSelesai:         db.Date(row.TanggalKursus.Time.Add(time.Duration(row.JumlahJam.Int32) * time.Hour)),
			Durasi:                 row.JumlahJam.Int32,
			Tahun:                  row.Tahun.Int16,
		}
	})

	total, err := s.repo.CountRiwayatPelatihanFungsional(ctx, pgtype.Text{String: nip, Valid: true})
	if err != nil {
		return nil, 0, fmt.Errorf("count list: %w", err)
	}

	return result, uint(total), nil
}

func (s *service) getBerkas(ctx context.Context, nip, id string) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatPelatihanFungsional(ctx, repo.GetBerkasRiwayatPelatihanFungsionalParams{
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
