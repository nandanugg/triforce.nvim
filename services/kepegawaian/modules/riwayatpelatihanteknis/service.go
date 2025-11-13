package riwayatpelatihanteknis

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
	ListRiwayatPelatihanTeknis(ctx context.Context, arg sqlc.ListRiwayatPelatihanTeknisParams) ([]sqlc.ListRiwayatPelatihanTeknisRow, error)
	CountRiwayatPelatihanTeknis(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatPelatihanTeknis(ctx context.Context, arg sqlc.GetBerkasRiwayatPelatihanTeknisParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPelatihanTeknis, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatPelatihanTeknis(ctx, sqlc.ListRiwayatPelatihanTeknisParams{
		PnsNip: pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPelatihanTeknis(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatPelatihanTeknisRow) riwayatPelatihanTeknis {
		var tanggalSelesai time.Time
		var tahun *int
		if row.TanggalKursus.Valid {
			tanggalSelesai = row.TanggalKursus.Time.Add(time.Duration(row.Durasi.Float64) * time.Hour)
			tahun = typeutil.ToPtr(row.TanggalKursus.Time.Year())
		}

		return riwayatPelatihanTeknis{
			ID:                     int64(row.ID),
			TipePelatihan:          tipePelatihan(row.TipeKursus.String),
			JenisPelatihan:         row.JenisKursus.String,
			NamaPelatihan:          row.NamaKursus.String,
			TanggalMulai:           db.Date(row.TanggalKursus.Time),
			TanggalSelesai:         db.Date(tanggalSelesai),
			Tahun:                  tahun,
			Durasi:                 row.Durasi,
			InstitusiPenyelenggara: row.InstitusiPenyelenggara.String,
			NomorSertifikat:        row.NoSertifikat.String,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatPelatihanTeknis(ctx, sqlc.GetBerkasRiwayatPelatihanTeknisParams{
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
