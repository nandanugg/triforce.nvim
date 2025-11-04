package riwayatpendidikan

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/api"
	"gitlab.com/wartek-id/matk/nexus/nexus-be/lib/typeutil"
	sqlc "gitlab.com/wartek-id/matk/nexus/nexus-be/services/kepegawaian/db/repository"
)

type repository interface {
	CountRiwayatPendidikan(ctx context.Context, nip pgtype.Text) (int64, error)
	ListRiwayatPendidikan(ctx context.Context, arg sqlc.ListRiwayatPendidikanParams) ([]sqlc.ListRiwayatPendidikanRow, error)
	GetBerkasRiwayatPendidikan(ctx context.Context, arg sqlc.GetBerkasRiwayatPendidikanParams) (pgtype.Text, error)
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	GetRefPendidikan(ctx context.Context, id string) (sqlc.GetRefPendidikanRow, error)
	GetRefTingkatPendidikan(ctx context.Context, id int32) (sqlc.GetRefTingkatPendidikanRow, error)

	CreateRiwayatPendidikan(ctx context.Context, arg sqlc.CreateRiwayatPendidikanParams) (int32, error)
	UpdateRiwayatPendidikan(ctx context.Context, arg sqlc.UpdateRiwayatPendidikanParams) (int64, error)
	UploadBerkasRiwayatPendidikan(ctx context.Context, arg sqlc.UploadBerkasRiwayatPendidikanParams) (int64, error)
	DeleteRiwayatPendidikan(ctx context.Context, arg sqlc.DeleteRiwayatPendidikanParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPendidikan, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatPendidikan(ctx, sqlc.ListRiwayatPendidikanParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPendidikan(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatPendidikanRow) riwayatPendidikan {
		return riwayatPendidikan{
			ID:                   row.ID,
			TingkatPendidikanID:  row.TingkatPendidikanID,
			JenjangPendidikan:    row.JenjangPendidikan.String,
			PendidikanID:         row.PendidikanID,
			Pendidikan:           row.Pendidikan.String,
			NamaSekolah:          row.NamaSekolah.String,
			TahunLulus:           row.TahunLulus,
			NomorIjazah:          row.NoIjazah.String,
			GelarDepan:           row.GelarDepan.String,
			GelarBelakang:        row.GelarBelakang.String,
			TugasBelajar:         labelStatusBelajar[row.TugasBelajar.Int16],
			KeteranganPendidikan: row.NegaraSekolah.String,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPendidikan(ctx, sqlc.GetBerkasRiwayatPendidikanParams{
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

func (s *service) create(ctx context.Context, nip string, params upsertParams) (int32, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	if err := s.validateReferences(ctx, params); err != nil {
		return 0, err
	}

	id, err := s.repo.CreateRiwayatPendidikan(ctx, sqlc.CreateRiwayatPendidikanParams{
		TingkatPendidikanID: pgtype.Int2{Int16: params.TingkatPendidikanID, Valid: true},
		PendidikanID:        pgtype.Text{String: typeutil.FromPtr(params.PendidikanID), Valid: params.PendidikanID != nil},
		NamaSekolah:         pgtype.Text{String: params.NamaSekolah, Valid: true},
		TahunLulus:          pgtype.Int2{Int16: params.TahunLulus, Valid: true},
		NoIjazah:            pgtype.Text{String: params.NomorIjazah, Valid: true},
		GelarDepan:          pgtype.Text{String: params.GelarDepan, Valid: params.GelarDepan != ""},
		GelarBelakang:       pgtype.Text{String: params.GelarBelakang, Valid: params.GelarBelakang != ""},
		NegaraSekolah:       pgtype.Text{String: params.NegaraSekolah, Valid: params.NegaraSekolah != ""},
		TugasBelajar:        params.TugasBelajar.toID(),
		PnsID:               pgtype.Text{String: pnsID, Valid: true},
		Nip:                 pgtype.Text{String: nip, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id int32, nip string, params upsertParams) (bool, error) {
	if err := s.validateReferences(ctx, params); err != nil {
		return false, err
	}

	affected, err := s.repo.UpdateRiwayatPendidikan(ctx, sqlc.UpdateRiwayatPendidikanParams{
		ID:                  id,
		Nip:                 pgtype.Text{String: nip, Valid: true},
		TingkatPendidikanID: pgtype.Int2{Int16: params.TingkatPendidikanID, Valid: true},
		PendidikanID:        pgtype.Text{String: typeutil.FromPtr(params.PendidikanID), Valid: params.PendidikanID != nil},
		NamaSekolah:         pgtype.Text{String: params.NamaSekolah, Valid: true},
		TahunLulus:          pgtype.Int2{Int16: params.TahunLulus, Valid: true},
		NoIjazah:            pgtype.Text{String: params.NomorIjazah, Valid: true},
		GelarDepan:          pgtype.Text{String: params.GelarDepan, Valid: params.GelarDepan != ""},
		GelarBelakang:       pgtype.Text{String: params.GelarBelakang, Valid: params.GelarBelakang != ""},
		NegaraSekolah:       pgtype.Text{String: params.NegaraSekolah, Valid: params.NegaraSekolah != ""},
		TugasBelajar:        params.TugasBelajar.toID(),
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id int32, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatPendidikan(ctx, sqlc.DeleteRiwayatPendidikanParams{
		ID:  id,
		Nip: pgtype.Text{String: nip, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int32, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatPendidikan(ctx, sqlc.UploadBerkasRiwayatPendidikanParams{
		ID:         id,
		Nip:        pgtype.Text{String: nip, Valid: true},
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo upload berkas: %w", err)
	}

	return affected > 0, nil
}

func (s *service) validateReferences(ctx context.Context, params upsertParams) error {
	var errs []error
	if _, err := s.repo.GetRefTingkatPendidikan(ctx, int32(params.TingkatPendidikanID)); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return fmt.Errorf("repo get tingkat pendidikan: %w", err)
		}
		errs = append(errs, errTingkatPendidikanNotFound)
	}

	if params.PendidikanID != nil {
		if _, err := s.repo.GetRefPendidikan(ctx, *params.PendidikanID); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return fmt.Errorf("repo get pendidikan: %w", err)
			}
			errs = append(errs, errPendidikanNotFound)
		}
	}

	if len(errs) > 0 {
		return api.NewMultiError(errs)
	}
	return nil
}
