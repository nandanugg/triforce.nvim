package riwayatpelatihanstruktural

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
	ListRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.ListRiwayatPelatihanStrukturalParams) ([]sqlc.ListRiwayatPelatihanStrukturalRow, error)
	CountRiwayatPelatihanStruktural(ctx context.Context, pnsNip pgtype.Text) (int64, error)
	GetBerkasRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.GetBerkasRiwayatPelatihanStrukturalParams) (pgtype.Text, error)
	GetPegawaiByNIP(ctx context.Context, nip string) (sqlc.GetPegawaiByNIPRow, error)

	CreateRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.CreateRiwayatPelatihanStrukturalParams) (string, error)
	UpdateRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.UpdateRiwayatPelatihanStrukturalParams) (int64, error)
	DeleteRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.DeleteRiwayatPelatihanStrukturalParams) (int64, error)
	UploadBerkasRiwayatPelatihanStruktural(ctx context.Context, arg sqlc.UploadBerkasRiwayatPelatihanStrukturalParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPelatihanStruktural, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatPelatihanStruktural(ctx, sqlc.ListRiwayatPelatihanStrukturalParams{
		PnsNip: pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPelatihanStruktural(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatPelatihanStrukturalRow) riwayatPelatihanStruktural {
		return riwayatPelatihanStruktural{
			ID:         row.ID,
			NamaDiklat: row.NamaDiklat.String,
			Tanggal:    db.Date(row.Tanggal.Time),
			Nomor:      row.Nomor.String,
			Lama:       row.Lama,
			Tahun:      row.Tahun,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id string) (string, []byte, error) {
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

func (s *service) create(ctx context.Context, nip string, params upsertParams) (string, error) {
	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		return "", errPegawaiNotFound
	}

	id, err := s.repo.CreateRiwayatPelatihanStruktural(ctx, sqlc.CreateRiwayatPelatihanStrukturalParams{
		PnsID:      pgtype.Text{String: pegawai.PnsID, Valid: true},
		PnsNip:     pgtype.Text{String: nip, Valid: true},
		PnsNama:    pegawai.Nama,
		NamaDiklat: pgtype.Text{String: params.NamaDiklat, Valid: true},
		Tanggal:    params.Tanggal.ToPgtypeDate(),
		Tahun:      pgtype.Int2{Int16: params.Tahun, Valid: true},
		Lama:       pgtype.Float4{Float32: params.Lama, Valid: true},
		Nomor:      pgtype.Text{String: params.Nomor, Valid: params.Nomor != ""},
	})
	if err != nil {
		return "", fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id, nip string, params upsertParams) (bool, error) {
	affected, err := s.repo.UpdateRiwayatPelatihanStruktural(ctx, sqlc.UpdateRiwayatPelatihanStrukturalParams{
		ID:         id,
		Nip:        nip,
		NamaDiklat: pgtype.Text{String: params.NamaDiklat, Valid: true},
		Tanggal:    params.Tanggal.ToPgtypeDate(),
		Tahun:      pgtype.Int2{Int16: params.Tahun, Valid: true},
		Lama:       pgtype.Float4{Float32: params.Lama, Valid: true},
		Nomor:      pgtype.Text{String: params.Nomor, Valid: params.Nomor != ""},
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatPelatihanStruktural(ctx, sqlc.DeleteRiwayatPelatihanStrukturalParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatPelatihanStruktural(ctx, sqlc.UploadBerkasRiwayatPelatihanStrukturalParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo upload berkas: %w", err)
	}

	return affected > 0, nil
}
