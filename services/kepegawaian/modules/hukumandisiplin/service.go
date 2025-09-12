package hukumandisiplin

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/db"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	ListRiwayatHukdis(ctx context.Context, arg sqlc.ListRiwayatHukdisParams) ([]sqlc.ListRiwayatHukdisRow, error)
	CountRiwayatHukdis(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatHukdis(ctx context.Context, arg sqlc.GetBerkasRiwayatHukdisParams) (pgtype.Text, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]hukumanDisiplin, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatHukdis(ctx, sqlc.ListRiwayatHukdisParams{
		PnsNip: pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatHukdis(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	data := make([]hukumanDisiplin, 0, len(rows))
	for _, row := range rows {
		data = append(data, hukumanDisiplin{
			ID:           row.ID,
			JenisHukuman: row.JenisHukuman.String,
			NomorSK:      row.SkNomor.String,
			TanggalSK:    db.Date(row.SkTanggal.Time),
			TanggalMulai: db.Date(row.TanggalMulaiHukuman.Time),
			TanggalAkhir: db.Date(row.TanggalAkhirHukuman.Time),
			MasaTahun:    int(row.MasaTahun.Int16),
			MasaBulan:    int(row.MasaBulan.Int16),
		})
	}
	return data, uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatHukdis(ctx, sqlc.GetBerkasRiwayatHukdisParams{
		PnsNip: pgNip,
		ID:     id,
	})
	if errors.Is(err, pgx.ErrNoRows) || len(res.String) == 0 {
		return "", nil, nil
	}
	if err != nil {
		return "", nil, fmt.Errorf("repo get berkas: %w", err)
	}

	return api.GetMimetypeAndDecodedData(res.String)
}
