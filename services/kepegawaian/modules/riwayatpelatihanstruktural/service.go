package riwayatpelatihanstruktural

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
	pnsNIP := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatPelatihanStruktural(ctx, sqlc.ListRiwayatPelatihanStrukturalParams{
		PnsNip: pnsNIP,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPelatihanStruktural(ctx, pnsNIP)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row sqlc.ListRiwayatPelatihanStrukturalRow) riwayatPelatihanStruktural {
		var tahun *int16
		if row.Tahun.Valid {
			tahun = &row.Tahun.Int16
		}

		var tanggalSelesai time.Time
		if row.Tanggal.Valid {
			tanggalSelesai = row.Tanggal.Time.Add(time.Duration(row.Lama.Float32) * time.Hour)
			if tahun == nil {
				tahun = typeutil.ToPtr(int16(row.Tanggal.Time.Year()))
			}
		}

		return riwayatPelatihanStruktural{
			ID:                     row.ID,
			JenisDiklat:            row.JenisDiklat.String,
			NamaDiklat:             row.NamaDiklat.String,
			InstitusiPenyelenggara: row.InstitusiPenyelenggara.String,
			NomorSertifikat:        row.Nomor.String,
			TanggalMulai:           db.Date(row.Tanggal.Time),
			TanggalSelesai:         db.Date(tanggalSelesai),
			Tahun:                  tahun,
			Durasi:                 row.Lama,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip, id string) (string, []byte, error) {
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
