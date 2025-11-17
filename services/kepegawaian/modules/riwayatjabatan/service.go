package riwayatjabatan

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
	ListRiwayatJabatan(ctx context.Context, arg repo.ListRiwayatJabatanParams) ([]repo.ListRiwayatJabatanRow, error)
	CountRiwayatJabatan(ctx context.Context, pnsNip string) (int64, error)
	GetBerkasRiwayatJabatan(ctx context.Context, arg repo.GetBerkasRiwayatJabatanParams) (pgtype.Text, error)
	GetPegawaiByNIP(ctx context.Context, nip string) (repo.GetPegawaiByNIPRow, error)
	GetRefJabatanByKode(ctx context.Context, kodeJabatan string) (repo.GetRefJabatanByKodeRow, error)
	GetRefJenisJabatan(ctx context.Context, id int32) (repo.GetRefJenisJabatanRow, error)
	GetUnitKerja(ctx context.Context, id string) (repo.GetUnitKerjaRow, error)

	CreateRiwayatJabatan(ctx context.Context, arg repo.CreateRiwayatJabatanParams) (int64, error)
	UpdateRiwayatJabatan(ctx context.Context, arg repo.UpdateRiwayatJabatanParams) (int64, error)
	DeleteRiwayatJabatan(ctx context.Context, arg repo.DeleteRiwayatJabatanParams) (int64, error)
	UploadBerkasRiwayatJabatan(ctx context.Context, arg repo.UploadBerkasRiwayatJabatanParams) (int64, error)
}

type service struct {
	repo repository
}

func newService(r repository) *service {
	return &service{repo: r}
}

func (s *service) list(ctx context.Context, nip string, limit, offset uint) ([]riwayatJabatan, int64, error) {
	data, err := s.repo.ListRiwayatJabatan(ctx, repo.ListRiwayatJabatanParams{
		Limit:  int32(limit),
		Offset: int32(offset),
		PnsNip: nip,
	})
	if err != nil {
		return nil, 0, fmt.Errorf("repo list: %w", err)
	}

	count, err := s.repo.CountRiwayatJabatan(ctx, nip)
	if err != nil {
		return nil, 0, fmt.Errorf("repo count: %w", err)
	}

	result := typeutil.Map(data, func(row repo.ListRiwayatJabatanRow) riwayatJabatan {
		return riwayatJabatan{
			ID:                      row.ID,
			JenisJabatanID:          row.JenisJabatanID,
			JenisJabatan:            row.JenisJabatan.String,
			NamaJabatan:             row.NamaJabatan.String,
			IDJabatan:               row.IDJabatan,
			TmtJabatan:              db.Date(row.TmtJabatan.Time),
			NoSk:                    row.NoSk.String,
			TanggalSk:               db.Date(row.TanggalSk.Time),
			SatuanKerjaID:           row.SatuanKerjaID,
			SatuanKerja:             row.SatuanKerja.String,
			UnitOrganisasiID:        row.UnitOrganisasiID,
			UnitOrganisasi:          row.UnitOrganisasi.String,
			StatusPlt:               row.StatusPlt.Bool,
			KelasJabatanID:          row.KelasJabatanID,
			KelasJabatan:            row.KelasJabatan.String,
			PeriodeJabatanStartDate: db.Date(row.PeriodeJabatanStartDate.Time),
			PeriodeJabatanEndDate:   db.Date(row.PeriodeJabatanEndDate.Time),
		}
	})

	return result, count, nil
}

