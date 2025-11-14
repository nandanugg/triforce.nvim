package riwayathukumandisiplin

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
	ListRiwayatHukumanDisiplin(ctx context.Context, arg sqlc.ListRiwayatHukumanDisiplinParams) ([]sqlc.ListRiwayatHukumanDisiplinRow, error)
	CountRiwayatHukumanDisiplin(ctx context.Context, nip pgtype.Text) (int64, error)
	GetBerkasRiwayatHukumanDisiplin(ctx context.Context, arg sqlc.GetBerkasRiwayatHukumanDisiplinParams) (pgtype.Text, error)
	CreateRiwayatHukumanDisiplin(ctx context.Context, arg sqlc.CreateRiwayatHukumanDisiplinParams) (int64, error)
	GetRefGolongan(ctx context.Context, id int32) (sqlc.GetRefGolonganRow, error)
	GetRefJenisHukuman(ctx context.Context, id int32) (sqlc.GetRefJenisHukumanRow, error)
	GetPegawaiByNIP(ctx context.Context, nip string) (sqlc.GetPegawaiByNIPRow, error)
	UpdateRiwayatHukumanDisiplin(ctx context.Context, arg sqlc.UpdateRiwayatHukumanDisiplinParams) (int64, error)
	DeleteRiwayatHukumanDisiplin(ctx context.Context, arg sqlc.DeleteRiwayatHukumanDisiplinParams) (int64, error)
	UploadBerkasRiwayatHukumanDisiplin(ctx context.Context, arg sqlc.UploadBerkasRiwayatHukumanDisiplinParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatHukumanDisiplin, uint, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	rows, err := s.repo.ListRiwayatHukumanDisiplin(ctx, sqlc.ListRiwayatHukumanDisiplinParams{
		PnsNip: pgNip,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatHukumanDisiplin(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	return typeutil.Map(rows, func(row sqlc.ListRiwayatHukumanDisiplinRow) riwayatHukumanDisiplin {
		return riwayatHukumanDisiplin{
			ID:                  row.ID,
			JenisHukuman:        row.JenisHukuman.String,
			JenisHukumanID:      row.JenisHukumanID.Int16,
			NamaGolongan:        row.NamaGolongan.String,
			GolonganID:          row.GolonganID.Int16,
			NamaPangkat:         row.NamaPangkat.String,
			NomorSK:             row.SkNomor.String,
			TanggalSK:           db.Date(row.SkTanggal.Time),
			TanggalMulai:        db.Date(row.TanggalMulaiHukuman.Time),
			TanggalAkhir:        db.Date(row.TanggalAkhirHukuman.Time),
			MasaTahun:           int(row.MasaTahun.Int16),
			MasaBulan:           int(row.MasaBulan.Int16),
			NomorPP:             row.NoPp.String,
			NomorSKPembatalan:   row.NoSkPembatalan.String,
			TanggalSKPembatalan: db.Date(row.TanggalSkPembatalan.Time),
		}
	}), uint(count), nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	res, err := s.repo.GetBerkasRiwayatHukumanDisiplin(ctx, sqlc.GetBerkasRiwayatHukumanDisiplinParams{
		PnsNip: pgNip,
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

type references struct {
	golongan     sqlc.GetRefGolonganRow
	jenisHukuman sqlc.GetRefJenisHukumanRow
	pegawai      sqlc.GetPegawaiByNIPRow
	tahun        int16
	bulan        int16
}

func (s *service) create(ctx context.Context, req adminCreateRequest) (int64, error) {
	ref, err := s.validateReferences(ctx, req.GolonganID, req.JenisHukumanID, req.NIP, req.TanggalMulai, req.TanggalAkhir)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return 0, errPegawaiNotFound
		}
		return 0, err
	}

	id, err := s.repo.CreateRiwayatHukumanDisiplin(ctx, sqlc.CreateRiwayatHukumanDisiplinParams{
		PnsID:               pgtype.Text{String: ref.pegawai.PnsID, Valid: true},
		PnsNip:              pgtype.Text{String: req.NIP, Valid: true},
		Nama:                ref.pegawai.Nama,
		NamaGolongan:        ref.golongan.Nama,
		NamaJenisHukuman:    ref.jenisHukuman.Nama,
		JenisHukumanID:      pgtype.Int2{Int16: int16(req.JenisHukumanID), Valid: true},
		GolonganID:          pgtype.Int2{Int16: int16(req.GolonganID), Valid: true},
		SkNomor:             pgtype.Text{String: req.NomorSK, Valid: true},
		SkTanggal:           req.TanggalSK.ToPgtypeDate(),
		TanggalMulaiHukuman: req.TanggalMulai.ToPgtypeDate(),
		MasaTahun:           pgtype.Int2{Int16: ref.tahun, Valid: true},
		MasaBulan:           pgtype.Int2{Int16: ref.bulan, Valid: true},
		TanggalAkhirHukuman: req.TanggalAkhir.ToPgtypeDate(),
		NoPp:                pgtype.Text{String: req.NomorPP, Valid: true},
		NoSkPembatalan:      pgtype.Text{String: req.NomorSKPembatalan, Valid: true},
		TanggalSkPembatalan: req.TanggalSKPembatalan.ToPgtypeDate(),
	})
	if err != nil {
		return 0, fmt.Errorf("[riwayat hukuman disiplin-create] repo create: %w", err)
	}
	return id, nil
}

func (s *service) update(ctx context.Context, req adminUpdateRequest) (bool, error) {
	ref, err := s.validateReferences(ctx, req.GolonganID, req.JenisHukumanID, req.NIP, req.TanggalMulai, req.TanggalAkhir)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return false, errPegawaiNotFound
		}
		return false, err
	}

	affected, err := s.repo.UpdateRiwayatHukumanDisiplin(ctx, sqlc.UpdateRiwayatHukumanDisiplinParams{
		ID:                  req.ID,
		Nip:                 req.NIP,
		GolonganID:          pgtype.Int2{Int16: int16(req.GolonganID), Valid: true},
		JenisHukumanID:      pgtype.Int2{Int16: int16(req.JenisHukumanID), Valid: true},
		NamaGolongan:        ref.golongan.Nama,
		NamaJenisHukuman:    ref.jenisHukuman.Nama,
		SkNomor:             pgtype.Text{String: req.NomorSK, Valid: true},
		SkTanggal:           req.TanggalSK.ToPgtypeDate(),
		TanggalMulaiHukuman: req.TanggalMulai.ToPgtypeDate(),
		MasaTahun:           pgtype.Int2{Int16: ref.tahun, Valid: true},
		MasaBulan:           pgtype.Int2{Int16: ref.bulan, Valid: true},
		TanggalAkhirHukuman: req.TanggalAkhir.ToPgtypeDate(),
		NoPp:                pgtype.Text{String: req.NomorPP, Valid: true},
		NoSkPembatalan:      pgtype.Text{String: req.NomorSKPembatalan, Valid: true},
		TanggalSkPembatalan: req.TanggalSKPembatalan.ToPgtypeDate(),
	})
	if err != nil {
		return false, fmt.Errorf("[riwayat hukuman disiplin-update] repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id int32, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatHukumanDisiplin(ctx, sqlc.DeleteRiwayatHukumanDisiplinParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("[riwayat hukuman disiplin-delete] repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int64, nip string, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatHukumanDisiplin(ctx, sqlc.UploadBerkasRiwayatHukumanDisiplinParams{
		ID:         id,
		Nip:        nip,
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("[riwayat hukuman disiplin-upload berkas] repo upload berkas: %w", err)
	}

	return affected > 0, nil
}

func (s *service) validateReferences(ctx context.Context, golonganID int32, jenisHukumanID int32, nip string, tanggalMulai, tanggalAkhir db.Date) (*references, error) {
	var errs []error
	ref := references{}

	golongan, err := s.repo.GetRefGolongan(ctx, golonganID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("[riwayat hukuman disiplin-validate references] repo get golongan: %w", err)
		}
		errs = append(errs, errGolonganNotFound)
	}
	ref.golongan = golongan

	jenisHukuman, err := s.repo.GetRefJenisHukuman(ctx, jenisHukumanID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("[riwayat hukuman disiplin-validate references] repo get jenis hukuman: %w", err)
		}
		errs = append(errs, errJenisHukumanNotFound)
	}
	ref.jenisHukuman = jenisHukuman

	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("[riwayat hukuman disiplin-validate references] repo get pegawai: %w", err)
		}
		return nil, errPegawaiNotFound
	}
	ref.pegawai = pegawai

	tahun, bulan, err := s.calculateMasaTahun(tanggalMulai, tanggalAkhir)
	if err != nil {
		errs = append(errs, errMasaHukumanTidakValid)
	}
	ref.tahun = tahun
	ref.bulan = bulan

	if len(errs) > 0 {
		return nil, api.NewMultiError(errs)
	}

	return &ref, nil
}

func (s *service) calculateMasaTahun(tanggalMulai, tanggalAkhir db.Date) (tahun, bulan int16, err error) {
	start, end := time.Time(tanggalMulai), time.Time(tanggalAkhir)

	if start.IsZero() || end.IsZero() {
		return 0, 0, errors.New("tanggal tidak valid")
	}
	if end.Before(start) {
		return 0, 0, errors.New("tanggal akhir tidak boleh sebelum tanggal mulai")
	}

	totalMonths := int(end.Year()-start.Year())*12 + int(end.Month()-start.Month())

	if end.Day() < start.Day() {
		totalMonths--
	}

	tahun = int16(totalMonths / 12)
	bulan = int16(totalMonths % 12)

	return tahun, bulan, nil
}
