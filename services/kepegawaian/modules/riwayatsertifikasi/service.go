package riwayatsertifikasi

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
	ListRiwayatSertifikasi(ctx context.Context, arg sqlc.ListRiwayatSertifikasiParams) ([]sqlc.ListRiwayatSertifikasiRow, error)
	CountRiwayatSertifikasi(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatSertifikasi(ctx context.Context, arg sqlc.GetBerkasRiwayatSertifikasiParams) (pgtype.Text, error)
	CreateRiwayatSertifikasi(ctx context.Context, arg sqlc.CreateRiwayatSertifikasiParams) (int64, error)
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	UpdateRiwayatSertifikasiByIDAndNIP(ctx context.Context, arg sqlc.UpdateRiwayatSertifikasiByIDAndNIPParams) (int64, error)
	DeleteRiwayatSertifikasiByIDAndNIP(ctx context.Context, arg sqlc.DeleteRiwayatSertifikasiByIDAndNIPParams) (int64, error)
	UpdateBerkasRiwayatSertifikasiByIDAndNIP(ctx context.Context, arg sqlc.UpdateBerkasRiwayatSertifikasiByIDAndNIPParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatSertifikasi, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatSertifikasi(ctx, sqlc.ListRiwayatSertifikasiParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatSertifikasi(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatSertifikasiRow) riwayatSertifikasi {
		return riwayatSertifikasi{
			ID:              row.ID,
			NamaSertifikasi: row.NamaSertifikasi.String,
			Tahun:           row.Tahun,
			Deskripsi:       row.Deskripsi.String,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatSertifikasi(ctx, sqlc.GetBerkasRiwayatSertifikasiParams{
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

func (s *service) create(ctx context.Context, params adminCreateRequest) (int64, error) {
	_, err := s.repo.GetPegawaiPNSIDByNIP(ctx, params.NIP)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errPegawaiNotFound
		}
		return 0, fmt.Errorf("[riwayatsertifikasi-create] repo get pegawai: %w", err)
	}

	id, err := s.repo.CreateRiwayatSertifikasi(ctx, sqlc.CreateRiwayatSertifikasiParams{
		Nip:             pgtype.Text{String: params.NIP, Valid: true},
		Tahun:           pgtype.Int8{Int64: int64(params.Tahun), Valid: true},
		NamaSertifikasi: pgtype.Text{String: params.NamaSertifikasi, Valid: true},
		Deskripsi:       pgtype.Text{String: params.Deskripsi, Valid: params.Deskripsi != ""},
	})
	if err != nil {
		return 0, fmt.Errorf("[riwayatsertifikasi-create] repo create: %w", err)
	}
	return id, nil
}

func (s *service) update(ctx context.Context, params adminUpdateRequest) (bool, error) {
	affected, err := s.repo.UpdateRiwayatSertifikasiByIDAndNIP(ctx, sqlc.UpdateRiwayatSertifikasiByIDAndNIPParams{
		ID:              params.ID,
		Nip:             params.NIP,
		Tahun:           pgtype.Int8{Int64: int64(params.Tahun), Valid: true},
		NamaSertifikasi: pgtype.Text{String: params.NamaSertifikasi, Valid: true},
		Deskripsi:       pgtype.Text{String: params.Deskripsi, Valid: params.Deskripsi != ""},
	})
	if err != nil {
		return false, fmt.Errorf("[riwayatsertifikasi-update] repo update: %w", err)
	}
	return affected == 1, nil
}

func (s *service) delete(ctx context.Context, id int64, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatSertifikasiByIDAndNIP(ctx, sqlc.DeleteRiwayatSertifikasiByIDAndNIPParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("[riwayatsertifikasi-delete] repo delete: %w", err)
	}
	return affected == 1, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int64, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UpdateBerkasRiwayatSertifikasiByIDAndNIP(ctx, sqlc.UpdateBerkasRiwayatSertifikasiByIDAndNIPParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("[riwayatsertifikasi-upload-berkas] repo update berkas: %w", err)
	}
	return affected == 1, nil
}
