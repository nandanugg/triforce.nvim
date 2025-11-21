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
	DeleteRiwayatPenghargaan(ctx context.Context, arg repo.DeleteRiwayatPenghargaanParams) (int64, error)
	CreateRiwayatPenghargaan(ctx context.Context, arg repo.CreateRiwayatPenghargaanParams) (int32, error)
	UpdateRiwayatPenghargaan(ctx context.Context, arg repo.UpdateRiwayatPenghargaanParams) (int64, error)
	UpdateRiwayatPenghargaanBerkas(ctx context.Context, arg repo.UpdateRiwayatPenghargaanBerkasParams) (int64, error)
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

func (s *service) uploadBerkas(ctx context.Context, id int32, nip string, base64 string) (bool, error) {
	res, err := s.repo.UpdateRiwayatPenghargaanBerkas(ctx, repo.UpdateRiwayatPenghargaanBerkasParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{Valid: true, String: base64},
	})
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return false, fmt.Errorf("repo upload berkas: %w", err)
	}

	if res == 0 {
		return false, nil
	}

	return true, nil
}

func (s *service) create(ctx context.Context, nip string, params upsertParams) (int32, error) {
	_, valid := s.validateJenisPenghargaan(params.JenisPenghargaan)
	if !valid {
		return 0, NewError(ErrJenisPenghargaanInvalid, params.JenisPenghargaan)
	}

	id, err := s.repo.CreateRiwayatPenghargaan(ctx, repo.CreateRiwayatPenghargaanParams{
		Nip:                  pgtype.Text{Valid: true, String: nip},
		NamaPenghargaan:      pgtype.Text{Valid: true, String: params.NamaPenghargaan},
		JenisPenghargaan:     pgtype.Text{Valid: true, String: params.JenisPenghargaan},
		DeskripsiPenghargaan: pgtype.Text{Valid: params.Deskripsi != "", String: params.Deskripsi},
		TanggalPenghargaan:   pgtype.Date{Valid: params.Tanggal.ToPgtypeDate().Valid, Time: params.Tanggal.ToPgtypeDate().Time},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id int32, nip string, params upsertParams) (bool, error) {
	_, valid := s.validateJenisPenghargaan(params.JenisPenghargaan)
	if !valid {
		return false, NewError(ErrJenisPenghargaanInvalid, params.JenisPenghargaan)
	}

	res, err := s.repo.UpdateRiwayatPenghargaan(ctx, repo.UpdateRiwayatPenghargaanParams{
		ID:                   id,
		Nip:                  pgtype.Text{Valid: true, String: nip},
		NamaPenghargaan:      pgtype.Text{Valid: true, String: params.NamaPenghargaan},
		JenisPenghargaan:     pgtype.Text{Valid: true, String: params.JenisPenghargaan},
		DeskripsiPenghargaan: pgtype.Text{Valid: params.Deskripsi != "", String: params.Deskripsi},
		TanggalPenghargaan:   pgtype.Date{Valid: params.Tanggal.ToPgtypeDate().Valid, Time: params.Tanggal.ToPgtypeDate().Time},
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	if res == 0 {
		return false, nil
	}

	return true, nil
}

func (s *service) delete(ctx context.Context, id int32, nip string) (bool, error) {
	res, err := s.repo.DeleteRiwayatPenghargaan(ctx, repo.DeleteRiwayatPenghargaanParams{
		ID: id, Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	if res == 0 {
		return false, nil
	}
	return true, nil
}
