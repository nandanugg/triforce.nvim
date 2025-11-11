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
	GetPegawaiPNSIDByNIP(ctx context.Context, nip string) (string, error)
	GetRefJenisDiklat(ctx context.Context, id int32) (sqlc.GetRefJenisDiklatRow, error)

	CreateRiwayatPelatihanSIASN(ctx context.Context, arg sqlc.CreateRiwayatPelatihanSIASNParams) (int64, error)
	UpdateRiwayatPelatihanSIASN(ctx context.Context, arg sqlc.UpdateRiwayatPelatihanSIASNParams) (int64, error)
	DeleteRiwayatPelatihanSIASN(ctx context.Context, arg sqlc.DeleteRiwayatPelatihanSIASNParams) (int64, error)
	UploadBerkasRiwayatPelatihanSIASN(ctx context.Context, arg sqlc.UploadBerkasRiwayatPelatihanSIASNParams) (int64, error)
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
			JenisDiklatID:          row.JenisDiklatID,
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

func (s *service) create(ctx context.Context, nip string, params upsertParams) (int64, error) {
	pnsID, err := s.repo.GetPegawaiPNSIDByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	ref, err := s.validateReferences(ctx, params)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateRiwayatPelatihanSIASN(ctx, sqlc.CreateRiwayatPelatihanSIASNParams{
		PnsOrangID:             pgtype.Text{String: pnsID, Valid: true},
		NipBaru:                pgtype.Text{String: nip, Valid: true},
		NamaDiklat:             pgtype.Text{String: params.NamaDiklat, Valid: true},
		JenisDiklatID:          pgtype.Int2{Int16: params.JenisDiklatID, Valid: true},
		JenisDiklat:            ref.jenisDiklat.nama,
		InstitusiPenyelenggara: pgtype.Text{String: params.InstitusiPenyelenggara, Valid: true},
		NoSertifikat:           pgtype.Text{String: params.NomorSertifikat, Valid: true},
		TanggalMulai:           params.TanggalMulai.ToPgtypeDate(),
		TanggalSelesai:         params.TanggalSelesai.ToPgtypeDate(),
		TahunDiklat:            pgtype.Int4{Int32: typeutil.FromPtr(params.Tahun), Valid: params.Tahun != nil},
		DurasiJam:              pgtype.Int4{Int32: params.Durasi, Valid: true},
	})
	if err != nil {
		return 0, fmt.Errorf("repo create: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, id int64, nip string, params upsertParams) (bool, error) {
	ref, err := s.validateReferences(ctx, params)
	if err != nil {
		return false, err
	}

	affected, err := s.repo.UpdateRiwayatPelatihanSIASN(ctx, sqlc.UpdateRiwayatPelatihanSIASNParams{
		ID:                     id,
		Nip:                    nip,
		NamaDiklat:             pgtype.Text{String: params.NamaDiklat, Valid: true},
		JenisDiklatID:          pgtype.Int2{Int16: params.JenisDiklatID, Valid: true},
		JenisDiklat:            ref.jenisDiklat.nama,
		InstitusiPenyelenggara: pgtype.Text{String: params.InstitusiPenyelenggara, Valid: true},
		NoSertifikat:           pgtype.Text{String: params.NomorSertifikat, Valid: true},
		TanggalMulai:           params.TanggalMulai.ToPgtypeDate(),
		TanggalSelesai:         params.TanggalSelesai.ToPgtypeDate(),
		TahunDiklat:            pgtype.Int4{Int32: typeutil.FromPtr(params.Tahun), Valid: params.Tahun != nil},
		DurasiJam:              pgtype.Int4{Int32: params.Durasi, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id int64, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatPelatihanSIASN(ctx, sqlc.DeleteRiwayatPelatihanSIASNParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int64, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatPelatihanSIASN(ctx, sqlc.UploadBerkasRiwayatPelatihanSIASNParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo upload berkas: %w", err)
	}

	return affected > 0, nil
}

type references struct {
	jenisDiklat struct {
		nama pgtype.Text
	}
}

func (s *service) validateReferences(ctx context.Context, params upsertParams) (*references, error) {
	var errs []error
	ref := references{}

	jenisDiklat, err := s.repo.GetRefJenisDiklat(ctx, int32(params.JenisDiklatID))
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get jenis diklat: %w", err)
		}
		errs = append(errs, errJenisDiklatNotFound)
	}
	ref.jenisDiklat.nama = jenisDiklat.JenisDiklat

	if len(errs) > 0 {
		return nil, api.NewMultiError(errs)
	}
	return &ref, nil
}
