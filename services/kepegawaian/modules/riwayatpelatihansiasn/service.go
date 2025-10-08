package riwayatpelatihansiasn

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
	ListRiwayatPelatihanSIASN(ctx context.Context, arg sqlc.ListRiwayatPelatihanSIASNParams) ([]sqlc.ListRiwayatPelatihanSIASNRow, error)
	CountRiwayatPelatihanSIASN(ctx context.Context, pnsNip pgtype.Text) (int64, error)
	GetBerkasRiwayatPelatihanSIASN(ctx context.Context, arg sqlc.GetBerkasRiwayatPelatihanSIASNParams) (pgtype.Text, error)
}
type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPelatihanSIASN, uint, error) {
	pnsNIP := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatPelatihanSIASN(ctx, sqlc.ListRiwayatPelatihanSIASNParams{
		NipBaru: pnsNIP,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPelatihanSIASN(ctx, pnsNIP)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row sqlc.ListRiwayatPelatihanSIASNRow) riwayatPelatihanSIASN {
		if !row.TahunDiklat.Valid && row.TanggalSelesai.Valid {
			row.TahunDiklat = pgtype.Int4{Int32: int32(row.TanggalSelesai.Time.Year()), Valid: true}
		}

		return riwayatPelatihanSIASN{
			ID:                     row.ID,
			JenisDiklat:            row.JenisDiklat.String,
			NamaDiklat:             row.NamaDiklat.String,
			InstitusiPenyelenggara: row.InstitusiPenyelenggara.String,
			NomorSertifikat:        row.NoSertifikat.String,
			TanggalMulai:           db.Date(row.TanggalMulai.Time),
			TanggalSelesai:         db.Date(row.TanggalSelesai.Time),
			Tahun:                  row.TahunDiklat,
			Durasi:                 row.DurasiJam,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatPelatihanSIASN(ctx, sqlc.GetBerkasRiwayatPelatihanSIASNParams{
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
