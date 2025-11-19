package riwayatpenugasan

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
	ListRiwayatPenugasan(ctx context.Context, arg sqlc.ListRiwayatPenugasanParams) ([]sqlc.ListRiwayatPenugasanRow, error)
	CountRiwayatPenugasan(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatPenugasan(ctx context.Context, arg sqlc.GetBerkasRiwayatPenugasanParams) (pgtype.Text, error)
	CreateRiwayatPenugasan(ctx context.Context, arg sqlc.CreateRiwayatPenugasanParams) (int32, error)
	UpdateRiwayatPenugasan(ctx context.Context, arg sqlc.UpdateRiwayatPenugasanParams) (int64, error)
	DeleteRiwayatPenugasan(ctx context.Context, arg sqlc.DeleteRiwayatPenugasanParams) (int64, error)
	UploadBerkasRiwayatPenugasan(ctx context.Context, arg sqlc.UploadBerkasRiwayatPenugasanParams) (int64, error)
	GetPegawaiByNIP(ctx context.Context, nip string) (sqlc.GetPegawaiByNIPRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatPenugasan, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatPenugasan(ctx, sqlc.ListRiwayatPenugasanParams{
		Nip:    pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatPenugasan(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(data, func(row sqlc.ListRiwayatPenugasanRow) riwayatPenugasan {
		isMenjabat := row.IsMenjabat.Bool || row.TanggalSelesai.Time.IsZero() || row.TanggalSelesai.Time.After(time.Now())
		return riwayatPenugasan{
			ID:               row.ID,
			TipeJabatan:      row.TipeJabatan.String,
			NameJabatan:      row.NamaJabatan.String,
			DeskripsiJabatan: row.DeskripsiJabatan.String,
			TanggalMulai:     db.Date(row.TanggalMulai.Time),
			TanggalSelesai:   db.Date(row.TanggalSelesai.Time),
			IsMenjabat:       isMenjabat,
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int32) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatPenugasan(ctx, sqlc.GetBerkasRiwayatPenugasanParams{
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

func (s *service) create(ctx context.Context, req adminCreateRequest) (int32, error) {
	_, err := s.repo.GetPegawaiByNIP(ctx, req.NIP)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return 0, errPegawaiNotFound
		}
		return 0, fmt.Errorf("[riwayatpenugasan-create] repo get pegawai: %w", err)
	}

	id, err := s.repo.CreateRiwayatPenugasan(ctx, sqlc.CreateRiwayatPenugasanParams{
		Nip:              pgtype.Text{String: req.NIP, Valid: true},
		TipeJabatan:      pgtype.Text{String: req.TipeJabatan, Valid: true},
		NamaJabatan:      pgtype.Text{String: req.NamaJabatan, Valid: true},
		DeskripsiJabatan: pgtype.Text{String: req.DeskripsiJabatan, Valid: req.DeskripsiJabatan != ""},
		TanggalMulai:     req.TanggalMulai.ToPgtypeDate(),
		TanggalSelesai:   req.TanggalSelesai.ToPgtypeDate(),
		IsMenjabat:       pgtype.Bool{Bool: req.IsMenjabat, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("[riwayatpenugasan-create] repo create: %w", err)
	}
	return id, nil
}

func (s *service) update(ctx context.Context, req adminUpdateRequest) (bool, error) {
	_, err := s.repo.GetPegawaiByNIP(ctx, req.NIP)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, errPegawaiNotFound
		}
		return false, fmt.Errorf("[riwayatpenugasan-update] repo get pegawai: %w", err)
	}

	affected, err := s.repo.UpdateRiwayatPenugasan(ctx, sqlc.UpdateRiwayatPenugasanParams{
		ID:               req.ID,
		Nip:              req.NIP,
		TipeJabatan:      pgtype.Text{String: req.TipeJabatan, Valid: true},
		NamaJabatan:      pgtype.Text{String: req.NamaJabatan, Valid: true},
		DeskripsiJabatan: pgtype.Text{String: req.DeskripsiJabatan, Valid: req.DeskripsiJabatan != ""},
		TanggalMulai:     req.TanggalMulai.ToPgtypeDate(),
		TanggalSelesai:   req.TanggalSelesai.ToPgtypeDate(),
		IsMenjabat:       pgtype.Bool{Bool: req.IsMenjabat, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("[riwayatpenugasan-update] repo update: %w", err)
	}
	return affected == 1, nil
}

func (s *service) delete(ctx context.Context, req adminDeleteRequest) (bool, error) {
	affected, err := s.repo.DeleteRiwayatPenugasan(ctx, sqlc.DeleteRiwayatPenugasanParams{
		ID:  req.ID,
		Nip: req.NIP,
	})
	if err != nil {
		return false, fmt.Errorf("[riwayatpenugasan-delete] repo delete: %w", err)
	}
	return affected == 1, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int32, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatPenugasan(ctx, sqlc.UploadBerkasRiwayatPenugasanParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("[riwayatpenugasan-upload-berkas] repo upload berkas: %w", err)
	}
	return affected == 1, nil
}
