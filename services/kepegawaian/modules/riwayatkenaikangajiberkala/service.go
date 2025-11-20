package riwayatkenaikangajiberkala

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
	ListRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.ListRiwayatKenaikanGajiBerkalaParams) ([]repo.ListRiwayatKenaikanGajiBerkalaRow, error)
	CountRiwayatKenaikanGajiBerkala(ctx context.Context, nipBaru pgtype.Text) (int64, error)
	GetBerkasRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.GetBerkasRiwayatKenaikanGajiBerkalaParams) (pgtype.Text, error)
	CreateRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.CreateRiwayatKenaikanGajiBerkalaParams) (int64, error)
	UpdateRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.UpdateRiwayatKenaikanGajiBerkalaParams) (int64, error)
	DeleteRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.DeleteRiwayatKenaikanGajiBerkalaParams) (int64, error)
	UploadBerkasRiwayatKenaikanGajiBerkala(ctx context.Context, arg repo.UploadBerkasRiwayatKenaikanGajiBerkalaParams) (int64, error)
	GetPegawaiByNIP(ctx context.Context, nip string) (repo.GetPegawaiByNIPRow, error)
	GetRefGolongan(ctx context.Context, id int32) (repo.GetRefGolonganRow, error)
	GetUnitKerja(ctx context.Context, id string) (repo.GetUnitKerjaRow, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatKenaikanGajiBerkala, int64, error) {
	pgNip := pgtype.Text{String: nip, Valid: true}
	data, err := s.repo.ListRiwayatKenaikanGajiBerkala(ctx, repo.ListRiwayatKenaikanGajiBerkalaParams{
		Limit:   int32(limit),
		Offset:  int32(offset),
		NipBaru: pgNip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatKenaikanGajiBerkala(ctx, pgNip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatKenaikanGajiBerkalaRow) riwayatKenaikanGajiBerkala {
		return riwayatKenaikanGajiBerkala{
			ID:                     row.ID,
			IDGolongan:             row.GolonganID,
			NamaGolongan:           row.GolonganNama.String,
			NamaGolonganPangkat:    row.GolonganNamaPangkat.String,
			NomorSK:                row.NoSk.String,
			TanggalSK:              db.Date(row.TanggalSk.Time),
			TMTGolongan:            db.Date(row.TmtGolongan.Time),
			MasaKerjaGolonganTahun: row.MasaKerjaGolonganTahun,
			MasaKerjaGolonganBulan: row.MasaKerjaGolonganBulan,
			TMTKenaikanGajiBerkala: db.Date(row.TmtKenaikanGajiBerkala.Time),
			GajiPokok:              row.GajiPokok,
			Jabatan:                row.Jabatan.String,
			TMTJabatan:             db.Date(row.TmtJabatan.Time),
			Pendidikan:             row.Pendidikan.String,
			TanggalLulus:           db.Date(row.TanggalLulus.Time),
			KantorPembayaran:       row.KantorPembayaran.String,
			UnitKerjaIndukID:       row.UnitKerjaIndukID.String,
			UnitKerjaInduk:         row.UnitKerjaInduk.String,
			Pejabat:                row.Pejabat.String,
		}
	})

	return result, count, nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatKenaikanGajiBerkala(ctx, repo.GetBerkasRiwayatKenaikanGajiBerkalaParams{
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

type references struct {
	golongan  repo.GetRefGolonganRow
	unitKerja repo.GetUnitKerjaRow
	pegawai   repo.GetPegawaiByNIPRow
}

func (s *service) create(ctx context.Context, arg adminCreateRequest) (int64, error) {
	ref, err := s.validateReferences(ctx, arg.NIP, arg.UnitKerjaIndukID, arg.GolonganID)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return 0, errPegawaiNotFound
		}
		return 0, err
	}

	id, err := s.repo.CreateRiwayatKenaikanGajiBerkala(ctx, repo.CreateRiwayatKenaikanGajiBerkalaParams{
		PegawaiID:                      pgtype.Int4{Int32: ref.pegawai.ID, Valid: true},
		GolonganID:                     pgtype.Int4{Int32: ref.golongan.ID, Valid: true},
		NGolRuang:                      ref.golongan.Nama,
		TmtGolongan:                    arg.TMTGolongan.ToPgtypeDate(),
		MasaKerjaGolonganTahun:         pgtype.Int2{Int16: arg.MasaKerjaGolonganTahun, Valid: true},
		MasaKerjaGolonganBulan:         pgtype.Int2{Int16: arg.MasaKerjaGolonganBulan, Valid: true},
		NoSk:                           pgtype.Text{String: arg.NomorSK, Valid: true},
		TanggalSk:                      arg.TanggalSK.ToPgtypeDate(),
		GajiPokok:                      pgtype.Int4{Int32: arg.GajiPokok, Valid: true},
		Jabatan:                        pgtype.Text{String: arg.Jabatan, Valid: arg.Jabatan != ""},
		TmtJabatan:                     arg.TMTJabatan.ToPgtypeDate(),
		TmtSk:                          arg.TMTKenaikanGajiBerkala.ToPgtypeDate(),
		PendidikanTerakhir:             pgtype.Text{String: arg.Pendidikan, Valid: arg.Pendidikan != ""},
		TanggalLulusPendidikanTerakhir: arg.TanggalLulus.ToPgtypeDate(),
		UnitKerjaIndukText:             ref.unitKerja.Nama,
		UnitKerjaIndukID:               pgtype.Text{String: arg.UnitKerjaIndukID, Valid: arg.UnitKerjaIndukID != ""},
		KantorPembayaran:               pgtype.Text{String: arg.KantorPembayaran, Valid: arg.KantorPembayaran != ""},
		Pejabat:                        pgtype.Text{String: arg.Pejabat, Valid: arg.Pejabat != ""},
		PegawaiNama:                    ref.pegawai.Nama,
		PegawaiNip:                     pgtype.Text{String: arg.NIP, Valid: true},
		TempatLahir:                    ref.pegawai.TempatLahir,
		TanggalLahir:                   ref.pegawai.TanggalLahir,
	})
	if err != nil {
		return 0, fmt.Errorf("repo create riwayat kenaikan gaji berkala: %w", err)
	}

	return id, nil
}

func (s *service) update(ctx context.Context, arg adminUpdateRequest) (bool, error) {
	ref, err := s.validateReferences(ctx, arg.NIP, arg.UnitKerjaIndukID, arg.GolonganID)
	if err != nil {
		if errors.Is(err, errPegawaiNotFound) {
			return false, errPegawaiNotFound
		}
		return false, err
	}

	affected, err := s.repo.UpdateRiwayatKenaikanGajiBerkala(ctx, repo.UpdateRiwayatKenaikanGajiBerkalaParams{
		ID:                             arg.ID,
		PegawaiID:                      pgtype.Int4{Int32: ref.pegawai.ID, Valid: true},
		GolonganID:                     pgtype.Int4{Int32: ref.golongan.ID, Valid: true},
		NGolRuang:                      ref.golongan.Nama,
		TmtGolongan:                    arg.TMTGolongan.ToPgtypeDate(),
		MasaKerjaGolonganTahun:         pgtype.Int2{Int16: arg.MasaKerjaGolonganTahun, Valid: true},
		MasaKerjaGolonganBulan:         pgtype.Int2{Int16: arg.MasaKerjaGolonganBulan, Valid: true},
		NoSk:                           pgtype.Text{String: arg.NomorSK, Valid: true},
		TanggalSk:                      arg.TanggalSK.ToPgtypeDate(),
		GajiPokok:                      pgtype.Int4{Int32: arg.GajiPokok, Valid: true},
		Jabatan:                        pgtype.Text{String: arg.Jabatan, Valid: arg.Jabatan != ""},
		TmtJabatan:                     arg.TMTJabatan.ToPgtypeDate(),
		TmtSk:                          arg.TMTKenaikanGajiBerkala.ToPgtypeDate(),
		PendidikanTerakhir:             pgtype.Text{String: arg.Pendidikan, Valid: arg.Pendidikan != ""},
		TanggalLulusPendidikanTerakhir: arg.TanggalLulus.ToPgtypeDate(),
		UnitKerjaIndukText:             ref.unitKerja.Nama,
		UnitKerjaIndukID:               pgtype.Text{String: arg.UnitKerjaIndukID, Valid: arg.UnitKerjaIndukID != ""},
		KantorPembayaran:               pgtype.Text{String: arg.KantorPembayaran, Valid: arg.KantorPembayaran != ""},
		Pejabat:                        pgtype.Text{String: arg.Pejabat, Valid: arg.Pejabat != ""},
	})
	if err != nil {
		return false, fmt.Errorf("repo update riwayat kenaikan gaji berkala: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, nip string, id int64) (bool, error) {
	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, errPegawaiNotFound
		}
		return false, fmt.Errorf("repo get pegawai: %w", err)
	}

	affected, err := s.repo.DeleteRiwayatKenaikanGajiBerkala(ctx, repo.DeleteRiwayatKenaikanGajiBerkalaParams{
		ID:        id,
		PegawaiID: pgtype.Int4{Int32: pegawai.ID, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo delete riwayat kenaikan gaji berkala: %w", err)
	}
	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int64, nip, fileBase64 string) (bool, error) {
	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, errPegawaiNotFound
		}
		return false, fmt.Errorf("repo get pegawai: %w", err)
	}

	affected, err := s.repo.UploadBerkasRiwayatKenaikanGajiBerkala(ctx, repo.UploadBerkasRiwayatKenaikanGajiBerkalaParams{
		ID:         id,
		PegawaiID:  pgtype.Int4{Int32: pegawai.ID, Valid: true},
		FileBase64: pgtype.Text{String: fileBase64, Valid: true},
	})
	if err != nil {
		return false, fmt.Errorf("repo upload berkas riwayat kenaikan gaji berkala: %w", err)
	}

	return affected > 0, nil
}

func (s *service) validateReferences(ctx context.Context, nip, unitKerjaIndukID string, golonganID int32) (*references, error) {
	var errs []error
	ref := references{}

	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get pegawai: %w", err)
		}
		return nil, errPegawaiNotFound
	}
	ref.pegawai = pegawai

	golongan, err := s.repo.GetRefGolongan(ctx, golonganID)
	if err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get ref golongan: %w", err)
		}
		errs = append(errs, errGolonganNotFound)
	}
	ref.golongan = golongan

	if unitKerjaIndukID != "" {
		unitKerja, err := s.repo.GetUnitKerja(ctx, unitKerjaIndukID)
		if err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("repo get unit kerja: %w", err)
			}
			errs = append(errs, errUnitKerjaNotFound)
		}
		ref.unitKerja = unitKerja
	}

	if len(errs) > 0 {
		return nil, api.NewMultiError(errs)
	}

	return &ref, nil
}