func (s *service) getBerkas(ctx context.Context, nip string, id int64) (string, []byte, error) {
	res, err := s.repo.GetBerkasRiwayatJabatan(ctx, repo.GetBerkasRiwayatJabatanParams{
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

func (s *service) create(ctx context.Context, nip string, params upsertParams) (int64, error) {
	pegawai, err := s.repo.GetPegawaiByNIP(ctx, nip)
	if err != nil {
		return 0, errPegawaiNotFound
	}

	ref, err := s.validateReferences(ctx, params)
	if err != nil {
		return 0, err
	}

	id, err := s.repo.CreateRiwayatJabatan(ctx, repo.CreateRiwayatJabatanParams{
		PnsID:                   pgtype.Text{String: pegawai.PnsID, Valid: true},
		PnsNip:                  pgtype.Text{String: nip, Valid: true},
		PnsNama:                 pegawai.Nama,
		JenisJabatanID:          pgtype.Int4{Int32: typeutil.FromPtr(params.JenisJabatanID), Valid: params.JenisJabatanID != nil},
		JenisJabatan:            ref.jenisJabatan.Nama,
		JabatanID:               pgtype.Text{String: params.JabatanID, Valid: true},
		NamaJabatan:             ref.jabatan.Nama,
		JabatanIDBkn:            ref.jabatan.KodeBkn,
		SatuanKerjaID:           pgtype.Text{String: params.SatuanKerjaID, Valid: true},
		UnorID:                  pgtype.Text{String: params.UnitOrganisasiID, Valid: true},
		UnorIDBkn:               pgtype.Text{String: params.UnitOrganisasiID, Valid: true},
		Unor:                    ref.unitOrganisasi.Nama,
		TmtJabatan:              params.TMTJabatan.ToPgtypeDate(),
		NoSk:                    pgtype.Text{String: params.NoSK, Valid: true},
		TanggalSk:               params.TanggalSK.ToPgtypeDate(),
		StatusPlt:               pgtype.Bool{Bool: typeutil.FromPtr(params.StatusPlt), Valid: params.StatusPlt != nil},
		PeriodeJabatanStartDate: params.PeriodeJabatanStartDate.ToPgtypeDate(),
		PeriodeJabatanEndDate:   params.PeriodeJabatanEndDate.ToPgtypeDate(),
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

	affected, err := s.repo.UpdateRiwayatJabatan(ctx, repo.UpdateRiwayatJabatanParams{
		ID:                      id,
		Nip:                     nip,
		JenisJabatanID:          pgtype.Int4{Int32: typeutil.FromPtr(params.JenisJabatanID), Valid: params.JenisJabatanID != nil},
		JenisJabatan:            ref.jenisJabatan.Nama,
		JabatanID:               pgtype.Text{String: params.JabatanID, Valid: true},
		NamaJabatan:             ref.jabatan.Nama,
		JabatanIDBkn:            ref.jabatan.KodeBkn,
		SatuanKerjaID:           pgtype.Text{String: params.SatuanKerjaID, Valid: true},
		UnorID:                  pgtype.Text{String: params.UnitOrganisasiID, Valid: true},
		UnorIDBkn:               pgtype.Text{String: params.UnitOrganisasiID, Valid: true},
		Unor:                    ref.unitOrganisasi.Nama,
		TmtJabatan:              params.TMTJabatan.ToPgtypeDate(),
		NoSk:                    pgtype.Text{String: params.NoSK, Valid: true},
		TanggalSk:               params.TanggalSK.ToPgtypeDate(),
		StatusPlt:               pgtype.Bool{Bool: typeutil.FromPtr(params.StatusPlt), Valid: params.StatusPlt != nil},
		PeriodeJabatanStartDate: params.PeriodeJabatanStartDate.ToPgtypeDate(),
		PeriodeJabatanEndDate:   params.PeriodeJabatanEndDate.ToPgtypeDate(),
	})
	if err != nil {
		return false, fmt.Errorf("repo update: %w", err)
	}

	return affected > 0, nil
}

func (s *service) delete(ctx context.Context, id int64, nip string) (bool, error) {
	affected, err := s.repo.DeleteRiwayatJabatan(ctx, repo.DeleteRiwayatJabatanParams{
		ID:  id,
		Nip: nip,
	})
	if err != nil {
		return false, fmt.Errorf("repo delete: %w", err)
	}

	return affected > 0, nil
}

func (s *service) uploadBerkas(ctx context.Context, id int64, nip, fileBase64 string) (bool, error) {
	affected, err := s.repo.UploadBerkasRiwayatJabatan(ctx, repo.UploadBerkasRiwayatJabatanParams{
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
	jenisJabatan   repo.GetRefJenisJabatanRow
	jabatan        repo.GetRefJabatanByKodeRow
	unitOrganisasi repo.GetUnitKerjaRow
}

func (s *service) validateReferences(ctx context.Context, params upsertParams) (*references, error) {
	var err error
	var errs []error
	ref := references{}

	if params.JenisJabatanID != nil {
		if ref.jenisJabatan, err = s.repo.GetRefJenisJabatan(ctx, *params.JenisJabatanID); err != nil {
			if !errors.Is(err, pgx.ErrNoRows) {
				return nil, fmt.Errorf("repo get jenis jabatan: %w", err)
			}
			errs = append(errs, errJenisJabatanNotFound)
		}
	}

	if ref.jabatan, err = s.repo.GetRefJabatanByKode(ctx, params.JabatanID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get jabatan: %w", err)
		}
		errs = append(errs, errJabatanNotFound)
	}

	if _, err = s.repo.GetUnitKerja(ctx, params.SatuanKerjaID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get satuan kerja: %w", err)
		}
		errs = append(errs, errSatuanKerjaNotFound)
	}

	if ref.unitOrganisasi, err = s.repo.GetUnitKerja(ctx, params.UnitOrganisasiID); err != nil {
		if !errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("repo get unit organisasi: %w", err)
		}
		errs = append(errs, errUnitOrganisasiNotFound)
	}

	if len(errs) > 0 {
		return nil, api.NewMultiError(errs)
	}
	return &ref, nil
}
