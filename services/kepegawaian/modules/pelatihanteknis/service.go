package pelatihanteknis

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repositoryInterface interface {
	ListPelatihanTeknis(ctx context.Context, arg sqlc.ListPelatihanTeknisParams) ([]sqlc.ListPelatihanTeknisRow, error)
	CountPelatihanTeknis(ctx context.Context, nip pgtype.Text) (int64, error)
}

type service struct {
	repo repositoryInterface
}

func newService(r repositoryInterface) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]pelatihanTeknis, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListPelatihanTeknis(ctx, sqlc.ListPelatihanTeknisParams{
		NipBaru: pgNip,
		Limit:   int32(limit),
		Offset:  int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountPelatihanTeknis(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	data := make([]pelatihanTeknis, 0, len(rows))
	for _, row := range rows {
		tanggalMulai := row.TanggalKursus.Time
		var tanggalSelesai time.Time

		durasi := int(row.Durasi)

		if row.TanggalKursus.Valid {
			tanggalSelesai = tanggalMulai.Add(time.Duration(durasi) * time.Hour)
		}

		tahun := func() *int {
			if row.Tahun.Valid {
				if v, err := row.Tahun.Int64Value(); err == nil {
					year := int(v.Int64)
					return &year
				}
			}
			return nil
		}()

		data = append(data, pelatihanTeknis{
			ID:                     int64(row.ID),
			TipePelatihan:          row.TipeKursus.String,
			JenisPelatihan:         row.JenisKursus.String,
			NamaPelatihan:          row.NamaKursus.String,
			TanggalMulai:           db.Date(tanggalMulai),
			TanggalSelesai:         db.Date(tanggalSelesai),
			Tahun:                  tahun,
			Durasi:                 durasi,
			InstitusiPenyelenggara: row.InstitusiPenyelenggara.String,
			NomorSertifikat:        row.NoSertifikat.String,
		})
	}
	return data, uint(count), nil
}
